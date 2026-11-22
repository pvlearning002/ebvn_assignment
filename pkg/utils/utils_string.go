package utils

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
)

type StringUtils interface {
	GenerateRandomString(length int) (string, error)
}

type stringUtils struct {
	cfg      config.Config
	cfgError config.ConfigError
	logger   LoggerUtils
}

func NewStringUtils(cfg config.Config, cfgError config.ConfigError, logger LoggerUtils) StringUtils {
	return &stringUtils{
		cfg:      cfg,
		cfgError: cfgError,
		logger:   logger,
	}
}

func (self *stringUtils) GenerateRandomString(length int) (string, error) {
	if length <= 0 {
		self.logger.WriteLogError(self.cfgError.ErrorLengthZero)
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
