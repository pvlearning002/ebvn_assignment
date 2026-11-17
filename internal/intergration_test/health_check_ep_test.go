package intergrationtest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/pvlearning002/ebvn_assignment/internal/api"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
)

// TestHealthCheckEndPoint tests the health check endpoint.
func TestHealthCheckEndPoint(t *testing.T) {
	// Initialize configuration.
	cfg := *config.CreateNewConfig()
	cfgStatus := *config.CreateNewConfigStatus()
	cfgError := *config.CreateNewConfigError()
	cfgPaths := *config.CreateNewConfigPaths()
	cfgHealthCheck := *config.CreateNewConfigHealthCheck()
	cfgParameter := *config.CreateNewConfigParameter()

	// Initialize a valid UUID for HealthCheckService.
	cfg.InstanceID = uuid.NewString()

	// Start an in-memory Redis server.
	miniRedis, err := miniredis.Run()
	require.NoError(t, err)
	defer miniRedis.Close()

	// Create Redis client.
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	defer redisClient.Close()

	// Initialize Gin engine.
	testAPI := api.NewEngine(
		&cfg,
		&cfgStatus,
		&cfgError,
		&cfgPaths,
		&cfgHealthCheck,
		&cfgParameter,
		redisClient,
	)

	testCases := []struct {
		name             string
		method           string
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:             "normal case",
			method:           http.MethodGet,
			expectedStatus:   http.StatusOK,
			expectedResponse: cfgHealthCheck.BookmarkService,
		},
		{
			name:             "wrong method case",
			method:           http.MethodPost,
			expectedStatus:   http.StatusNotFound,
			expectedResponse: "404 page not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()

			req := httptest.NewRequest(
				tc.method,
				cfgPaths.PathGetHealthCheck,
				nil,
			)

			testAPI.ServeHTTP(recorder, req)

			// Verify HTTP status code.
			assert.Equal(
				t,
				tc.expectedStatus,
				recorder.Code,
			)

			// Verify response body.
			assert.Contains(
				t,
				recorder.Body.String(),
				tc.expectedResponse,
			)
		})
	}
}
