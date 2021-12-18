package cameraClient

import (
	"image"
	"strconv"
)

type Dimension interface {
	Width() int
	Height() int
}

type dimension struct {
	width  int
	height int
}

func DimensionCacheKey(dim Dimension) string {
	return strconv.Itoa(dim.Width()) + "x" + strconv.Itoa(dim.Height())
}

func DimensionOfImage(img image.Image) Dimension {
	if img == nil {
		return dimension{0, 0}
	}

	bounds := img.Bounds()
	return dimension{
		width:  bounds.Dx(),
		height: bounds.Dy(),
	}
}

func (d dimension) Width() int {
	return d.width
}

func (d dimension) Height() int {
	return d.height
}
