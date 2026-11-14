package handler

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/internal/model"
	"github.com/pvlearning002/ebvn_assignment/internal/service"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockShortenURL implements service.ShortenURL.
type MockShortenURL struct {
	mock.Mock
}

func (m *MockShortenURL) CreateShortenLink(
	ctx context.Context,
	url string,
) (*model.ShortenLinkOutput, error) {

	args := m.Called(ctx, url)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.ShortenLinkOutput), args.Error(1)
}

// Compile-time interface check.
var _ service.ShortenURL = (*MockShortenURL)(nil)

// createTestRedis creates a Redis client using miniredis.
func createTestRedis(t *testing.T) (*miniredis.Miniredis, *redis.Client) {
	t.Helper()

	miniRedis, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	t.Cleanup(func() {
		client.Close()
		miniRedis.Close()
	})

	return miniRedis, client
}

// createTestHandler creates a handler for testing.
func createTestHandler(
	svc service.ShortenURL,
	client *redis.Client,
) *shortenURLHandler {

	cfg := config.Config{}

	cfgError := config.ConfigError{
		ErrorTitle:           "error",
		ErrorRedisConnection: "Redis connection error",
		ErrorInvalidInput:    "Invalid input",
		ErrorInternalServer:  "Internal server error",
	}

	return NewShortenURLHandler(svc, cfg, cfgError, client)
}

// TestCreateShortenLink tests CreateShortenLink.
func TestCreateShortenLink(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("redis_connection_error", func(t *testing.T) {

		miniRedis, client := createTestRedis(t)

		// Stop Redis to simulate connection failure.
		miniRedis.Close()

		mockService := new(MockShortenURL)
		handler := createTestHandler(mockService, client)

		req := httptest.NewRequest(
			http.MethodPost,
			"/v1/links/shorten",
			bytes.NewBufferString(`{"url":"https://google.com"}`),
		)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = req

		handler.CreateShortenLink(ctx)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		assert.Empty(t, mockService.Calls)
	})

	t.Run("invalid_input", func(t *testing.T) {

		_, client := createTestRedis(t)

		mockService := new(MockShortenURL)
		handler := createTestHandler(mockService, client)

		req := httptest.NewRequest(
			http.MethodPost,
			"/v1/links/shorten",
			bytes.NewBufferString(`{"url":`),
		)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = req

		handler.CreateShortenLink(ctx)

		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		assert.Empty(t, mockService.Calls)
	})

	t.Run("service_error", func(t *testing.T) {

		_, client := createTestRedis(t)

		mockService := new(MockShortenURL)

		mockService.On(
			"CreateShortenLink",
			mock.Anything,
			"https://google.com",
		).Return(
			nil,
			errors.New("service error"),
		).Once()

		handler := createTestHandler(mockService, client)

		req := httptest.NewRequest(
			http.MethodPost,
			"/v1/links/shorten",
			bytes.NewBufferString(`{"url":"https://google.com"}`),
		)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = req

		handler.CreateShortenLink(ctx)

		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {

		_, client := createTestRedis(t)

		mockService := new(MockShortenURL)

		expectedOutput := &model.ShortenLinkOutput{}

		mockService.On(
			"CreateShortenLink",
			mock.Anything,
			"https://google.com",
		).Return(
			expectedOutput,
			nil,
		).Once()

		handler := createTestHandler(mockService, client)

		req := httptest.NewRequest(
			http.MethodPost,
			"/v1/links/shorten",
			bytes.NewBufferString(`{"url":"https://google.com"}`),
		)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(recorder)
		ctx.Request = req

		handler.CreateShortenLink(ctx)

		assert.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "application/json; charset=utf-8",
			recorder.Header().Get("Content-Type"))

		mockService.AssertExpectations(t)
	})
}
