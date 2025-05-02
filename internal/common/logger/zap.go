package logger

import (
	"go.uber.org/zap"
)

var Sugared *zap.SugaredLogger

func init() {
	log := zap.Must(zap.NewDevelopment())
	Sugared = log.Sugar()
}
