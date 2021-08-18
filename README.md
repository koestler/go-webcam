# go-webcams

This daemon conntets to multiple network cameras and serves all their images as one central server. It can authenticate
at the cameras to fetch images and serve the images in aggregated views to different users.

The tool can also connect to an MQTT server to publish the health state of each webcam and the url where the image can
be fetched.

The tool consists of the following components:

* **httpServer**: *serves* the images as well as a simple frontend to the clients
* **cameraClient**: connects to a camera and *receives* images
* **mqttClient**: connects to a MQTT Server and *send* messages

## Basic Usage

```
Usage:
  go-webcam [-c <path to yaml config file>]

Application Options:
      --version     Print the build version and timestamp
  -c, --config=     Config File in yaml format (default: ./config.yaml)
      --cpuprofile= write cpu profile to <file>
      --memprofile= write memory profile to <file>

Help Options:
  -h, --help        Show this help message
```

## Config

The Configuration is stored in one yaml file. There are mandatory fields and there are optional fields which have a
default value.

### Complete, explained example

```yaml
Version: 0                                                 # mandatory, version is always 0 (reserved for later use)
LogConfig: True                                            # optional, default False, outputs the configuration including defaults on startup
LogWorkerStart: True                                       # optional, default False, write log for starting / stoping of workers
LogMqttDebug: False                                        # optional, default False, enable debug output of the mqtt module
HttpServer:                                                # optional, default Disabled, start the http server
  Bind: 0.0.0.0                                            # optional, default ::1 (ipv6 loopback)
  Port: 80                                                 # optional, default 8042
  LogRequests: True                                        # optional, default False, log all requests to stdout

MqttClients:                                               # mandatory, a list of MQTT servers to connect to
  0-piegn-mosquitto:                                       # mandatory, an arbitrary name used in log outputs and for reference in the converters section
    Broker: "tcp://mqtt.exampel.com:1883"                  # mandatory, the address / port of the server
    User: Bob                                              # optional, if given used for login
    Password: Jeir2Jie4zee                                 # optional, if given used for login
    ClientId: "config-tester"                              # optional, default go-webcam, client-id sent to the server
    Qos: 2                                                 # optional, default 0, QOS-level used for subscriptions
    AvailabilityTopic: test/%Prefix%tele/%clientId%/LWT    # optional, if given, a message with Online/Offline will be published on connect/disconnect
                                                           # supported placeholders:
                                                           # - %Prefix$   : as specified in this config section
                                                           # - %clientId% : as specified in this config section
    TopicPrefix: piegn/                                    # optional, default empty
    LogMessages: False                                     # optional, default False, logs all received messages

  1-local-mosquitto:                                       # optional, a second MQTT erver
    Broker: "tcp://172.17.0.5:1883"                        # optional, the second MQTT servers broker...

CameraClients:
  0-cam-east:
    Host: 192.168.8.32
    Implementation: ubnt
    User: ubnt
    Password: 1234

  1-cam-north:
    Host: 192.168.8.33
    Implementation: ubnt
    User: ubnt
    Password: abcde

Views:                                                     # mandatory, a list of Views that shall be available
  public:
    Route: /
    Cameras:
      - cam0
    ResolutionMaxWidth: 1024
    ResolutionMaxHeight: 768
  private:
    Route: /all/
    Cameras:
      - cam0
      - cam1
    ResolutionMaxWidth: 1920
    ResolutionMaxHeight: 1080
    Htaccess:
      - "user:password"
```  

## Local Production Build
```
docker build -f docker/Dockerfile -t go-webcam .
docker run --rm --name go-webcam -p 127.0.0.1:8043:8043 -v "$(pwd)"/config.yaml:/app/config.yaml:ro go-webcam
```

## Dockerhub Production Build
```
docker buildx build --platform linux/arm64 -f docker/Dockerfile -t koestler/go-webcam .
docker buildx build --platform linux/amd64 -f docker/Dockerfile -t koestler/go-webcam .
docker push koestler/go-webcam

```

# License

MIT License