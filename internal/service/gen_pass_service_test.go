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
		expectError    bool
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
		{
			name:           "zero length",
			passwordLength: 0,
			expectError:    true,
		},
		{
			name:           "negative length",
			passwordLength: -1,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := *config.CreateNewConfig()
			cfgError := *config.CreateNewConfigError()
			// Update password length for each test case.
			cfg.PasswordLengthDefault = tc.passwordLength
			service := &genPassService{
				cfg:      cfg,
				cfgError: cfgError,
			}

			result, err := service.GenerateRandomString(tc.passwordLength)

			if tc.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
				assert.Equal(t, cfgError.ErrorLengthZero, err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Len(t, result, tc.passwordLength)
			assert.NotContains(t, result, "-")
			assert.NotEmpty(t, result)
		})
	}
}

func TestGeneratePassword_InvalidLength(t *testing.T) {
	cfg := *config.CreateNewConfig()
	cfgError := *config.CreateNewConfigError()

	cfg.PasswordLengthDefault = 0

	service := &genPassService{
		cfg:      cfg,
		cfgError: cfgError,
	}

	result, err := service.GeneratePassword()

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.Equal(t, cfgError.ErrorLengthZero, err.Error())
}
