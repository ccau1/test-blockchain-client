package eth

import (
	"github.com/go-chi/chi/v5"
)

func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/block-number", c_getBlockNumber)

	r.Get("/block-by-number/{blockNumber}", c_getBlockByNumber)

	return r
}
