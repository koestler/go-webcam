package cameraClient

import (
	"github.com/google/uuid"
	"image"
	"time"
)

type rawImageReadRequest struct {
	response chan rawImageReadResponse
}

type rawImageReadResponse struct {
	img  image.Image
	uuid uuid.UUID
	err  error
}

func (rr rawImageReadResponse) GetImg() image.Image {
	return rr.img
}

func (rr rawImageReadResponse) GetUuid() uuid.UUID {
	return rr.uuid
}

func (c *Client) handleRawImageReadRequest(request rawImageReadRequest) {
	// fetch new image every RefreshInterval
	if time.Now().After(c.lastFetched.Add(c.Config().RefreshInterval())) {
		c.lastUuid = uuid.New()
		c.lastFetched = time.Now()
		c.lastImg, c.lastErr = c.ubntGetRawImage()
	}

	// return current image / uuid / error
	request.response <- rawImageReadResponse{
		img:  c.lastImg,
		uuid: c.lastUuid,
		err:  c.lastErr,
	}
}