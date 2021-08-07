package cameraClient

import (
	"github.com/disintegration/imaging"
	"image"
)

func (c *Client) GetResizedImage(dim Dimension) (img image.Image, err error) {
	rawImg, err := c.GetRawImage()
	if err != nil {
		return nil, err
	}
	return imageResize(rawImg, dim.Width(), dim.Height()), nil
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

	if inputWidth == width || inputHeight == height {
		return inpImg
	}

	return imaging.Resize(inpImg, width, height, imaging.Box)
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

