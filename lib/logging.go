package lib

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func CreateProductionLogger(component string) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"component": component,
		},
	}

	return zap.Must(config.Build())
}

func CreateChildLogger(logger *zap.Logger, service string) *zap.SugaredLogger {
	childLogger := logger.With(
		zap.String("service", service),
	)

	return childLogger.Sugar()
}
