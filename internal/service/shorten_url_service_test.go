package service

import (
	"context"
	"errors"
	"testing"

	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock GenPass
type MockGenPass struct {
	mock.Mock
}

func (m *MockGenPass) GenerateRandomString(length int) (string, error) {
	args := m.Called(length)
	return args.String(0), args.Error(1)
}

func (m *MockGenPass) GeneratePassword() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// Mock URLStorage
type MockURLStorage struct {
	mock.Mock
}

func (m *MockURLStorage) GetURL(
	ctx context.Context,
	key string,
) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *MockURLStorage) StoreURL(
	ctx context.Context,
	key string,
	url string,
) error {
	args := m.Called(ctx, key, url)
	return args.Error(0)
}

// Compile-time interface checks
var _ GenPass = (*MockGenPass)(nil)
var _ repository.URLStorage = (*MockURLStorage)(nil)

// Create test service
func createTestShortenURLService(
	storage *MockURLStorage,
	genPass *MockGenPass,
) ShortenURL {
	cfg := config.Config{
		KeycodeLengthDefault: 6,
		StringEmpty:          "",
	}

	return NewShortenURLService(
		storage,
		genPass,
		cfg,
		config.ConfigError{},
	)
}

// Test CreateShortenLink
func TestCreateShortenLink(t *testing.T) {
	ctx := context.Background()
	originalURL := "https://www.google.com"

	t.Run("success", func(t *testing.T) {
		mockStorage := new(MockURLStorage)
		mockGenPass := new(MockGenPass)

		service := createTestShortenURLService(
			mockStorage,
			mockGenPass,
		)

		mockGenPass.On(
			"GenerateRandomString",
			6,
		).Return("abc123", nil).Once()

		mockStorage.On(
			"GetURL",
			mock.Anything,
			"abc123",
		).Return("", nil).Once()

		mockStorage.On(
			"StoreURL",
			mock.Anything,
			"abc123",
			originalURL,
		).Return(nil).Once()

		result, err := service.CreateShortenLink(
			ctx,
			originalURL,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "abc123", result.KeyCode)
		assert.Equal(
			t,
			"Shorten URL generated successfully",
			result.Message,
		)

		mockGenPass.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})

	t.Run("generate password error", func(t *testing.T) {
		mockStorage := new(MockURLStorage)
		mockGenPass := new(MockGenPass)

		service := createTestShortenURLService(
			mockStorage,
			mockGenPass,
		)

		expectedError := errors.New("generate random string failed")

		mockGenPass.On(
			"GenerateRandomString",
			6,
		).Return("", expectedError).Once()

		result, err := service.CreateShortenLink(
			ctx,
			originalURL,
		)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.NotNil(t, result)
		assert.Equal(t, "", result.KeyCode)
		assert.Equal(t, expectedError.Error(), result.Message)

		mockGenPass.AssertExpectations(t)
		mockStorage.AssertNotCalled(
			t,
			"StoreURL",
			mock.Anything,
			mock.Anything,
			mock.Anything,
		)
	})

	t.Run("storage error", func(t *testing.T) {
		mockStorage := new(MockURLStorage)
		mockGenPass := new(MockGenPass)

		service := createTestShortenURLService(
			mockStorage,
			mockGenPass,
		)

		expectedError := errors.New("redis store failed")

		mockGenPass.On(
			"GenerateRandomString",
			6,
		).Return("abc123", nil).Once()

		mockStorage.On(
			"GetURL",
			mock.Anything,
			"abc123",
		).Return("", nil).Once()

		mockStorage.On(
			"StoreURL",
			mock.Anything,
			"abc123",
			originalURL,
		).Return(expectedError).Once()

		result, err := service.CreateShortenLink(
			ctx,
			originalURL,
		)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.NotNil(t, result)
		assert.Equal(t, "", result.KeyCode)
		assert.Equal(t, expectedError.Error(), result.Message)

		mockGenPass.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})

	t.Run("generate new key when key already exists", func(t *testing.T) {
		mockStorage := new(MockURLStorage)
		mockGenPass := new(MockGenPass)

		service := createTestShortenURLService(
			mockStorage,
			mockGenPass,
		)

		mockGenPass.On(
			"GenerateRandomString",
			6,
		).Return("abc123", nil).Once()

		mockStorage.On(
			"GetURL",
			mock.Anything,
			"abc123",
		).Return("https://existing.com", nil).Once()

		mockGenPass.On(
			"GenerateRandomString",
			6,
		).Return("xyz789", nil).Once()

		mockStorage.On(
			"GetURL",
			mock.Anything,
			"xyz789",
		).Return("", nil).Once()

		mockStorage.On(
			"StoreURL",
			mock.Anything,
			"xyz789",
			originalURL,
		).Return(nil).Once()

		result, err := service.CreateShortenLink(
			ctx,
			originalURL,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "xyz789", result.KeyCode)
		assert.Equal(
			t,
			"Shorten URL generated successfully",
			result.Message,
		)

		mockGenPass.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})

	t.Run("generate password error during retry", func(t *testing.T) {
		mockStorage := new(MockURLStorage)
		mockGenPass := new(MockGenPass)

		service := createTestShortenURLService(
			mockStorage,
			mockGenPass,
		)

		expectedError := errors.New("generate random string failed")

		mockGenPass.On(
			"GenerateRandomString",
			6,
		).Return("abc123", nil).Once()

		mockStorage.On(
			"GetURL",
			mock.Anything,
			"abc123",
		).Return("https://existing.com", nil).Once()

		mockGenPass.On(
			"GenerateRandomString",
			6,
		).Return("", expectedError).Once()

		result, err := service.CreateShortenLink(
			ctx,
			originalURL,
		)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.NotNil(t, result)
		assert.Equal(t, "", result.KeyCode)
		assert.Equal(t, expectedError.Error(), result.Message)

		mockGenPass.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})

	t.Run("GetURL error is ignored", func(t *testing.T) {
		mockStorage := new(MockURLStorage)
		mockGenPass := new(MockGenPass)

		service := createTestShortenURLService(
			mockStorage,
			mockGenPass,
		)

		mockGenPass.On(
			"GenerateRandomString",
			6,
		).Return("abc123", nil).Once()

		mockStorage.On(
			"GetURL",
			mock.Anything,
			"abc123",
		).Return("", errors.New("redis get failed")).Once()

		mockStorage.On(
			"StoreURL",
			mock.Anything,
			"abc123",
			originalURL,
		).Return(nil).Once()

		result, err := service.CreateShortenLink(
			ctx,
			originalURL,
		)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "abc123", result.KeyCode)

		mockGenPass.AssertExpectations(t)
		mockStorage.AssertExpectations(t)
	})
}
