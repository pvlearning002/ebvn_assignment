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
	"github.com/pvlearning002/ebvn_assignment/internal/app/model"
	"github.com/pvlearning002/ebvn_assignment/internal/app/service"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockShortenURL implements service.ShortenURL.
type MockShortenURL struct {
	mock.Mock
}

// Compile-time interface check.
var _ service.ShortenURL = (*MockShortenURL)(nil)

// CreateShortenLink mocks the CreateShortenLink service method.
func (m *MockShortenURL) CreateShortenLink(
	ctx context.Context,
	url string,
) (*model.ShortenLinkOutput, error) {

	args := m.Called(ctx, url)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	output, ok := args.Get(0).(*model.ShortenLinkOutput)
	if !ok {
		return nil, args.Error(1)
	}

	return output, args.Error(1)
}

// GetUrlByCode mocks the GetUrlByCode service method.
func (m *MockShortenURL) GetUrlByCode(
	ctx context.Context,
	keyCode string,
) (string, error) {

	args := m.Called(ctx, keyCode)
	return args.String(0), args.Error(1)
}

// createTestRedis creates a Redis client using miniredis.
func createTestRedis(
	t *testing.T,
) (*miniredis.Miniredis, *redis.Client) {

	t.Helper()

	miniRedis, err := miniredis.Run()
	require.NoError(t, err)

	client := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	t.Cleanup(func() {
		_ = client.Close()
		miniRedis.Close()
	})

	return miniRedis, client
}

// createTestHandler creates a handler for testing.
func createTestHandler(
	svc service.ShortenURL,
	client *redis.Client,
) *shortenURLHandler {

	cfg := *config.CreateNewConfig()
	cfgError := *config.CreateNewConfigError()
	cfgParameter := *config.CreateNewConfigParameter()

	return NewShortenURLHandler(
		svc,
		cfg,
		cfgError,
		cfgParameter,
		client,
	)
}

// createShortenLinkRequest creates an HTTP request for shortening a URL.
func createShortenLinkRequest(body string) *http.Request {

	req := httptest.NewRequest(
		http.MethodPost,
		"/v1/links/shorten",
		bytes.NewBufferString(body),
	)

	req.Header.Set("Content-Type", "application/json")

	return req
}

// createTestContext creates a Gin test context.
func createTestContext(
	req *http.Request,
) (*gin.Context, *httptest.ResponseRecorder) {

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = req

	return ctx, recorder
}

// TestCreateShortenLink tests CreateShortenLink.
func TestCreateShortenLink(t *testing.T) {

	gin.SetMode(gin.TestMode)

	t.Run("redis_connection_error", func(t *testing.T) {

		miniRedis, client := createTestRedis(t)

		// Close Redis to simulate connection failure.
		miniRedis.Close()

		mockService := new(MockShortenURL)
		handler := createTestHandler(mockService, client)

		req := createShortenLinkRequest(
			`{"url":"https://google.com","expiration_time":3600}`,
		)

		ctx, recorder := createTestContext(req)

		handler.CreateShortenLink(ctx)

		assert.Equal(
			t,
			http.StatusInternalServerError,
			recorder.Code,
		)

		assert.Empty(t, mockService.Calls)
	})

	t.Run("invalid_input", func(t *testing.T) {

		_, client := createTestRedis(t)

		mockService := new(MockShortenURL)
		handler := createTestHandler(mockService, client)

		// Invalid JSON.
		req := createShortenLinkRequest(
			`{"url":`,
		)

		ctx, recorder := createTestContext(req)

		handler.CreateShortenLink(ctx)

		assert.Equal(
			t,
			http.StatusBadRequest,
			recorder.Code,
		)

		assert.Empty(t, mockService.Calls)
	})

	t.Run("missing_expiration_time", func(t *testing.T) {

		_, client := createTestRedis(t)

		mockService := new(MockShortenURL)
		handler := createTestHandler(mockService, client)

		// Missing required expiration_time.
		req := createShortenLinkRequest(
			`{"url":"https://google.com"}`,
		)

		ctx, recorder := createTestContext(req)

		handler.CreateShortenLink(ctx)

		assert.Equal(
			t,
			http.StatusBadRequest,
			recorder.Code,
		)

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

		req := createShortenLinkRequest(
			`{"url":"https://google.com","expiration_time":3600}`,
		)

		ctx, recorder := createTestContext(req)

		handler.CreateShortenLink(ctx)

		assert.Equal(
			t,
			http.StatusInternalServerError,
			recorder.Code,
		)

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

		req := createShortenLinkRequest(
			`{"url":"https://google.com","expiration_time":3600}`,
		)

		ctx, recorder := createTestContext(req)

		handler.CreateShortenLink(ctx)

		assert.Equal(
			t,
			http.StatusOK,
			recorder.Code,
		)

		assert.Equal(
			t,
			"application/json; charset=utf-8",
			recorder.Header().Get("Content-Type"),
		)

		mockService.AssertExpectations(t)
	})
}

// TestRedirectToOriginalURL tests RedirectToOriginalURL.
func TestRedirectToOriginalURL(t *testing.T) {

	gin.SetMode(gin.TestMode)

	testCases := []struct {
		name           string
		keyCode        string
		mockURL        string
		mockError      error
		expectedStatus int
		expectedURL    string
		shouldCallMock bool
	}{
		{
			name:           "empty_key_code",
			keyCode:        "",
			expectedStatus: http.StatusBadRequest,
			shouldCallMock: false,
		},
		{
			name:           "service_error",
			keyCode:        "abc123",
			mockError:      errors.New("redis connection error"),
			expectedStatus: http.StatusInternalServerError,
			shouldCallMock: true,
		},
		{
			name:           "url_not_found",
			keyCode:        "abc123",
			mockURL:        "",
			expectedStatus: http.StatusNotFound,
			shouldCallMock: true,
		},
		{
			name:           "redirect_success",
			keyCode:        "abc123",
			mockURL:        "https://www.google.com",
			expectedStatus: http.StatusMovedPermanently,
			expectedURL:    "https://www.google.com",
			shouldCallMock: true,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {

			mockService := new(MockShortenURL)

			if tc.shouldCallMock {
				mockService.On(
					"GetUrlByCode",
					mock.Anything,
					tc.keyCode,
				).Return(
					tc.mockURL,
					tc.mockError,
				).Once()
			}

			cfg := *config.CreateNewConfig()
			cfgError := *config.CreateNewConfigError()
			cfgParameter := *config.CreateNewConfigParameter()

			handler := &shortenURLHandler{
				shortenURLService: mockService,
				cfg:               cfg,
				cfgError:          cfgError,
				cfgParameter:      cfgParameter,
			}

			router := gin.New()

			if tc.keyCode == "" {
				router.GET("/", handler.RedirectToOriginalURL)
			} else {
				router.GET("/:code", handler.RedirectToOriginalURL)
			}

			request := httptest.NewRequest(
				http.MethodGet,
				"/"+tc.keyCode,
				nil,
			)

			recorder := httptest.NewRecorder()

			router.ServeHTTP(recorder, request)

			assert.Equal(
				t,
				tc.expectedStatus,
				recorder.Code,
			)

			if tc.expectedURL != "" {
				assert.Equal(
					t,
					tc.expectedURL,
					recorder.Header().Get("Location"),
				)
			}

			mockService.AssertExpectations(t)
		})
	}
}
