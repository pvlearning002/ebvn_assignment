package service

import (
	"strings"

	"github.com/google/uuid"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
)

// mockery --name GenPass --filename gen_pass_service_mock.go
// GenPass defines the interface for generating passwords.
type GenPass interface {
	GeneratePassword() (string, error)
}

// genPassService is the implementation of the GenPass interface.
type genPassService struct {
	cfg config.Config
}

// NewGenPass creates a new instance of the genPassService with the provided configuration.
func NewGenPass(cfg config.Config) GenPass {
	return &genPassService{
		cfg: cfg,
	}
}

// GeneratePassword generates a random password using UUIDs and returns it as a string.
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
