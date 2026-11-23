package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pvlearning002/ebvn_assignment/internal/app/model"
	"github.com/pvlearning002/ebvn_assignment/internal/app/service"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type healthCheckHandlerServiceMock struct{ mock.Mock }

func (m *healthCheckHandlerServiceMock) GenerateHealthCheck() (*model.HealthCheckOutput, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	result, ok := args.Get(0).(*model.HealthCheckOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return result, args.Error(1)
}

var _ service.HealthCheck = (*healthCheckHandlerServiceMock)(nil)

func TestNewHealthCheck(t *testing.T) {
	got := NewHealthCheck(
		new(healthCheckHandlerServiceMock),
		config.Config{},
		utils.NewResponseUtils(config.ConfigError{}),
		utils.NewLoggerUtils(),
	)
	require.NotNil(t, got)
	assert.IsType(t, &healthCheckHandler{}, got)
}

func TestHealthCheckHandler_GenerateHealthCheck_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	instanceID := uuid.New()
	want := &model.HealthCheckOutput{
		Message:          "OK",
		ServiceName:      "bookmark-service",
		InstanceID:       instanceID,
		IsRedisConnected: true,
	}
	serviceMock := new(healthCheckHandlerServiceMock)
	serviceMock.On("GenerateHealthCheck").Return(want, nil).Once()
	handler := NewHealthCheck(serviceMock, config.Config{}, utils.NewResponseUtils(config.ConfigError{}), utils.NewLoggerUtils())
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/health", nil)

	handler.GenerateHealthCheck(ctx)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.JSONEq(t, `{"message":"OK","service_name":"bookmark-service","instance_id":"`+instanceID.String()+`","is_redis_connected":true}`, recorder.Body.String())
	serviceMock.AssertExpectations(t)
}

func TestHealthCheckHandler_GenerateHealthCheck_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	serviceMock := new(healthCheckHandlerServiceMock)
	serviceMock.On("GenerateHealthCheck").Return(nil, errors.New("health check failed")).Once()
	cfgError := config.ConfigError{
		ErrorTitle:          "Error",
		ErrorInternalServer: "Internal Server Error",
	}
	handler := NewHealthCheck(serviceMock, config.Config{}, utils.NewResponseUtils(cfgError), utils.NewLoggerUtils())
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/health", nil)

	handler.GenerateHealthCheck(ctx)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.JSONEq(t, `{"error_title":"Error","message":"Internal Server Error"}`, recorder.Body.String())
	serviceMock.AssertExpectations(t)
}
