package intergrationtest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	"github.com/pvlearning002/ebvn_assignment/internal/api"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
)

type shortenLinkResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Title string `json:"title"`
}

// setupTestServer initializes the API server with a real Redis repository.
func setupTestServer(t *testing.T) (*miniredis.Miniredis, *redis.Client, http.Handler) {
	t.Helper()

	miniRedis := miniredis.RunT(t)

	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	t.Cleanup(func() {
		_ = redisClient.Close()
	})

	cfg := config.CreateNewConfig()
	cfgStatus := config.CreateNewConfigStatus()
	cfgError := config.CreateNewConfigError()
	cfgPaths := config.CreateNewConfigPaths()
	cfgHealthCheck := config.CreateNewConfigHealthCheck()
	cfgParameter := config.CreateNewConfigParameter()

	router := api.NewEngine(
		cfg,
		cfgStatus,
		cfgError,
		cfgPaths,
		cfgHealthCheck,
		cfgParameter,
		redisClient,
	)

	return miniRedis, redisClient, router
}

// createShortenLinkRequest creates a valid HTTP request.
func createShortenLinkRequest() *http.Request {
	body := `{
		"expiration_time": 3600,
		"url": "https://www.google.com"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/links/shorten",
		bytes.NewBufferString(body),
	)

	req.Header.Set("Content-Type", "application/json")

	return req
}

// parseShortenLinkResponse parses the API response.
func parseShortenLinkResponse(t *testing.T, body []byte) shortenLinkResponse {
	t.Helper()

	var response shortenLinkResponse

	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to parse response: %v. Body: %s", err, string(body))
	}

	return response
}

func TestCreateShortenLink_Integration(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectCode     bool
	}{
		{
			name: "Success",
			requestBody: `{
				"expiration_time": 3600,
				"url": "https://www.google.com"
			}`,
			expectedStatus: http.StatusOK,
			expectCode:     true,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"url":}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Missing expiration time",
			requestBody: `{
				"url": "https://www.google.com"
			}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid expiration time",
			requestBody: `{
				"expiration_time": 100,
				"url": "https://www.google.com"
			}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid URL",
			requestBody: `{
				"expiration_time": 3600,
				"url": "invalid-url"
			}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, router := setupTestServer(t)

			req := httptest.NewRequest(
				http.MethodPost,
				"/v1/links/shorten",
				bytes.NewBufferString(tt.requestBody),
			)

			req.Header.Set("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Fatalf(
					"Expected status %d, got %d. Response: %s",
					tt.expectedStatus,
					recorder.Code,
					recorder.Body.String(),
				)
			}

			if tt.expectCode {
				response := parseShortenLinkResponse(t, recorder.Body.Bytes())

				if response.Code == "" {
					t.Error("Expected non-empty shorten URL code")
				}

				if response.Message == "" {
					t.Error("Expected non-empty response message")
				}
			}
		})
	}
}

func TestCreateShortenLink_RedisConnectionError_Integration(t *testing.T) {
	miniRedis, _, router := setupTestServer(t)

	// Stop Redis to simulate a connection failure.
	miniRedis.Close()

	req := createShortenLinkRequest()

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf(
			"Expected status %d, got %d. Response: %s",
			http.StatusInternalServerError,
			recorder.Code,
			recorder.Body.String(),
		)
	}
}

func TestRedirectToOriginalURL_Integration(t *testing.T) {
	_, _, router := setupTestServer(t)

	// First create a shortened URL.
	req := createShortenLinkRequest()

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"Failed to create shortened URL. Status: %d, Response: %s",
			recorder.Code,
			recorder.Body.String(),
		)
	}

	response := parseShortenLinkResponse(t, recorder.Body.Bytes())

	if response.Code == "" {
		t.Fatal("Shorten URL code is empty")
	}

	// Redirect using the generated code.
	redirectReq := httptest.NewRequest(
		http.MethodGet,
		"/v1/links/redirect/"+response.Code,
		nil,
	)

	redirectRecorder := httptest.NewRecorder()

	router.ServeHTTP(redirectRecorder, redirectReq)

	if redirectRecorder.Code != http.StatusMovedPermanently {
		t.Fatalf(
			"Expected status %d, got %d. Response: %s",
			http.StatusMovedPermanently,
			redirectRecorder.Code,
			redirectRecorder.Body.String(),
		)
	}

	location := redirectRecorder.Header().Get("Location")

	if location != "https://www.google.com" {
		t.Errorf(
			"Expected redirect URL %s, got %s",
			"https://www.google.com",
			location,
		)
	}
}

func TestRedirectToOriginalURL_NotFound_Integration(t *testing.T) {
	_, _, router := setupTestServer(t)

	req := httptest.NewRequest(
		http.MethodGet,
		"/v1/links/redirect/nonexistent123",
		nil,
	)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf(
			"Expected status %d, got %d. Response: %s",
			http.StatusNotFound,
			recorder.Code,
			recorder.Body.String(),
		)
	}
}
