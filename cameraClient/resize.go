package cameraClient

import (
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
)

func (c *Client) GetResizedImage(dim Dimension) (img image.Image, err error) {
	imgReader, err := c.ubntGetRawImageReader()
	if err != nil {
		return nil, err
	}
	defer imgReader.Close()

	img, err = jpeg.Decode(imgReader)
	if err != nil {
		return nil, err
	}
	return imageResize(img, dim.Width(), dim.Height()), nil
}

func imageResize(inpImg image.Image, requestedWidth, requestedHeight int) (oupImg image.Image) {
	inputWidth := inpImg.Bounds().Max.X - inpImg.Bounds().Min.X
	inputHeight := inpImg.Bounds().Max.Y - inpImg.Bounds().Min.Y

	var width, height int

	if requestedWidth * inputHeight / inputWidth < requestedHeight {
		width = min(inputWidth, requestedWidth)
		height = 0
	} else {
		width = 0
		height = min(inputHeight, requestedHeight)
	}

	return imaging.Resize(inpImg, width, height, imaging.Box)
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

