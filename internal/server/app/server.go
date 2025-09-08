package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mkolibaba/metrics/internal/common/log"
	"github.com/mkolibaba/metrics/internal/common/rsa"
	"github.com/mkolibaba/metrics/internal/server/config"
	"github.com/mkolibaba/metrics/internal/server/http/router"
	"github.com/mkolibaba/metrics/internal/server/storage/inmemory"
	"github.com/mkolibaba/metrics/internal/server/storage/jsonfile"
	"github.com/mkolibaba/metrics/internal/server/storage/postgres"
	"github.com/mkolibaba/metrics/migrations"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	stdlog "log"
	"net/http"
	"os/signal"
	"syscall"
)

var noopFn = func() {}

// Run инициализирует и запускает сервер метрик.
// Функция настраивает конфигурацию, логгер, хранилище (в памяти, файловое или БД)
// и HTTP-роутер, а затем запускает сервер.
// Работает до прерывания.
func Run() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer stop()

	cfg, err := config.LoadServerConfig()
	if err != nil {
		stdlog.Fatalf("error creating config: %v", err)
	}

	logger := log.New()

	db, err := createDB(cfg.DatabaseDSN)
	if err != nil {
		logger.Fatalf("error creating db: %v", err)
	}
	defer db.Close()

	store, closeFn, err := createStore(ctx, cfg, db, logger)
	if err != nil {
		logger.Fatalf("error creating store: %v", err)
	}
	defer closeFn()

	decryptor := rsa.NopDecryptor
	if cfg.CryptoKey != "" {
		decryptor, err = rsa.NewDecryptor(cfg.CryptoKey)
		if err != nil {
			logger.Fatalf("error creating decryptor: %v", err)
		}
	}

	r := router.New(store, db, cfg.Key, logger, decryptor)

	if err := runServer(ctx, cfg, r, logger); err != nil {
		logger.Fatalf("error running server: %v", err)
	}
}

func runServer(
	ctx context.Context,
	cfg *config.ServerConfig,
	r chi.Router,
	logger *zap.SugaredLogger,
) error {
	logger.Infof("running server on %s", cfg.ServerAddress)

	server := http.Server{Addr: cfg.ServerAddress, Handler: r}

	shutdown := make(chan struct{})

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			return fmt.Errorf("error shutting down server: %w", err)
		}
		close(shutdown)
		return nil
	})

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	<-shutdown

	logger.Info("server stopped")

	return g.Wait()
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
