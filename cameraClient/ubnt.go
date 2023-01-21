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
	httpClient        *http.Client
	authenticated     bool
	firstAttemptError string
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
		authenticated:     false,
		firstAttemptError: "",
	}
}

func (c *Client) ubntLogin(force bool) (err error) {
	if force {
		c.ubnt.authenticated = false
	}

	if c.ubnt.authenticated {
		return
	}

	start := time.Now()

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
		log.Printf("cameraClient[%s]: ubntLogin failed: %s", c.Name(), err)
		return
	}

	if res.StatusCode != 200 {
		log.Printf("cameraClient[%s]: ubntLogin failed, code=%d, addr=%s, requestBody=%s", c.Name(), res.StatusCode, addr, bodyJson)
		return fmt.Errorf("got code %d from camera during ubntLogin", res.StatusCode)
	}

	// we are authenticated
	c.ubnt.authenticated = true
	if c.Config().LogDebug() {
		log.Printf("cameraClient[%s]: ubntLogin successful, took=%s", c.Name(), time.Since(start))
	}

	return nil
}

func (c *Client) ubntGetRawImage() (img []byte, err error) {
	// ubntLogin
	err = c.ubntLogin(false)
	if err != nil {
		return
	}

	start := time.Now()

	// create address
	imageUrl, err := url.Parse("https://" + path.Join(c.Config().Address(), "snap.jpeg"))
	if err != nil {
		return
	}

	// first attempt
	res, err := c.ubnt.httpClient.Get(imageUrl.String())
	if err != nil {
		errMsg := fmt.Sprintf("cameraClient[%s]: fetch failed: %s", c.Name(), err)
		// only print the same connect error once and not on every retry
		// this is prevents the logs from filling up when a host is unavailable
		if c.ubnt.firstAttemptError != errMsg {
			log.Println(errMsg)
		}
		c.ubnt.firstAttemptError = errMsg
		return
	}
	c.ubnt.firstAttemptError = ""

	// second attempt; try to relogin if unauthorized is returned
	if res.StatusCode == http.StatusUnauthorized {
		// unauthorized -> reloggin
		err = c.ubntLogin(true)
		if err != nil {
			return
		}

		res, err = c.ubnt.httpClient.Get(imageUrl.String())
		if err != nil {
			log.Printf("cameraClient[%s]: fetch failed: %s", c.Name(), err)
			return
		}

	}

	// handle errors
	if res.StatusCode != http.StatusOK {
		if c.Config().LogDebug() {
			log.Printf("cameraClient[%s]: fetch failed with code %d, error: %s", c.Name(), res.StatusCode, res.Body)
		}
		return nil, fmt.Errorf("got code %v from camera when fetching a snapshot", res.StatusCode)
	}

	// read and return body
	defer res.Body.Close()
	img, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if c.Config().LogDebug() {
		log.Printf("cameraClient[%s]: ubntImage fetched, took=%.3f", c.Name(), time.Since(start).Seconds())
	}
	return img, nil
}
