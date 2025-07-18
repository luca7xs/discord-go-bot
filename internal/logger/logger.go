package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// Init inicializa o logger com configurações coloridas e legíveis
func Init() error {
	config := zap.NewDevelopmentConfig()

	// Configuração para output colorido e legível
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	var err error
	Log, err = config.Build()
	if err != nil {
		return err
	}

	return nil
}

// Sync força a escrita de logs pendentes
func Sync() {
	_ = Log.Sync()
}
