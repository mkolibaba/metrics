package app

import (
	"github.com/mkolibaba/metrics/internal/common/logger"
	"github.com/mkolibaba/metrics/internal/server/config"
	"github.com/mkolibaba/metrics/internal/server/http/router"
	"github.com/mkolibaba/metrics/internal/server/storage/jsonfile"
	"net/http"
	"strings"
)

func Run() {
	cfg := config.MustLoadServerConfig()

	store := mustCreateFileStorage(cfg)
	defer store.Close()
	r := router.New(store)

	logger.Sugared.Infof("Running server on %s", serverAddress)
	if err := http.ListenAndServe(serverAddress, r); err != nil {
	logger.Sugared.Infof("running server on %s", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		logger.Sugared.Fatal(err)
	}
}

func mustCreateFileStorage(cfg *config.ServerConfig) *jsonfile.FileStorage {
	store, err := jsonfile.NewFileStorage(cfg.FileStoragePath, cfg.StoreInterval, cfg.Restore)
	if err != nil {
		logger.Sugared.Fatalf("error creating file storage: %v", err)
	}
	return store
}
