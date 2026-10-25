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
}

func NewGenPass(genPassSvc service.GenPass) GenPass {
	return &genPassHandler{
		genPassSvc: genPassSvc,
	}
}

func (g *genPassHandler) GeneratePassword(c *gin.Context) {
	pass, err := g.genPassSvc.GeneratePassword(config.PASSWORD_LENGTH_DEFAULT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"password": pass})
}
