package log

import (
	"go.uber.org/zap"
)

func New() *zap.SugaredLogger {
	return zap.Must(zap.NewDevelopment()).Sugar()
}
