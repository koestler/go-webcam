package cameraClient

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"image/jpeg"
	"log"
	"time"
)

type resizedImageRequest struct {
	refreshInterval time.Duration
	dim             Dimension
}

type resizedImageReadRequest struct {
	resizedImageRequest
	response chan *sizedCameraPicture
}

type resizedImageComputeResponse struct {
	cacheKey     string
	resizedImage *sizedCameraPicture
}

func (c *Client) resizedImageRoutine() {
	for {
		select {
		case readRequest := <-c.resizedImageReadRequestChannel:
			c.handleResizedImageReadRequest(readRequest)
		case computeResponse := <-c.resizeComputeResponseChannel:
			c.handleResizeComputeResponse(computeResponse)
		}
	}
}

func (c *Client) handleResizedImageReadRequest(request resizedImageReadRequest) {
	cacheKey := request.computeCacheKey()
	c.resizeCache.purgeExpired(-c.Config().ExpireEarly())
	if cp, ok := c.resizeCache[cacheKey]; ok {
		log.Printf("cameraClient[%s]: resize image cache HIT, cacheKey=%s", c.Name(), cacheKey)
		request.response <- cp
	} else {
		log.Printf("cameraClient[%s]: resize image cache MISS, cacheKey=%s", c.Name(), cacheKey)
		if responses, ok := c.resizeWaitingResponses[cacheKey]; ok {
			log.Printf("cameraClient[%s]: resizeWaitingResponses HIT, cacheKey=%s", c.Name(), cacheKey)
			c.resizeWaitingResponses[cacheKey] = append(responses, request.response)
		} else {
			log.Printf("cameraClient[%s]: resizeWaitingResponses MISS, cacheKey=%s", c.Name(), cacheKey)
			c.resizeWaitingResponses[cacheKey] = []chan *sizedCameraPicture{request.response}
			go c.resizeOperation(request.resizedImageRequest)
		}
	}
}

func (c *Client) handleResizeComputeResponse(response resizedImageComputeResponse) {
	// response to all pending read requests
	for _, r := range c.resizeWaitingResponses[response.cacheKey] {
		r <- response.resizedImage
	}
	delete(c.resizeWaitingResponses, response.cacheKey)

	// add new image to cache
	c.resizeCache[response.cacheKey] = response.resizedImage
}

func (c *Client) resizeOperation(request resizedImageRequest) {
	delayedImg := c.GetDelayedImage(request.refreshInterval)

	log.Printf("cameraClient[%s]: resizeOperation(%s) start", c.Name(), request.computeCacheKey())

	var img []byte
	var imgDim dimension
	err := delayedImg.Err()
	if err == nil {
		dim := request.dim
		img, imgDim, err = imageResize(delayedImg.Img(), dim.Width(), dim.Height())
		time.Sleep(2 * time.Second)
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

	log.Printf("cameraClient[%s]: resizeOperation(%s) finish", c.Name(), request.computeCacheKey())
	c.resizeComputeResponseChannel <- resizedImageComputeResponse{
		request.computeCacheKey(),
		resizedImage,
	}
}

func (request resizedImageRequest) computeCacheKey() string {
	return fmt.Sprintf("%s-%s", request.refreshInterval.String(), dimensionCacheKey(request.dim))
}

func imageResize(inpImg []byte, requestedWidth, requestedHeight int) (oupImg []byte, oupDim dimension, err error) {
	if len(inpImg) == 0 {
		return inpImg, oupDim, nil
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
		return inpImg, oupDim, nil
	}

	resizedImp := imaging.Resize(decodedImg, width, height, imaging.Box)

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	err = jpeg.Encode(w, resizedImp, &jpeg.Options{Quality: 90})
	if err != nil {
		return
	}

	return b.Bytes(), dimensionFromBound(resizedImp.Bounds()), nil
}

func minInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}
