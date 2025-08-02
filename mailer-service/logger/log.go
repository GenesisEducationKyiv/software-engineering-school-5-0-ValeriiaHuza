package logger

import (
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func InitZapLogger() error {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		return err
	}
	return nil
}

func Sync() {
	if logger != nil {
		_ = logger.Sync()
	}
}

func GetLogger() *zap.Logger {
	return logger
}

func InitTestLogger() {
	logger = zap.NewExample()
}
