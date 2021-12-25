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
	ExpireEarly() time.Duration
	LogDebug() bool
}

type Client struct {
	// configuration
	config Config

	ubnt    ubntState
	raw     rawState
	delayed delayedState
	resize  resizeState
}

func RunClient(config Config) (*Client, error) {
	client := &Client{
		config:  config,
		ubnt:    createUbntState(),
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
	// todo: implement proper shutdown
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
