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
	if time.Now().After(c.raw.expires) {
		rawImg, err := c.ubntGetRawImage()

		now := time.Now()

		c.raw = cameraPicture{
			img:     rawImg,
			fetched: now,
			expires: now.Add(c.Config().RefreshInterval()),
			uuid:    uuid.New().String(),
			err:     err,
		}

		c.purgeResizeCache()
	}

	request.response <- &c.raw
}
