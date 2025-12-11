package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getCore(level zap.AtomicLevel, isDev bool) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	var encoder zapcore.Encoder
	if isDev {
		encoder = zapcore.NewConsoleEncoder(developmentCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(productionCfg)
	}

	return zapcore.NewTee(
		zapcore.NewCore(encoder, stdout, level),
		zapcore.NewCore(encoder, file, level),
	)
}

func getAtomicLevel(logLevel string) zap.AtomicLevel {
	var level zapcore.Level
	if err := level.Set(logLevel); err != nil {
		level = zapcore.InfoLevel
	}

	return zap.NewAtomicLevelAt(level)
}
