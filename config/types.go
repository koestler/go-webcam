package config

import "time"

type Config struct {
	Version        int                 `yaml:"Version"`        // must be 0
	MqttClients    []*MqttClientConfig `yaml:"MqttClient"`     // mandatory: at least 1 must be defined
	Cameras        []*CameraConfig     `yaml:"Cameras"`        // mandatory: at least 1 must be defined
	Views          []*ViewConfig       `yaml:"Views"`          // mandatory: at least 1 must be defined
	HttpServer     HttpServerConfig    `yaml:"HttpServer"`     // optional: default Disabled
	LogConfig      bool                `yaml:"LogConfig"`      // optional: default False
	LogWorkerStart bool                `yaml:"LogWorkerStart"` // optional: default False
	LogMqttDebug   bool                `yaml:"LogMqttDebug"`   // optional: default False
}

type MqttClientConfig struct {
	name              string // defined automatically by map key
	broker            string // mandatory
	user              string // optional: default empty
	password          string // optional: default empty
	clientId          string // optional: default go-webcam
	qos               byte   // optional: default 0, must be 0, 1, 2
	availabilityTopic string // optional: default %Prefix%tele/%clientId%/LWT
	topicPrefix       string // optional: default empty
	logMessages       bool   // optional: default False
}

type CameraConfig struct {
	name            string        // defined automatically by map key
	address         string        // mandatory
	user            string        // optional: default empty
	password        string        // optional: default empty
	refreshInterval time.Duration // optional: default 200ms
}

type ViewConfig struct {
	name                string        // defined automatically by map key
	route               string        // mandatory: must end with a '/'
	cameras             []string      // mandatory: a list of cameraClient naems
	resolutionMaxWidth  int           // optional: defaults to 3840
	resolutionMaxHeight int           // optional: defaults  2160
	refreshInterval     time.Duration // optional: default 1m
}

type HttpServerConfig struct {
	enabled     bool   // defined automatically if HttpServer section exists
	bind        string // optional: defaults to ::1 (ipv6 loopback)
	port        int    // optional: defaults to 8043
	logRequests bool   // optional:  default False
}

// Read structs are given to yaml for decoding and are slightly less exact in types
type configRead struct {
	Version        *int                    `yaml:"Version"`
	MqttClients    mqttClientConfigReadMap `yaml:"MqttClients"`
	Cameras        cameraConfigReadMap     `yaml:"Cameras"`
	Views          viewConfigReadMap       `yaml:"Views"`
	HttpServer     *httpServerConfigRead   `yaml:"HttpServer"`
	LogConfig      *bool                   `yaml:"LogConfig"`
	LogWorkerStart *bool                   `yaml:"LogWorkerStart"`
	LogMqttDebug   *bool                   `yaml:"LogMqttDebug"`
}

type mqttClientConfigRead struct {
	Broker            string  `yaml:"Broker"`
	User              string  `yaml:"User"`
	Password          string  `yaml:"Password"`
	ClientId          string  `yaml:"ClientId"`
	Qos               *byte   `yaml:"Qos"`
	AvailabilityTopic *string `yaml:"AvailabilityTopic"`
	TopicPrefix       string  `yaml:"TopicPrefix"`
	LogMessages       *bool   `yaml:"LogMessages"`
}

type mqttClientConfigReadMap map[string]mqttClientConfigRead

type cameraConfigRead struct {
	Address         string `yaml:"Address"`
	User            string `yaml:"User"`
	Password        string `yaml:"Password"`
	RefreshInterval string `yaml:"RefreshInterval"`
}

type cameraConfigReadMap map[string]cameraConfigRead

type viewConfigRead struct {
	Route               string   `yaml:"Route"`
	Cameras             []string `yaml:"Cameras"`
	ResolutionMaxWidth  *int     `yaml:"ResolutionMaxWidth"`
	ResolutionMaxHeight *int     `yaml:"ResolutionMaxHeight"`
	RefreshInterval     string   `yaml:"RefreshInterval"`
}

type viewConfigReadMap map[string]viewConfigRead

type httpServerConfigRead struct {
	Bind        string `yaml:"Bind"`
	Port        *int   `yaml:"Port"`
	LogRequests *bool  `yaml:"LogRequests"`
}
