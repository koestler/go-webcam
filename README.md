# go-webcams

[![Docker Image CI](https://github.com/koestler/go-webcam/actions/workflows/docker-image.yml/badge.svg)](https://github.com/koestler/go-webcam/actions/workflows/docker-image.yml)

The goal of this project is to show still images of IP security cameras on the web.

This daemon is a http server, that fetches images from cameras, caches and scales those images
and serves them in a simple [REST-Api](#api).

There is a [frontend](#frontend) that shows
all the available images and can also autoplay them (reload when new image is available).
Specific resolutions of a webcam can be made available publicly or only behind a login
(eg. high resolution on when logged in, low resolution thumbnail for everyone).

The daemon reads a [configuration](#configuration) written in YAML to configure:
* Project title etc.
* What cameras are available and how to acces them.
* How often a specific camera or a resized image is computed.
* If the image / resolution is public or only available for certain users.
* Some technical stuff.

The tool is written with relatively low cpu-resources of embedded systems
(eg. a [Raspberry Pi](https://www.raspberrypi.com/) or a [PC Engines APU Board](https://pcengines.ch/))
in mind. It therefore supports an in-memory caching mechanism to only fetch images once per timeframe,
to never redo rescaling and therefore is able to serve hundreds of clients on such an embedded system.

It uses [HTTP Cache-Control headers](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control) 
and [HTTP conditional requests](https://developer.mozilla.org/en-US/docs/Web/HTTP/Conditional_requests)
to minimize bandwidth-requirements. This is especially useful when the webcam and the computer running
this software is connected over a slow (eg. LTE) connection to a webserver running as a
[reverse-proxy](https://en.wikipedia.org/wiki/Reverse_proxy).
The webserver can than deliver the same image to many clients while it is only transmitted
once over the slow connection. This even works for images only available behind a login.
In this case, images are sent only once to the reverse proxy
and only small redirect / authentication requests are sent for each user.

## Basic Usage

### Running in docker-compose
I use [docker-compose](https://docs.docker.com/compose/) to deploy this tool.
This has the advantage, that autostart on boot, running as non-root-user, monitoring
and log-rotation are all taken care of by docker.
There is a precompiled docker image (approx. 10 MB) on [dockerhub](https://hub.docker.com/r/koestler/go-webcam),
which I normally use for deployment.
An example docker-compose is here:

```yaml
# documentation/docker-compose.yml

version: "3"
services:
  go-webcam:
    restart: always
    image: ghcr.io/koestler/go-webcam:v1
    volumes:
      - ${PWD}/config.yaml:/config.yaml:ro
      # - ${PWD}/auth.passwd:/auth.passwd:ro
    ports:
      - "80:8043"
```

Setup like this:
```bash
# create configuration files
mkdir -p /srv/dc/webcam
cd /srv/dc/webcam
curl https://raw.githubusercontent.com/koestler/go-webcam/main/documentation/docker-compose.yml -o docker-compose.yml
curl https://raw.githubusercontent.com/koestler/go-webcam/main/documentation/config.yaml -o config.yaml
# edit config.yaml

# create htpasswd file
sudo apt-get install apache2-utils
htpasswd -c auth.passwd user

# start it
docker-compose up -d
```

### Run using docker
```bash
docker run --rm --name go-webcam \
  -p 80:8043 \
  -v "$(pwd)"/config.yaml:/config.yaml:ro \
  -v "$(pwd)"/auth.passwd:/auth.passwd:ro \
  koestler/go-webcam
```

### Basic Usage of the binary
```txt
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

Return Codes:
```txt
ExitSuccess          = 0
ExitDueToCmdOptions  = 1
ExitDueToConfig      = 2
ExitDueToModuleStart = 3
```

## Frontend
The frontend is a client-side application based on [React](https://reactjs.org/).
It is developed in a separate [repository](https://github.com/koestler/js-webcam)
but normally bundled into the build of this project.

## Config
The Configuration is stored in one yaml file. This file is only read once when the server is started.
Restart the backend whenever you change something.
There are mandatory fields and there are optional fields which have a default value.
Whenever a mandatory field is missing or an invalid value is given, the backend refuses to start.

### Complete, explained example

```yaml
# documentation/config.yaml

Version: 0
ProjectTitle: Configurable Title of Project
LogConfig: True                                            # optional, default False, outputs the configuration including defaults on startup
LogWorkerStart: True

Auth:
  HtaccessFile: ./auth.passwd

HttpServer:
  Bind: 0.0.0.0                                            # optional, default ::1 (ipv6 loopback)
  Port: 8043                                               # optional, default 8043
  LogRequests: True
  LogAuth: True                                            # optional, default False, log when login is successful / fails

Cameras:
  0-cam-east:
    Address: rtsps://192.168.1.100:7441/DGGXXX3487348?enableSrtp
    RefreshInterval: 10s

  1-cam-north:
    Address: rtsps://192.168.1.101:7441/DGGXXX3487348?enableSrtp
    RefreshInterval: 10s


Views:
  - Name: low
    Title: Low Resolution
    Cameras:
      - Name: 0-cam-east
        title: Camera East
      - Name: 1-cam-north
        Title: Camera North
    ResolutionMaxWidth: 480
    RefreshInterval: 2s
  - Name: highres
    Title: High Resolution
    Cameras:
      - Name: 0-cam-east
        Title: Camera East
      - Name: 1-cam-north
        Title: Camera North
    ResolutionMaxWidth: 1024
    RefreshInterval: 2s
    AllowedUsers:
      - tester0
```

### Minimalistic example
The following example shows the minimal configuration that only specifies the mandatory configuration
fields. Start with this one and override the defaults.

```yaml
# documentation/minimal-config.yaml

Version: 0

HttpServer:
  Port: 8043                                               # optional, default 8043

Cameras:
  0-cam-east:
    Address: 192.168.8.63
    User: ubnt
    Password: my-password-1234
    RefreshInterval: 10s

Views:
  - Name: raw
    Title: Full Resolution
    Cameras:
      - Name: 0-cam-east
        Title: Camera East

```

### JwtSecret
The JwtSecret is optional. When it is missing, a random secret is generated on every startup of the
backend. This causes all users to be logged out whenever the backend is restarted.
To avoid this, it can be fixed via the configuration.
Most easily, you start the backend the first time with `LogDebug: True`
and copy the randomly generated secret into the configuration file.

### HashSecret
In order to allow reverse proxies to cache images even when authentication is used, all authenticated
images requests are redirected to an unauthenticated, unguessable URL of the image.
This URL includes a hash of the uuid uf the image plus the `HashSecret`. To avoid making
all cached images obsolete after a restart of the backend, this secret should be hardcoded.
Most easily, you start the backend the first time with `LogDebug: True`
and copy the randomly generated secret into the configuration file.

## Cameras

### Unifi
Login to the Unifi Protect controller and in the camera settings "Enable Secure RTSPS Output" and copy the
returned URL into the `Address` field of the camera configuration.

The requests to the cameras are encrypted however validation of the camera's
server certificate is always skipped.

## Authentication
The user/password database is stored in a single file in the format of the apache `htpasswd` tool.
The file can is reloaded automatically.

Use `htpasswd` to generate password files like this:
```bash
sudo apt install apache2-utils
htpasswd -c auth.passwd username
```

## Local Development

### Install dependencies
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### Compile and run on host
```bash
go generate && go build && ./go-webcam
```

### Compile and run inside docker
```bash
docker build -f docker/Dockerfile -t go-webcam .
docker run --rm --name go-webcam -p 127.0.0.1:8043:8043 \
  -v "$(pwd)"/config.yaml:/config.yaml:ro \
  -v "$(pwd)"/auth.passwd:/auth.passwd:ro \
  go-webcam
```

### Frontend development
Whe developing on the frontend, the backend can be configured to be a reverse-proxy
instead of serving a static frontend build.

```yaml
HttpServer:
  FrontendProxy: "http://127.0.0.1:3000/"
```

### Update README.md
```bash
npx embedme README.md
```

## Production build
### Install dependencies
Buildx must be installed: https://docs.docker.com/buildx/working-with-buildx/
On Linux, binfmt_misc needs to be installed and a builder needs to be created:
```bash
docker run --privileged --rm tonistiigi/binfmt --install all
docker buildx create --name mbuilder
docker buildx use mbuilder
```

### Local Production Build
Build:
```bash
docker buildx build --load --platform linux/amd64 -f docker/Dockerfile -t ghcr.io/koestler/go-webcam .
docker buildx build --load --platform linux/arm64 -f docker/Dockerfile -t ghcr.io/koestler/go-webcam .
docker buildx build --load --platform linux/arm/v7 -f docker/Dockerfile -t ghcr.io/koestler/go-webcam .
```

Test:
```bash
docker run --rm --name go-webcam -p 127.0.0.1:8043:8043 -v "$(pwd)"/config.yaml:/app/config.yaml:ro ghcr.io/koestler/go-webcam
```

### Dockerhub Production amd64/arm64 Build
This is for testing only. Production builds are generated by Github Actions.
```bash
docker buildx build --push --platform linux/arm64,linux/amd64 -f docker/Dockerfile -t koestler/go-webcam:tdev .
```

## License
[MIT License](LICENSE)

## Contributing
This is a private project currently maintained by one person only. Therfore only the cameras my friends
and I own are supported. However I'm happy to receive bug reports via github.
I'm also happy to merge pull requests (eg. support for other cameras) and change this section
to give credits to others.