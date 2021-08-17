package config

import "time"

func (c Config) Version() int {
	return c.version
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

func (c Config) LogMqttDebug() bool {
	return c.logMqttDebug
}

func (c Config) ProjectTitle() string {
	return c.projectTitle
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

func (c ViewConfig) Name() string {
	return c.name
}

func (c ViewConfig) Title() string {
	return c.title
}

func (c ViewConfig) Cameras() []string {
	return c.cameras
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

func (c HttpServerConfig) EnableProxyFrontend() bool {
	return len(c.proxyFrontend) > 0
}

func (c HttpServerConfig) ProxyFrontend() string {
	return c.proxyFrontend
}
