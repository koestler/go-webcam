package mqttClient

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"log"
	"os"
	"strings"
)

type Client struct {
	cfg        Config
	mqttClient mqtt.Client
	shutdown   chan struct{}
}

type Config interface {
	Name() string
	Broker() string
	User() string
	Password() string
	ClientId() string
	Qos() byte
	TopicPrefix() string
	AvailabilityTopic() string
	LogDebug() bool
}

func RunClient(cfg Config) (*Client, error) {
	// configure client and start connection
	opts := mqtt.NewClientOptions().
		AddBroker(cfg.Broker()).
		SetClientID(cfg.ClientId()).
		SetOrderMatters(false).
		SetCleanSession(true) // use clean, non-persistent session since we only publish

	if user := cfg.User(); len(user) > 0 {
		opts.SetUsername(user)
	}
	if password := cfg.Password(); len(password) > 0 {
		opts.SetPassword(password)
	}

	// setup availability topic using will
	if availabilityTopic := getAvailabilityTopic(cfg); len(availabilityTopic) > 0 {
		opts.SetWill(availabilityTopic, "offline", cfg.Qos(), true)

		// publish availability after each connect
		opts.SetOnConnectHandler(func(client mqtt.Client) {
			client.Publish(availabilityTopic, cfg.Qos(), true, "online")
		})
	}

	mqtt.ERROR = log.New(os.Stdout, "", 0)
	if cfg.LogDebug() {
		mqtt.DEBUG = log.New(os.Stdout, "", 0)
	}

	mqttClient := mqtt.NewClient(opts)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("connect failed: %s", token.Error())
	}

	clientStruct := Client{
		cfg:        cfg,
		mqttClient: mqttClient,
		shutdown:   make(chan struct{}),
	}

	return &clientStruct, nil
}

func (c *Client) Config() Config {
	return c.cfg
}

func (c *Client) Shutdown() {
	close(c.shutdown)

	// publish availability offline
	if availabilityTopic := getAvailabilityTopic(c.cfg); len(availabilityTopic) > 0 {
		c.mqttClient.Publish(availabilityTopic, c.cfg.Qos(), true, "offline")
	}

	c.mqttClient.Disconnect(1000)
	log.Printf("mqttClient[%s]: shutdown completed", c.cfg.Name())
}

func getAvailabilityTopic(cfg Config) string {
	return replaceTemplate(cfg.AvailabilityTopic(), cfg)
}

func replaceTemplate(template string, cfg Config) (r string) {
	r = strings.Replace(template, "%Prefix%", cfg.TopicPrefix(), 1)
	r = strings.Replace(r, "%ClientId%", cfg.ClientId(), 1)
	return
}
