package cameraClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"time"
)

type Client struct {
	config     Config
	httpClient *http.Client
	lastAuth time.Time
}



type Config interface {
	Name() string
	Address() string
	User() string
	Password() string
	RefreshInterval() time.Duration
}

func RunClient(config Config) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &Client{
		config: config,
		httpClient: &http.Client{
			Jar: jar,
			// this tool is designed to serve cameras running on the local network
			// -> us a relatively short timeout
			Timeout: 10 * time.Second,
		},
	}, nil
}

func (c *Client) Shutdown() {}

func (c *Client) Name() string {
	return c.config.Name()
}

func (c *Client) Config() Config {
	return c.config
}

func (c *Client) login() (err error) {
	if time.Now().Before(c.lastAuth.Add(12 * time.Hour)) {
		// nothing todo, login was done recently
		return
	}

	// create address
	loginUrl, err := url.Parse("http://" + path.Join(c.Config().Address(), "api/1.1/login"))
	if err != nil {
		return
	}

	// create payload
	body := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		c.Config().User(),
		c.Config().Password(),
	}

	bodyJson, err := json.Marshal(body)
	if err != nil {
		return
	}

	resp, err := c.httpClient.Post(loginUrl.String(), "application/json", bytes.NewBuffer(bodyJson))
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("got code %v from camera", resp.StatusCode)
	}

	// set auth time to now
	c.lastAuth = time.Now()
	log.Printf("cameraClient[%s]: login successful", c.Name())

	return nil
}

func (c *Client) GetRawImage() (image []byte, err error) {
	err = c.login()
	if err != nil {
		return
	}

	return nil, errors.New("image not available for camera " + c.Name())
}
