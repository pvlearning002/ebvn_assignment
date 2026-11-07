package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/pvlearning002/ebvn_assignment/docs"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/internal/handler"
	service "github.com/pvlearning002/ebvn_assignment/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type engine struct {
	app            *gin.Engine
	cfg            *config.Config
	cfgStatus      *config.ConfigStatus
	cfgError       *config.ConfigError
	cfgPaths       *config.ConfigPaths
	cfgHealthCheck *config.ConfigHealthCheck
}

func NewEngine(cfg *config.Config, cfgStatus *config.ConfigStatus, cfgError *config.ConfigError, cfgPaths *config.ConfigPaths, cfgHealthCheck *config.ConfigHealthCheck) Engine {
	app := &engine{
		app:            gin.Default(),
		cfg:            cfg,
		cfgStatus:      cfgStatus,
		cfgError:       cfgError,
		cfgPaths:       cfgPaths,
		cfgHealthCheck: cfgHealthCheck,
	}
	app.initRoutes()
	return app
}

// Start starts the HTTP server on the specified port.
func (self *engine) Start() error {
	return self.app.Run(":" + strconv.Itoa(self.cfg.AppPort))
}

// ServeHTTP implements the http.Handler interface, allowing the engine to be used as an HTTP handler.
func (self *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self.app.ServeHTTP(w, r)
}

// initRoutes initializes the routes for the API engine.
func (self *engine) initRoutes() {
	// Initialize gen pass service
	genPassService := service.NewGenPass(*self.cfg)
	// Initialize gen pass handler
	genPassHandler := handler.NewGenPass(genPassService, *self.cfg, *self.cfgStatus, *self.cfgError)
	// Register gen pass route
	// self.app.GET(self.cfgPaths.PathSwagger, ginSwagger.WrapHandler(swaggerFiles.Handler))
	self.app.GET(self.cfgPaths.PathSwagger, ginSwagger.WrapHandler(swaggerFiles.Handler))
	self.app.GET(self.cfgPaths.PathGetGenPass, genPassHandler.GeneratePassword)

	// Initialize health check service
	healthCheckService := service.NewHealthCheck(*self.cfgStatus, *self.cfgHealthCheck)
	// Initialize health check handler
	healthCheckHandler := handler.NewHealthCheck(healthCheckService, *self.cfg, *self.cfgError)
	// Register health check route
	self.app.GET(self.cfgPaths.PathGetHealthCheck, healthCheckHandler.GenerateHealthCheck)
}
