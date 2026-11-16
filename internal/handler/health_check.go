package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/internal/service"
)

type HealthCheck interface {
	GenerateHealthCheck(c *gin.Context)
}

type healthCheckHandler struct {
	svc      service.HealthCheck
	cfg      config.Config
	cfgError config.ConfigError
}

func NewHealthCheck(svc service.HealthCheck, cfg config.Config, cfgError config.ConfigError) HealthCheck {
	return &healthCheckHandler{svc: svc, cfg: cfg, cfgError: cfgError}
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
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{self.cfgError.ErrorTitle: self.cfgError.ErrorInternalServer})
		return
	}
	c.JSON(http.StatusOK, modelHealthCheck)
}
