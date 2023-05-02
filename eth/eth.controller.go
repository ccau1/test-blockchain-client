package eth

import (
	"strconv"
	"net/http"
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/go-chi/chi/v5"

	"github.com/ccau1/test-blockchain-client/chain_base"
	"github.com/ccau1/test-blockchain-client/utils"
)

var DEFAULT_JSON_VERSION = "2.0"

var CallPolygonRPC = chain_base.CallPolygonRPCFactory("eth")

func c_getBlockNumber(rw http.ResponseWriter, r *http.Request) {
	// get param fields
	paramId := chi.URLParam(r, "id")
	blockNumberId, err := strconv.Atoi(paramId)
	if err != nil {
		render.Render(rw, r, utils.ErrInvalidRequest("Cannot parse param [id]: " + paramId))
		return
	}
	// get query fields
	jsonrpc := r.URL.Query().Get("jsonrpc")
	if jsonrpc == "" {
		jsonrpc = DEFAULT_JSON_VERSION
	}
	// Encode the data
	body := &GetBlockNumberBody{
		JSONRPC: jsonrpc,
		Method: "eth_blockNumber",
		ID: blockNumberId,
	}
	requestBodyByte, _ := json.Marshal(body)
	// fetch result from polygonRPC
	result, err := CallPolygonRPC(requestBodyByte)
	if err != nil {
		render.Render(rw, r, utils.ErrServerError("PolygonRPC request error: " + err.Error()))
		return
	}
	// return result
	rw.Write(result)
}

func c_getBlockByNumber(rw http.ResponseWriter, r *http.Request) {
	// get param fields
	paramId := chi.URLParam(r, "id")
	blockNumberId, err := strconv.Atoi(paramId)
	if err != nil {
		render.Render(rw, r, utils.ErrInvalidRequest("Cannot parse param [id]: " + paramId))
		return
	}
	// get query fields
	jsonrpc := r.URL.Query().Get("jsonrpc")
	if jsonrpc == "" {
		jsonrpc = DEFAULT_JSON_VERSION
	}
	// Encode the data
	body := &GetBlockByNumberBody{
		JSONRPC: jsonrpc,
		Method: "eth_getBlockByNumber",
		ID: blockNumberId,
		Params: []interface{}{
			"0x134e82a",
			true,
		},
	}
	requestBodyByte, _ := json.Marshal(body)
	// fetch result from polygonRPC
	result, err := CallPolygonRPC(requestBodyByte)
	if err != nil {
		render.Render(rw, r, utils.ErrServerError("PolygonRPC request error: " + err.Error()))
		return
	}
	// return result
	rw.Write(result)
}