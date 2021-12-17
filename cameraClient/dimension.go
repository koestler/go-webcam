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

func dimensionCacheKey(dim Dimension) string {
	return strconv.Itoa(dim.Width()) + "x" + strconv.Itoa(dim.Height())
}

func dimensionFromBound(bounds image.Rectangle) dimension {
	return dimension{
		width:  bounds.Max.X - bounds.Min.X,
		height: bounds.Max.Y - bounds.Min.Y,
	}
}

func (d dimension) Width() int {
	return d.width
}

func (d dimension) Height() int {
	return d.height
}
