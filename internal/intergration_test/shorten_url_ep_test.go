package intergrationtest_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/pvlearning002/ebvn_assignment/internal/api"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestCreateShortenLinkEndPoint(t *testing.T) {
	cfg := *config.CreateNewConfig()
	cfgStatus := *config.CreateNewConfigStatus()
	cfgError := *config.CreateNewConfigError()
	cfgPaths := *config.CreateNewConfigPaths()
	cfgHealthCheck := *config.CreateNewConfigHealthCheck()

	miniRedis := miniredis.RunT(t)
	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})
	defer redisClient.Close()

	testCases := []struct {
		name             string
		setupTestHTTP    func(api api.Engine) *httptest.ResponseRecorder
		expectedStatus   int
		expectedResponse string
	}{
		{
			name: "normal case",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				recorder := httptest.NewRecorder()

				req := httptest.NewRequest(
					http.MethodPost,
					cfgPaths.PathPostShortenURL,
					strings.NewReader(`{"url":"https://google.com"}`),
				)

				req.Header.Set("Content-Type", "application/json")

				api.ServeHTTP(recorder, req)
				return recorder
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: `"code"`,
		},
		{
			name: "wrong method case",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				recorder := httptest.NewRecorder()

				req := httptest.NewRequest(
					http.MethodGet,
					cfgPaths.PathPostShortenURL,
					strings.NewReader(`{"url":"https://google.com"}`),
				)

				api.ServeHTTP(recorder, req)
				return recorder
			},
			expectedStatus:   http.StatusNotFound,
			expectedResponse: "404 page not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testAPI := api.NewEngine(
				&cfg,
				&cfgStatus,
				&cfgError,
				&cfgPaths,
				&cfgHealthCheck,
				redisClient,
			)

			recorder := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponse)
		})
	}
}
