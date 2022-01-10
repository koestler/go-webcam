package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

const NameRegexp = "^[a-zA-Z0-9\\-]{1,32}$"

var nameMatcher = regexp.MustCompile(NameRegexp)

func ReadConfigFile(exe, source string) (config Config, err []error) {
	yamlStr, e := ioutil.ReadFile(source)
	if e != nil {
		return config, []error{fmt.Errorf("cannot read configuration: %v; use see `%s --help`", err, exe)}
	}

	return ReadConfig(yamlStr)
}

func ReadConfig(yamlStr []byte) (config Config, err []error) {
	var configRead configRead

	yamlStr = []byte(os.ExpandEnv(string(yamlStr)))
	e := yaml.Unmarshal(yamlStr, &configRead)
	if e != nil {
		return config, []error{fmt.Errorf("cannot parse yaml: %s", e)}
	}

	return configRead.TransformAndValidate()
}

func (c Config) PrintConfig() (err error) {
	newYamlStr, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("cannot encode yaml again: %s", err)
	}

	log.Print("config: use the following config:")
	for _, line := range strings.Split(string(newYamlStr), "\n") {
		log.Print("config: ", line)
	}
	return nil
}

func (c configRead) TransformAndValidate() (ret Config, err []error) {
	var e []error
	ret.auth, e = c.Auth.TransformAndValidate()
	err = append(err, e...)

	ret.mqttClients, e = c.MqttClients.TransformAndValidate()
	err = append(err, e...)

	ret.cameras, e = c.Cameras.TransformAndValidate()
	err = append(err, e...)

	ret.views, e = c.Views.TransformAndValidate(ret.cameras)
	err = append(err, e...)

	ret.httpServer, e = c.HttpServer.TransformAndValidate()
	err = append(err, e...)

	if c.Version == nil {
		err = append(err, fmt.Errorf("Version must be defined. Use Version=0."))
	} else {
		ret.version = *c.Version
		if ret.version != 0 {
			err = append(err, fmt.Errorf("Version=%d is not supported.", ret.version))
		}
	}

	if len(c.ProjectTitle) > 0 {
		ret.projectTitle = c.ProjectTitle
	} else {
		ret.projectTitle = "go-webcam"
	}

	if c.LogConfig != nil && *c.LogConfig {
		ret.logConfig = true
	}

	if c.LogWorkerStart != nil && *c.LogWorkerStart {
		ret.logWorkerStart = true
	}

	if c.LogDebug != nil && *c.LogDebug {
		ret.logDebug = true
	}

	return
}

func (c *authConfigRead) TransformAndValidate() (ret AuthConfig, err []error) {
	ret.enabled = false
	ret.jwtValidityPeriod = time.Hour

	if randString, e := randomString(64); err == nil {
		ret.jwtSecret = []byte(randString)
	} else {
		err = append(err, fmt.Errorf("Auth->JwtSecret: error while generating random secret: %s", e))
	}

	if c == nil {
		return
	}

	ret.enabled = true

	if c.JwtSecret != nil {
		if len(*c.JwtSecret) < 32 {
			err = append(err, fmt.Errorf("Auth->JwtSecret must be empty ot >= 32 chars"))
		} else {
			ret.jwtSecret = []byte(*c.JwtSecret)
		}
	}

	if len(c.JwtValidityPeriod) < 1 {
		// use default
	} else if authJwtValidityPeriod, e := time.ParseDuration(c.JwtValidityPeriod); e != nil {
		err = append(err, fmt.Errorf("Auth->JwtValidityPeriod='%s' parse error: %s",
			c.JwtValidityPeriod, e,
		))
	} else if authJwtValidityPeriod < 0 {
		err = append(err, fmt.Errorf("Auth->JwtValidityPeriod='%s' must be positive",
			c.JwtValidityPeriod,
		))
	} else {
		ret.jwtValidityPeriod = authJwtValidityPeriod
	}

	if c.HtaccessFile != nil && len(*c.HtaccessFile) > 0 {
		if info, e := os.Stat(*c.HtaccessFile); e != nil {
			err = append(err, fmt.Errorf("Auth->HtaccessFile='%s' cannot open file. error: %s",
				*c.HtaccessFile, e,
			))
		} else if info.IsDir() {
			err = append(err, fmt.Errorf("Auth->HtaccessFile='%s' must be a file, not a directory",
				*c.HtaccessFile,
			))
		}

		ret.htaccessFile = *c.HtaccessFile
	}

	return
}

