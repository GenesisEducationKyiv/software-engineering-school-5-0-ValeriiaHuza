package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	logger *zap.SugaredLogger
}

func NewLogger() (*Logger, error) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderCfg)

	writer := zapcore.Lock(os.Stdout)

	logLevel := zapcore.InfoLevel
	core := zapcore.NewCore(encoder, writer, logLevel)

	// Wrap the core with sampling
	sampledCore := zapcore.NewSamplerWithOptions(
		core,
		// Sampling window (1 second here)
		// meaning it resets counters every second
		time.Second,
		100, // first 100 logs per second are logged
		10,  // then 1 every 10 is logged
	)

	// Create a new logger using the sampled core
	zapLogger := zap.New(sampledCore,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{
		logger: zapLogger.Sugar(),
	}, nil
}

func (l *Logger) Info(msg string, keysAndValues ...any) {
	l.logger.Infow(msg, keysAndValues...)
}

func (l *Logger) Error(msg string, keysAndValues ...any) {
	l.logger.Errorw(msg, keysAndValues...)
}

func (l *Logger) Debug(msg string, keysAndValues ...any) {
	l.logger.Debugw(msg, keysAndValues...)
}

func (l *Logger) Sync() error {
	return l.logger.Sync()
}
