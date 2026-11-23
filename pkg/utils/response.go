package utils

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
)

type ResponseMessage struct {
	Title   string `json:"error_title"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type ResponseUtils interface {
	GetInputFieldError(err error) ResponseMessage
	GetErrorInternalServer() ResponseMessage
	GetErrorInvalidInput() ResponseMessage
	GetErrorUrlNotFound() ResponseMessage
}

type responseUtils struct {
	cfgError config.ConfigError
}

func NewResponseUtils(cfgError config.ConfigError) ResponseUtils {
	return &responseUtils{
		cfgError: cfgError,
	}
}

func (self *responseUtils) GetInputFieldError(err error) ResponseMessage {
	if is_ok := errors.As(err, &validator.ValidationErrors{}); !is_ok {
		return self.GetErrorInvalidInput()
	}
	var errs []string
	for _, err := range err.(validator.ValidationErrors) {
		errs = append(errs, err.Field()+" is invalid ("+err.Tag()+")")
	}
	return ResponseMessage{
		Title:   self.cfgError.ErrorTitle,
		Message: self.cfgError.ErrorInvalidInput,
		Details: errs,
	}
}

func (self *responseUtils) GetErrorInternalServer() ResponseMessage {
	return ResponseMessage{
		Title:   self.cfgError.ErrorTitle,
		Message: self.cfgError.ErrorInternalServer,
	}
}

func (self *responseUtils) GetErrorInvalidInput() ResponseMessage {
	return ResponseMessage{
		Title:   self.cfgError.ErrorTitle,
		Message: self.cfgError.ErrorInvalidInput,
	}
}

func (self *responseUtils) GetErrorUrlNotFound() ResponseMessage {
	return ResponseMessage{
		Title:   self.cfgError.ErrorTitle,
		Message: self.cfgError.ErrorUrlNotFound,
	}
}
