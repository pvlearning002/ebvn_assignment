package utils_log

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestNewCustomLogger(t *testing.T) {
	logger := NewCustomLogger()

	if logger == nil {
		t.Fatal("NewCustomLogger() returned nil")
	}

	_, ok := logger.(*customLogger)
	if !ok {
		t.Error("NewCustomLogger() did not return *customLogger")
	}
}

func TestSetLogLevel(t *testing.T) {
	originalLevel := zerolog.GlobalLevel()
	defer zerolog.SetGlobalLevel(originalLevel)

	logger := NewCustomLogger()

	tests := []struct {
		name          string
		level         string
		expectedLevel zerolog.Level
	}{
		{
			name:          "Debug level",
			level:         "debug",
			expectedLevel: zerolog.DebugLevel,
		},
		{
			name:          "Info level",
			level:         "info",
			expectedLevel: zerolog.InfoLevel,
		},
		{
			name:          "Warn level",
			level:         "warn",
			expectedLevel: zerolog.WarnLevel,
		},
		{
			name:          "Error level",
			level:         "error",
			expectedLevel: zerolog.ErrorLevel,
		},
		{
			name:          "Fatal level",
			level:         "fatal",
			expectedLevel: zerolog.FatalLevel,
		},
		{
			name:          "Disabled level",
			level:         "disabled",
			expectedLevel: zerolog.Disabled,
		},
		{
			name:          "Invalid level",
			level:         "invalid",
			expectedLevel: zerolog.InfoLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.SetLogLevel(tt.level)

			if got := zerolog.GlobalLevel(); got != tt.expectedLevel {
				t.Errorf(
					"SetLogLevel(%q) = %v, want %v",
					tt.level,
					got,
					tt.expectedLevel,
				)
			}
		})
	}
}

func TestGetPathFunction(t *testing.T) {
	logger := NewCustomLogger()

	tests := []struct {
		name string
		skip []int
	}{
		{
			name: "Default skip",
			skip: nil,
		},
		{
			name: "Custom skip",
			skip: []int{1},
		},
		{
			name: "Custom skip 2",
			skip: []int{2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := logger.GetPathFunction(tt.skip...)

			if result == "" {
				t.Error("GetPathFunction() returned empty string")
			}

			if result == "unknown" {
				t.Error("GetPathFunction() returned unknown")
			}
		})
	}
}

func TestGetPathFunction_InvalidSkip(t *testing.T) {
	logger := NewCustomLogger()

	result := logger.GetPathFunction(999999)

	if result != "unknown" {
		t.Errorf(
			"GetPathFunction(999999) = %q, want %q",
			result,
			"unknown",
		)
	}
}

func TestBuildLogMessage(t *testing.T) {
	logger := NewCustomLogger()

	tests := []struct {
		name    string
		message string
		skip    []int
	}{
		{
			name:    "Default skip",
			message: "test info message",
			skip:    nil,
		},
		{
			name:    "Custom skip",
			message: "test error message",
			skip:    []int{1},
		},
		{
			name:    "Empty message",
			message: "",
			skip:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := logger.BuildLogMessage(tt.message, tt.skip...)

			if !strings.HasPrefix(result, "func: ") {
				t.Errorf("BuildLogMessage() missing function prefix: %q", result)
			}

			if !strings.Contains(result, " - log: ") {
				t.Errorf("BuildLogMessage() missing log separator: %q", result)
			}

			if !strings.HasSuffix(result, tt.message) {
				t.Errorf(
					"BuildLogMessage() = %q, want suffix %q",
					result,
					tt.message,
				)
			}
		})
	}
}

func TestWriteLogInfo(t *testing.T) {
	testWriteLog(t, "info", func(logger CustomLogger) {
		logger.WriteLogInfo("test info message", 1)
	})
}

func TestWriteLogError(t *testing.T) {
	testWriteLog(t, "error", func(logger CustomLogger) {
		logger.WriteLogError("test error message", 1)
	})
}

func TestWriteLogDebug(t *testing.T) {
	testWriteLog(t, "debug", func(logger CustomLogger) {
		logger.WriteLogDebug("test debug message", 1)
	})
}

func TestWriteLogWarning(t *testing.T) {
	testWriteLog(t, "warn", func(logger CustomLogger) {
		logger.WriteLogWarning("test warning message", 1)
	})
}

func testWriteLog(
	t *testing.T,
	level string,
	writeFunc func(CustomLogger),
) {
	t.Helper()

	originalLogger := log.Logger
	originalLevel := zerolog.GlobalLevel()

	var output bytes.Buffer

	log.Logger = zerolog.New(&output)
	zerolog.SetGlobalLevel(zerolog.TraceLevel)

	defer func() {
		log.Logger = originalLogger
		zerolog.SetGlobalLevel(originalLevel)
	}()

	logger := NewCustomLogger()
	writeFunc(logger)

	result := output.String()

	if result == "" {
		t.Fatal("WriteLog() produced no output")
	}

	if !strings.Contains(result, `"level":"`+level+`"`) {
		t.Errorf(
			"WriteLog() output = %q, want level %q",
			result,
			level,
		)
	}

	if !strings.Contains(result, "test ") {
		t.Errorf("WriteLog() output does not contain message: %q", result)
	}

	if !strings.Contains(result, "func: ") {
		t.Errorf("WriteLog() output does not contain function path: %q", result)
	}
}
