package cameraClient

import (
	"time"
)

type Config interface {
	Name() string
	Address() string
	User() string
	Password() string
	RefreshInterval() time.Duration
	PreemptiveFetch() time.Duration
	ExpireEarly() time.Duration
	LogDebug() bool
}

type Client struct {
	// configuration
	config Config

	rtsp    rtspState
	raw     rawState
	delayed delayedState
	resize  resizeState
}

func RunClient(config Config) (*Client, error) {
	client := &Client{
		config:  config,
		rtsp:    createRtspState(),
		raw:     createRawState(),
		delayed: createDelayedState(),
		resize:  createResizeState(),
	}

	go client.rawImageRoutine()
	go client.delayedImageRoutine()
	go client.resizedImageRoutine()

	return client, nil
}

func (c *Client) Shutdown() {
	// send shutdown signals
	close(c.raw.shutdown)
	close(c.delayed.shutdown)
	close(c.resize.shutdown)
	// wait for all 3 go routines to send the closed signal
	<-c.raw.closed
	<-c.delayed.closed
	<-c.resize.closed
	c.rtsp.Close()
}

func (c *Client) Name() string {
	return c.config.Name()
}

func (c *Client) Config() Config {
	return c.config
}

func (c *Client) GetRawImage() *cameraPicture {
	response := make(chan *cameraPicture)
	c.raw.readRequestChannel <- rawImageReadRequest{response}
	return <-response
}

func (c *Client) GetDelayedImage(refreshInterval time.Duration) *cameraPicture {
	response := make(chan *cameraPicture)
	c.delayed.readRequestChannel <- delayedImageReadRequest{refreshInterval, response}
	return <-response
}

func (c *Client) GetResizedImage(refreshInterval time.Duration, dim Dimension, jpgQuality int) *cameraPicture {
	response := make(chan *cameraPicture)
	c.resize.readRequestChannel <- resizedImageReadRequest{
		resizedImageRequest{refreshInterval, dim, jpgQuality},
		response}
	return <-response
}
