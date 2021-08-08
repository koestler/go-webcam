package config

func (c Config) MarshalYAML() (interface{}, error) {
	return configRead{
		Version: &c.Version,
		MqttClients: func() mqttClientConfigReadMap {
			mqttClients := make(mqttClientConfigReadMap, len(c.MqttClients))
			for _, c := range c.MqttClients {
				mqttClients[c.Name()] = c.convertToRead()
			}
			return mqttClients
		}(),
		Cameras: func() cameraConfigReadMap {
			cameras := make(cameraConfigReadMap, len(c.Cameras))
			for _, c := range c.Cameras {
				cameras[c.Name()] = c.convertToRead()
			}
			return cameras
		}(),
		Views: func() viewConfigReadMap {
			views := make(viewConfigReadMap, len(c.Views))
			for _, c := range c.Views {
				views[c.Name()] = c.convertToRead()
			}
			return views
		}(),
		HttpServer: func() *httpServerConfigRead {
			if !c.HttpServer.Enabled() {
				return nil
			}
			r := c.HttpServer.convertToRead()
			return &r
		}(),
		LogConfig:      &c.LogConfig,
		LogWorkerStart: &c.LogWorkerStart,
		LogMqttDebug:   &c.LogMqttDebug,
	}, nil
}

func (c MqttClientConfig) convertToRead() mqttClientConfigRead {
	return mqttClientConfigRead{
		Broker:            c.broker,
		User:              c.user,
		Password:          c.password,
		ClientId:          c.clientId,
		Qos:               &c.qos,
		AvailabilityTopic: &c.availabilityTopic,
		TopicPrefix:       c.topicPrefix,
		LogMessages:       &c.logMessages,
	}
}

func (c CameraConfig) convertToRead() cameraConfigRead {
	return cameraConfigRead{
		Address:         c.address,
		User:            c.user,
		Password:        c.password,
		RefreshInterval: c.refreshInterval.String(),
	}
}

func (c ViewConfig) convertToRead() viewConfigRead {
	return viewConfigRead{
		Route:               c.route,
		Cameras:             c.cameras,
		ResolutionMaxWidth:  &c.resolutionMaxWidth,
		ResolutionMaxHeight: &c.resolutionMaxHeight,
		RefreshInterval:     c.refreshInterval.String(),
	}
}

func (c HttpServerConfig) convertToRead() httpServerConfigRead {
	return httpServerConfigRead{
		Bind:        c.bind,
		Port:        &c.port,
		LogRequests: &c.logRequests,
	}
}
