package service

import (
	"context"
	"errors"
	"testing"

	"github.com/pvlearning002/ebvn_assignment/internal/app/repository"
	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/pkg/utils"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type shortenServiceStringUtilsMock struct {
	mock.Mock
	generate func(int) (string, error)
}

func (m *shortenServiceStringUtilsMock) GenerateRandomString(length int) (string, error) {
	if m.generate != nil {
		return m.generate(length)
	}
	args := m.Called(length)
	return args.String(0), args.Error(1)
}

type shortenServiceStorageMock struct{ mock.Mock }

func (m *shortenServiceStorageMock) GetURL(ctx context.Context, key string) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}
func (m *shortenServiceStorageMock) StoreURL(ctx context.Context, key, value string) error {
	return m.Called(ctx, key, value).Error(0)
}

type shortenServiceLoggerMock struct{ mock.Mock }

func (m *shortenServiceLoggerMock) SetLogLevel(level string) { m.Called(level) }
func (m *shortenServiceLoggerMock) GetPathFunction(skip ...int) string {
	return m.Called(skip).String(0)
}
func (m *shortenServiceLoggerMock) BuildLogMessage(message string, skip ...int) string {
	return m.Called(message, skip).String(0)
}
func (m *shortenServiceLoggerMock) WriteLogInfo(message string, skip ...int) { m.Called(message, skip) }
func (m *shortenServiceLoggerMock) WriteLogError(message string, skip ...int) {
	m.Called(message, skip)
}
func (m *shortenServiceLoggerMock) WriteLogDebug(message string, skip ...int) {
	m.Called(message, skip)
}
func (m *shortenServiceLoggerMock) WriteLogWarning(message string, skip ...int) {
	m.Called(message, skip)
}

var _ utils.StringUtils = (*shortenServiceStringUtilsMock)(nil)
var _ utils.LoggerUtils = (*shortenServiceLoggerMock)(nil)
var _ repository.URLStorage = (*shortenServiceStorageMock)(nil)

func newShortenServiceForTest(storage *shortenServiceStorageMock, stringUtils *shortenServiceStringUtilsMock, logger *shortenServiceLoggerMock) ShortenURL {
	return NewShortenURLService(storage,
		config.Config{KeycodeLengthDefault: 6, StringEmpty: ""}, logger, stringUtils)
}

func TestNewShortenURLService(t *testing.T) {
	got := newShortenServiceForTest(new(shortenServiceStorageMock), new(shortenServiceStringUtilsMock), new(shortenServiceLoggerMock))
	require.NotNil(t, got)
	assert.IsType(t, &shortenURLService{}, got)
}

func TestCreateShortenLink_Success(t *testing.T) {
	ctx := context.Background()
	storage := new(shortenServiceStorageMock)
	stringUtils := new(shortenServiceStringUtilsMock)
	stringUtils.On("GenerateRandomString", 6).Return("abc123", nil).Once()
	storage.On("GetURL", ctx, "abc123").Return("", nil).Once()
	storage.On("StoreURL", ctx, "abc123", "https://example.com").Return(nil).Once()

	got, err := newShortenServiceForTest(storage, stringUtils, new(shortenServiceLoggerMock)).CreateShortenLink(ctx, "https://example.com")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "abc123", got.KeyCode)
	assert.Equal(t, "Shorten URL generated successfully", got.Message)
	stringUtils.AssertExpectations(t)
	storage.AssertExpectations(t)
}

func TestCreateShortenLink_RetriesCollision(t *testing.T) {
	ctx := context.Background()
	keys := []string{"taken1", "free22"}
	stringUtils := &shortenServiceStringUtilsMock{generate: func(length int) (string, error) {
		require.Equal(t, 6, length)
		if len(keys) == 0 {
			t.Fatal("generated more keys than expected")
		}
		key := keys[0]
		keys = keys[1:]
		return key, nil
	}}
	storage := new(shortenServiceStorageMock)
	storage.On("GetURL", ctx, "taken1").Return("https://used.example", nil).Once()
	storage.On("GetURL", ctx, "free22").Return("", nil).Once()
	storage.On("StoreURL", ctx, "free22", "https://example.com").Return(nil).Once()

	got, err := newShortenServiceForTest(storage, stringUtils, new(shortenServiceLoggerMock)).CreateShortenLink(ctx, "https://example.com")
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "free22", got.KeyCode)
	storage.AssertExpectations(t)
}

