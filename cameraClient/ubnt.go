package cameraClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
)

func (c *Client) ubntLogin(force bool) (err error) {
	if force {
		c.authenticated = false
	}

	if c.authenticated {
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
		return fmt.Errorf("got code %d from camera during ubntLogin", res.StatusCode)
	}

	// we are authenticated
	c.authenticated = true
	log.Printf("cameraClient[%s]: ubntLogin successful", c.Name())

	return nil
}

func (c *Client) ubntGetRawImage() (img []byte, err error) {
	// ubntLogin
	err = c.ubntLogin(false)
	if err != nil {
		return
	}

	// create address
	imageUrl, err := url.Parse("http://" + path.Join(c.Config().Address(), "snap.jpeg"))
	if err != nil {
		return
	}

	// first attempt
	res, err := c.httpClient.Get(imageUrl.String())
	if err != nil {
		return
	}

	// second attempt; try to relogin if unauthorized is returned
	if res.StatusCode == http.StatusUnauthorized {
		// unauthorized -> reloggin
		err = c.ubntLogin(true)
		if err != nil {
			return
		}

		res, err = c.httpClient.Get(imageUrl.String())
		if err != nil {
			return
		}

	}

	// handle errors
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got code %v from camera when fetching a snapshot", res.StatusCode)
	}

	// read and return body
	defer res.Body.Close()
	img, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return img, nil
}
