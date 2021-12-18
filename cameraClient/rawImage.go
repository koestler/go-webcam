package cameraClient

import (
	"bytes"
	"github.com/google/uuid"
	"image"
	"image/jpeg"
	"log"
	"time"
)

type rawImageReadRequest struct {
	response chan *cameraPicture
}

func (c *Client) rawImageRoutine() {
	for {
		readRequest := <-c.rawImageReadRequestChannel
		// per camera, maximum one fetch operation is in progress
		c.handleRawImageReadRequest(readRequest)
	}
}

func (c *Client) handleRawImageReadRequest(request rawImageReadRequest) {
	// fetch new image every RefreshInterval
	if c.raw.Expired(-c.Config().ExpireEarly()) {
		log.Printf("cameraClient[%s]: raw image cache MISS", c.Name())

		log.Printf("cameraClient[%s]: GetRawImage start", c.Name())

		rawImg, err := c.ubntGetRawImage()
		time.Sleep(time.Second)

		var decodedRawImg image.Image

		t := time.Now()
		if err == nil {
			decodedRawImg, err = jpeg.Decode(bytes.NewReader(rawImg))
		}
		log.Printf("jpg decode took: %s", time.Since(t))
		log.Printf("cameraClient[%s]: GetRawImage finis", c.Name())

		now := time.Now()
		c.raw = cameraPicture{
			jpgImg:     rawImg,
			decodedImg: decodedRawImg,
			fetched:    now,
			expires:    now.Add(c.Config().RefreshInterval()),
			uuid:       uuid.New().String(),
			err:        err,
		}
	} else {
		log.Printf("cameraClient[%s]: raw image cache HIT", c.Name())
	}

	request.response <- &c.raw
}
