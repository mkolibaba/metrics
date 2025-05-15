package app

import (
	"github.com/mkolibaba/metrics/internal/common/log"
	"github.com/mkolibaba/metrics/internal/server/config"
	"github.com/mkolibaba/metrics/internal/server/http/router"
	"github.com/mkolibaba/metrics/internal/server/storage/jsonfile"
	"go.uber.org/zap"
	"net/http"
)

func Run() {
	cfg := config.MustLoadServerConfig()

	logger := log.New()

	store := mustCreateFileStorage(cfg, logger)
	defer store.Close()
	r := router.New(store, logger)

	logger.Infof("running server on %s", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		logger.Fatal(err)
	}
}

func mustCreateFileStorage(cfg *config.ServerConfig, logger *zap.SugaredLogger) *jsonfile.FileStorage {
	store, err := jsonfile.NewFileStorage(cfg.FileStoragePath, cfg.StoreInterval, cfg.Restore, logger)
	if err != nil {
		logger.Fatalf("error creating file storage: %v", err)
	}
	return store
}
