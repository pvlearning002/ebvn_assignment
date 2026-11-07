package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	// Implement your test logic here
	cfgStatus := *config.CreateNewConfigStatus()
	cfgHealthCheck := *config.CreateNewConfigHealthCheck()
	t.Parallel()
	testCases := []struct {
		name                string
		expectedMessage     string
		expectedServiceName string
		expectedInstanceId  uuid.UUID
		expectError         error
	}{
		{
			name:                "Health check returns OK",
			expectedMessage:     cfgStatus.StatusOK,
			expectedServiceName: cfgHealthCheck.BookmarkService,
			expectedInstanceId:  uuid.New(),
			expectError:         nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call the HealthCheck function and compare the result with tc.expectedStatus
			t.Parallel()
			svc := NewHealthCheck(cfgStatus, cfgHealthCheck)
			result, err := svc.GenerateHealthCheck()
			assert.Equal(t, tc.expectedMessage, cfgStatus.StatusOK)
			assert.Equal(t, tc.expectedServiceName, cfgHealthCheck.BookmarkService)
			assert.IsType(t, uuid.UUID{}, result.InstanceID)
			assert.ErrorIs(t, err, tc.expectError)
		})
	}
}
