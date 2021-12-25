package config

func (c Config) MarshalYAML() (interface{}, error) {

	return configRead{
		Version: &c.version,
		Auth: func() *authConfigRead {
			if !c.auth.enabled {
				return nil
			}
			r := c.auth.convertToRead()
			return &r
		}(),
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
		Views: func() viewConfigReadList {
			views := make(viewConfigReadList, len(c.views))
			for _, c := range c.views {
				views = append(views, c.convertToRead())
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
		LogDebug:       &c.logDebug,
		ProjectTitle:   c.projectTitle,
	}, nil
}

func (c AuthConfig) convertToRead() authConfigRead {
	jwtSecret := string(c.jwtSecret)
	return authConfigRead{
		JwtSecret:         &jwtSecret,
		JwtValidityPeriod: c.jwtValidityPeriod.String(),
		HtaccessFile:      &c.htaccessFile,
	}
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

func (c ViewCameraConfig) convertToRead() viewCameraConfigRead {
	return viewCameraConfigRead{
		Title: c.title,
	}
}

func (c ViewConfig) convertToRead() viewConfigRead {
	return viewConfigRead{
		Name:  c.name,
		Title: c.title,
		Cameras: func() viewCameraConfigReadMap {
			views := make(viewCameraConfigReadMap, len(c.cameras))
			for _, c := range c.cameras {
				views[c.name] = c.convertToRead()
			}
			return views
		}(),
		ResolutionMaxWidth:  &c.resolutionMaxWidth,
		ResolutionMaxHeight: &c.resolutionMaxHeight,
		JpgQuality:          &c.jpgQuality,
		RefreshInterval:     c.refreshInterval.String(),
		Autoplay:            &c.autoplay,
		AllowedUsers:        mapKeys(c.allowedUsers),
		Hidden:              &c.hidden,
	}
}

func mapKeys(m map[string]struct{}) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func (c HttpServerConfig) convertToRead() httpServerConfigRead {
	frontendProxy := ""
	if c.frontendProxy != nil {
		frontendProxy = c.frontendProxy.String()
	}

	return httpServerConfigRead{
		Bind:            c.bind,
		Port:            &c.port,
		LogRequests:     &c.logRequests,
		EnableDocs:      &c.enableDocs,
		FrontendProxy:   frontendProxy,
		FrontendPath:    c.frontendPath,
		FrontendExpires: c.frontendExpires.String(),
		ConfigExpires:   c.configExpires.String(),
		HashTimeout:     c.hashTimeout.String(),
		HashSecret:      &c.hashSecret,
	}
}
