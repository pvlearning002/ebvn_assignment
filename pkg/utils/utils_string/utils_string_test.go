package utils_string

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/pvlearning002/ebvn_assignment/internal/config"
)

// Mock CustomLogger
type MockCustomLogger struct {
	mock.Mock
}

func (m *MockCustomLogger) SetLogLevel(levelStr string) {
	m.Called(levelStr)
}

func (m *MockCustomLogger) GetPathFunction(skip ...int) string {
	args := m.Called(skip)
	return args.String(0)
}

func (m *MockCustomLogger) BuildLogMessage(message string, skip ...int) string {
	args := m.Called(message, skip)
	return args.String(0)
}

func (m *MockCustomLogger) WriteLogInfo(message string, skip ...int) {
	m.Called(message, skip)
}

func (m *MockCustomLogger) WriteLogError(message string, skip ...int) {
	m.Called(message, skip)
}

func (m *MockCustomLogger) WriteLogDebug(message string, skip ...int) {
	m.Called(message, skip)
}

func (m *MockCustomLogger) WriteLogWarning(message string, skip ...int) {
	m.Called(message, skip)
}

// Test NewUtilsString
func TestNewUtilsString(t *testing.T) {
	cfg := config.Config{}
	cfgError := config.ConfigError{}
	logger := &MockCustomLogger{}

	result := NewUtilsString(cfg, cfgError, logger)

	assert.NotNil(t, result)
	assert.IsType(t, &utilsString{}, result)
}

// Test GenerateRandomString
func TestGenerateRandomString(t *testing.T) {
	tests := []struct {
		name          string
		length        int
		expectedError bool
	}{
		{
			name:          "Length is zero",
			length:        0,
			expectedError: true,
		},
		{
			name:          "Length is negative",
			length:        -1,
			expectedError: true,
		},
		{
			name:          "Length is 1",
			length:        1,
			expectedError: false,
		},
		{
			name:          "Length is 12",
			length:        12,
			expectedError: false,
		},
		{
			name:          "Length is 32",
			length:        32,
			expectedError: false,
		},
		{
			name:          "Length is 64",
			length:        64,
			expectedError: false,
		},
		{
			name:          "Length is 100",
			length:        100,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.Config{
				Dash:        "-",
				StringEmpty: "",
			}

			cfgError := config.ConfigError{
				ErrorLengthZero: "length must be greater than zero",
			}

			logger := &MockCustomLogger{}

			if tt.expectedError {
				logger.On(
					"WriteLogError",
					cfgError.ErrorLengthZero,
					mock.Anything,
				).Return().Once()
			}

			utils := NewUtilsString(cfg, cfgError, logger)

			result, err := utils.(*utilsString).GenerateRandomString(tt.length)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, cfgError.ErrorLengthZero, err.Error())
				assert.Empty(t, result)
				logger.AssertExpectations(t)
				return
			}

			assert.NoError(t, err)
			assert.NotEmpty(t, result)
			assert.Equal(t, tt.length, len(result))
			assert.NotContains(t, result, "-")
			assert.True(t, len(result) == tt.length)
			logger.AssertExpectations(t)
		})
	}
}

// Test GenerateRandomString with multiple generated values
func TestGenerateRandomString_Unique(t *testing.T) {
	cfg := config.Config{
		Dash:        "-",
		StringEmpty: "",
	}

	cfgError := config.ConfigError{
		ErrorLengthZero: "length must be greater than zero",
	}

	logger := &MockCustomLogger{}

	utils := NewUtilsString(cfg, cfgError, logger)

	first, err1 := utils.(*utilsString).GenerateRandomString(32)
	second, err2 := utils.(*utilsString).GenerateRandomString(32)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, first, second)
}

// Test GenerateRandomString returns hexadecimal characters
func TestGenerateRandomString_ValidCharacters(t *testing.T) {
	cfg := config.Config{
		Dash:        "-",
		StringEmpty: "",
	}

	cfgError := config.ConfigError{}
	logger := &MockCustomLogger{}

	utils := NewUtilsString(cfg, cfgError, logger)

	result, err := utils.(*utilsString).GenerateRandomString(64)

	assert.NoError(t, err)
	assert.Len(t, result, 64)

	for _, char := range result {
		assert.True(
			t,
			(char >= '0' && char <= '9') ||
				(char >= 'a' && char <= 'f'),
			"Invalid character: %c",
			char,
		)
	}

	assert.True(t, strings.ToLower(result) == result)
}
