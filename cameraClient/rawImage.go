package cameraClient

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"
	"time"

	"github.com/google/uuid"
)

type rawState struct {
	readRequestChannel chan rawImageReadRequest

	// img image
	img cameraPicture

	preemptiveTickerRunning bool
	preemptiveTicker        *time.Ticker

	// shutdown handling
	shutdown chan struct{}
	closed   chan struct{}
}

type rawImageReadRequest struct {
	response chan *cameraPicture
}

func createRawState() rawState {
	// create a stopped ticker
	ticker := time.NewTicker(time.Hour)
	ticker.Stop()

	return rawState{
		readRequestChannel:      make(chan rawImageReadRequest, 16),
		preemptiveTickerRunning: false,
		preemptiveTicker:        ticker,
		shutdown:                make(chan struct{}),
		closed:                  make(chan struct{}),
	}
}

func (c *Client) rawImageRoutine() {
	defer close(c.raw.closed)

	cfg := c.Config()

	c.startPreemptiveTicker()
	defer c.raw.preemptiveTicker.Stop()

	// track when the last non-preemptive fetch was made
	lastFetch := time.Now()

	for {
		select {
		case readRequest := <-c.raw.readRequestChannel:
			// per camera, maximum one fetch operation is in progress
			c.handleRawImageReadRequest(readRequest)
			lastFetch = time.Now()

			// check if preemptive fetch needs to be started
			c.startPreemptiveTicker()
		case <-c.raw.preemptiveTicker.C:
			if cfg.LogDebug() {
				log.Printf("cameraClient[%s]: preemptive fetch", c.Name())
			}
			// trigger a new camera fetch whenever ticker goes off
			c.fetchImage()

			// check if preemptive fetch needs to be stopped
			if lastFetch.Add(cfg.PreemptiveFetch()).Before(time.Now()) {
				c.stopPreemptiveTicker()
			}
		case <-c.raw.shutdown:
			return
		}
	}
}

func (c *Client) isPreemptiveFetchEnabled() bool {
	cfg := c.Config()
	return cfg.RefreshInterval() > (50*time.Millisecond) && cfg.PreemptiveFetch() > cfg.RefreshInterval()
}

func (c *Client) startPreemptiveTicker() {
	if !c.isPreemptiveFetchEnabled() {
		// not enabled, do not start
		return
	}

	if c.raw.preemptiveTickerRunning {
		// already running, do not restart
		return
	}

	if c.Config().LogDebug() {
		log.Printf("cameraClient[%s]: start preemptive fetch", c.Name())
	}

	c.raw.preemptiveTickerRunning = true
	c.raw.preemptiveTicker.Reset(c.Config().RefreshInterval())
}

func (c *Client) stopPreemptiveTicker() {
	if !c.raw.preemptiveTickerRunning {
		// not running, do not stop
		return
	}

	if c.Config().LogDebug() {
		log.Printf("cameraClient[%s]: stop preemptive fetch", c.Name())
	}

	c.raw.preemptiveTickerRunning = false
	c.raw.preemptiveTicker.Stop()
}

func (c *Client) handleRawImageReadRequest(request rawImageReadRequest) {
	// fetch new image every RefreshInterval
	if c.raw.img.Expired(-c.Config().ExpireEarly()) {
		if c.Config().LogDebug() {
			log.Printf("cameraClient[%s]: raw image cache MISS", c.Name())
		}
		c.fetchImage()
	} else if c.Config().LogDebug() {
		log.Printf("cameraClient[%s]: raw image cache HIT, expiresIn=%s", c.Name(), time.Until(c.raw.img.expires))
	}

	request.response <- &c.raw.img
}

func (c *Client) fetchImage() {
	start0 := time.Now()

	rawImg, err := c.rtspGetRawImage()
	if err != nil {
		log.Printf("cameraClient[%s]: failed to fetch raw image: %v", c.Name(), err)
	}

	var decodedRawImg image.Image

	start1 := time.Now()

	if err == nil {
		decodedRawImg, err = jpeg.Decode(bytes.NewReader(rawImg))
	}

	now := time.Now()
	c.raw.img = cameraPicture{
		jpgImg:     rawImg,
		decodedImg: decodedRawImg,
		fetched:    now,
		expires:    now.Add(c.Config().RefreshInterval()),
		uuid:       uuid.New().String(),
		err:        err,
	}

	if c.Config().LogDebug() && decodedRawImg != nil {
		log.Printf(
			"cameraClient[%s]: raw image fetched, took=%.3fs, total=%.3fs, dim=%s",
			c.Name(),
			time.Since(start1).Seconds(),
			time.Since(start0).Seconds(),
			DimensionCacheKey(DimensionOfImage(decodedRawImg)),
		)
	}
}
