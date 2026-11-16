package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/internal/model"
	services_mocks "github.com/pvlearning002/ebvn_assignment/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

type healthCheckMockAdapter struct {
	mock *services_mocks.HealthCheck
	cfg  *config.Config
}

func (a healthCheckMockAdapter) GenerateHealthCheck() (*model.HealthCheckOutput, error) {
	return a.mock.GenerateHealthCheck()
}

func TestHealthCheckHandler(t *testing.T) {
	// Implement your test logic here
	t.Parallel()
	cfg := *config.CreateNewConfig()
	cfgStatus := *config.CreateNewConfigStatus()
	cfgError := *config.CreateNewConfigError()
	cfgPaths := *config.CreateNewConfigPaths()
	cfgHealthCheck := *config.CreateNewConfigHealthCheck()
	instanceID := uuid.New()
	errExpected := errors.New(string(cfgError.ErrorInternalServer))
	testCases := []struct {
		name             string
		setupRequest     func(c *gin.Context)
		setupMockService func(ctx context.Context) *services_mocks.HealthCheck
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "success case",
			setupRequest: func(c *gin.Context) {
				c.Request = httptest.NewRequest(
					http.MethodGet,
					cfgPaths.PathGetHealthCheck,
					nil,
				)
			},
			setupMockService: func(ctx context.Context) *services_mocks.HealthCheck {
				mockService := &services_mocks.HealthCheck{}

				mockService.On("GenerateHealthCheck").
					Return(&model.HealthCheckOutput{
						Message:     cfgStatus.StatusOK,
						ServiceName: cfgHealthCheck.BookmarkService,
						InstanceID:  instanceID,
					}, nil)

				return mockService
			},
			expectedStatus: http.StatusOK,
			expectedResponse: `{"message":"` + cfgStatus.StatusOK +
				`","service_name":"` + cfgHealthCheck.BookmarkService +
				`","instance_id":"` + instanceID.String() + `"}`,
		},
		{
			name: "error case 01",
			setupRequest: func(c *gin.Context) {
				c.Request = httptest.NewRequest(
					http.MethodGet,
					cfgPaths.PathGetHealthCheck,
					nil,
				)
			},
			setupMockService: func(ctx context.Context) *services_mocks.HealthCheck {
				mockService := &services_mocks.HealthCheck{}

				mockService.On("GenerateHealthCheck").
					Return(nil, errExpected)

				return mockService
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: `{"` + cfgError.ErrorTitle + `":"` +
				cfgError.ErrorInternalServer + `"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Implement your test logic here
			rec := httptest.NewRecorder()
			ginCtx, _ := gin.CreateTestContext(rec)
			tc.setupRequest(ginCtx)
			mockService := tc.setupMockService(context.Background())
			handler := NewHealthCheck(healthCheckMockAdapter{mock: mockService, cfg: &cfg}, cfg, cfgError)
			handler.GenerateHealthCheck(ginCtx)
			if rec.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rec.Code)
			}
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, tc.expectedResponse, rec.Body.String())
		})
	}
}
