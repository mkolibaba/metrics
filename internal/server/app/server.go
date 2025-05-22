package app

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mkolibaba/metrics/internal/common/log"
	"github.com/mkolibaba/metrics/internal/server/config"
	"github.com/mkolibaba/metrics/internal/server/http/router"
	"github.com/mkolibaba/metrics/internal/server/storage/jsonfile"
	"go.uber.org/zap"
	"net/http"
)

func Run() {
	cfg := mustCreateConfig()

	logger := log.New()

	store := mustCreateFileStorage(cfg, logger)
	defer store.Close()

	db := mustCreateDB(cfg.DatabaseDSN, logger)
	defer db.Close()

	r := router.New(store, db, logger)

	logger.Infof("running server on %s", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		logger.Fatal(err)
	}
}

// TODO
func mustCreateConfig() *config.ServerConfig {
	return config.MustLoadServerConfig()
}

func mustCreateFileStorage(cfg *config.ServerConfig, logger *zap.SugaredLogger) *jsonfile.FileStorage {
	store, err := jsonfile.NewFileStorage(cfg.FileStoragePath, cfg.StoreInterval, cfg.Restore, logger)
	if err != nil {
		logger.Fatalf("error creating file storage: %v", err)
	}
	return store
}

func mustCreateDB(databaseDSN string, logger *zap.SugaredLogger) *sql.DB {
	db, err := sql.Open("pgx", databaseDSN)
	if err != nil {
		logger.Fatalf("error creating db: %v", err)
	}
	return db
}
