package http

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/mkolibaba/metrics/internal/server/config"
	"github.com/mkolibaba/metrics/internal/server/http/router"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type Server struct {
	r      chi.Router
	logger *zap.SugaredLogger
}

func NewServer(
	store router.MetricsStorage,
	db *sql.DB,
	cfg *config.ServerConfig,
	logger *zap.SugaredLogger,
) (*Server, error) {
	r, err := router.New(store, db, cfg, logger)
	if err != nil {
		return nil, err
	}

	return &Server{
		r:      r,
		logger: logger,
	}, nil
}

func (s *Server) Start(ctx context.Context, addr string) error {
	s.logger.Infof("running http server on %s", addr)

	server := http.Server{Addr: addr, Handler: s.r}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		<-ctx.Done()
		if err := server.Shutdown(context.Background()); err != nil {
			return fmt.Errorf("error shutting down server: %w", err)
		}
		return nil
	})

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	s.logger.Info("server stopped")

	return g.Wait()
}
