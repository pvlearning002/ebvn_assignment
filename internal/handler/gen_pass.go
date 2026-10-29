package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/internal/service"
)

type GenPass interface {
	GeneratePassword(c *gin.Context)
}

type genPassHandler struct {
	genPassSvc service.GenPass
	cfg        config.Config
}

func NewGenPass(genPassSvc service.GenPass, cfg config.Config) GenPass {
	return &genPassHandler{
		genPassSvc: genPassSvc,
		cfg:        cfg,
	}
}

func (g *genPassHandler) GeneratePassword(c *gin.Context) {
	pwd, err := g.genPassSvc.GeneratePassword()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{g.cfg.ErrorTitle: g.cfg.ErrorInternalServer})
		return
	}
	c.JSON(http.StatusOK, gin.H{g.cfg.PasswordTitle: pwd})
}
