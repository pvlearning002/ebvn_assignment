package utils_log

import (
	"runtime"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type CustomLogger interface {
	SetLogLevel(levelStr string)
	GetPathFunction(skip ...int) string
	BuildLogMessage(message string, skip ...int) string
	WriteLogInfo(message string, skip ...int)
	WriteLogError(message string, skip ...int)
	WriteLogDebug(message string, skip ...int)
	WriteLogWarning(message string, skip ...int)
}

type customLogger struct {
}

func NewCustomLogger() CustomLogger {
	return &customLogger{}
}

func (self *customLogger) SetLogLevel(levelStr string) {
	level := zerolog.NoLevel
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
}

func (self *customLogger) GetPathFunction(skip ...int) string {
	default_skip := 3
	if len(skip) > 0 {
		default_skip = skip[0]
	}
	pc, _, _, ok := runtime.Caller(default_skip)
	result := "unknown"
	if ok {
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			result = fn.Name()
		}
	}
	return result
}

// Added more function path into log
func (self *customLogger) BuildLogMessage(message string, skip ...int) string {
	caller := self.GetPathFunction(skip...)
	output := "func: " + caller + " - log: " + message
	return output
}

func (self *customLogger) WriteLogInfo(message string, skip ...int) {
	log.Info().Msg(self.BuildLogMessage(message, skip...))
}

func (self *customLogger) WriteLogError(message string, skip ...int) {
	log.Error().Msg(self.BuildLogMessage(message, skip...))
}

func (self *customLogger) WriteLogDebug(message string, skip ...int) {
	log.Debug().Msg(self.BuildLogMessage(message, skip...))
}

func (self *customLogger) WriteLogWarning(message string, skip ...int) {
	log.Warn().Msg(self.BuildLogMessage(message, skip...))
}
