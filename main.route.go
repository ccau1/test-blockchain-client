package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/ccau1/test-blockchain-client/eth"
)

func GenRoute() *chi.Mux {
	// init router
	r := chi.NewRouter()

	// mount routers
	r.Mount("/eth", eth.Router())

	return r
}