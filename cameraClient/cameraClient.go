package cameraClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
	"io"
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
	lastAuth   time.Time
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

	res, err := c.httpClient.Post(loginUrl.String(), "application/json", bytes.NewBuffer(bodyJson))
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("got code %d from camera during login", res.StatusCode)
	}

	// set auth time to now
	c.lastAuth = time.Now()
	log.Printf("cameraClient[%s]: login successful", c.Name())

	return nil
}

func (c *Client) GetRawImageReader() (imgReader io.ReadCloser, err error) {
	// login
	err = c.login()
	if err != nil {
		return
	}

	// create address
	imageUrl, err := url.Parse("http://" + path.Join(c.Config().Address(), "snap.jpeg"))
	if err != nil {
		return
	}

	res, err := c.httpClient.Get(imageUrl.String())
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("got code %v from camera when fetching a snapshot", res.StatusCode)
	}

	log.Printf("cameraClient[%s]: image fetched", c.Name())

	return res.Body, nil
}

func (c *Client) GetRawImage() (img []byte, err error) {
	imgReader, err := c.GetRawImageReader()
	if err != nil {
		return nil, err
	}
	defer imgReader.Close()
	return io.ReadAll(imgReader)
}

func (c *Client) GetImage(dim Dimension) (img image.Image, err error) {
	imgReader, err := c.GetRawImageReader()
	if err != nil {
		return nil, err
	}
	defer imgReader.Close()

	img, err = jpeg.Decode(imgReader)
	if err != nil {
		return nil, err
	}
	return imageResize(img, dim.Width(), dim.Height()), nil
}

func imageResize(inpImg image.Image, requestedWidth, requestedHeight int) (oupImg image.Image) {
	inputWidth := inpImg.Bounds().Max.X - inpImg.Bounds().Min.X
	inputHeight := inpImg.Bounds().Max.Y - inpImg.Bounds().Min.Y

	var width, height int

	if requestedWidth * inputHeight / inputWidth < requestedHeight {
		width = min(inputWidth, requestedWidth)
		height = 0
	} else {
		width = 0
		height = min(inputHeight, requestedHeight)
	}

	return imaging.Resize(inpImg, width, height, imaging.Box)
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
