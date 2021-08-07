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

type Dimension interface {
	Width() int
	Height() int
}

type Client struct {
	// configuration
	config Config

	// fetching
	httpClient *http.Client
	lastAuth   time.Time

	// interfacing
	rawImageReadRequestChannel     chan rawImageReadRequest
	resizedImageReadRequestChannel chan resizedImageReadRequest

	// raw image
	raw cameraPicture

	// resized image cache
	resizeCache map[string]*cameraPicture
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
		resizeCache:                    make(map[string]*cameraPicture),
	}

	go client.rawImageRoutine()
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
	c.rawImageReadRequestChannel <- rawImageReadRequest{
		response: response,
	}
	return <-response
}

func (c *Client) GetResizedImage(dim Dimension) CameraPicture {
	response := make(chan *cameraPicture)
	c.resizedImageReadRequestChannel <- resizedImageReadRequest{
		dim:      dim,
		response: response,
	}
	return <-response
}
