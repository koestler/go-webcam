package config

import (
	"net/url"
	"time"
)

type Config struct {
	version               int                 `yaml:"Version"`           // must be 0
	mqttClients           []*MqttClientConfig `yaml:"MqttClient"`        // mandatory: at least 1 must be defined
	cameras               []*CameraConfig     `yaml:"Cameras"`           // mandatory: at least 1 must be defined
	views                 []*ViewConfig       `yaml:"Views"`             // mandatory: at least 1 must be defined
	httpServer            HttpServerConfig    `yaml:"HttpServer"`        // optional: default Disabled
	logConfig             bool                `yaml:"LogConfig"`         // optional: default False
	logWorkerStart        bool                `yaml:"LogWorkerStart"`    // optional: default False
	logMqttDebug          bool                `yaml:"LogMqttDebug"`      // optional: default False
	projectTitle          string              `yaml:"ProjectTitle"`      // optional: default go-webcam
	authJwtSecret         string              `yaml:"JwtSecret"`         // optional: default new random string on startup
	authJwtValidityPeriod time.Duration       `yaml:"JwtValidityPeriod"` // optional: default 1h
	authHtaccessFile      string              `yaml:"authHtaccessFile"`  // optional: default no valid users
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

type ViewCameraConfig struct {
	name  string // defined automatically by map key
	title string // mandatory: a nice title for the frontend
}

type ViewConfig struct {
	name                string              // defined automatically by map key
	title               string              // mandatory: a nice title for the frontend
	cameras             []*ViewCameraConfig // mandatory: a list of cameraClient names
	resolutionMaxWidth  int                 // optional: defaults to 3840
	resolutionMaxHeight int                 // optional: defaults  2160
	refreshInterval     time.Duration       // optional: default 1m
	autoplay            bool                // optional: default false
	allowedUsers        []string
}

type HttpServerConfig struct {
	enabled       bool     // defined automatically if HttpServer section exists
	bind          string   // optional: defaults to ::1 (ipv6 loopback)
	port          int      // optional: defaults to 8043
	logRequests   bool     // optional: default False
	enableDocs    bool     // optional: default False
	frontendProxy *url.URL // optional: default deactivated; otherwise an address of the frontend dev-server
	frontendPath  string   // optional: default deactivated; otherwise set to a path where the frontend build is located
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
	ProjectTitle   string                  `yaml:"ProjectTitle"`
	AuthJwtSecret         *string              `yaml:"JwtSecret"`
	AuthJwtValidityPeriod string     `yaml:"JwtValidityPeriod"`
	AuthHtaccessFile      *string              `yaml:"authHtaccessFile"`

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

type viewCameraConfigRead struct {
	Title string `yaml:"Title"`
}

type viewCameraConfigReadMap map[string]viewCameraConfigRead

type viewConfigRead struct {
	Title               string                  `yaml:"Title"`
	Cameras             viewCameraConfigReadMap `yaml:"Cameras"`
	ResolutionMaxWidth  *int                    `yaml:"ResolutionMaxWidth"`
	ResolutionMaxHeight *int                    `yaml:"ResolutionMaxHeight"`
	RefreshInterval     string                  `yaml:"RefreshInterval"`
	Autoplay            *bool                   `yaml:"Autoplay"`
	AllowedUsers        []string                `yaml:"AllowedUsers"`
}

type viewConfigReadMap map[string]viewConfigRead

type httpServerConfigRead struct {
	Bind          string `yaml:"Bind"`
	Port          *int   `yaml:"Port"`
	LogRequests   *bool  `yaml:"LogRequests"`
	EnableDocs    *bool  `yaml:"EnableDocs"`
	FrontendProxy string `yaml:"FrontendProxy"`
	FrontendPath  string `yaml:"FrontendPath"`
}
