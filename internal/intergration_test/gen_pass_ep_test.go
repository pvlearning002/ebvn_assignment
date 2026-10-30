package intergrationtest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pvlearning002/ebvn_assignment/internal/api"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestGenPassEndPoint(t *testing.T) {
	t.Parallel()
	cfg := *config.CreateNewConfig()
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
				req := httptest.NewRequest(http.MethodGet, cfg.PathGetGenPass, nil)
				api.ServeHTTP(recorder, req)
				return recorder
			},
			expectedStatus:   http.StatusOK,
			expectedResponse: cfg.PasswordTitle,
		}, {
			name: "wrong method case",
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				recorder := httptest.NewRecorder()
				req := httptest.NewRequest(http.MethodPost, cfg.PathGetGenPass, nil)
				api.ServeHTTP(recorder, req)
				return recorder
			},
			expectedStatus:   http.StatusNotFound,
			expectedResponse: "404 page not found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			testAPI := api.NewEngine(&cfg)
			recorder := tc.setupTestHTTP(testAPI)
			assert.Equal(t, tc.expectedStatus, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponse)
		})
	}
}
