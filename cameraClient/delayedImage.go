package cameraClient

import (
	"log"
	"time"
)

type delayedImageReadRequest struct {
	refreshInterval time.Duration
	dim             Dimension
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
	dim := request.dim

	cacheKey := dimensionCacheKey(dim) + "-" + refreshInterval.String()

	c.delayedCache.purgeExpired()

	if cp, ok := c.delayedCache[cacheKey]; ok {
		log.Printf("cameraClient[%s]: delayed image cache HIT, cacheKey=%s", c.Name(), cacheKey)
		request.response <- cp
	} else {
		resizedImg := c.GetResizedImage(dim)

		delayedImage := &cameraPicture{
			img:     resizedImg.Img(),
			fetched: resizedImg.Fetched(),
			expires: laterTime(resizedImg.Expires(), resizedImg.Fetched().Add(refreshInterval)),
			uuid:    resizedImg.Uuid(),
			err:     resizedImg.Err(),
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
