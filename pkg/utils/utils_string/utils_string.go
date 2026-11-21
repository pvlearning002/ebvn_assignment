package utils_string

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/pkg/utils/utils_log"
)

type UtilsString interface {
	GenerateRandomString(length int) (string, error)
}

type utilsString struct {
	cfg      config.Config
	cfgError config.ConfigError
	logger   utils_log.CustomLogger
}

func NewUtilsString(cfg config.Config, cfgError config.ConfigError, logger utils_log.CustomLogger) UtilsString {
	return &utilsString{
		cfg:      cfg,
		cfgError: cfgError,
		logger:   logger,
	}
}

func (self *utilsString) GenerateRandomString(length int) (string, error) {
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
