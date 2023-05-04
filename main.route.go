package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/ccau1/test-blockchain-client/proxy"
)

func GenRoute() *chi.Mux {
	// init router
	r := chi.NewRouter()

	// mount routers
	r.Mount("/", proxy.Router())

	return r
}