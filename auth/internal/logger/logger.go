package logger

import (
	"go.uber.org/zap"
)

var globalLogger *zap.Logger

func Init(isDev bool, level string) {
	core := getCore(getAtomicLevel(level), isDev)
	globalLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func Debug(msg string, fields ...zap.Field) {
	globalLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	globalLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	globalLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	globalLogger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	globalLogger.Fatal(msg, fields...)
}

func WithOptions(opts ...zap.Option) *zap.Logger {
	return globalLogger.WithOptions(opts...)
}
