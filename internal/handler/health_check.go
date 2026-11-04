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
	svc service.HealthCheck
	cfg config.Config
}

func NewHealthCheck(svc service.HealthCheck, cfg config.Config) HealthCheck {
	return &healthCheckHandler{svc: svc, cfg: cfg}
}

func (self *healthCheckHandler) GenerateHealthCheck(c *gin.Context) {
	modelHealthCheck, err := self.svc.GenerateHealthCheck()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{self.cfg.ErrorTitle: self.cfg.ErrorInternalServer})
		return
	}
	c.JSON(http.StatusOK, modelHealthCheck)
}
