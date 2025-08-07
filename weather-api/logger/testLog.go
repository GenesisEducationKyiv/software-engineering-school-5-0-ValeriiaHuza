package logger

import (
	"bytes"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type TestBuffer interface {
	Logs() string
}

type testSink struct {
	mux sync.Mutex
	buf bytes.Buffer
}

func (s *testSink) Write(p []byte) (n int, err error) {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.buf.Write(p)
}

func (s *testSink) Logs() string {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.buf.String()
}

func NewTestLogger() (*Logger, TestBuffer) {
	sink := &testSink{}

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        "", // Omit time
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "", // Omit caller
		FunctionKey:    "", // Omit function
		MessageKey:     "msg",
		StacktraceKey:  "", // Omit stack
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     nil, // omit
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.AddSync(sink),
		zapcore.DebugLevel,
	)

	logger := zap.New(core).Sugar()

	return &Logger{logger: logger}, sink
}
