package cameraClient

import (
	"time"
)

type CameraPicture interface {
	Img() []byte
	Fetched() time.Time
	Uuid() string
	Err() error
}

type cameraPicture struct {
	img     []byte
	fetched time.Time
	uuid    string
	err     error
}

func (cp cameraPicture) Img() []byte {
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
