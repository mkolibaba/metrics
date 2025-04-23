package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/mkolibaba/metrics/internal/server/http/handlers/list"
	"github.com/mkolibaba/metrics/internal/server/http/handlers/read"
	"github.com/mkolibaba/metrics/internal/server/http/handlers/update"
)

type MetricsStorage interface {
	list.AllMetricsGetter
	read.MetricsGetter
	update.MetricsUpdater
}

func New(store MetricsStorage) chi.Router {
	r := chi.NewRouter()

	r.Get("/", list.New(store))
	r.Get("/value/{type}/{name}", read.New(store))
	r.Post("/update/{type}/{name}/{value}", update.New(store))

	return r
}
