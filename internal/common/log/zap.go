package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New создает и настраивает новый экземпляр zap.SugaredLogger
// для использования в приложении.
func New() (*zap.SugaredLogger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}
