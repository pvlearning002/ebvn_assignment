package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	// cSpell:ignore pvlearning ebvn
	"github.com/pvlearning002/ebvn_assignment/docs"
	_ "github.com/pvlearning002/ebvn_assignment/docs"
	"github.com/pvlearning002/ebvn_assignment/internal/app/handler"
	repository "github.com/pvlearning002/ebvn_assignment/internal/app/repository"
	service "github.com/pvlearning002/ebvn_assignment/internal/app/service"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/pkg/utils"
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
	cfgParameter   *config.ConfigParameter
	stringUtils    *utils.StringUtils
	redisUtils     *utils.RedisClientUtils
	response       *utils.ResponseUtils
	logger         *utils.LoggerUtils
}

func NewEngine(cfg *config.Config,
	cfgStatus *config.ConfigStatus,
	cfgError *config.ConfigError,
	cfgPaths *config.ConfigPaths,
	cfgHealthCheck *config.ConfigHealthCheck,
	cfgParameter *config.ConfigParameter,
	stringUtils *utils.StringUtils,
	redisUtils *utils.RedisClientUtils,
	response *utils.ResponseUtils,
	logger *utils.LoggerUtils) Engine {
	app := &engine{
		app:            gin.Default(),
		cfg:            cfg,
		cfgStatus:      cfgStatus,
		cfgError:       cfgError,
		cfgPaths:       cfgPaths,
		cfgHealthCheck: cfgHealthCheck,
		cfgParameter:   cfgParameter,
		stringUtils:    stringUtils,
		redisUtils:     redisUtils,
		response:       response,
		logger:         logger,
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
	docs.SwaggerInfo.BasePath = self.cfg.BasePath
	// Initialize gen pass service
	genPassService := service.NewGenPass(*self.cfg, *self.logger, *self.stringUtils)
	// Initialize gen pass handler
	genPassHandler := handler.NewGenPass(genPassService, *self.cfg, *self.logger, *self.response)

	// Initialize health check service
	healthCheckService := service.NewHealthCheck(*self.cfg, *self.cfgStatus, *self.cfgHealthCheck, *self.redisUtils)
	// Initialize health check handler
	healthCheckHandler := handler.NewHealthCheck(healthCheckService, *self.cfg, *self.response, *self.logger)

	urlStorage := repository.NewURLStorage(*self.redisUtils)
	// Initialize shorten URL service
	shortenURLService := service.NewShortenURLService(urlStorage, *self.cfg, *self.logger, *self.stringUtils)
	// Initialize shorten URL handler
	shortenURLHandler := handler.NewShortenURLHandler(shortenURLService, *self.cfg, *self.cfgError, *self.response, *self.cfgParameter, *self.logger)
	// Register Swagger route
	self.app.GET(self.cfgPaths.PathSwagger, ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Register gen pass route
	self.app.GET(self.cfgPaths.PathGetGenPass, genPassHandler.GeneratePassword)
	// Register health check route
	self.app.GET(self.cfgPaths.PathGetHealthCheck, healthCheckHandler.GenerateHealthCheck)
	// Register shorten URL route
	self.app.POST(self.cfgPaths.PathPostShortenURL, shortenURLHandler.CreateShortenLink)
	// Register redirect route
	self.app.GET(self.cfgPaths.PathGetURL, shortenURLHandler.RedirectToOriginalURL)
}
