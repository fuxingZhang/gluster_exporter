package logger

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/prometheus/common/promlog"
)

var logger log.Logger

func Init(config *promlog.Config) {
	logger = promlog.New(config)
}

func Logger() log.Logger {
	return logger
}

func Debug(keyvals ...any) {
	level.Debug(logger).Log(keyvals...)
}

func Warn(keyvals ...any) {
	level.Warn(logger).Log(keyvals...)
}

func Info(keyvals ...any) {
	level.Info(logger).Log(keyvals...)
}

func Error(keyvals ...any) {
	level.Error(logger).Log(keyvals...)
}
