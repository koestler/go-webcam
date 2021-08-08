package cameraClient

import (
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Config interface {
	Name() string
	Address() string
	User() string
	Password() string
	RefreshInterval() time.Duration
}

type Client struct {
	// configuration
	config Config

	// interfacing
	rawImageReadRequestChannel     chan rawImageReadRequest
	resizedImageReadRequestChannel chan resizedImageReadRequest
	delayedImageReadRequestChannel chan delayedImageReadRequest

	// fetching
	httpClient    *http.Client
	authenticated bool

	// raw image
	raw cameraPicture

	// resized image cache
	resizeCache  cameraPictureMap
	delayedCache cameraPictureMap
}

func RunClient(config Config) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	client := &Client{
		config: config,
		httpClient: &http.Client{
			Jar: jar,
			// this tool is designed to serve cameras running on the local network
			// -> us a relatively short timeout
			Timeout: 10 * time.Second,
		},
		rawImageReadRequestChannel:     make(chan rawImageReadRequest, 16),
		resizedImageReadRequestChannel: make(chan resizedImageReadRequest, 16),
		delayedImageReadRequestChannel: make(chan delayedImageReadRequest, 16),
		resizeCache:                    make(cameraPictureMap),
		delayedCache:                   make(cameraPictureMap),
	}

	go client.rawImageRoutine()
	go client.delayedImageRoutine()
	go client.resizedImageRoutine()

	return client, nil
}

func (c *Client) Shutdown() {}

func (c *Client) Name() string {
	return c.config.Name()
}

func (c *Client) Config() Config {
	return c.config
}

func (c *Client) GetRawImage() CameraPicture {
	response := make(chan *cameraPicture)
	c.rawImageReadRequestChannel <- rawImageReadRequest{response}
	return <-response
}

func (c *Client) GetDelayedImage(refreshInterval time.Duration) CameraPicture {
	response := make(chan *cameraPicture)
	c.delayedImageReadRequestChannel <- delayedImageReadRequest{refreshInterval, response}
	return <-response
}

func (c *Client) GetResizedImage(refreshInterval time.Duration, dim Dimension) CameraPicture {
	response := make(chan *cameraPicture)
	c.resizedImageReadRequestChannel <- resizedImageReadRequest{refreshInterval, dim, response}
	return <-response
}
