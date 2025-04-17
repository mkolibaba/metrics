package main

import (
	"github.com/mkolibaba/metrics/internal/http/handlers"
	"github.com/mkolibaba/metrics/internal/storage/inmemory"
	"log"
	"net/http"
)

func main() {
	store := inmemory.NewMemStorage()

	mux := http.NewServeMux()
	mux.HandleFunc(handlers.RouteUpdate, handlers.NewUpdateHandler(store))

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
