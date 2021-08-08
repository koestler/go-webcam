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

type cameraPicture struct {
	img     []byte
	fetched time.Time
	expires time.Time
	uuid    string
	err     error
}

type cameraPictureMap map[string]*cameraPicture

func (cp cameraPicture) Img() []byte {
	return cp.img
}

func (cp cameraPicture) Fetched() time.Time {
	return cp.fetched
}

func (cp cameraPicture) Expires() time.Time {
	return cp.expires
}

func (cp cameraPicture) Expired() bool {
	return time.Now().After(cp.expires)
}

func (cp cameraPicture) Uuid() string {
	return cp.uuid
}

func (cp cameraPicture) Err() error {
	return cp.err
}

func (m cameraPictureMap) purgeExpired() {
	for k, e := range m {
		if e.Expired() {
			delete(m, k)
		}
	}
}
