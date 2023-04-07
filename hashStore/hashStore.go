package hashStore

import (
	"github.com/koestler/go-webcam/cameraClient"
	"time"
)

type HashStore struct {
	config Config

	shutdown chan struct{}
	closed   chan struct{}

	setChannel chan setRequest
	getChannel chan getRequest

	storage map[string]value
}

type Config interface {
	HashTimeout() time.Duration
}

type setRequest struct {
	hash     string
	cp       cameraClient.CameraPicture
	response chan struct{}
}

type getRequest struct {
	hash     string
	response chan cameraClient.CameraPicture
}

type value struct {
	cp      cameraClient.CameraPicture
	touched time.Time
}

func Run(config Config) *HashStore {
	h := &HashStore{
		config:     config,
		shutdown:   make(chan struct{}),
		closed:     make(chan struct{}),
		setChannel: make(chan setRequest, 16),
		getChannel: make(chan getRequest, 16),
		storage:    make(map[string]value),
	}

	go h.worker()

	return h
}

func (h *HashStore) Shutdown() {
	// send remaining points
	close(h.shutdown)
	// wait for worker to shut down
	<-h.closed
}

func (h *HashStore) Set(hash string, cp cameraClient.CameraPicture) {
	response := make(chan struct{})
	h.setChannel <- setRequest{hash, cp, response}
	<-response
}

func (h *HashStore) Get(hash string) cameraClient.CameraPicture {
	response := make(chan cameraClient.CameraPicture)
	h.getChannel <- getRequest{hash, response}
	return <-response
}

func (h *HashStore) Config() Config {
	return h.config
}

func (h *HashStore) worker() {
	defer close(h.closed)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case setRequest := <-h.setChannel:
			if v, ok := h.storage[setRequest.hash]; ok {
				v.touched = time.Now()
			} else {
				h.storage[setRequest.hash] = value{
					cp:      setRequest.cp,
					touched: time.Now(),
				}
			}
			close(setRequest.response)
		case getRequest := <-h.getChannel:
			if v, ok := h.storage[getRequest.hash]; ok {
				getRequest.response <- v.cp
			} else {
				getRequest.response <- nil
			}
		case <-ticker.C:
			now := time.Now()
			for k, v := range h.storage {
				if v.touched.Add(h.config.HashTimeout()).Before(now) {
					delete(h.storage, k)
				}
			}
		case <-h.shutdown:
			return // shutdown
		}
	}
}
