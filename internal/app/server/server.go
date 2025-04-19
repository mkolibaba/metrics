package server

import (
	"github.com/mkolibaba/metrics/internal/config"
	"github.com/mkolibaba/metrics/internal/http/router"
	"github.com/mkolibaba/metrics/internal/storage/inmemory"
	"log"
	"net/http"
)

func Run() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	serverAddress := config.LoadServerConfig().ServerAddress

	store := inmemory.NewMemStorage()
	r := router.New(store)

	log.Printf("Running server on %s", serverAddress)
	if err := http.ListenAndServe(serverAddress, r); err != nil {
		log.Fatal(err)
	}
}
