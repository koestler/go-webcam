# build backend
FROM golang:1.16-alpine as go-builder

RUN apk add git

COPY . /app
WORKDIR /app

RUN go get -u github.com/swaggo/swag/cmd/swag
RUN go generate

RUN  VERSION=`git symbolic-ref -q --short HEAD || git describe --tags --exact-match` \
     BUILD_TIME=`date -Is` \
     GOOS=linux \
     go build -ldflags="-s -w -X main.buildVersion=$VERSION -X main.buildTime=$BUILD_TIME" \
     -o /go-webcam

# build frontend
FROM node:16 as js-builder

WORKDIR /app
RUN git clone https://github.com/koestler/js-webcam.git
RUN cd js-webcam && npm install && npm run build

# build final image
FROM alpine:latest

RUN  addgroup -S app && adduser -S app -G app
USER app

WORKDIR /app
COPY --from=go-builder /go-webcam ./
COPY --from=js-builder /app/js-webcam/build ./frontend-build

CMD ["/app/go-webcam", "-c", "/app/config.yaml"]

