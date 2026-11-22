package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pvlearning002/ebvn_assignment/internal/app/service"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/pkg/utils"
)

// GenPass defines the interface for handling password generation requests.
type GenPass interface {
	GeneratePassword(c *gin.Context)
}

// genPassHandler is the implementation of the GenPass interface.
type genPassHandler struct {
	genPassSvc service.GenPass
	cfg        config.Config
	cfgError   config.ConfigError
	logger     utils.LoggerUtils
}

// NewGenPass creates a new instance of the genPassHandler with the provided GenPass service and configuration.
func NewGenPass(genPassSvc service.GenPass, cfg config.Config, cfgError config.ConfigError, logger utils.LoggerUtils) GenPass {
	return &genPassHandler{
		genPassSvc: genPassSvc,
		cfg:        cfg,
		cfgError:   cfgError,
		logger:     logger,
	}
}

// @Summary Generate Random Password
// @Description GeneratePassword handles the HTTP request for generating a password.
// @Tags Password
// @Produces application/json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /genpass [get]
func (self *genPassHandler) GeneratePassword(c *gin.Context) {
	pwd, err := self.genPassSvc.GeneratePassword()
	if err != nil {
		self.logger.WriteLogError(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{self.cfgError.ErrorTitle: self.cfgError.ErrorInternalServer})
		return
	}
	c.JSON(http.StatusOK, gin.H{self.cfg.PasswordTitle: pwd})
}
