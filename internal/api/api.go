package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/internal/handler"
	service "github.com/pvlearning002/ebvn_assignment/internal/service"
)

type Engine interface {
	Start() error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type engine struct {
	app *gin.Engine
	cfg *config.Config
}

func NewEngine(cfg *config.Config) Engine {
	app := &engine{
		app: gin.Default(),
		cfg: cfg,
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
	genPassHandler := handler.NewGenPass(genPassService, *self.cfg)
	// Register gen pass route
	self.app.GET(self.cfg.PathGetGenPass, genPassHandler.GeneratePassword)

	// // Initialize health check service
	// healthCheckService := service.NewHealthCheck()
	// // Initialize health check handler
	// healthCheckHandler := handler.NewHealthCheck(healthCheckService)
	// // Register health check route
	// e.app.GET(config.PATH_GET_HEALTH_CHECK, healthCheckHandler.GenerateHealthCheck)
}