func (c *httpServerConfigRead) TransformAndValidate() (ret HttpServerConfig, err []error) {
	ret.enabled = false
	ret.bind = "[::1]"
	ret.port = 8043

	if randString, e := randomString(64); err == nil {
		ret.hashSecret = randString
	} else {
		err = append(err, fmt.Errorf("HttpServerConfig->HashSecret: error while generating random secret: %s", e))
	}

	if c == nil {
		return
	}

	ret.enabled = true

	if len(c.Bind) > 0 {
		ret.bind = c.Bind
	}

	if c.Port != nil {
		ret.port = *c.Port
	}

	if c.LogRequests != nil && *c.LogRequests {
		ret.logRequests = true
	}

	ret.enableDocs = true
	if c.EnableDocs != nil && !*c.EnableDocs {
		ret.enableDocs = false
	}

	if len(c.FrontendProxy) > 0 {
		u, parseError := url.Parse(c.FrontendProxy)
		if parseError == nil {
			ret.frontendProxy = u
		} else {
			err = append(err, fmt.Errorf("HttpServerConfig->FrontendProxy must not be empty (=disabled) or a valid URL, err: %s", parseError))
		}
	}

	if len(c.FrontendPath) > 0 {
		ret.frontendPath = c.FrontendPath
	} else {
		ret.frontendPath = "frontend-build"
	}

	if len(c.FrontendExpires) < 1 {
		// use default 5min
		ret.frontendExpires = 5 * time.Minute
	} else if frontendExpires, e := time.ParseDuration(c.FrontendExpires); e != nil {
		err = append(err, fmt.Errorf("HttpServerConfig->FrontendExpires='%s' parse error: %s", c.FrontendExpires, e))
	} else if frontendExpires < 0 {
		err = append(err, fmt.Errorf("HttpServerConfig->FrontendExpires='%s' must be positive", c.FrontendExpires))
	} else {
		ret.frontendExpires = frontendExpires
	}

	if len(c.ConfigExpires) < 1 {
		// use default 1min
		ret.configExpires = 1 * time.Minute
	} else if configExpires, e := time.ParseDuration(c.ConfigExpires); e != nil {
		err = append(err, fmt.Errorf("HttpServerConfig->ConfigExpires='%s' parse error: %s", c.ConfigExpires, e))
	} else if configExpires < 0 {
		err = append(err, fmt.Errorf("HttpServerConfig->ConfigExpires='%s' must be positive", c.ConfigExpires))
	} else {
		ret.configExpires = configExpires
	}

	if len(c.HashTimeout) < 1 {
		// use default 10s
		ret.hashTimeout = 10 * time.Second
	} else if hashTimeout, e := time.ParseDuration(c.HashTimeout); e != nil {
		err = append(err, fmt.Errorf("HttpServerConfig->HashTimeout='%s' parse error: %s", c.HashTimeout, e))
	} else if hashTimeout < 0 {
		err = append(err, fmt.Errorf("HttpServerConfig->HashTimeout='%s' must be positive", c.HashTimeout))
	} else {
		ret.hashTimeout = hashTimeout
	}

	if len(c.ImageEarlyExpire) < 1 {
		// use default 10s
		ret.imageEarlyExpire = time.Second
	} else if imageEarlyExpire, e := time.ParseDuration(c.ImageEarlyExpire); e != nil {
		err = append(err, fmt.Errorf("HttpServerConfig->ImageEarlyExpire='%s' parse error: %s", c.ImageEarlyExpire, e))
	} else if imageEarlyExpire < 0 {
		err = append(err, fmt.Errorf("HttpServerConfig->ImageEarlyExpire='%s' must be positive", c.ImageEarlyExpire))
	} else {
		ret.imageEarlyExpire = imageEarlyExpire
	}

	if c.HashSecret != nil {
		if len(*c.HashSecret) < 32 {
			err = append(err, fmt.Errorf("HashSecret must be empty ot >= 32 chars"))
		} else {
			ret.hashSecret = *c.HashSecret
		}
	}

	return
}

