package service

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/pvlearning002/ebvn_assignment/internal/app/model"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/pkg/utils"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type healthCheckRedisUtilsMock struct {
	containerRunning bool
}

func (m *healthCheckRedisUtilsMock) NewRedisClient(config.ConfigRedis) (*redis.Client, error) {
	return nil, nil
}
func (m *healthCheckRedisUtilsMock) DockerInspect() ([]byte, error) { return nil, nil }
func (m *healthCheckRedisUtilsMock) IsRedisConnect(context.Context) (bool, error) {
	return false, nil
}
func (m *healthCheckRedisUtilsMock) IsRedisContainerRunning() bool {
	return m.containerRunning
}
func (m *healthCheckRedisUtilsMock) StoreStringValue(context.Context, string, string) error {
	return nil
}
func (m *healthCheckRedisUtilsMock) GetStringValue(context.Context, string) (string, error) {
	return "", nil
}

var _ utils.RedisClientUtils = (*healthCheckRedisUtilsMock)(nil)

func TestNewHealthCheck(t *testing.T) {
	got := NewHealthCheck(
		config.Config{},
		config.ConfigStatus{},
		config.ConfigHealthCheck{},
		&healthCheckRedisUtilsMock{},
	)
	require.NotNil(t, got)
	assert.IsType(t, &healthCheckService{}, got)
}

func TestHealthCheckService_GenerateHealthCheck(t *testing.T) {
	validID := uuid.New()
	tests := []struct {
		name           string
		instanceID     string
		redisRunning   bool
		wantOutput     *model.HealthCheckOutput
		wantParseError bool
	}{
		{
			name:         "valid UUID and Redis running",
			instanceID:   validID.String(),
			redisRunning: true,
			wantOutput: &model.HealthCheckOutput{
				Message: "OK", ServiceName: "bookmark-service",
				InstanceID: validID, IsRedisConnected: true,
			},
		},
		{
			name:         "valid UUID and Redis stopped",
			instanceID:   validID.String(),
			redisRunning: false,
			wantOutput: &model.HealthCheckOutput{
				Message: "OK", ServiceName: "bookmark-service",
				InstanceID: validID, IsRedisConnected: false,
			},
		},
		{name: "invalid UUID", instanceID: "invalid-uuid", wantParseError: true},
		{name: "empty UUID", instanceID: "", wantParseError: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := NewHealthCheck(
				config.Config{InstanceID: tc.instanceID},
				config.ConfigStatus{StatusOK: "OK"},
				config.ConfigHealthCheck{BookmarkService: "bookmark-service"},
				&healthCheckRedisUtilsMock{containerRunning: tc.redisRunning},
			)

			got, err := service.GenerateHealthCheck()
			if tc.wantParseError {
				require.Error(t, err)
				assert.Nil(t, got)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantOutput, got)
		})
	}
}
