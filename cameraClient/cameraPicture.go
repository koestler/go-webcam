package cameraClient

import (
	"image"
	"time"
)

type CameraPicture interface {
	Img() image.Image
	Fetched() time.Time
	Uuid() string
	Err() error
}

type cameraPicture struct {
	img     image.Image
	fetched time.Time
	uuid    string
	err     error
}

func (cp cameraPicture) Img() image.Image {
	return cp.img
}

func (cp cameraPicture) Fetched() time.Time {
	return cp.fetched
}

func (cp cameraPicture) Uuid() string {
	return cp.uuid
}

func (cp cameraPicture) Err() error {
	return cp.err
}
