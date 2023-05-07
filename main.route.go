package main

import (
	"net/http"

	"github.com/ccau1/test-blockchain-client/proxy"
	"github.com/go-chi/chi/v5"
)

func GenRoute() *chi.Mux {
	// init router
	r := chi.NewRouter()

	// mount routers
	r.Mount("/", proxy.Router())

	// define default health-check call
	r.Get("/health", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("healthy"))
	})

	return r
}