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
		var img []byte
		err := rawImg.Err()
		if err == nil {
			img, err = imageResize(rawImg.Img(), dim.Width(), dim.Height())
		}

		resizedImage := &cameraPicture{
			img:     img,
			fetched: rawImg.Fetched(),
			expires: rawImg.Expires(),
			uuid:    rawImg.Uuid(),
			err:     err,
		}
		log.Printf("cameraClient[%s]: image resized, cacheKey=%s", c.Name(), cacheKey)

		request.response <- resizedImage
		c.resizeCache[cacheKey] = resizedImage
	}
}
