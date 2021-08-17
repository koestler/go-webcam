package config

func (c Config) MarshalYAML() (interface{}, error) {
	return configRead{
		Version: &c.version,
		MqttClients: func() mqttClientConfigReadMap {
			mqttClients := make(mqttClientConfigReadMap, len(c.mqttClients))
			for _, c := range c.mqttClients {
				mqttClients[c.name] = c.convertToRead()
			}
			return mqttClients
		}(),
		Cameras: func() cameraConfigReadMap {
			cameras := make(cameraConfigReadMap, len(c.cameras))
			for _, c := range c.cameras {
				cameras[c.name] = c.convertToRead()
			}
			return cameras
		}(),
		Views: func() viewConfigReadMap {
			views := make(viewConfigReadMap, len(c.views))
			for _, c := range c.views {
				views[c.name] = c.convertToRead()
			}
			return views
		}(),
		HttpServer: func() *httpServerConfigRead {
			if !c.httpServer.enabled {
				return nil
			}
			r := c.httpServer.convertToRead()
			return &r
		}(),
		LogConfig:      &c.logConfig,
		LogWorkerStart: &c.logWorkerStart,
		LogMqttDebug:   &c.logMqttDebug,
		ProjectTitle:   c.projectTitle,
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
		Title:               c.title,
		Cameras:             c.cameras,
		ResolutionMaxWidth:  &c.resolutionMaxWidth,
		ResolutionMaxHeight: &c.resolutionMaxHeight,
		RefreshInterval:     c.refreshInterval.String(),
		Autoplay:            &c.autoplay,
	}
}

func (c HttpServerConfig) convertToRead() httpServerConfigRead {
	return httpServerConfigRead{
		Bind:          c.bind,
		Port:          &c.port,
		LogRequests:   &c.logRequests,
		EnableDocs:    &c.enableDocs,
		ProxyFrontend: c.proxyFrontend,
	}
}
