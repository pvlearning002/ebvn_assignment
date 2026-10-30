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
}

// NewGenPass creates a new instance of the genPassHandler with the provided GenPass service and configuration.
func NewGenPass(genPassSvc service.GenPass, cfg config.Config) GenPass {
	return &genPassHandler{
		genPassSvc: genPassSvc,
		cfg:        cfg,
	}
}

// GeneratePassword handles the HTTP request for generating a password.
func (self *genPassHandler) GeneratePassword(c *gin.Context) {
	pwd, err := self.genPassSvc.GeneratePassword()
	if err != nil || pwd == self.cfg.StringEmpty {
		c.JSON(http.StatusInternalServerError, gin.H{self.cfg.ErrorTitle: self.cfg.ErrorInternalServer})
		return
	}
	c.JSON(http.StatusOK, gin.H{self.cfg.PasswordTitle: pwd})
}
