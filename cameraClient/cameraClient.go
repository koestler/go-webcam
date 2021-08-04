package cameraClient

import (
	"github.com/pkg/errors"
	"time"
)

type Client struct {
	config      Config
	lastFetched time.Time
}

type Config interface {
	Name() string
	Address() string
	User() string
	Password() string
	RefreshInterval() time.Duration
}

func RunClient(config Config) (*Client, error) {
	return &Client{
		config:      config,
		lastFetched: time.Time{},
	}, nil
}

func (c *Client) Shutdown() {}

func (c *Client) Name() string {
	return c.config.Name()
}

func (c *Client) Config() Config {
	return c.config
}

func (c *Client) GetRawImage() ([]byte, error) {
	return nil, errors.New("image not available for camera " + c.Name())
}
