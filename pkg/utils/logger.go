package utils

import (
	"runtime"

	"github.com/rs/zerolog/log"
)

func BuildLogMessage(message string, skip ...int) string {
	default_skip := 2
	if len(skip) > 0 {
		default_skip = skip[0]
	}
	pc, _, _, ok := runtime.Caller(default_skip)
	caller := "unknown"
	if ok {
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			caller = fn.Name()
		}
	}
	output := "func: " + caller + " - log: " + message
	return output
}

func WriteLogInfo(message string, skip ...int) {
	log.Info().Msg(BuildLogMessage(message, skip...))
}

func WriteLogError(message string, skip ...int) {
	log.Error().Msg(BuildLogMessage(message, skip...))
}

func WriteLogDebug(message string, skip ...int) {
	log.Debug().Msg(BuildLogMessage(message, skip...))
}

func WriteLogWarning(message string, skip ...int) {
	log.Warn().Msg(BuildLogMessage(message, skip...))
}
