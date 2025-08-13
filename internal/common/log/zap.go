package log

import (
	"go.uber.org/zap"
)

// New создает и настраивает новый экземпляр zap.SugaredLogger
// для использования в приложении.
func New() *zap.SugaredLogger {
	return zap.Must(zap.NewDevelopment()).Sugar()
}
