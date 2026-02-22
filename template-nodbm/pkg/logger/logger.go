package logger

import (

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(level string) (*zap.Logger, error) {
	var config zap.Config

	if level == "debug" {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// Parse level string
	atomicLevel, err := zap.ParseAtomicLevel(level)
	if err == nil {
		config.Level = atomicLevel
	}

	config.OutputPaths = []string{"stdout"}
	return config.Build()
}
