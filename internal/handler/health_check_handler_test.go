package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/internal/model"
	services_mocks "github.com/pvlearning002/ebvn_assignment/internal/service/services_mocks"
	"github.com/stretchr/testify/assert"
)

type healthCheckMockAdapter struct {
	mock *services_mocks.HealthCheck
	cfg  *config.Config
}

func (a healthCheckMockAdapter) GenerateHealthCheck() (*model.HealthCheck, error) {
	return a.mock.GenerateHealthCheck()
}

func TestHealthCheckHandler(t *testing.T) {
	// Implement your test logic here
	t.Parallel()
	cfg := *config.CreateNewConfig()
	instanceID := uuid.New()
	testCases := []struct {
		name             string
		setupRequest     func(c *gin.Context)
		setupMockService func(ctx context.Context) *services_mocks.HealthCheck
		expectedStatus   int
		expectedResponse bool
	}{
		{
			name: "normal case",
			setupRequest: func(c *gin.Context) {
				// Setup your request here
				c.Request = httptest.NewRequest(http.MethodGet, cfg.PathGetHealthCheck, nil)
			},
			setupMockService: func(ctx context.Context) *services_mocks.HealthCheck {
				mockService := &services_mocks.HealthCheck{}
				// Setup your mock service here
				mockService.On("GenerateHealthCheck").Return(&model.HealthCheck{Message: cfg.StatusOK, ServiceName: cfg.BookmarkService, InstanceID: instanceID}, nil)
				return mockService
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: true,
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
			handler := NewHealthCheck(healthCheckMockAdapter{mock: mockService, cfg: &cfg}, cfg)
			handler.GenerateHealthCheck(ginCtx)
			if rec.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rec.Code)
			}
			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String() == `{"message":"`+cfg.StatusOK+`","service_name":"`+cfg.BookmarkService+`","instance_id":"`+instanceID.String()+`"}`)
		})
	}
}
