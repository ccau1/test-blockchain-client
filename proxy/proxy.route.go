package proxy

import (
	"github.com/go-chi/chi/v5"
)

func Router() chi.Router {
	r := chi.NewRouter()

	r.Get("/{chainType}/block-number", c_getBlockNumber)

	r.Get("/{chainType}/block-by-number/{blockNumber}", c_getBlockByNumber)

	return r
}
