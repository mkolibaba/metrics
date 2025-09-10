package router

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/mkolibaba/metrics/internal/common/rsa"
	"github.com/mkolibaba/metrics/internal/server/http/handlers/list"
	"github.com/mkolibaba/metrics/internal/server/http/handlers/ping"
	"github.com/mkolibaba/metrics/internal/server/http/handlers/read"
	"github.com/mkolibaba/metrics/internal/server/http/handlers/update"
	"github.com/mkolibaba/metrics/internal/server/http/middleware"
	"go.uber.org/zap"
	"net"
)

type MetricsStorage interface {
	list.AllMetricsGetter
	read.MetricsGetter
	update.MetricsUpdater
}

// New создает и настраивает новый chi.Router.
// Роутер включает в себя middleware для логирования, подписи, сжатия
// и шифровки, а также регистрирует обработчики для всех эндпоинтов сервера метрик.
func New(
	store MetricsStorage,
	db *sql.DB,
	hashKey string,
	subnet *net.IPNet,
	logger *zap.SugaredLogger,
	decryptor rsa.Decryptor,
) chi.Router {
	r := chi.NewRouter()

	// middleware
	r.Use(middleware.Logger(logger))
	if subnet != nil {
		r.Use(middleware.Subnet(subnet))
	}
	r.Use(middleware.Decryptor(decryptor, logger))
	if hashKey != "" {
		r.Use(middleware.Hash(hashKey, logger))
	}
	r.Use(middleware.Compressor(logger))
	jsonContentTypeMiddleware := middleware.ContentType("application/json")

	// read
	r.Get("/", list.New(store, logger))
	r.Get("/value/{type}/{name}", read.New(store, logger))
	r.With(jsonContentTypeMiddleware).Post("/value/", read.NewJSON(store))

	// update
	updateAPI := update.NewAPI(store, logger)
	r.Post("/update/{type}/{name}/{value}", updateAPI.HandlePlain)
	r.With(jsonContentTypeMiddleware).Post("/update/", updateAPI.HandleJSON)
	r.With(jsonContentTypeMiddleware).Post("/updates/", updateAPI.HandleJSONBatch)

	// other
	r.Get("/ping", ping.New(db))

	return r
}
