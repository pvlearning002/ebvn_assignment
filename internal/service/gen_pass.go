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

func (self *genPassService) GeneratePassword() (string, error) {
	result := strings.ReplaceAll(uuid.New().String(), self.cfg.Dash, self.cfg.StringEmpty)
	pass_len := len(result)
	for pass_len < self.cfg.PasswordLengthDefault {
		result += strings.ReplaceAll(uuid.New().String(), self.cfg.Dash, self.cfg.StringEmpty)
		pass_len = len(result)
	}
	result = result[:self.cfg.PasswordLengthDefault]
	return result, nil
}
