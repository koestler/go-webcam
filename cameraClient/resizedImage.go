package cameraClient

import (
	"log"
	"strconv"
)

type resizedImageReadRequest struct {
	dim      Dimension
	response chan *cameraPicture
}

func (c *Client) resizedImageRoutine() {
	for {
		readRequest := <-c.resizedImageReadRequestChannel
		c.handleResizedImageReadRequest(readRequest)
	}
}

func (c *Client) purgeResizeCache() {
	c.resizeCache = make(map[string]*cameraPicture)
}

func (c *Client) handleResizedImageReadRequest(request resizedImageReadRequest) {
	dim := request.dim
	cacheKey := strconv.Itoa(dim.Width()) + "x" + strconv.Itoa(dim.Height())

	rawImg := c.GetRawImage()
	if cp, ok := c.resizeCache[cacheKey]; ok && cp.Uuid() == rawImg.Uuid() {
		request.response <- cp
	} else {
		resizedImage := &cameraPicture{
			img:     imageResize(rawImg.Img(), dim.Width(), dim.Height()),
			fetched: rawImg.Fetched(),
			uuid:    rawImg.Uuid(),
			err:     rawImg.Err(),
		}
		log.Printf("cameraClient[%s]: image resized", c.Name())

		request.response <- resizedImage
		c.resizeCache[cacheKey] = resizedImage
	}
}
