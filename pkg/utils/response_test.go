package utils

import (
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
)

func TestNewResponseUtils(t *testing.T) {
	cfgError := config.ConfigError{
		ErrorInternalServer: "internal server error",
		ErrorInvalidInput:   "invalid input",
	}

	got := NewResponseUtils(cfgError)
	if got == nil {
		t.Fatal("NewResponseUtils() returned nil")
	}
	if _, ok := got.(*responseUtils); !ok {
		t.Fatalf("NewResponseUtils() returned %T, want *responseUtils", got)
	}
}

func TestInputFieldError_NonValidationError(t *testing.T) {
	cfgError := config.ConfigError{ErrorInternalServer: "internal server error"}
	utils := NewResponseUtils(cfgError).(*responseUtils)

	got := utils.InputFieldError(errors.New("database unavailable"))
	if got.Message != cfgError.ErrorInternalServer {
		t.Errorf("Message = %q, want %q", got.Message, cfgError.ErrorInternalServer)
	}
	if got.Details != nil {
		t.Errorf("Details = %#v, want nil", got.Details)
	}
}

func TestInputFieldError_ValidationErrors(t *testing.T) {
	type input struct {
		Email string `validate:"required,email"`
		Age   int    `validate:"gte=18"`
	}

	err := validator.New().Struct(input{Email: "not-an-email", Age: 16})
	if err == nil {
		t.Fatal("validator.Struct() returned nil, want validation errors")
	}

	cfgError := config.ConfigError{ErrorInvalidInput: "invalid input"}
	utils := NewResponseUtils(cfgError).(*responseUtils)
	got := utils.InputFieldError(err)

	if got.Message != cfgError.ErrorInvalidInput {
		t.Errorf("Message = %q, want %q", got.Message, cfgError.ErrorInvalidInput)
	}

	details, ok := got.Details.([]string)
	if !ok {
		t.Fatalf("Details has type %T, want []string", got.Details)
	}
	want := []string{"Email is invalid (email)", "Age is invalid (gte)"}
	if len(details) != len(want) {
		t.Fatalf("Details = %#v, want %#v", details, want)
	}
	for i := range want {
		if details[i] != want[i] {
			t.Errorf("Details[%d] = %q, want %q", i, details[i], want[i])
		}
	}
}
