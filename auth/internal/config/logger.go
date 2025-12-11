package config

import (
	"os"

	"github.com/pkg/errors"
)

const (
	logLevelEnvName = "LOG_LEVEL"
	envName         = "ENV"
)

type LoggerConfig interface {
	Level() string
	IsDev() bool
}

type loggerConfig struct {
	level string
	isDev bool
}

func NewLoggerConfig() (LoggerConfig, error) {
	level := os.Getenv(logLevelEnvName)
	if len(level) == 0 {
		level = "info" // default level
	}

	env := os.Getenv(envName)
	if len(env) == 0 {
		return nil, errors.New("env name not found")
	}

	isDev := false
	if env == "local" || env == "dev" {
		isDev = true
	}

	return &loggerConfig{
		level: level,
		isDev: isDev,
	}, nil
}

func (cfg *loggerConfig) Level() string {
	return cfg.level
}

func (cfg *loggerConfig) IsDev() bool {
	return cfg.isDev
}
