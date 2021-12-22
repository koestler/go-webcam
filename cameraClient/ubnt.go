package cameraClient

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path"
	"time"
)

type ubntState struct {
	httpClient    *http.Client
	authenticated bool
}

func createUbntState() ubntState {
	jar, _ := cookiejar.New(nil)

	return ubntState{
		httpClient: &http.Client{
			Jar: jar,
			// this tool is designed to serve cameras running on the local network
			// -> us a relatively short timeout
			Timeout: 10 * time.Second,

			// ubnt cameras don't use valid certificates
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
		authenticated: false,
	}
}

func (c *Client) ubntLogin(force bool) (err error) {
	if force {
		c.ubnt.authenticated = false
	}

	if c.ubnt.authenticated {
		return
	}

	// create address
	addr := "https://" + path.Join(c.Config().Address(), "api/1.1/login")
	loginUrl, err := url.Parse(addr)
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

	res, err := c.ubnt.httpClient.Post(loginUrl.String(), "application/json", bytes.NewBuffer(bodyJson))
	if err != nil {
		return
	}

	if res.StatusCode != 200 {
		log.Printf("cameraClient[%s]: ubntLogin failed, code=%d, addr=%s, requestBody=%s", c.Name(), res.StatusCode, addr, bodyJson)
		return fmt.Errorf("got code %d from camera during ubntLogin", res.StatusCode)
	}

	// we are authenticated
	c.ubnt.authenticated = true
	if c.Config().LogDebug() {
		log.Printf("cameraClient[%s]: ubntLogin successful", c.Name())
	}

	return nil
}

func (c *Client) ubntGetRawImage() (img []byte, err error) {
	// ubntLogin
	err = c.ubntLogin(false)
	if err != nil {
		return
	}

	// create address
	imageUrl, err := url.Parse("https://" + path.Join(c.Config().Address(), "snap.jpeg"))
	if err != nil {
		return
	}

	// first attempt
	res, err := c.ubnt.httpClient.Get(imageUrl.String())
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

		res, err = c.ubnt.httpClient.Get(imageUrl.String())
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
