package cameraClient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/url"
	"path"
	"time"
)

func (c *Client) ubntLogin() (err error) {
	if time.Now().Before(c.lastAuth.Add(12 * time.Hour)) {
		// nothing todo, ubntLogin was done recently
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

	// set auth time to now
	c.lastAuth = time.Now()
	log.Printf("cameraClient[%s]: ubntLogin successful", c.Name())

	return nil
}

func (c *Client) ubntGetRawImage() (img image.Image, err error) {
	// ubntLogin
	err = c.ubntLogin()
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
	defer res.Body.Close()

	decodedImg, err := jpeg.Decode(res.Body)
	if err != nil {
		return nil, err
	}
	return decodedImg, nil
}
