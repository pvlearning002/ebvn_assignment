package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/pvlearning002/ebvn_assignment/internal/app/model"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthCheckService_GenerateHealthCheck(t *testing.T) {
	validUUID := uuid.New()

	tests := []struct {
		name          string
		instanceID    string
		expectedError bool
		expected      *model.HealthCheckOutput
	}{
		{
			name:          "Success - valid UUID",
			instanceID:    validUUID.String(),
			expectedError: false,
			expected: &model.HealthCheckOutput{
				Message:     "OK",
				ServiceName: "BookmarkService",
				InstanceID:  validUUID,
			},
		},
		{
			name:          "Failed - invalid UUID",
			instanceID:    "invalid-uuid",
			expectedError: true,
			expected:      nil,
		},
		{
			name:          "Failed - empty UUID",
			instanceID:    "",
			expectedError: true,
			expected:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{
				InstanceID: tt.instanceID,
			}

			cfgStatus := config.ConfigStatus{
				StatusOK: "OK",
			}

			cfgHealthCheck := config.ConfigHealthCheck{
				BookmarkService: "BookmarkService",
			}

			service := NewHealthCheck(
				cfg,
				cfgStatus,
				cfgHealthCheck,
			)

			result, err := service.GenerateHealthCheck()

			if tt.expectedError {
				require.Error(t, err)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)

				assert.Equal(t, tt.expected.Message, result.Message)
				assert.Equal(t, tt.expected.ServiceName, result.ServiceName)
				assert.Equal(t, tt.expected.InstanceID, result.InstanceID)
			}
		})
	}
}
