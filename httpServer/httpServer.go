package httpServer

import (
	"context"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/koestler/go-webcam/cameraClient"
	"github.com/koestler/go-webcam/config"
	"github.com/koestler/go-webcam/hashStore"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type HttpServer struct {
	config Config
	server *http.Server
}

type Environment struct {
	ProjectTitle             string
	Views                    []*config.ViewConfig
	Auth                     config.AuthConfig
	CameraClientPoolInstance *cameraClient.ClientPool
	HashStorage              *hashStore.HashStore
}

type Config interface {
	Bind() string
	Port() int
	LogRequests() bool
	LogDebug() bool
	LogConfig() bool
	EnableDocs() bool
	FrontendProxy() *url.URL
	FrontendPath() string
	GetViewNames() []string
	FrontendExpires() time.Duration
	ConfigExpires() time.Duration
}

func Run(config Config, env *Environment) (httpServer *HttpServer) {
	gin.SetMode("release")
	engine := gin.New()
	if config.LogRequests() {
		engine.Use(gin.Logger())
	}
	engine.Use(gin.Recovery())
	engine.Use(gzip.Gzip(gzip.BestCompression))
	engine.Use(authJwtMiddleware(env))

	if config.EnableDocs() {
		setupSwaggerDocs(engine, config)
	}
	addApiV0Routes(engine, config, env)
	setupFrontend(engine, config)

	server := &http.Server{
		Addr:    config.Bind() + ":" + strconv.Itoa(config.Port()),
		Handler: engine,
	}

	go func() {
		if config.LogDebug() {
			log.Printf("httpServer: listening on %v", server.Addr)
		}
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

func addApiV0Routes(r *gin.Engine, config Config, env *Environment) {
	v0 := r.Group("/api/v0/")
	setupConfig(v0, config, env)
	setupLogin(v0, config, env)
	setupImagesByHash(v0, config, env)
	setupImages(v0, config, env)
}
