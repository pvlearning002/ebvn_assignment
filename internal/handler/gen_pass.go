package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/internal/service"
)

// GenPass defines the interface for handling password generation requests.
type GenPass interface {
	GeneratePassword(c *gin.Context)
}

// genPassHandler is the implementation of the GenPass interface.
type genPassHandler struct {
	genPassSvc service.GenPass
	cfg        config.Config
	cfgStatus  config.ConfigStatus
	cfgError   config.ConfigError
}

// NewGenPass creates a new instance of the genPassHandler with the provided GenPass service and configuration.
func NewGenPass(genPassSvc service.GenPass, cfg config.Config, cfgStatus config.ConfigStatus, cfgError config.ConfigError) GenPass {
	return &genPassHandler{
		genPassSvc: genPassSvc,
		cfg:        cfg,
		cfgStatus:  cfgStatus,
		cfgError:   cfgError,
	}
}

// @Summary Generate Random Password
// @Description GeneratePassword handles the HTTP request for generating a password.
// @Tags Password
// @Produces application/json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /generate-password [get]
func (self *genPassHandler) GeneratePassword(c *gin.Context) {
	pwd, err := self.genPassSvc.GeneratePassword()
	if err != nil || len(pwd) != self.cfg.PasswordLengthDefault {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{self.cfgStatus.StatusOK: self.cfgError.ErrorInternalServer})
		return
	}
	c.JSON(http.StatusOK, gin.H{self.cfg.PasswordTitle: pwd})
}
