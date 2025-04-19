package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/mkolibaba/metrics/internal/http/handlers/list"
	"github.com/mkolibaba/metrics/internal/http/handlers/read"
	"github.com/mkolibaba/metrics/internal/http/handlers/update"
	"github.com/mkolibaba/metrics/internal/storage"
)

func New(store storage.MetricsStorage) chi.Router {
	r := chi.NewRouter()

	r.Get("/", list.New(store))
	r.Get("/value/{type}/{name}", read.New(store))
	r.Post("/update/{type}/{name}/{value}", update.New(store))

	return r
}
