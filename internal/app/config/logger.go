package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetupLogger() {
	// Init zap log
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.DisableStacktrace = true

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	logger, _ := config.Build()
	zap.ReplaceGlobals(logger)
}
