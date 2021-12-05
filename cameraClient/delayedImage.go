package cameraClient

import (
	"log"
	"time"
)

type delayedImageReadRequest struct {
	refreshInterval time.Duration
	response        chan *cameraPicture
}

func (c *Client) delayedImageRoutine() {
	for {
		readRequest := <-c.delayedImageReadRequestChannel
		c.handleDelayedImageReadRequest(readRequest)
	}
}

func (c *Client) handleDelayedImageReadRequest(request delayedImageReadRequest) {
	refreshInterval := request.refreshInterval

	cacheKey := refreshInterval.String()

	// expire images 50ms early
	// this ensures that always a new image is fetched during periodic reloads with a jitter of up to 50ms
	c.delayedCache.purgeExpired(-50 * time.Millisecond)

	if cp, ok := c.delayedCache[cacheKey]; ok {
		log.Printf("cameraClient[%s]: delayed image cache HIT, cacheKey=%s", c.Name(), cacheKey)
		request.response <- cp
	} else {
		rawImg := c.GetRawImage()

		delayedImage := &cameraPicture{
			img:     rawImg.Img(),
			fetched: rawImg.Fetched(),
			expires: laterTime(rawImg.Expires(), rawImg.Fetched().Add(refreshInterval)),
			uuid:    rawImg.Uuid(),
			err:     rawImg.Err(),
		}

		log.Printf("cameraClient[%s]: delayed image cache MISS, updated, cacheKey=%s", c.Name(), cacheKey)

		request.response <- delayedImage
		c.delayedCache[cacheKey] = delayedImage
	}
}

func laterTime(x, y time.Time) time.Time {
	if x.After(y) {
		return x
	}
	return y
}
