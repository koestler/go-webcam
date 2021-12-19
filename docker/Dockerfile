# build backend
FROM golang:1.16-alpine as go-builder

RUN apk add git upx

RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . /app
WORKDIR /app

RUN go get
RUN go generate

RUN  VERSION=`git symbolic-ref -q --short HEAD || git describe --tags --exact-match` \
     BUILD_TIME=`date -Is` \
     CGO_ENABLED=0 \
     GOOS=linux \
     go build -ldflags="-s -w -X main.buildVersion=$VERSION -X main.buildTime=$BUILD_TIME" \
     -o /go-webcam

RUN upx /go-webcam

# used to update /etc/passwd
RUN addgroup -S app && adduser -S app -G app

# build frontend
FROM node:16 as js-builder
ENV GENERATE_SOURCEMAP=false

WORKDIR /app
RUN git clone https://github.com/koestler/js-webcam.git
RUN cd js-webcam && npm install && npm run build

# build final image
FROM scratch
USER app
COPY --from=go-builder /go-webcam           /go-webcam
COPY --from=go-builder /etc/group           /etc/group
COPY --from=go-builder /etc/passwd          /etc/passwd
COPY --from=js-builder /app/js-webcam/build /frontend-build
CMD ["/go-webcam", "-c", "/config.yaml"]
