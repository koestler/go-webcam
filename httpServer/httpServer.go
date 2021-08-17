package httpServer

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/koestler/go-webcam/cameraClient"
	"github.com/koestler/go-webcam/config"
	"log"
	"net/http"
	"strconv"
	"time"
)

type HttpServer struct {
	config Config
	server *http.Server
}

type Environment struct {
	Views                    []*config.ViewConfig
	CameraClientPoolInstance *cameraClient.ClientPool
}

type Config interface {
	Bind() string
	Port() int
	LogRequests() bool
}

func Run(config Config, env *Environment) (httpServer *HttpServer) {
	engine := gin.New()
	if config.LogRequests() {
		engine.Use(gin.Logger())
	}
	engine.Use(gin.Recovery())

	setupSwaggerDocs(engine, config)
	addRoutes(engine, env)

	server := &http.Server{
		Addr:    config.Bind() + ":" + strconv.Itoa(config.Port()),
		Handler: engine,
	}

	go func() {
		log.Printf("httpServer: listening on %v", server.Addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("httpServer: stopped due to error: %s", err)
		}
	}()

	return &HttpServer{
		config: config,
		server: server,
	}
}

func (s *HttpServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Printf("httpServer: graceful shutdown failed: %s", err)
	}
}
