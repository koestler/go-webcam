package cameraClient

import "strconv"

type Dimension interface {
	Width() int
	Height() int
}

func dimensionCacheKey(dim Dimension) string {
	return strconv.Itoa(dim.Width()) + "x" + strconv.Itoa(dim.Height())
}
