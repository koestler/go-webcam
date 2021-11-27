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
		return config, []error{fmt.Errorf("cannot parse yaml: %s", err)}
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

	if c.LogConfig != nil && *c.LogConfig {
		ret.logConfig = true
	}

	if c.LogWorkerStart != nil && *c.LogWorkerStart {
		ret.logWorkerStart = true
	}

	if c.LogMqttDebug != nil && *c.LogMqttDebug {
		ret.logMqttDebug = true
	}

	if len(c.ProjectTitle) > 0 {
		ret.projectTitle = c.ProjectTitle
	} else {
		ret.projectTitle = "go-webcam"
	}

	return
}

func (c *httpServerConfigRead) TransformAndValidate() (ret HttpServerConfig, err []error) {
	ret.enabled = false
	ret.bind = "[::1]"
	ret.port = 8043

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

	if c.EnableDocs != nil && *c.EnableDocs {
		ret.enableDocs = true
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
		// use default 0
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

func (c viewConfigReadMap) getOrderedKeys() (ret []string) {
	ret = make([]string, len(c))
	i := 0
	for k := range c {
		ret[i] = k
		i++
	}
	sort.Strings(ret)
	return
}

func (c viewConfigReadMap) TransformAndValidate(cameras []*CameraConfig) (ret []*ViewConfig, err []error) {
	if len(c) < 1 {
		return ret, []error{fmt.Errorf("Views section must no be empty.")}
	}

	ret = make([]*ViewConfig, len(c))
	j := 0
	for _, name := range c.getOrderedKeys() {
		r, e := c[name].TransformAndValidate(name, cameras)
		ret[j] = &r
		err = append(err, e...)
		j++
	}
	return
}

func (c viewConfigRead) TransformAndValidate(
	name string,
	cameras []*CameraConfig,
) (ret ViewConfig, err []error) {
	ret = ViewConfig{
		name:         name,
		title:        c.Title,
		allowedUsers: c.AllowedUsers,
	}

	if !nameMatcher.MatchString(ret.name) {
		err = append(err, fmt.Errorf("Views->Name='%s' does not match %s", ret.name, NameRegexp))
	}

	if len(c.Title) < 1 {
		err = append(err, fmt.Errorf("Views->%s->Title must not be empty", name))
	}

	{
		var camerasErr []error
		ret.cameras, camerasErr = c.Cameras.TransformAndValidate(cameras)
		for ce := range camerasErr {
			err = append(err, fmt.Errorf("Views->%s: %s", name, ce))
		}
	}

	if c.ResolutionMaxWidth == nil {
		ret.resolutionMaxWidth = 3840
	} else if *c.ResolutionMaxWidth > 0 {
		ret.resolutionMaxWidth = *c.ResolutionMaxWidth
	} else {
		err = append(err, fmt.Errorf("Views->%s->ResolutionMaxWidth=%d but must be a positive integer", name, *c.ResolutionMaxWidth))
	}

	if c.ResolutionMaxHeight == nil {
		ret.resolutionMaxHeight = 3840
	} else if *c.ResolutionMaxHeight > 0 {
		ret.resolutionMaxHeight = *c.ResolutionMaxHeight
	} else {
		err = append(err, fmt.Errorf("Views->%s->ResolutionMaxHeight=%d but must be a positive integer", name, *c.ResolutionMaxHeight))
	}

	if len(c.RefreshInterval) < 1 {
		// use default 0
		ret.refreshInterval = time.Minute
	} else if refreshInterval, e := time.ParseDuration(c.RefreshInterval); e != nil {
		err = append(err, fmt.Errorf("viewConfig->%s->RefreshInterval='%s' parse error: %s",
			name, c.RefreshInterval, e,
		))
	} else if refreshInterval < 0 {
		err = append(err, fmt.Errorf("viewConfig->%s->RefreshInterval='%s' must be positive",
			name, c.RefreshInterval,
		))
	} else {
		ret.refreshInterval = refreshInterval
	}

	if c.Autoplay != nil && *c.Autoplay {
		ret.autoplay = true
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
