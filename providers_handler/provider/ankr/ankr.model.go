package ankr

import (
	"github.com/ccau1/test-blockchain-client/provider_accounts_handler"
	"github.com/ccau1/test-blockchain-client/providers_handler/provider"
	"github.com/ccau1/test-blockchain-client/utils"
)

type ProviderAccountsHandler = provider_accounts_handler.ProviderAccountsHandler

type ChainBlock = provider.ChainBlock

var Log = utils.Log

type AnkrGetBlockNumberBody struct {
	JSONRPC    	string `json:"jsonrpc"`
	Method 			string `json:"method"`
	ID 					int `json:"id"`
	Params			[]interface{} `json:"params"`
}

type AnkrCallResponse[T any] struct {
	JSONRPC    	string `json:"jsonrpc"`
	ID 					int `json:"id"`
	Result    	T `json:"result"`
	Error				*AnkrCallResponseError `json:"error"`
}

type AnkrCallResponseError struct {
	Code				int `json:"code"`
	Message			string `json:"message"`
}
