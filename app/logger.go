package app

import (
	"os"

	"github.com/humaniq/hmq-finance-price-feed/pkg/blogger"
)

var loggerInstance blogger.Logger

func InitDefaultLogger(level uint8) {
	loggerInstance = blogger.NewLog([]func(next blogger.BufferedHandler) blogger.BufferedHandler{
		blogger.LogLevelFilter(level),
	},
		blogger.NewIOWriterRouter(os.Stdout, os.Stderr, os.Stderr, true))
}
func InitLogger(logger blogger.Logger) {
	loggerInstance = logger
}

func Logger() blogger.Logger {
	if loggerInstance == nil {
		InitDefaultLogger(blogger.LevelDefault)
	}
	return loggerInstance
}
