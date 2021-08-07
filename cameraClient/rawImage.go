package cameraClient

import (
	"github.com/google/uuid"
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
	if time.Now().After(c.raw.fetched.Add(c.Config().RefreshInterval())) {
		rawImg, err := c.ubntGetRawImage()

		c.raw = cameraPicture{
			img:     rawImg,
			fetched: time.Now(),
			uuid:    uuid.New().String(),
			err:     err,
		}

		c.purgeResizeCache()
	}

	request.response <- &c.raw
}
