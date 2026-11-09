package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/stretchr/testify/assert"
)

type mockShortenURLService struct {
	key       string
	err       error
	calledURL string
}

func (m *mockShortenURLService) CreateShortenLink(
	ctx context.Context,
	url string,
) (string, error) {
	m.calledURL = url
	return m.key, m.err
}

func TestCreateShortenLink(t *testing.T) {
	gin.SetMode(gin.TestMode)

	cfg := config.Config{
		KeyTitle: "key",
	}

	cfgError := config.ConfigError{
		ErrorTitle:          "message",
		ErrorInvalidInput:   "Invalid input",
		ErrorInternalServer: "Internal server error",
	}

	testCases := []struct {
		name           string
		requestBody    string
		mockKey        string
		mockError      error
		expectedStatus int
		expectedBody   string
		expectedURL    string
	}{
		{
			name:           "success",
			requestBody:    `{"url":"https://google.com"}`,
			mockKey:        "abc123",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"key":"abc123"}`,
			expectedURL:    "https://google.com",
		},
		{
			name:           "invalid input",
			requestBody:    `{"url":`,
			mockKey:        "",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"Invalid input"}`,
			expectedURL:    "",
		},
		{
			name:           "internal server error",
			requestBody:    `{"url":"https://google.com"}`,
			mockKey:        "",
			mockError:      errors.New("database connection error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"Internal server error"}`,
			expectedURL:    "https://google.com",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockService := &mockShortenURLService{
				key: tc.mockKey,
				err: tc.mockError,
			}

			handler := NewShortenURLHandler(
				mockService,
				cfg,
				cfgError,
			)

			router := gin.New()
			router.POST("/v1/links/shorten", handler.CreateShortenLink)

			req := httptest.NewRequest(
				http.MethodPost,
				"/v1/links/shorten",
				strings.NewReader(tc.requestBody),
			)
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatus, rec.Code)
			assert.JSONEq(t, tc.expectedBody, rec.Body.String())
			assert.Equal(t, tc.expectedURL, mockService.calledURL)
		})
	}
}
