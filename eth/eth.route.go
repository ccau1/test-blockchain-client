package eth

import (
	"github.com/go-chi/chi/v5"
)

func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/block-number/{id}", c_getBlockNumber)

	r.Get("/block-by-number/{id}", c_getBlockByNumber)

	return r
}
