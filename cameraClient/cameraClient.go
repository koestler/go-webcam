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
	readRequestChannel chan readRequest

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

type readResponse struct {
	img  image.Image
	uuid uuid.UUID
	err  error
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

func (rr readResponse) GetImg() image.Image {
	return rr.img
}

func (rr readResponse) GetUuid() uuid.UUID {
	return rr.uuid
}

type readRequest struct {
	response chan readResponse
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
		readRequestChannel: make(chan readRequest, 32),
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
	response := make(chan readResponse)
	c.readRequestChannel <- readRequest{
		response: response,
	}
	r := <-response
	return r, r.err
}

func (c *Client) mainRoutine() {
	for {
		select {
		case readRequest := <-c.readRequestChannel:
			// fetch new image every RefreshInterval
			if time.Now().After(c.lastFetched.Add(c.Config().RefreshInterval())) {
				c.lastUuid = uuid.New()
				c.lastFetched = time.Now()
				c.lastImg, c.lastErr = c.ubntGetRawImage()
			}

			// return current image / uuid / error
			readRequest.response <- readResponse{
				img:  c.lastImg,
				uuid: c.lastUuid,
				err:  c.lastErr,
			}
		}
	}
}
