package eth

import (
	"strconv"
	"net/http"
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/go-chi/chi/v5"

	"github.com/ccau1/test-blockchain-client/utils"
	"github.com/ccau1/test-blockchain-client/providers_handler"
)

type GetNextProviderFilter = providers_handler.GetNextProviderFilter

var DEFAULT_JSON_VERSION = "2.0"

var providersHandler = &providers_handler.ProvidersHandler{}

func c_getBlockNumber(rw http.ResponseWriter, r *http.Request) {
	// get query fields
	jsonrpc := r.URL.Query().Get("jsonrpc")
	if jsonrpc == "" {
		jsonrpc = DEFAULT_JSON_VERSION
	}
	batchIdStr := r.URL.Query().Get("batchId")
	batchId := 1
	if batchIdStr != "" {
		batchIdParsed, err := strconv.Atoi(batchIdStr)
		if err != nil {
			render.Render(rw, r, utils.ErrInvalidRequest("Cannot parse param [id]: " + batchIdStr))
			return
		}
		batchId = batchIdParsed
	}
	// Encode the data
	body := &GetBlockNumberBody{
		JSONRPC: jsonrpc,
		Method: "eth_blockNumber",
		ID: batchId,
	}
	requestBodyByte, _ := json.Marshal(body)
	// get provider to handle this call
	provider := providersHandler.GetNextProvider(GetNextProviderFilter{ ChainType: "eth" })
	// fetch result from polygonRPC
	result, err := (*provider).Call("eth", requestBodyByte)
	if err != nil {
		render.Render(rw, r, utils.ErrServerError("PolygonRPC request error: " + err.Error()))
		return
	}
	// return result
	rw.Write(result)
}

func c_getBlockByNumber(rw http.ResponseWriter, r *http.Request) {
	// get param fields
	paramId := chi.URLParam(r, "blockNumber")
	blockNumberId, err := strconv.Atoi(paramId)
	if err != nil {
		render.Render(rw, r, utils.ErrInvalidRequest("Cannot parse param [blockNumber]: " + paramId))
		return
	}
	// get query fields
	jsonrpc := r.URL.Query().Get("jsonrpc")
	if jsonrpc == "" {
		jsonrpc = DEFAULT_JSON_VERSION
	}
	batchIdStr := r.URL.Query().Get("batchId")
	batchId := 1
	if batchIdStr != "" {
		batchIdParsed, err := strconv.Atoi(batchIdStr)
		if err != nil {
			render.Render(rw, r, utils.ErrInvalidRequest("Cannot parse param [id]: " + batchIdStr))
			return
		}
		batchId = batchIdParsed
	}
	// Encode the data
	body := &GetBlockByNumberBody{
		JSONRPC: jsonrpc,
		Method: "eth_getBlockByNumber",
		ID: batchId,
		Params: []interface{}{
			blockNumberId,
			true,
		},
	}
	requestBodyByte, _ := json.Marshal(body)
	// get provider to handle this call
	provider := providersHandler.GetNextProvider(GetNextProviderFilter{ ChainType: "eth" })
	// fetch result from provider call
	result, err := (*provider).Call("eth", requestBodyByte)
	if err != nil {
		render.Render(rw, r, utils.ErrServerError("PolygonRPC request error: " + err.Error()))
		return
	}
	// return result
	rw.Write(result)
}