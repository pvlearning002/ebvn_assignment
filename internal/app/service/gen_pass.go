package service

import (
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/pkg/utils/utils_log"
	"github.com/pvlearning002/ebvn_assignment/pkg/utils/utils_string"
)

// mockery --name GenPass --filename gen_pass_service_mock.go
// GenPass defines the interface for generating passwords.
type GenPass interface {
	GeneratePassword() (string, error)
}

// genPassService is the implementation of the GenPass interface.
type genPassService struct {
	cfg           config.Config
	logger        utils_log.CustomLogger
	custom_string utils_string.UtilsString
}

// NewGenPass creates a new instance of the genPassService with the provided configuration.
func NewGenPass(cfg config.Config, logger utils_log.CustomLogger, custom_string utils_string.UtilsString) GenPass {
	return &genPassService{
		cfg:           cfg,
		logger:        logger,
		custom_string: custom_string,
	}
}

// GeneratePassword generates a random password using UUIDs and returns it as a string.
func (self *genPassService) GeneratePassword() (string, error) {
	result, err := self.custom_string.GenerateRandomString(self.cfg.PasswordLengthDefault)
	if err != nil {
		self.logger.WriteLogError(err.Error())
		return self.cfg.StringEmpty, err
	}
	return result, nil
}
