package router

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/mkolibaba/metrics/internal/server/http/handlers/list"
	"github.com/mkolibaba/metrics/internal/server/http/handlers/ping"
	"github.com/mkolibaba/metrics/internal/server/http/handlers/read"
	"github.com/mkolibaba/metrics/internal/server/http/handlers/update"
	"github.com/mkolibaba/metrics/internal/server/http/middleware"
	"go.uber.org/zap"
)

type MetricsStorage interface {
	list.AllMetricsGetter
	read.MetricsGetter
	update.MetricsUpdater
}

func New(store MetricsStorage, db *sql.DB, logger *zap.SugaredLogger) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.Logger(logger))
	r.Use(middleware.Compressor(logger))

	r.Get("/", list.New(store, logger))
	r.Get("/value/{type}/{name}", read.New(store, logger))
	r.Post("/update/{type}/{name}/{value}", update.New(store))

	r.Post("/value/", read.NewJSON(store))
	r.Post("/update/", update.NewJSON(store, logger))
	r.Post("/updates/", update.NewJSONBatch(store, logger))

	r.Get("/ping", ping.New(db))

	return r
}
