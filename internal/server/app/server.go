package app

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mkolibaba/metrics/internal/common/log"
	"github.com/mkolibaba/metrics/internal/server/config"
	"github.com/mkolibaba/metrics/internal/server/grpc"
	"github.com/mkolibaba/metrics/internal/server/http"
	"github.com/mkolibaba/metrics/internal/server/http/router"
	"github.com/mkolibaba/metrics/internal/server/storage/inmemory"
	"github.com/mkolibaba/metrics/internal/server/storage/jsonfile"
	"github.com/mkolibaba/metrics/internal/server/storage/postgres"
	"github.com/mkolibaba/metrics/migrations"
	"go.uber.org/zap"
	stdlog "log"
	"os/signal"
	"syscall"
)

type Server interface {
	Start(ctx context.Context, addr string) error
}

var noopFn = func() {}

// Run инициализирует и запускает сервер метрик.
// Функция настраивает конфигурацию, логгер, хранилище (в памяти, файловое или БД)
// и HTTP-роутер, а затем запускает сервер.
// Работает до прерывания.
func Run() {
	// context
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	// config
	cfg, err := config.LoadServerConfig()
	if err != nil {
		stdlog.Fatalf("error creating config: %v", err)
	}

	// logger
	logger := log.New()

	// database
	db, err := createDB(cfg.DatabaseDSN)
	if err != nil {
		logger.Fatalf("error creating db: %v", err)
	}
	defer db.Close()

	// storage
	store, closeFn, err := createStore(ctx, cfg, db, logger)
	if err != nil {
		logger.Fatalf("error creating store: %v", err)
	}
	defer closeFn()

	// server
	var server Server
	if cfg.UseGRPC {
		server = grpc.NewServer(store, logger)
	} else {
		server, err = http.NewServer(store, db, cfg, logger)
		if err != nil {
			logger.Fatalf("error creating server: %v", err)
		}
	}

	if err := server.Start(ctx, cfg.ServerAddress); err != nil {
		logger.Fatalf("server error: %v", err)
	}
}

func createDB(databaseDSN string) (*sql.DB, error) {
	return sql.Open("pgx", databaseDSN)
}

func createStore(
	ctx context.Context,
	cfg *config.ServerConfig,
	db *sql.DB,
	logger *zap.SugaredLogger,
) (router.MetricsStorage, func(), error) {
	if cfg.DatabaseDSN != "" {
		store := postgres.New(db)
		if err := migrations.Run(ctx, db, migrations.AppServer); err != nil {
			return nil, nil, fmt.Errorf("error during db migration: %w", err)
		}
		return store, noopFn, nil
	}

	if cfg.FileStoragePath != "" {
		store, err := jsonfile.NewFileStorage(cfg.FileStoragePath, cfg.StoreInterval, cfg.Restore, logger)
		if err != nil {
			return nil, nil, fmt.Errorf("error creating file storage: %w", err)
		}
		return store, store.Close, nil
	}

	return inmemory.NewMemStorage(), noopFn, nil
}
