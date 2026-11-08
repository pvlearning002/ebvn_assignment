package service

import (
	"testing"

	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGeneratePassword(t *testing.T) {
	testCases := []struct {
		name           string
		passwordLength int
	}{
		{
			name:           "password length less than UUID length",
			passwordLength: 12,
		},
		{
			name:           "password length equal to UUID length",
			passwordLength: 32,
		},
		{
			name:           "password length greater than UUID length",
			passwordLength: 36,
		},
		{
			name:           "password length requires multiple UUIDs",
			passwordLength: 64,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := *config.CreateNewConfig()
			// Update password length for each test case.
			cfg.PasswordLengthDefault = tc.passwordLength
			service := &genPassService{
				cfg: cfg,
			}

			result, err := service.GeneratePassword()

			assert.NoError(t, err)
			assert.Equal(t, tc.passwordLength, len(result))
			assert.NotContains(t, result, "-")
			assert.NotEmpty(t, result)
		})
	}
}
