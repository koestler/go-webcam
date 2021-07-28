package main

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/koestler/go-webcam/config"
	"github.com/koestler/go-webcam/mqttClient"
	"log"
	"os"
)

func runMqttClient(
	cfg *config.Config,
	initiateShutdown chan<- error,
) map[string]*mqttClient.MqttClient {
	mqtt.ERROR = log.New(os.Stdout, "MqttDebugLog: ", log.LstdFlags)
	if cfg.LogMqttDebug {
		mqtt.DEBUG = log.New(os.Stdout, "MqttDebugLog: ", log.LstdFlags)
	}

	mqttClientInstances := make(map[string]*mqttClient.MqttClient)

	for _, mqttClientConfig := range cfg.MqttClients {
		if cfg.LogWorkerStart {
			log.Printf(
				"mqttClient[%s]: start: Broker='%s', ClientId='%s'",
				mqttClientConfig.Name(), mqttClientConfig.Broker(), mqttClientConfig.ClientId(),
			)
		}

		if client, err := mqttClient.Run(mqttClientConfig); err != nil {
			log.Printf("mqttClient[%s]: start failed: %s", mqttClientConfig.Name(), err)
		} else {
			mqttClientInstances[mqttClientConfig.Name()] = client
			log.Printf("mqttClient[%s]: started", mqttClientConfig.Name())
		}
	}

	return mqttClientInstances
}
