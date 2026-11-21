package utils_response

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
)

type ResponseMessage struct {
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

type CustomResponse interface {
}

type customResponse struct {
	cfgError config.ConfigError
}

func NewCustomResponse(cfgError config.ConfigError) CustomResponse {
	return &customResponse{
		cfgError: cfgError,
	}
}

func (self *customResponse) InputFieldError(err error) ResponseMessage {
	if is_ok := errors.As(err, &validator.ValidationErrors{}); !is_ok {
		return ResponseMessage{
			Message: self.cfgError.ErrorInternalServer,
		}
	}
	var errs []string
	for _, err := range err.(validator.ValidationErrors) {
		errs = append(errs, err.Field()+" is invalid ("+err.Tag()+")")
	}
	return ResponseMessage{
		Message: self.cfgError.ErrorInvalidInput,
		Details: errs,
	}
}
