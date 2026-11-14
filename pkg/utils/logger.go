package utils

import (
	"runtime"

	"github.com/rs/zerolog/log"
)

func BuildLogMessage(message string) string {
	pc, _, _, ok := runtime.Caller(2)
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

func WriteLogInfo(message string) {
	log.Info().Msg(BuildLogMessage(message))
}

func WriteLogError(message string) {
	log.Error().Msg(BuildLogMessage(message))
}

func WriteLogDebug(message string) {
	log.Debug().Msg(BuildLogMessage(message))
}

func WriteLogWarning(message string) {
	log.Warn().Msg(BuildLogMessage(message))
}
