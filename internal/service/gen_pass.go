package service

import (
	"strings"

	"github.com/google/uuid"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
)

type GenPass interface {
	GeneratePassword() (string, error)
}

type genPassService struct {
	cfg config.Config
}

func NewGenPass(cfg config.Config) GenPass {
	return &genPassService{
		cfg: cfg,
	}
}

func (g *genPassService) GeneratePassword() (string, error) {
	stringEmpty := ""
	dash := g.cfg.Dash
	result := strings.ReplaceAll(uuid.New().String(), dash, stringEmpty)
	pass_len := len(result)
	for pass_len < g.cfg.PasswordLengthDefault {
		result += strings.ReplaceAll(uuid.New().String(), dash, stringEmpty)
		pass_len = len(result)
	}
	result = result[:g.cfg.PasswordLengthDefault]
	return result, nil
}
