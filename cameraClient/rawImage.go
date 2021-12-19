package cameraClient

import (
	"bytes"
	"github.com/google/uuid"
	"image"
	"image/jpeg"
	"log"
	"time"
)

type rawState struct {
	readRequestChannel chan rawImageReadRequest

	// img image
	img cameraPicture
}

type rawImageReadRequest struct {
	response chan *cameraPicture
}

func createRawState() rawState {
	return rawState{
		readRequestChannel: make(chan rawImageReadRequest, 16),
	}
}

func (c *Client) rawImageRoutine() {
	for {
		readRequest := <-c.raw.readRequestChannel
		// per camera, maximum one fetch operation is in progress
		c.handleRawImageReadRequest(readRequest)
	}
}

func (c *Client) handleRawImageReadRequest(request rawImageReadRequest) {
	// fetch new image every RefreshInterval
	if c.raw.img.Expired(-c.Config().ExpireEarly()) {
		if c.Config().LogDebug() {
			log.Printf("cameraClient[%s]: raw image cache MISS", c.Name())
		}

		rawImg, err := c.ubntGetRawImage()
		var decodedRawImg image.Image

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
	} else if c.Config().LogDebug() {
		log.Printf("cameraClient[%s]: raw image cache HIT, expiresIn=%s", c.Name(), time.Until(c.raw.img.expires))
	}

	request.response <- &c.raw.img
}
