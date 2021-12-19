# build backend
FROM golang:1.16 as go-builder

RUN apt-get update && apt-get install -y git upx
RUN go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app
COPY . /app

RUN go generate

RUN  VERSION=`git symbolic-ref -q --short HEAD || git describe --tags --exact-match` \
     BUILD_TIME=`date -Is` \
     CGO_ENABLED=0 \
     GOOS=linux \
     go build -ldflags="-s -w -X main.buildVersion=$VERSION -X main.buildTime=$BUILD_TIME" \
     -o /go-webcam

RUN upx /go-webcam

# used to generate /etc/passwd /etc/group
RUN groupadd -g 1000 app && \
    useradd -r -u 1000 -g app app
RUN chown app:app /go-webcam

# use already compiled frontend
FROM koestler/js-webcam AS js-builder

# build final image
FROM scratch
USER app
COPY --from=go-builder /go-webcam      /go-webcam
COPY --from=go-builder /etc/group      /etc/group
COPY --from=go-builder /etc/passwd     /etc/passwd
COPY --from=js-builder /frontend-build /frontend-build
CMD ["/go-webcam", "-c", "/config.yaml"]
