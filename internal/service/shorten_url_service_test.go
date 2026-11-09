package service

import (
	"context"
	"errors"
	"testing"

	config "github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/internal/repository"
	"github.com/stretchr/testify/assert"
	testifymock "github.com/stretchr/testify/mock"
)

// Mock GenPass
type MockKeyGenerator struct {
	testifymock.Mock
}

func (m *MockKeyGenerator) GeneratePassword() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// Mock URLStorage
type MockURLStorage struct {
	testifymock.Mock
}

func (m *MockURLStorage) StoreURL(
	ctx context.Context,
	key string,
	url string,
) error {
	args := m.Called(ctx, key, url)
	return args.Error(0)
}

func (m *MockURLStorage) GetURL(
	ctx context.Context,
	key string,
) (string, error) {
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

// Compile-time interface checks
var _ GenPass = (*MockKeyGenerator)(nil)
var _ repository.URLStorage = (*MockURLStorage)(nil)

func TestCreateShortenLink(t *testing.T) {
	cfg := *config.CreateNewConfig()
	ctx := context.Background()
	generateError := errors.New("generate password failed")
	storageError := errors.New("store URL failed")
	key_test := "abc123"
	url_test := "https://google.com"

	testCases := []struct {
		name      string
		setupMock func(
			*MockKeyGenerator,
			*MockURLStorage,
		)
		expectedResult string
		expectedError  error
	}{
		{
			name: "success",
			setupMock: func(
				mockGenerator *MockKeyGenerator,
				mockStorage *MockURLStorage,
			) {
				mockGenerator.
					On("GeneratePassword").
					Return(key_test, nil).
					Once()

				mockStorage.
					On("GetURL", ctx, key_test).
					Return("", nil).
					Once()

				mockStorage.
					On("StoreURL", ctx, key_test, url_test).
					Return(nil).
					Once()
			},
			expectedResult: key_test,
			expectedError:  nil,
		},
		{
			name: "generate password error",
			setupMock: func(
				mockGenerator *MockKeyGenerator,
				mockStorage *MockURLStorage,
			) {
				mockGenerator.
					On("GeneratePassword").
					Return("", generateError).
					Once()
			},
			expectedResult: cfg.StringEmpty,
			expectedError:  generateError,
		},
		{
			name: "store URL error",
			setupMock: func(
				mockGenerator *MockKeyGenerator,
				mockStorage *MockURLStorage,
			) {
				mockGenerator.
					On("GeneratePassword").
					Return(key_test, nil).
					Once()

				// Bổ sung mock GetURL
				mockStorage.
					On("GetURL", ctx, key_test).
					Return("", nil).
					Once()

				// Mock StoreURL trả về lỗi
				mockStorage.
					On("StoreURL", ctx, key_test, url_test).
					Return(storageError).
					Once()
			},
			expectedResult: cfg.StringEmpty,
			expectedError:  storageError,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockGenerator := new(MockKeyGenerator)
			mockStorage := new(MockURLStorage)

			tc.setupMock(mockGenerator, mockStorage)

			service := &shortenURLService{
				urlStorage: mockStorage,
				genPass:    mockGenerator,
				cfg:        cfg,
			}

			result, err := service.CreateShortenLink(
				ctx,
				url_test,
			)

			assert.Equal(t, tc.expectedResult, result)

			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}

			mockGenerator.AssertExpectations(t)
			mockStorage.AssertExpectations(t)
		})
	}
}
