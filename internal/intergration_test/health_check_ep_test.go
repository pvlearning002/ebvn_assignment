package intergrationtest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheckEndPoint(t *testing.T) {
	_, engine, _, paths, healthCfg := setupIntegrationAPI(t)
	tests := []struct {
		name         string
		method       string
		wantStatus   int
		wantBodyPart string
	}{
		{name: "health endpoint", method: http.MethodGet, wantStatus: http.StatusOK, wantBodyPart: healthCfg.BookmarkService},
		{name: "wrong method", method: http.MethodPost, wantStatus: http.StatusNotFound, wantBodyPart: "404 page not found"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(tc.method, paths.PathGetHealthCheck, nil)
			engine.ServeHTTP(recorder, request)

			assert.Equal(t, tc.wantStatus, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.wantBodyPart)
		})
	}
}
