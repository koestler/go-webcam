package config

import (
	"net/url"
	"time"
)

func (c Config) Version() int {
	return c.version
}

func (c Config) Auth() AuthConfig {
	return c.auth
}

func (c Config) MqttClients() []*MqttClientConfig {
	return c.mqttClients
}

func (c Config) Cameras() []*CameraConfig {
	return c.cameras
}

func (c Config) Views() []*ViewConfig {
	return c.views
}

func (c Config) HttpServer() HttpServerConfig {
	return c.httpServer
}

func (c Config) LogConfig() bool {
	return c.logConfig
}

func (c Config) LogWorkerStart() bool {
	return c.logWorkerStart
}

func (c Config) LogDebug() bool {
	return c.logDebug
}

func (c Config) ProjectTitle() string {
	return c.projectTitle
}

func (c AuthConfig) Enabled() bool {
	return c.enabled
}

func (c AuthConfig) JwtSecret() []byte {
	return c.jwtSecret
}

func (c AuthConfig) JwtValidityPeriod() time.Duration {
	return c.jwtValidityPeriod
}

func (c AuthConfig) HtaccessFile() string {
	return c.htaccessFile
}

func (c MqttClientConfig) Name() string {
	return c.name
}

func (c MqttClientConfig) Broker() string {
	return c.broker
}

func (c MqttClientConfig) User() string {
	return c.user
}

func (c MqttClientConfig) Password() string {
	return c.password
}

func (c MqttClientConfig) ClientId() string {
	return c.clientId
}

func (c MqttClientConfig) Qos() byte {
	return c.qos
}

func (c MqttClientConfig) AvailabilityTopic() string {
	return c.availabilityTopic
}

func (c MqttClientConfig) TopicPrefix() string {
	return c.topicPrefix
}

func (c MqttClientConfig) LogMessages() bool {
	return c.logMessages
}

func (c CameraConfig) Name() string {
	return c.name
}

func (c CameraConfig) Address() string {
	return c.address
}

func (c CameraConfig) User() string {
	return c.user
}

func (c CameraConfig) Password() string {
	return c.password
}

func (c CameraConfig) RefreshInterval() time.Duration {
	return c.refreshInterval
}

func (c CameraConfig) ExpireEarly() time.Duration {
	// expire images 100ms early
	// this ensures that always a new image is fetched during periodic reloads with a jitter of up to 100ms
	return 100 * time.Millisecond
}

func (c ViewCameraConfig) Name() string {
	return c.name
}

func (c ViewCameraConfig) Title() string {
	return c.title
}

func (c ViewConfig) Name() string {
	return c.name
}

func (c ViewConfig) Title() string {
	return c.title
}

func (c ViewConfig) Cameras() []*ViewCameraConfig {
	return c.cameras
}

func (c ViewConfig) CameraNames() []string {
	names := make([]string, len(c.cameras))
	for i, camera := range c.cameras {
		names[i] = camera.Name()
	}
	return names
}

func (c ViewConfig) ResolutionMaxWidth() int {
	return c.resolutionMaxWidth
}

func (c ViewConfig) ResolutionMaxHeight() int {
	return c.resolutionMaxHeight
}

func (c ViewConfig) RefreshInterval() time.Duration {
	return c.refreshInterval
}

func (c ViewConfig) Autoplay() bool {
	return c.autoplay
}

func (c ViewConfig) IsAllowed(user string) bool {
	_, ok := c.allowedUsers[user]
	return ok
}

func (c ViewConfig) IsPublic() bool {
	return len(c.allowedUsers) == 0
}

func (c ViewConfig) Hidden() bool {
	return c.hidden
}

func (c HttpServerConfig) Enabled() bool {
	return c.enabled
}

func (c HttpServerConfig) Bind() string {
	return c.bind
}

func (c HttpServerConfig) Port() int {
	return c.port
}

func (c HttpServerConfig) LogRequests() bool {
	return c.logRequests
}

func (c HttpServerConfig) EnableDocs() bool {
	return c.enableDocs
}

func (c HttpServerConfig) FrontendProxy() *url.URL {
	return c.frontendProxy
}

func (c HttpServerConfig) FrontendPath() string {
	return c.frontendPath
}

func (c HttpServerConfig) FrontendExpires() time.Duration {
	return c.frontendExpires
}

func (c HttpServerConfig) ConfigExpires() time.Duration {
	return c.configExpires
}

func (c HttpServerConfig) HashTimeout() time.Duration {
	return c.hashTimeout
}

func (c Config) GetViewNames() (ret []string) {
	ret = []string{}
	for _, v := range c.Views() {
		ret = append(ret, v.Name())
	}
	return
}
