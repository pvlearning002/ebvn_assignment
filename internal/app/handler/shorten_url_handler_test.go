package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pvlearning002/ebvn_assignment/internal/app/model"
	"github.com/pvlearning002/ebvn_assignment/internal/app/service"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type shortenHandlerServiceMock struct{ mock.Mock }

func (m *shortenHandlerServiceMock) CreateShortenLink(ctx context.Context, url string) (*model.ShortenLinkOutput, error) {
	args := m.Called(ctx, url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	result, ok := args.Get(0).(*model.ShortenLinkOutput)
	if !ok {
		return nil, args.Error(1)
	}
	return result, args.Error(1)
}
func (m *shortenHandlerServiceMock) GetUrlByCode(ctx context.Context, code string) (string, error) {
	args := m.Called(ctx, code)
	return args.String(0), args.Error(1)
}

type shortenHandlerLoggerMock struct{ mock.Mock }

func (m *shortenHandlerLoggerMock) SetLogLevel(level string) { m.Called(level) }
func (m *shortenHandlerLoggerMock) GetPathFunction(skip ...int) string {
	return m.Called(skip).String(0)
}
func (m *shortenHandlerLoggerMock) BuildLogMessage(message string, skip ...int) string {
	return m.Called(message, skip).String(0)
}
func (m *shortenHandlerLoggerMock) WriteLogInfo(message string, skip ...int) {
	m.Called(message, skip)
}
func (m *shortenHandlerLoggerMock) WriteLogError(message string, skip ...int) {
	m.Called(message, skip)
}
func (m *shortenHandlerLoggerMock) WriteLogDebug(message string, skip ...int) {
	m.Called(message, skip)
}
func (m *shortenHandlerLoggerMock) WriteLogWarning(message string, skip ...int) {
	m.Called(message, skip)
}

var _ service.ShortenURL = (*shortenHandlerServiceMock)(nil)
var _ utils.LoggerUtils = (*shortenHandlerLoggerMock)(nil)

func newShortenHandlerForTest(svc *shortenHandlerServiceMock, logger *shortenHandlerLoggerMock) *shortenURLHandler {
	return NewShortenURLHandler(
		svc,
		config.Config{StringEmpty: ""},
		config.ConfigError{
			ErrorTitle:          "Error",
			ErrorInvalidInput:   "Invalid input",
			ErrorInternalServer: "Internal Server Error",
			ErrorInvalidCode:    "Invalid code",
			ErrorUrlNotFound:    "URL not found",
		},
		utils.NewResponseUtils(config.ConfigError{
			ErrorTitle:          "Error",
			ErrorInvalidInput:   "Invalid input",
			ErrorInternalServer: "Internal Server Error",
			ErrorUrlNotFound:    "URL not found",
		}),
		config.ConfigParameter{Code: "code"},
		logger,
	)
}

func TestNewShortenURLHandler(t *testing.T) {
	got := newShortenHandlerForTest(new(shortenHandlerServiceMock), new(shortenHandlerLoggerMock))
	require.NotNil(t, got)
	assert.IsType(t, &shortenURLHandler{}, got)
}

func TestShortenURLHandler_CreateShortenLink(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name           string
		body           string
		serviceOutput  *model.ShortenLinkOutput
		serviceError   error
		wantStatus     int
		wantBody       string
		wantLogMessage string
	}{
		{
			name:          "success",
			body:          `{"expiration_time":500,"url":"https://example.com"}`,
			serviceOutput: &model.ShortenLinkOutput{KeyCode: "abc123", Message: "created"},
			wantStatus:    http.StatusOK,
			wantBody:      `{"code":"abc123","message":"created"}`,
		},
		{
			name:           "invalid JSON",
			body:           `{"url":`,
			wantStatus:     http.StatusBadRequest,
			wantBody:       `{"error_title":"Error","message":"Invalid input"}`,
			wantLogMessage: "invalid request",
		},
		{
			name:           "service error",
			body:           `{"expiration_time":500,"url":"https://example.com"}`,
			serviceError:   errors.New("store failed"),
			wantStatus:     http.StatusInternalServerError,
			wantBody:       `{"error_title":"Error","message":"Internal Server Error"}`,
			wantLogMessage: "store failed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			serviceMock := new(shortenHandlerServiceMock)
			logger := new(shortenHandlerLoggerMock)
			if strings.Contains(tc.body, "https://example.com") {
				serviceMock.On("CreateShortenLink", mock.Anything, "https://example.com").Return(tc.serviceOutput, tc.serviceError).Once()
			}
			if tc.wantLogMessage == "store failed" {
				logger.On("WriteLogError", tc.wantLogMessage, mock.Anything).Return().Once()
			} else if tc.name == "invalid JSON" {
				// Gin's JSON bind error text varies by version; assert the logger call without matching its message.
				logger.On("WriteLogError", mock.Anything, mock.Anything).Return().Once()
			}

			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", strings.NewReader(tc.body))
			ctx.Request.Header.Set("Content-Type", "application/json")
			newShortenHandlerForTest(serviceMock, logger).CreateShortenLink(ctx)

			assert.Equal(t, tc.wantStatus, recorder.Code)
			assert.JSONEq(t, tc.wantBody, recorder.Body.String())
			serviceMock.AssertExpectations(t)
			logger.AssertExpectations(t)
		})
	}
}

func TestShortenURLHandler_RedirectToOriginalURL(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name           string
		code           string
		url            string
		serviceError   error
		wantStatus     int
		wantLocation   string
		wantBody       string
		wantLogMessage string
	}{
		{name: "redirect", code: "abc123", url: "https://example.com", wantStatus: http.StatusMovedPermanently, wantLocation: "https://example.com"},
		{name: "missing code", code: "", wantStatus: http.StatusBadRequest, wantBody: `{"error_title":"Error","message":"Invalid input"}`, wantLogMessage: "Invalid code"},
		{name: "URL not found", code: "missing", url: "", wantStatus: http.StatusNotFound, wantBody: `{"error_title":"Error","message":"URL not found"}`},
		{name: "service error", code: "abc123", serviceError: errors.New("lookup failed"), wantStatus: http.StatusInternalServerError, wantBody: `{"error_title":"Error","message":"Internal Server Error"}`, wantLogMessage: "lookup failed"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			serviceMock := new(shortenHandlerServiceMock)
			logger := new(shortenHandlerLoggerMock)
			if tc.code != "" {
				serviceMock.On("GetUrlByCode", mock.Anything, tc.code).Return(tc.url, tc.serviceError).Once()
			}
			if tc.wantLogMessage != "" {
				logger.On("WriteLogError", tc.wantLogMessage, mock.Anything).Return().Once()
			}

			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/"+tc.code, nil)
			ctx.Params = gin.Params{{Key: "code", Value: tc.code}}
			newShortenHandlerForTest(serviceMock, logger).RedirectToOriginalURL(ctx)

			assert.Equal(t, tc.wantStatus, recorder.Code)
			if tc.wantLocation != "" {
				assert.Equal(t, tc.wantLocation, recorder.Header().Get("Location"))
			} else {
				assert.JSONEq(t, tc.wantBody, recorder.Body.String())
			}
			serviceMock.AssertExpectations(t)
			logger.AssertExpectations(t)
		})
	}
}
