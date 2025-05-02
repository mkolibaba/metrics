package app

import (
	"github.com/mkolibaba/metrics/internal/common/logger"
	"github.com/mkolibaba/metrics/internal/server/config"
	"github.com/mkolibaba/metrics/internal/server/http/router"
	"github.com/mkolibaba/metrics/internal/server/storage/inmemory"
	"net/http"
	"strings"
)

func Run() {
	cfg := config.MustLoadServerConfig()
	// TODO: почему gzip для text/html не работает с localhost:port?
	serverAddress := strings.TrimPrefix(cfg.ServerAddress, "localhost")

	store := inmemory.NewMemStorage()
	r := router.New(store)

	logger.Sugared.Infof("Running server on %s", serverAddress)
	if err := http.ListenAndServe(serverAddress, r); err != nil {
		logger.Sugared.Fatal(err)
	}
}
