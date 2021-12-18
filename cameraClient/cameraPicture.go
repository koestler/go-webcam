package cameraClient

import (
	"image"
	"time"
)

type CameraPicture interface {
	JpgImg() []byte
	DecodedImg() image.Image
	Fetched() time.Time
	Expires() time.Time
	Expired(delay time.Duration) bool
	Uuid() string
	Err() error
}

type cameraPicture struct {
	jpgImg     []byte
	decodedImg image.Image
	fetched    time.Time
	expires    time.Time
	uuid       string
	err        error
}

type cameraPictureMap map[string]*cameraPicture

func (cp cameraPicture) JpgImg() []byte {
	return cp.jpgImg
}

func (cp cameraPicture) DecodedImg() image.Image {
	return cp.decodedImg
}

func (cp cameraPicture) Dimension() Dimension {
	return DimensionOfImage(cp.decodedImg)
}

func (cp cameraPicture) Fetched() time.Time {
	return cp.fetched
}

func (cp cameraPicture) Expires() time.Time {
	return cp.expires
}

func (cp cameraPicture) Expired(delay time.Duration) bool {
	return time.Now().Add(-delay).After(cp.expires)
}

func (cp cameraPicture) Uuid() string {
	return cp.uuid
}

func (cp cameraPicture) Err() error {
	return cp.err
}

func (m cameraPictureMap) purgeExpired(delay time.Duration) {
	for k, e := range m {
		if e.Expired(delay) {
			delete(m, k)
		}
	}
}
