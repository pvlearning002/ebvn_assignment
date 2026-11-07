package service

import (
	"testing"

	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGeneratePassword(t *testing.T) {
	cfg := *config.CreateNewConfig()
	testCases := []struct {
		name           string
		expectedLength int
		expectError    error
	}{
		{
			name:           "success",
			expectedLength: cfg.PasswordLengthDefault,
			expectError:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testSvc := NewGenPass(cfg)
			pass, err := testSvc.GeneratePassword()
			assert.NotNil(t, pass)
			assert.Equal(t, tc.expectError, err)
			assert.Equal(t, tc.expectedLength, len(pass))
		})
	}
}
