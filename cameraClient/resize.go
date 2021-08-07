package cameraClient

import (
	"github.com/disintegration/imaging"
	"image"
)

func imageResize(inpImg image.Image, requestedWidth, requestedHeight int) (oupImg image.Image) {
	if inpImg == nil {
		return nil
	}

	inputWidth := inpImg.Bounds().Max.X - inpImg.Bounds().Min.X
	inputHeight := inpImg.Bounds().Max.Y - inpImg.Bounds().Min.Y

	var width, height int

	if requestedWidth*inputHeight/inputWidth < requestedHeight {
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
