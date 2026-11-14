package service

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
)

// mockery --name GenPass --filename gen_pass_service_mock.go
// GenPass defines the interface for generating passwords.
type GenPass interface {
	GenerateRandomString(length int) (string, error)
	GeneratePassword() (string, error)
}

// genPassService is the implementation of the GenPass interface.
type genPassService struct {
	cfg      config.Config
	cfgError config.ConfigError
}

// NewGenPass creates a new instance of the genPassService with the provided configuration.
func NewGenPass(cfg config.Config, cfgError config.ConfigError) GenPass {
	return &genPassService{
		cfg:      cfg,
		cfgError: cfgError,
	}
}

func (self *genPassService) GenerateRandomString(length int) (string, error) {
	if length <= 0 {
		return self.cfg.StringEmpty, errors.New(self.cfgError.ErrorLengthZero)
	}

	result := strings.ReplaceAll(uuid.New().String(), self.cfg.Dash, self.cfg.StringEmpty)
	pass_len := len(result)

	for pass_len < length {
		result += strings.ReplaceAll(uuid.New().String(), self.cfg.Dash, self.cfg.StringEmpty)
		pass_len = len(result)
	}

	result = result[:length]
	return result, nil
}

// GeneratePassword generates a random password using UUIDs and returns it as a string.
func (self *genPassService) GeneratePassword() (string, error) {
	result, err := self.GenerateRandomString(self.cfg.PasswordLengthDefault)
	if err != nil {
		return self.cfg.StringEmpty, err
	}
	return result, nil
}
