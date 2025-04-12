package logger

import (
	"go.uber.org/zap"
)

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
}

type Field struct {
	Key   string
	Value interface{}
}

type loggerImpl struct {
	zap *zap.Logger
}

func New(level string) Logger {
	cfg := zap.NewProductionConfig()

	switch level {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		cfg.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	zapLogger, err := cfg.Build()
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}

	return &loggerImpl{zap: zapLogger}
}

func (l *loggerImpl) Debug(msg string, fields ...Field) {
	l.zap.Debug(msg, toZapFields(fields)...)
}

func (l *loggerImpl) Info(msg string, fields ...Field) {
	l.zap.Info(msg, toZapFields(fields)...)
}

func (l *loggerImpl) Warn(msg string, fields ...Field) {
	l.zap.Warn(msg, toZapFields(fields)...)
}

func (l *loggerImpl) Error(msg string, fields ...Field) {
	l.zap.Error(msg, toZapFields(fields)...)
}

func (l *loggerImpl) Fatal(msg string, fields ...Field) {
	l.zap.Fatal(msg, toZapFields(fields)...)
}

func toZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, f := range fields {
		zapFields[i] = zap.Any(f.Key, f.Value)
	}
	return zapFields
}
