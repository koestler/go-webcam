package cameraClient

import (
	"time"
)

type CameraPicture interface {
	Img() []byte
	Fetched() time.Time
	Expires() time.Time
	Uuid() string
	Err() error
}

type SizedCameraPicture interface {
	CameraPicture
	Dimension() dimension
}

type cameraPicture struct {
	img     []byte
	fetched time.Time
	expires time.Time
	uuid    string
	err     error
}

type sizedCameraPicture struct {
	cameraPicture
	dimension dimension
}

type cameraPictureMap map[string]*cameraPicture
type sizedCameraPictureMap map[string]*sizedCameraPicture

func (cp cameraPicture) Img() []byte {
	return cp.img
}

func (cp sizedCameraPicture) Dimension() dimension {
	return cp.dimension
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

func (m sizedCameraPictureMap) purgeExpired(delay time.Duration) {
	for k, e := range m {
		if e.Expired(delay) {
			delete(m, k)
		}
	}
}
