package logger

import (
	"go.uber.org/zap"
)

var log *zap.SugaredLogger

func init() {
	logger, _ := zap.NewProduction()
	log = logger.Sugar()
}

func Get() *zap.SugaredLogger {
	return log
}

func Info(msg string, keysAndValues ...interface{}) {
	log.Infow(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...interface{}) {
	log.Errorw(msg, keysAndValues...)
}

func Debug(msg string, keysAndValues ...interface{}) {
	log.Debugw(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...interface{}) {
	log.Warnw(msg, keysAndValues...)
}
