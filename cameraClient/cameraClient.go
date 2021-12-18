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
	ExpireEarly() time.Duration
}

type Client struct {
	// configuration
	config Config

	// fetching
	httpClient    *http.Client
	authenticated bool

	// processing
	raw     rawState
	delayed delayedState
	resize  resizeState
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
		raw:     createRawState(),
		delayed: createDelayedState(),
		resize:  createResizeState(),
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
	c.raw.readRequestChannel <- rawImageReadRequest{response}
	return <-response
}

func (c *Client) GetDelayedImage(refreshInterval time.Duration) CameraPicture {
	response := make(chan *cameraPicture)
	c.delayed.readRequestChannel <- delayedImageReadRequest{refreshInterval, response}
	return <-response
}

func (c *Client) GetResizedImage(refreshInterval time.Duration, dim Dimension) CameraPicture {
	response := make(chan *cameraPicture)
	c.resize.readRequestChannel <- resizedImageReadRequest{
		resizedImageRequest{refreshInterval, dim},
		response}
	return <-response
}
