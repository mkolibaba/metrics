package app

import (
	"github.com/mkolibaba/metrics/internal/server/config"
	"github.com/mkolibaba/metrics/internal/server/http/router"
	"github.com/mkolibaba/metrics/internal/server/storage/inmemory"
	"log"
	"net/http"
)

func Run() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	serverAddress := config.MustLoadServerConfig().ServerAddress

	store := inmemory.NewMemStorage()
	r := router.New(store)

	log.Printf("Running server on %s", serverAddress)
	if err := http.ListenAndServe(serverAddress, r); err != nil {
		log.Fatal(err)
	}
}
