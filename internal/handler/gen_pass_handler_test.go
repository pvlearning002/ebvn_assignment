package handler

import (
	"context"
	http "net/http"
	"net/http/httptest"
	"testing"

	gin "github.com/gin-gonic/gin"
	config "github.com/pvlearning002/ebvn_assignment/internal/config"
	mocks "github.com/pvlearning002/ebvn_assignment/internal/service/mocks"
	assert "github.com/stretchr/testify/assert"
)

func TestGenPassHandler_GeneratePassword(t *testing.T) {
	// Test implementation
	t.Parallel()
	cfg := *config.CreateNewConfig()
	cfgPaths := *config.CreateNewConfigPaths()
	cfgStatus := *config.CreateNewConfigStatus()
	cfgError := *config.CreateNewConfigError()
	pwd_test := "123456789012"
	testCases := []struct {
		name             string
		setupRequest     func(ctx *gin.Context)
		setupMockService func(ctx context.Context) *mocks.GenPass
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "service success",
			setupRequest: func(ctx *gin.Context) {
				// Setup request context if needed
				ctx.Request = httptest.NewRequest(http.MethodGet, cfgPaths.PathGetGenPass, nil)
			},
			setupMockService: func(ctx context.Context) *mocks.GenPass {
				mockSvc := mocks.NewGenPass(t)
				mockSvc.On("GeneratePassword").Return(pwd_test, nil).Once()
				return mockSvc
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"` + cfg.PasswordTitle + `":"` + pwd_test + `"}`,
		},
		{
			name: "service failed",
			setupRequest: func(ctx *gin.Context) {
				// Setup request context if needed
				ctx.Request = httptest.NewRequest(http.MethodGet, cfgPaths.PathGetGenPass, nil)
			},
			setupMockService: func(ctx context.Context) *mocks.GenPass {
				mockSvc := mocks.NewGenPass(t)
				mockSvc.On("GeneratePassword").Return(cfg.StringEmpty, assert.AnError).Once()
				return mockSvc
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: `{"` + cfgStatus.StatusOK + `":"` + cfgError.ErrorInternalServer + `"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			tc.setupRequest(ctx)
			mockSvc := tc.setupMockService(ctx)
			handler := NewGenPass(mockSvc, cfg, cfgStatus, cfgError)
			handler.GeneratePassword(ctx)
			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.Equal(t, tc.expectedResponse, recorder.Body.String())
		})
	}
}