func TestCreateShortenLink_GenerationError(t *testing.T) {
	ctx := context.Background()
	wantErr := errors.New("key generation failed")
	stringUtils := new(shortenServiceStringUtilsMock)
	stringUtils.On("GenerateRandomString", 6).Return("", wantErr).Once()
	logger := new(shortenServiceLoggerMock)
	logger.On("WriteLogError", wantErr.Error(), mock.Anything).Return().Twice()
	storage := new(shortenServiceStorageMock)

	got, err := newShortenServiceForTest(storage, stringUtils, logger).CreateShortenLink(ctx, "https://example.com")
	require.ErrorIs(t, err, wantErr)
	require.NotNil(t, got)
	assert.Empty(t, got.KeyCode)
	assert.Equal(t, wantErr.Error(), got.Message)
	storage.AssertNotCalled(t, "StoreURL", mock.Anything, mock.Anything, mock.Anything)
	logger.AssertExpectations(t)
}

func TestCreateShortenLink_GenerationErrorDuringRetry(t *testing.T) {
	ctx := context.Background()
	wantErr := errors.New("key generation failed during retry")
	generationCount := 0
	stringUtils := &shortenServiceStringUtilsMock{generate: func(length int) (string, error) {
		require.Equal(t, 6, length)
		generationCount++
		if generationCount == 1 {
			return "taken1", nil
		}
		return "", wantErr
	}}
	storage := new(shortenServiceStorageMock)
	storage.On("GetURL", ctx, "taken1").Return("https://used.example", nil).Once()
	logger := new(shortenServiceLoggerMock)
	// generateKey logs the retry error, then CreateShortenLink logs it again.
	logger.On("WriteLogError", wantErr.Error(), mock.Anything).Return().Twice()

	got, err := newShortenServiceForTest(storage, stringUtils, logger).CreateShortenLink(ctx, "https://example.com")
	require.ErrorIs(t, err, wantErr)
	require.NotNil(t, got)
	assert.Empty(t, got.KeyCode)
	assert.Equal(t, wantErr.Error(), got.Message)
	assert.Equal(t, 2, generationCount)
	storage.AssertNotCalled(t, "StoreURL", mock.Anything, mock.Anything, mock.Anything)
	storage.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestCreateShortenLink_StoreError(t *testing.T) {
	ctx := context.Background()
	wantErr := errors.New("redis store failed")
	storage := new(shortenServiceStorageMock)
	stringUtils := new(shortenServiceStringUtilsMock)
	stringUtils.On("GenerateRandomString", 6).Return("abc123", nil).Once()
	storage.On("GetURL", ctx, "abc123").Return("", nil).Once()
	storage.On("StoreURL", ctx, "abc123", "https://example.com").Return(wantErr).Once()
	logger := new(shortenServiceLoggerMock)
	logger.On("WriteLogError", wantErr.Error(), mock.Anything).Return().Once()

	got, err := newShortenServiceForTest(storage, stringUtils, logger).CreateShortenLink(ctx, "https://example.com")
	require.ErrorIs(t, err, wantErr)
	require.NotNil(t, got)
	assert.Empty(t, got.KeyCode)
	assert.Equal(t, wantErr.Error(), got.Message)
	logger.AssertExpectations(t)
	storage.AssertExpectations(t)
}

func TestGetUrlByCode(t *testing.T) {
	ctx := context.Background()
	tests := []struct {
		name      string
		value     string
		getErr    error
		wantValue string
		wantErr   error
		wantLog   bool
	}{
		{name: "found", value: "https://example.com", wantValue: "https://example.com"},
		{name: "missing", getErr: redis.Nil, wantLog: true},
		{name: "storage error", getErr: errors.New("redis unavailable"), wantErr: errors.New("redis unavailable"), wantLog: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			storage := new(shortenServiceStorageMock)
			storage.On("GetURL", ctx, "abc123").Return(tc.value, tc.getErr).Once()
			logger := new(shortenServiceLoggerMock)
			if tc.wantLog {
				logger.On("WriteLogError", tc.getErr.Error(), mock.Anything).Return().Once()
			}

			got, err := newShortenServiceForTest(storage, new(shortenServiceStringUtilsMock), logger).GetUrlByCode(ctx, "abc123")
			assert.Equal(t, tc.wantValue, got)
			if tc.wantErr != nil {
				require.Error(t, err)
				assert.EqualError(t, err, tc.wantErr.Error())
			} else {
				require.NoError(t, err)
			}
			storage.AssertExpectations(t)
			logger.AssertExpectations(t)
		})
	}
}
