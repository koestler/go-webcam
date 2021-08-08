package config

import "time"

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