func (c mqttClientConfigReadMap) getOrderedKeys() (ret []string) {
	ret = make([]string, len(c))
	i := 0
	for k := range c {
		ret[i] = k
		i++
	}
	sort.Strings(ret)
	return
}

func (c mqttClientConfigReadMap) TransformAndValidate() (ret []*MqttClientConfig, err []error) {
	ret = make([]*MqttClientConfig, len(c))
	j := 0
	for _, name := range c.getOrderedKeys() {
		r, e := c[name].TransformAndValidate(name)
		ret[j] = &r
		err = append(err, e...)
		j++
	}
	return
}

func (c mqttClientConfigRead) TransformAndValidate(name string) (ret MqttClientConfig, err []error) {
	ret = MqttClientConfig{
		name:        name,
		broker:      c.Broker,
		user:        c.User,
		password:    c.Password,
		clientId:    c.ClientId,
		topicPrefix: c.TopicPrefix,
	}

	if !nameMatcher.MatchString(ret.name) {
		err = append(err, fmt.Errorf("MqttClientConfig->Name='%s' does not match %s", ret.name, NameRegexp))
	}

	if len(ret.broker) < 1 {
		err = append(err, fmt.Errorf("MqttClientConfig->%s->Broker must not be empty", name))
	}
	if len(ret.clientId) < 1 {
		ret.clientId = "go-webcam"
	}
	if c.Qos == nil {
		ret.qos = 0
	} else if *c.Qos == 0 || *c.Qos == 1 || *c.Qos == 2 {
		ret.qos = *c.Qos
	} else {
		err = append(err, fmt.Errorf("MqttClientConfig->%s->Qos=%d but must be 0, 1 or 2", name, *c.Qos))
	}

	if c.AvailabilityTopic == nil {
		// use default
		ret.availabilityTopic = "%Prefix%tele/%clientId%/LWT"
	} else {
		ret.availabilityTopic = *c.AvailabilityTopic
	}

	if c.LogMessages != nil && *c.LogMessages {
		ret.logMessages = true
	}

	return
}

func (c cameraConfigReadMap) getOrderedKeys() (ret []string) {
	ret = make([]string, len(c))
	i := 0
	for k := range c {
		ret[i] = k
		i++
	}
	sort.Strings(ret)
	return
}

func (c cameraConfigReadMap) TransformAndValidate() (ret []*CameraConfig, err []error) {
	if len(c) < 1 {
		return ret, []error{fmt.Errorf("Cameras section must no be empty")}
	}

	ret = make([]*CameraConfig, len(c))
	j := 0
	for _, name := range c.getOrderedKeys() {
		r, e := c[name].TransformAndValidate(name)
		ret[j] = &r
		err = append(err, e...)
		j++
	}
	return
}

func (c cameraConfigRead) TransformAndValidate(name string) (ret CameraConfig, err []error) {
	ret = CameraConfig{
		name:     name,
		address:  c.Address,
		user:     c.User,
		password: c.Password,
	}

	if !nameMatcher.MatchString(ret.name) {
		err = append(err, fmt.Errorf("CameraConfig->Name='%s' does not match %s", ret.name, NameRegexp))
	}

	if len(ret.address) < 1 {
		err = append(err, fmt.Errorf("CameraConfig->%s->Address must not be empty", name))
	}

	if len(c.RefreshInterval) < 1 {
		// use default 200ms
		ret.refreshInterval = 200 * time.Millisecond
	} else if refreshInterval, e := time.ParseDuration(c.RefreshInterval); e != nil {
		err = append(err, fmt.Errorf("CameraConfig->%s->RefreshInterval='%s' parse error: %s",
			name, c.RefreshInterval, e,
		))
	} else if refreshInterval < 0 {
		err = append(err, fmt.Errorf("CameraConfig->%s->RefreshInterval='%s' must be positive",
			name, c.RefreshInterval,
		))
	} else {
		ret.refreshInterval = refreshInterval
	}

	return
}

func (c viewConfigReadList) TransformAndValidate(cameras []*CameraConfig) (ret []*ViewConfig, err []error) {
	if len(c) < 1 {
		return ret, []error{fmt.Errorf("Views section must no be empty.")}
	}

	ret = make([]*ViewConfig, len(c))
	j := 0
	for _, cr := range c {
		r, e := cr.TransformAndValidate(cameras)

		// check for duplicate name
		for i := 0; i < j; i++ {
			if r.Name() == ret[i].Name() {
				err = append(err, fmt.Errorf("Views->Name='%s': name must be unique", r.Name()))
			}
		}

		ret[j] = &r
		err = append(err, e...)
		j++
	}

	return
}

