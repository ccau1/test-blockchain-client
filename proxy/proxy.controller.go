package proxy

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/ccau1/test-blockchain-client/providers_handler"
	"github.com/ccau1/test-blockchain-client/utils"
)

type GetNextProviderFilter = providers_handler.GetNextProviderFilter

var DEFAULT_JSON_VERSION = "2.0"

var providersHandler = &providers_handler.ProvidersHandler{}

func c_getBlockNumber(rw http.ResponseWriter, r *http.Request) {
	// get param fields
	chainType := chi.URLParam(r, "chainType")
	// get provider to handle this call
	provider, err := providersHandler.GetNextProvider(GetNextProviderFilter{ ChainType: chainType })
	if err != nil {
		render.Render(rw, r, utils.ErrServerError("PolygonRPC request error: " + err.Error()))
		return
	}
	// fetch result from polygonRPC
	result, err := (*provider).GetLatestBlockNumber(chainType)
	if err != nil {
		render.Render(rw, r, utils.ErrServerError("PolygonRPC request error: " + err.Error()))
		return
	}
	// return result
	rw.Write([]byte(result))
}

func c_getBlockByNumber(rw http.ResponseWriter, r *http.Request) {
	// get param fields
	chainType := chi.URLParam(r, "chainType")
	blockNumber := chi.URLParam(r, "blockNumber")
	// blockNumberId, err := strconv.Atoi(paramId)
	// if err != nil {
	// 	render.Render(rw, r, utils.ErrInvalidRequest("Cannot parse param [blockNumber]: " + paramId))
	// 	return
	// }
	// get provider to handle this call
	provider, err := providersHandler.GetNextProvider(GetNextProviderFilter{ ChainType: chainType })
	if err != nil {
		render.Render(rw, r, utils.ErrServerError("PolygonRPC request error: " + err.Error()))
		return
	}
	// fetch result from provider call
	result, err := (*provider).GetByBlockNumber(chainType, blockNumber)
	if err != nil {
		render.Render(rw, r, utils.ErrServerError("PolygonRPC request error: " + err.Error()))
		return
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		render.Render(rw, r, utils.ErrServerError("PolygonRPC request result error: " + err.Error()))
		return
	}

	// return result
	rw.Write(resultBytes)
}