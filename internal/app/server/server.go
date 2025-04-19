package server

import (
	"github.com/mkolibaba/metrics/internal/http/router"
	"github.com/mkolibaba/metrics/internal/storage/inmemory"
	"log"
	"net/http"
)

func Run() {
	store := inmemory.NewMemStorage()

	r := router.New(store)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
