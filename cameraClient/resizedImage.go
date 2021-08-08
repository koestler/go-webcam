package cameraClient

import (
	"bufio"
	"bytes"
	"github.com/disintegration/imaging"
	"image/jpeg"
	"log"
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

func (c *Client) handleResizedImageReadRequest(request resizedImageReadRequest) {
	dim := request.dim
	cacheKey := dimensionCacheKey(dim)

	c.resizeCache.purgeExpired()

	if cp, ok := c.resizeCache[cacheKey]; ok {
		log.Printf("cameraClient[%s]: resize image cache HIT, cacheKey=%s", c.Name(), cacheKey)
		request.response <- cp
	} else {
		rawImg := c.GetRawImage()

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
		log.Printf("cameraClient[%s]: resize image cache MISS, computed, cacheKey=%s", c.Name(), cacheKey)

		request.response <- resizedImage
		c.resizeCache[cacheKey] = resizedImage
	}
}

func imageResize(inpImg []byte, requestedWidth, requestedHeight int) (oupImg []byte, err error) {
	if len(inpImg) == 0 {
		return inpImg, nil
	}

	decodedImg, err := jpeg.Decode(bytes.NewReader(inpImg))
	if err != nil {
		return
	}

	bounds := decodedImg.Bounds()
	inputWidth := bounds.Max.X - bounds.Min.X
	inputHeight := bounds.Max.Y - bounds.Min.Y

	var width, height int

	if requestedWidth*inputHeight/inputWidth < requestedHeight {
		width = minInt(inputWidth, requestedWidth)
		height = 0
	} else {
		width = 0
		height = minInt(inputHeight, requestedHeight)
	}

	if inputWidth == width || inputHeight == height {
		return inpImg, nil
	}

	resizedImp := imaging.Resize(decodedImg, width, height, imaging.Box)

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	err = jpeg.Encode(w, resizedImp, &jpeg.Options{Quality: 90})
	if err != nil {
		return
	}

	return b.Bytes(), nil
}

func minInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}
