FROM golang:1.16-alpine as builder

RUN apk add git

COPY . /app
WORKDIR /app

RUN  VERSION=`git symbolic-ref -q --short HEAD || git describe --tags --exact-match` \
     BUILD_TIME=`date -Is` \
     GOOS=linux \
     go build -ldflags="-s -w -X main.buildVersion=$VERSION -X main.buildTime=$BUILD_TIME" \
     -o /go-webcam

FROM alpine:latest

RUN  addgroup -S app && adduser -S app -G app
USER app

WORKDIR /app
COPY --from=builder /go-webcam ./
ENV GIN_MODE=release
CMD ["/app/go-webcam", "-c", "/app/config.yaml"]

