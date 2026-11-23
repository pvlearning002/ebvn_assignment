package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pvlearning002/ebvn_assignment/internal/app/service"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type genPassHandlerServiceMock struct{ mock.Mock }

func (m *genPassHandlerServiceMock) GeneratePassword() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

var _ service.GenPass = (*genPassHandlerServiceMock)(nil)

func TestNewGenPass(t *testing.T) {
	got := NewGenPass(
		new(genPassHandlerServiceMock),
		config.Config{},
		utils.NewLoggerUtils(),
		utils.NewResponseUtils(config.ConfigError{}),
	)
	require.NotNil(t, got)
	assert.IsType(t, &genPassHandler{}, got)
}

func TestGenPassHandler_GeneratePassword_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	serviceMock := new(genPassHandlerServiceMock)
	serviceMock.On("GeneratePassword").Return("pass-123", nil).Once()
	handler := NewGenPass(
		serviceMock,
		config.Config{PasswordTitle: "password"},
		utils.NewLoggerUtils(),
		utils.NewResponseUtils(config.ConfigError{}),
	)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/genpass", nil)

	handler.GeneratePassword(ctx)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.JSONEq(t, `{"password":"pass-123"}`, recorder.Body.String())
	serviceMock.AssertExpectations(t)
}

func TestGenPassHandler_GeneratePassword_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	serviceMock := new(genPassHandlerServiceMock)
	wantErr := errors.New("password generation failed")
	serviceMock.On("GeneratePassword").Return("", wantErr).Once()
	handler := NewGenPass(
		serviceMock,
		config.Config{},
		utils.NewLoggerUtils(),
		utils.NewResponseUtils(config.ConfigError{
			ErrorTitle:          "Error",
			ErrorInternalServer: "internal server error",
		}),
	)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodGet, "/genpass", nil)

	handler.GeneratePassword(ctx)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.JSONEq(t, `{"error_title":"Error","message":"internal server error"}`, recorder.Body.String())
	serviceMock.AssertExpectations(t)
}
