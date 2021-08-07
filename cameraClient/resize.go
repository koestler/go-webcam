package cameraClient

import (
	"bufio"
	"bytes"
	"github.com/disintegration/imaging"
	"image/jpeg"
)

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
		width = min(inputWidth, requestedWidth)
		height = 0
	} else {
		width = 0
		height = min(inputHeight, requestedHeight)
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

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
