package cameraClient

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
	"log"
	"time"
)

type resizeState struct {
	readRequestChannel chan resizedImageReadRequest

	cache                  cameraPictureMap
	computeResponseChannel chan resizedImageComputeResponse
	waitingResponses       map[string][]chan *cameraPicture
}

type resizedImageRequest struct {
	refreshInterval time.Duration
	dim             Dimension
}

type resizedImageReadRequest struct {
	resizedImageRequest
	response chan *cameraPicture
}

type resizedImageComputeResponse struct {
	cacheKey     string
	resizedImage *cameraPicture
}

func createResizeState() resizeState {
	return resizeState{
		readRequestChannel:     make(chan resizedImageReadRequest, 16),
		cache:                  make(cameraPictureMap),
		computeResponseChannel: make(chan resizedImageComputeResponse, 16),
		waitingResponses:       make(map[string][]chan *cameraPicture),
	}
}

func (c *Client) resizedImageRoutine() {
	for {
		select {
		case readRequest := <-c.resize.readRequestChannel:
			c.handleResizedImageReadRequest(readRequest)
		case computeResponse := <-c.resize.computeResponseChannel:
			c.handleResizeComputeResponse(computeResponse)
		}
	}
}

func (c *Client) handleResizedImageReadRequest(request resizedImageReadRequest) {
	cacheKey := request.computeCacheKey()
	c.resize.cache.purgeExpired(-c.Config().ExpireEarly())
	if cp, ok := c.resize.cache[cacheKey]; ok {
		log.Printf("cameraClient[%s]: resize image cache HIT, cacheKey=%s, expiresIn=%s", c.Name(), cacheKey, time.Until(cp.expires))
		request.response <- cp
	} else {
		log.Printf("cameraClient[%s]: resize image cache MISS, cacheKey=%s", c.Name(), cacheKey)
		if responses, ok := c.resize.waitingResponses[cacheKey]; ok {
			log.Printf("cameraClient[%s]: waitingResponses HIT, cacheKey=%s", c.Name(), cacheKey)
			c.resize.waitingResponses[cacheKey] = append(responses, request.response)
		} else {
			log.Printf("cameraClient[%s]: waitingResponses MISS, cacheKey=%s", c.Name(), cacheKey)
			c.resize.waitingResponses[cacheKey] = []chan *cameraPicture{request.response}
			go c.resizeOperation(request.resizedImageRequest)
		}
	}
}

func (c *Client) handleResizeComputeResponse(response resizedImageComputeResponse) {
	// response to all pending read requests
	for _, r := range c.resize.waitingResponses[response.cacheKey] {
		r <- response.resizedImage
	}
	delete(c.resize.waitingResponses, response.cacheKey)

	// add new image to cache
	c.resize.cache[response.cacheKey] = response.resizedImage
}

func (c *Client) resizeOperation(request resizedImageRequest) {
	delayedImg := c.GetDelayedImage(request.refreshInterval)

	log.Printf("cameraClient[%s]: resizeOperation(%s) start", c.Name(), request.computeCacheKey())

	var oupJpgImg []byte
	var oupDecodedImg image.Image
	err := delayedImg.Err()
	if err == nil {
		oupJpgImg, oupDecodedImg, err = imageResize(delayedImg.JpgImg(), delayedImg.DecodedImg(), request.dim)
	}

	resizedImage := &cameraPicture{
		jpgImg:     oupJpgImg,
		decodedImg: oupDecodedImg,
		fetched:    delayedImg.Fetched(),
		expires:    delayedImg.Expires(),
		uuid:       delayedImg.Uuid(),
		err:        err,
	}

	log.Printf("cameraClient[%s]: resizeOperation(%s) finish", c.Name(), request.computeCacheKey())
	c.resize.computeResponseChannel <- resizedImageComputeResponse{
		request.computeCacheKey(),
		resizedImage,
	}
}

func (request resizedImageRequest) computeCacheKey() string {
	return fmt.Sprintf("%s-%s", request.refreshInterval.String(), DimensionCacheKey(request.dim))
}

func imageResize(
	inpJpgImg []byte, inpDecodedImg image.Image, requestedDim Dimension,
) (oupJpgImg []byte, oupDecodedImg image.Image, err error) {
	if inpDecodedImg == nil {
		return inpJpgImg, inpDecodedImg, nil
	}

	inpDim := DimensionOfImage(inpDecodedImg)
	var width, height int
	if requestedDim.Width()*inpDim.Height()/inpDim.Width() < requestedDim.Height() {
		width = minInt(inpDim.Width(), requestedDim.Width())
		height = 0
	} else {
		width = 0
		height = minInt(inpDim.Height(), requestedDim.Height())
	}

	if inpDim.Width() == width || inpDim.Height() == height {
		return inpJpgImg, inpDecodedImg, nil
	}

	t := time.Now()
	resizedImp := imaging.Resize(inpDecodedImg, width, height, imaging.Box)
	log.Printf("jpg resize took: %s", time.Since(t))

	var b bytes.Buffer
	w := bufio.NewWriter(&b)

	t = time.Now()
	err = jpeg.Encode(w, resizedImp, &jpeg.Options{Quality: 90})
	log.Printf("jpg encode took: %s", time.Since(t))

	if err != nil {
		return
	}

	return b.Bytes(), resizedImp, nil
}

func minInt(x, y int) int {
	if x < y {
		return x
	}
	return y
}
