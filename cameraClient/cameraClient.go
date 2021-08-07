package cameraClient

import (
	"github.com/google/uuid"
	"image"
	"net/http"
	"net/http/cookiejar"
	"time"
)

type Client struct {
	// configuration
	config Config

	// fetching
	httpClient *http.Client
	lastAuth   time.Time

	// interfacing
	rawImageReadRequestChannel chan rawImageReadRequest

	// image cache
	lastFetched time.Time
	lastUuid    uuid.UUID
	lastImg     image.Image
	lastErr     error
}

type CameraPicture interface {
	GetImg() image.Image
	GetUuid() uuid.UUID
}

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
		rawImageReadRequestChannel: make(chan rawImageReadRequest, 32),
	}

	go client.mainRoutine()

	return client, nil
}

func (c *Client) Shutdown() {}

func (c *Client) Name() string {
	return c.config.Name()
}

func (c *Client) Config() Config {
	return c.config
}

func (c *Client) GetRawImage() (cameraPicture CameraPicture, err error) {
	response := make(chan rawImageReadResponse)
	c.rawImageReadRequestChannel <- rawImageReadRequest{
		response: response,
	}
	r := <-response
	return r, r.err
}

func (c *Client) mainRoutine() {
	for {
		select {
		case readRequest := <-c.rawImageReadRequestChannel:
			c.handleRawImageReadRequest(readRequest)
		}
	}
}
