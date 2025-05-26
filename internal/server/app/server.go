package app

import (
	"context"
	"database/sql"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mkolibaba/metrics/internal/common/log"
	"github.com/mkolibaba/metrics/internal/server/config"
	"github.com/mkolibaba/metrics/internal/server/http/router"
	"github.com/mkolibaba/metrics/internal/server/storage/jsonfile"
	"github.com/mkolibaba/metrics/internal/server/storage/postgres"
	"github.com/mkolibaba/metrics/migrations"
	"go.uber.org/zap"
	stdlog "log"
	"net/http"
)

func Run() {
	ctx := context.Background()

	cfg := mustCreateConfig()

	logger := log.New()

	db := mustCreateDB(cfg.DatabaseDSN, logger)
	defer db.Close()

	var r chi.Router
	if cfg.DatabaseDSN != "" {
		store := postgres.New(db, logger)
		runMigrations(ctx, db, logger)

		r = router.New(store, db, logger)
	} else {
		store := mustCreateFileStorage(cfg, logger)
		defer store.Close()

		r = router.New(store, db, logger)
	}

	logger.Infof("running server on %s", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		logger.Fatal(err)
	}
}

func mustCreateConfig() *config.ServerConfig {
	cfg, err := config.LoadServerConfig()
	if err != nil {
		stdlog.Fatalf("error creating config: %v", err)
	}
	return cfg
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

func runMigrations(ctx context.Context, db *sql.DB, logger *zap.SugaredLogger) {
	err := migrations.Run(ctx, db, migrations.AppServer)
	if err != nil {
		logger.Fatalf("error during db migration: %v", err)
	}
}
