package config

import (
	"net/url"
	"time"
)

type Config struct {
	version        int                 `yaml:"Version"`        // must be 0
	projectTitle   string              `yaml:"ProjectTitle"`   // optional: default go-webcam
	auth           AuthConfig          `yaml:"Auth"`           // optional: default Disabled
	mqttClients    []*MqttClientConfig `yaml:"MqttClient"`     // mandatory: at least 1 must be defined
	cameras        []*CameraConfig     `yaml:"Cameras"`        // mandatory: at least 1 must be defined
	views          []*ViewConfig       `yaml:"Views"`          // mandatory: at least 1 must be defined
	httpServer     HttpServerConfig    `yaml:"HttpServer"`     // optional: default Disabled
	logConfig      bool                `yaml:"LogConfig"`      // optional: default False
	logWorkerStart bool                `yaml:"LogWorkerStart"` // optional: default False
	logDebug       bool                `yaml:"LogDebug"`       // optional: default False
}

type AuthConfig struct {
	enabled           bool          // defined automatically if Auth section exists
	jwtSecret         []byte        `yaml:"JwtSecret"`         // optional: default new random string on startup
	jwtValidityPeriod time.Duration `yaml:"JwtValidityPeriod"` // optional: default 1h
	htaccessFile      string        `yaml:"HtaccessFile"`      // optional: default no valid users
	logAuth           bool          `yaml:"LogAuth"`           // optional: default False
}

type MqttClientConfig struct {
	name              string // defined automatically by map key
	broker            string // mandatory
	user              string // optional: default empty
	password          string // optional: default empty
	clientId          string // optional: default go-webcam-UUID
	qos               byte   // optional: default 1, must be 0, 1, 2
	availabilityTopic string // optional: default %Prefix%tele/%ClientId%/status
	topicPrefix       string // optional: default empty
	logDebug          bool   // optional: default False
}

type CameraConfig struct {
	name            string        // defined automatically by map key
	address         string        // mandatory
	user            string        // optional: default empty
	password        string        // optional: default empty
	refreshInterval time.Duration // optional: default 200ms
	preemptiveFetch time.Duration // optional: default 2 x refreshInterval
}

type ViewCameraConfig struct {
	name  string // defined automatically by map key
	title string // mandatory: a nice title for the frontend
}

type ViewConfig struct {
	name                string              // mandatory: A technical name used in the URLs
	title               string              // mandatory: a nice title for the frontend
	cameras             []*ViewCameraConfig // mandatory: a list of cameraClient names
	resolutionMaxWidth  int                 // optional: defaults to 3840
	resolutionMaxHeight int                 // optional: defaults  2160
	jpgQuality          int                 // optional: default 85
	refreshInterval     time.Duration       // optional: default 1m
	autoplay            bool                // optional: default false
	allowedUsers        map[string]struct{} // optional: if empty: view is public; otherwise only allowed to listed users
	hidden              bool                // optional: if true, view is not shown in menu unless logged in
}

type HttpServerConfig struct {
	enabled          bool          // defined automatically if HttpServer section exists
	bind             string        // optional: defaults to ::1 (ipv6 loopback)
	port             int           // optional: defaults to 8043
	logRequests      bool          // optional: default False
	enableDocs       bool          // optional: default True
	frontendProxy    *url.URL      // optional: default deactivated; otherwise an address of the frontend dev-server
	frontendPath     string        // optional: default "frontend-build"; otherwise set to a path where the frontend build is located
	frontendExpires  time.Duration // optional: default 5min; what cache-control header to sent for static frontend files
	configExpires    time.Duration // optional: default 1min; what cache-control header to sent for static frontend files
	hashTimeout      time.Duration // optional: default 10s; for how long, after a redirect to a imageByHash is made, the entry is stored
	imageEarlyExpire time.Duration // optional: default 2s; s-maxage of images is computed ad expiry - imageEarlyExpire;
	hashSecret       string        // optional: default random string on startup
}

// Read structs are given to yaml for decoding and are slightly less exact in types
type configRead struct {
	Version        *int                    `yaml:"Version"`
	ProjectTitle   string                  `yaml:"ProjectTitle"`
	Auth           *authConfigRead         `yaml:"Auth"`
	MqttClients    mqttClientConfigReadMap `yaml:"MqttClients"`
	Cameras        cameraConfigReadMap     `yaml:"Cameras"`
	Views          viewConfigReadList      `yaml:"Views"`
	HttpServer     *httpServerConfigRead   `yaml:"HttpServer"`
	LogConfig      *bool                   `yaml:"LogConfig"`
	LogWorkerStart *bool                   `yaml:"LogWorkerStart"`
	LogDebug       *bool                   `yaml:"LogDebug"`
}

type authConfigRead struct {
	JwtSecret         *string `yaml:"JwtSecret"`
	JwtValidityPeriod string  `yaml:"JwtValidityPeriod"`
	HtaccessFile      *string `yaml:"HtaccessFile"`
	LogAuth           *bool   `yaml:"LogAuth"`
}

type mqttClientConfigRead struct {
	Broker            string  `yaml:"Broker"`
	User              string  `yaml:"User"`
	Password          string  `yaml:"Password"`
	ClientId          string  `yaml:"ClientId"`
	Qos               *byte   `yaml:"Qos"`
	AvailabilityTopic *string `yaml:"AvailabilityTopic"`
	TopicPrefix       string  `yaml:"TopicPrefix"`
	LogDebug          *bool   `yaml:"LogDebug"`
}

type mqttClientConfigReadMap map[string]mqttClientConfigRead

type cameraConfigRead struct {
	Address         string `yaml:"Address"`
	User            string `yaml:"User"`
	Password        string `yaml:"Password"`
	RefreshInterval string `yaml:"RefreshInterval"`
	PreemptiveFetch string `yaml:"PreemptiveFetch"`
}

type cameraConfigReadMap map[string]cameraConfigRead

type viewCameraConfigRead struct {
	Name  string `yaml:"Name"`
	Title string `yaml:"Title"`
}

type viewCameraConfigReadList []viewCameraConfigRead

type viewConfigRead struct {
	Name                string                   `yaml:"Name"`
	Title               string                   `yaml:"Title"`
	Cameras             viewCameraConfigReadList `yaml:"Cameras"`
	ResolutionMaxWidth  *int                     `yaml:"ResolutionMaxWidth"`
	ResolutionMaxHeight *int                     `yaml:"ResolutionMaxHeight"`
	JpgQuality          *int                     `yaml:"JpgQuality"`
	RefreshInterval     string                   `yaml:"RefreshInterval"`
	Autoplay            *bool                    `yaml:"Autoplay"`
	AllowedUsers        []string                 `yaml:"AllowedUsers"`
	Hidden              *bool                    `yaml:"Hidden"`
}

type viewConfigReadList []viewConfigRead

type httpServerConfigRead struct {
	Bind             string  `yaml:"Bind"`
	Port             *int    `yaml:"Port"`
	LogRequests      *bool   `yaml:"LogRequests"`
	EnableDocs       *bool   `yaml:"EnableDocs"`
	FrontendProxy    string  `yaml:"FrontendProxy"`
	FrontendPath     string  `yaml:"FrontendPath"`
	FrontendExpires  string  `yaml:"FrontendExpires"`
	ConfigExpires    string  `yaml:"ConfigExpires"`
	HashTimeout      string  `yaml:"HashTimeout"`
	ImageEarlyExpire string  `yaml:"ImageEarlyExpire"`
	HashSecret       *string `yaml:"HashSecret"`
}
