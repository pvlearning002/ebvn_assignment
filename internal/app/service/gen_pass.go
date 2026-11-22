package service

import (
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/pkg/utils"
)

// mockery --name GenPass --filename gen_pass_service_mock.go
// GenPass defines the interface for generating passwords.
type GenPass interface {
	GeneratePassword() (string, error)
}

// genPass is the implementation of the GenPass interface.
type genPass struct {
	cfg         config.Config
	logger      utils.LoggerUtils
	stringUtils utils.StringUtils
}

// NewGenPass creates a new instance of the genPassService with the provided configuration.
func NewGenPass(
	cfg config.Config,
	logger utils.LoggerUtils,
	stringUtils utils.StringUtils) GenPass {
	return &genPass{
		cfg:         cfg,
		logger:      logger,
		stringUtils: stringUtils,
	}
}

// GeneratePassword generates a random password using UUIDs and returns it as a string.
func (self *genPass) GeneratePassword() (string, error) {
	result, err := self.stringUtils.GenerateRandomString(self.cfg.PasswordLengthDefault)
	if err != nil {
		self.logger.WriteLogError(err.Error())
		return self.cfg.StringEmpty, err
	}
	return result, nil
}
