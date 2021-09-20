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
	// expire images 50ms early
	// this ensures that always a new image is fetched during periodic reloads with a jitter of up to 50ms
	return time.Now().Add(50 * time.Millisecond).After(cp.expires)
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
