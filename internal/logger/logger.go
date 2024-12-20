package logger

import (
	"github.com/LI-SeNyA-vE/Gofermat/internal/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("15:04:05"))
}

func getCustomLoggerConfig(level string) (*zap.Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder
	encoderConfig.TimeKey = "time"

	config := zap.Config{
		Level:            lvl,
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return config.Build()
}

func Initialize(level string) error {
	logger, err := getCustomLoggerConfig(level)
	if err != nil {
		return err
	}
	global.Logger = logger.Sugar() // Преобразуем в SugaredLogger для удобства
	return nil
}
