package intergrationtest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type shortenLinkResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func sendShortenRequest(engine interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}, body string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	engine.ServeHTTP(recorder, request)
	return recorder
}

func parseShortenResponse(t *testing.T, body []byte) shortenLinkResponse {
	t.Helper()
	var result shortenLinkResponse
	require.NoError(t, json.Unmarshal(body, &result), "response body: %s", string(body))
	return result
}

func TestCreateShortenLink_Integration(t *testing.T) {
	tests := []struct {
		name             string
		body             string
		wantStatus       int
		wantNonEmptyCode bool
	}{
		{
			name:             "valid request",
			body:             `{"expiration_time":3600,"url":"https://www.google.com"}`,
			wantStatus:       http.StatusOK,
			wantNonEmptyCode: true,
		},
		{name: "invalid JSON", body: `{"url":}`, wantStatus: http.StatusBadRequest},
		{name: "missing expiration", body: `{"url":"https://www.google.com"}`, wantStatus: http.StatusBadRequest},
		{name: "expiration below minimum", body: `{"expiration_time":100,"url":"https://www.google.com"}`, wantStatus: http.StatusBadRequest},
		{name: "invalid URL", body: `{"expiration_time":3600,"url":"not-a-url"}`, wantStatus: http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, engine, _, _, _ := setupIntegrationAPI(t)
			recorder := sendShortenRequest(engine, tc.body)
			assert.Equal(t, tc.wantStatus, recorder.Code, recorder.Body.String())
			if tc.wantNonEmptyCode {
				result := parseShortenResponse(t, recorder.Body.Bytes())
				assert.NotEmpty(t, result.Code)
				assert.NotEmpty(t, result.Message)
			}
		})
	}
}

func TestCreateShortenLink_RedisConnectionError_Integration(t *testing.T) {
	redisServer, engine, _, _, _ := setupIntegrationAPI(t)
	redisServer.Close()

	recorder := sendShortenRequest(engine, `{"expiration_time":3600,"url":"https://www.google.com"}`)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code, recorder.Body.String())
}

func TestRedirectToOriginalURL_Integration(t *testing.T) {
	_, engine, _, _, _ := setupIntegrationAPI(t)
	recorder := sendShortenRequest(engine, `{"expiration_time":3600,"url":"https://www.google.com"}`)
	require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())
	result := parseShortenResponse(t, recorder.Body.Bytes())
	require.NotEmpty(t, result.Code)

	redirectRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/links/redirect/"+result.Code, nil)
	engine.ServeHTTP(redirectRecorder, request)

	assert.Equal(t, http.StatusMovedPermanently, redirectRecorder.Code)
	assert.Equal(t, "https://www.google.com", redirectRecorder.Header().Get("Location"))
}

func TestRedirectToOriginalURL_NotFound_Integration(t *testing.T) {
	_, engine, _, _, _ := setupIntegrationAPI(t)
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/v1/links/redirect/nonexistent123", nil)
	engine.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusNotFound, recorder.Code, recorder.Body.String())
}
