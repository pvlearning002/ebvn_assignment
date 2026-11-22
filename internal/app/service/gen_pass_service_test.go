package service

import (
	"errors"
	"testing"

	"github.com/pvlearning002/ebvn_assignment/internal/config"
	"github.com/pvlearning002/ebvn_assignment/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type genPassStringUtilsMock struct{ mock.Mock }

func (m *genPassStringUtilsMock) GenerateRandomString(length int) (string, error) {
	args := m.Called(length)
	return args.String(0), args.Error(1)
}

type genPassLoggerMock struct{ mock.Mock }

func (m *genPassLoggerMock) SetLogLevel(level string) { m.Called(level) }
func (m *genPassLoggerMock) GetPathFunction(skip ...int) string {
	return m.Called(skip).String(0)
}
func (m *genPassLoggerMock) BuildLogMessage(message string, skip ...int) string {
	return m.Called(message, skip).String(0)
}
func (m *genPassLoggerMock) WriteLogInfo(message string, skip ...int) {
	m.Called(message, skip)
}
func (m *genPassLoggerMock) WriteLogError(message string, skip ...int) {
	m.Called(message, skip)
}
func (m *genPassLoggerMock) WriteLogDebug(message string, skip ...int) {
	m.Called(message, skip)
}
func (m *genPassLoggerMock) WriteLogWarning(message string, skip ...int) {
	m.Called(message, skip)
}

var _ utils.StringUtils = (*genPassStringUtilsMock)(nil)
var _ utils.LoggerUtils = (*genPassLoggerMock)(nil)

func TestNewGenPass(t *testing.T) {
	got := NewGenPass(config.Config{}, new(genPassLoggerMock), new(genPassStringUtilsMock))
	require.NotNil(t, got)
	assert.IsType(t, &genPass{}, got)
}

func TestGeneratePassword_Success(t *testing.T) {
	const passwordLength = 16
	const password = "a1b2c3d4e5f6g7h8"

	stringUtils := new(genPassStringUtilsMock)
	stringUtils.On("GenerateRandomString", passwordLength).Return(password, nil).Once()
	service := NewGenPass(
		config.Config{PasswordLengthDefault: passwordLength},
		nil,
		stringUtils,
	)

	got, err := service.GeneratePassword()
	require.NoError(t, err)
	assert.Equal(t, password, got)
	stringUtils.AssertExpectations(t)
}

func TestGeneratePassword_GenerationError(t *testing.T) {
	const passwordLength = 12
	wantErr := errors.New("random string generation failed")
	stringUtils := new(genPassStringUtilsMock)
	stringUtils.On("GenerateRandomString", passwordLength).Return("partial", wantErr).Once()
	logger := new(genPassLoggerMock)
	logger.On("WriteLogError", wantErr.Error(), mock.Anything).Return().Once()
	service := NewGenPass(config.Config{
		PasswordLengthDefault: passwordLength,
		StringEmpty:           "",
	}, logger, stringUtils)

	got, err := service.GeneratePassword()
	require.ErrorIs(t, err, wantErr)
	assert.Empty(t, got)
	stringUtils.AssertExpectations(t)
	logger.AssertExpectations(t)
}
