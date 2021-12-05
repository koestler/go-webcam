package cameraClient

import (
	"bufio"
	"bytes"
	"github.com/disintegration/imaging"
	"image/jpeg"
	"log"
	"time"
)

type resizedImageReadRequest struct {
	refreshInterval time.Duration
	dim             Dimension
	response        chan *sizedCameraPicture
}

func (c *Client) resizedImageRoutine() {
	for {
		readRequest := <-c.resizedImageReadRequestChannel
		c.handleResizedImageReadRequest(readRequest)
	}
}

func (c *Client) handleResizedImageReadRequest(request resizedImageReadRequest) {
	refreshInterval := request.refreshInterval
	dim := request.dim
	cacheKey := refreshInterval.String() + "-" + dimensionCacheKey(dim)

	// expire images 50ms early
	// this ensures that always a new image is fetched during periodic reloads with a jitter of up to 50ms
	c.resizeCache.purgeExpired(-50 * time.Millisecond)

	if cp, ok := c.resizeCache[cacheKey]; ok {
		log.Printf("cameraClient[%s]: resize image cache HIT, cacheKey=%s", c.Name(), cacheKey)
		request.response <- cp
	} else {
		delayedImg := c.GetDelayedImage(refreshInterval)

		var img []byte
		var imgDim dimension
		err := delayedImg.Err()
		if err == nil {
			img, imgDim, err = imageResize(delayedImg.Img(), dim.Width(), dim.Height())
		}

		resizedImage := &sizedCameraPicture{
			cameraPicture: cameraPicture{
				img:     img,
				fetched: delayedImg.Fetched(),
				expires: delayedImg.Expires(),
				uuid:    delayedImg.Uuid(),
				err:     err,
			},
			dimension: imgDim,
		}

		log.Printf("cameraClient[%s]: resize image cache MISS, computed, cacheKey=%s", c.Name(), cacheKey)

		request.response <- resizedImage
		c.resizeCache[cacheKey] = resizedImage
	}
}

func imageResize(inpImg []byte, requestedWidth, requestedHeight int) (oupImg []byte, oupDim dimension, err error) {
	if len(inpImg) == 0 {
		return inpImg, oupDim,nil
	}

	decodedImg, err := jpeg.Decode(bytes.NewReader(inpImg))
	if err != nil {
		return
	}

	inpDim := dimensionFromBound(decodedImg.Bounds())

	var width, height int
	if requestedWidth*inpDim.height/inpDim.width < requestedHeight {
		width = minInt(inpDim.width, requestedWidth)
		height = 0
	} else {
		width = 0
		height = minInt(inpDim.height, requestedHeight)
	}

	if inpDim.width == width || inpDim.height == height {
		return inpImg, oupDim,nil
	}

	resizedImp := imaging.Resize(decodedImg, width, height, imaging.Box)

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	err = jpeg.Encode(w, resizedImp, &jpeg.Options{Quality: 90})
	if err != nil {
		return
	}

	return b.Bytes(), dimensionFromBound(resizedImp.Bounds()),nil
}

func minInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}
