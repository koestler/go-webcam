package cameraClient

import (
	"log"
	"time"
)

type delayedState struct {
	readRequestChannel chan delayedImageReadRequest
	cache              cameraPictureMap

	// shutdown handling
	shutdown chan struct{}
	closed   chan struct{}
}

type delayedImageReadRequest struct {
	refreshInterval time.Duration
	response        chan *cameraPicture
}

func createDelayedState() delayedState {
	return delayedState{
		readRequestChannel: make(chan delayedImageReadRequest, 16),
		cache:              make(cameraPictureMap),
		shutdown:           make(chan struct{}),
		closed:             make(chan struct{}),
	}
}

func (c *Client) delayedImageRoutine() {
	defer close(c.delayed.closed)
	for {
		select {
		case readRequest := <-c.delayed.readRequestChannel:
			c.handleDelayedImageReadRequest(readRequest)
		case <-c.delayed.shutdown:
			return
		}
	}
}

func (c *Client) handleDelayedImageReadRequest(request delayedImageReadRequest) {
	refreshInterval := request.refreshInterval
	cacheKey := refreshInterval.String()

	c.delayed.cache.purgeExpired(-c.Config().ExpireEarly())

	if cp, ok := c.delayed.cache[cacheKey]; ok {
		if c.Config().LogDebug() {
			log.Printf(
				"cameraClient[%s]: delayed image cache HIT, cacheKey=%s, expiresIn=%s",
				c.Name(), cacheKey,
				time.Until(cp.expires),
			)
		}
		request.response <- cp
	} else {
		if c.Config().LogDebug() {
			log.Printf("cameraClient[%s]: delayed image cache MISS, cacheKey=%s", c.Name(), cacheKey)
		}

		rawImg := c.GetRawImage()

		delayedImage := &cameraPicture{
			jpgImg:     rawImg.JpgImg(),
			decodedImg: rawImg.DecodedImg(),
			fetched:    rawImg.Fetched(),
			expires:    laterTime(rawImg.Expires(), rawImg.Fetched().Add(refreshInterval)),
			uuid:       rawImg.Uuid(),
			err:        rawImg.Err(),
		}

		request.response <- delayedImage
		c.delayed.cache[cacheKey] = delayedImage
	}
}

func laterTime(x, y time.Time) time.Time {
	if x.After(y) {
		return x
	}
	return y
}
