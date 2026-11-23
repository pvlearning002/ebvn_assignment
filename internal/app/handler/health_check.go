package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pvlearning002/ebvn_assignment/internal/app/service"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/pkg/utils"
)

type HealthCheck interface {
	GenerateHealthCheck(c *gin.Context)
}

type healthCheckHandler struct {
	svc      service.HealthCheck
	cfg      config.Config
	response utils.ResponseUtils
	logger   utils.LoggerUtils
}

func NewHealthCheck(svc service.HealthCheck,
	cfg config.Config,
	response utils.ResponseUtils,
	logger utils.LoggerUtils) HealthCheck {
	return &healthCheckHandler{svc: svc,
		cfg:      cfg,
		response: response,
		logger:   logger}
}

// @Summary Generate Health Check
// @Description Generate Health Check handles the HTTP request for generating a health check.
// @Tags Health
// @Produces application/json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /health [get]
func (self *healthCheckHandler) GenerateHealthCheck(c *gin.Context) {
	modelHealthCheck, err := self.svc.GenerateHealthCheck()
	if err != nil {
		self.logger.WriteLogError(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, self.response.GetErrorInternalServer())
		return
	}
	c.JSON(http.StatusOK, modelHealthCheck)
}
