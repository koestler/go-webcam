package cameraClient

import (
	"github.com/google/uuid"
	"log"
	"time"
)

type rawImageReadRequest struct {
	response chan *cameraPicture
}

func (c *Client) rawImageRoutine() {
	for {
		readRequest := <-c.rawImageReadRequestChannel
		c.handleRawImageReadRequest(readRequest)
	}
}

func (c *Client) handleRawImageReadRequest(request rawImageReadRequest) {
	// fetch new image every RefreshInterval
	if c.raw.Expired() {
		rawImg, err := c.ubntGetRawImage()

		log.Printf("cameraClient[%s]: raw image cache MISS, fetched", c.Name())

		now := time.Now()

		c.raw = cameraPicture{
			img:     rawImg,
			fetched: now,
			expires: now.Add(c.Config().RefreshInterval()),
			uuid:    uuid.New().String(),
			err:     err,
		}
	} else {
		log.Printf("cameraClient[%s]: raw image cache HIT", c.Name())
	}

	request.response <- &c.raw
}