func (c viewConfigRead) TransformAndValidate(cameras []*CameraConfig) (ret ViewConfig, err []error) {
	ret = ViewConfig{
		name:         c.Name,
		title:        c.Title,
		allowedUsers: make(map[string]struct{}),
		hidden:       false,
	}

	if !nameMatcher.MatchString(ret.name) {
		err = append(err, fmt.Errorf("Views->Name='%s' does not match %s", ret.name, NameRegexp))
	}

	if len(c.Title) < 1 {
		err = append(err, fmt.Errorf("Views->%s->Title must not be empty", c.Name))
	}

	{
		var camerasErr []error
		ret.cameras, camerasErr = c.Cameras.TransformAndValidate(cameras)
		for _, ce := range camerasErr {
			err = append(err, fmt.Errorf("Views->%s: %s", c.Name, ce))
		}
	}

	if c.ResolutionMaxWidth == nil {
		ret.resolutionMaxWidth = 3840
	} else if *c.ResolutionMaxWidth > 0 {
		ret.resolutionMaxWidth = *c.ResolutionMaxWidth
	} else {
		err = append(err, fmt.Errorf("Views->%s->ResolutionMaxWidth=%d but must be a positive integer", c.Name, *c.ResolutionMaxWidth))
	}

	if c.ResolutionMaxHeight == nil {
		ret.resolutionMaxHeight = 3840
	} else if *c.ResolutionMaxHeight > 0 {
		ret.resolutionMaxHeight = *c.ResolutionMaxHeight
	} else {
		err = append(err, fmt.Errorf("Views->%s->ResolutionMaxHeight=%d but must be a positive integer", c.Name, *c.ResolutionMaxHeight))
	}

	if c.JpgQuality == nil {
		ret.jpgQuality = 85
	} else if *c.JpgQuality > 0 && *c.JpgQuality <= 100 {
		ret.jpgQuality = *c.JpgQuality
	} else {
		err = append(err, fmt.Errorf("Views->%s->JpgQuality=%d but must be >0 and <= 100", c.Name, *c.JpgQuality))
	}

	if len(c.RefreshInterval) < 1 {
		// use default 0
		ret.refreshInterval = time.Minute
	} else if refreshInterval, e := time.ParseDuration(c.RefreshInterval); e != nil {
		err = append(err, fmt.Errorf("viewConfig->%s->RefreshInterval='%s' parse error: %s",
			c.Name, c.RefreshInterval, e,
		))
	} else if refreshInterval < 0 {
		err = append(err, fmt.Errorf("viewConfig->%s->RefreshInterval='%s' must be positive",
			c.Name, c.RefreshInterval,
		))
	} else {
		ret.refreshInterval = refreshInterval
	}

	if c.Autoplay != nil && *c.Autoplay {
		ret.autoplay = true
	}

	for _, user := range c.AllowedUsers {
		ret.allowedUsers[user] = struct{}{}
	}

	if c.Hidden != nil && *c.Hidden {
		ret.hidden = true
	}

	return
}

func (c viewCameraConfigReadMap) TransformAndValidate(cameras []*CameraConfig) (ret []*ViewCameraConfig, err []error) {
	if len(c) < 1 {
		return ret, []error{fmt.Errorf("Cameras section must no be empty.")}
	}

	ret = make([]*ViewCameraConfig, len(c))
	j := 0
	for name, camera := range c {
		r, e := camera.TransformAndValidate(name, cameras)
		ret[j] = &r
		err = append(err, e...)
		j++
	}
	return

}

func (c viewCameraConfigRead) TransformAndValidate(
	name string,
	cameras []*CameraConfig,
) (ret ViewCameraConfig, err []error) {
	if !cameraExists(name, cameras) {
		err = append(err, fmt.Errorf("Camera='%s' is not defined", name))
	}

	ret = ViewCameraConfig{
		name:  name,
		title: c.Title,
	}

	return
}

func cameraExists(cameraName string,
	cameras []*CameraConfig) bool {
	for _, client := range cameras {
		if cameraName == client.name {
			return true
		}
	}
	return false
}
