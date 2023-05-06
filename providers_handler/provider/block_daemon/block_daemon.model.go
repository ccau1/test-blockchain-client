package provider

import (
	"github.com/ccau1/test-blockchain-client/provider_accounts_handler"
	"github.com/ccau1/test-blockchain-client/providers_handler/provider"
	"github.com/ccau1/test-blockchain-client/utils"
)

type ProviderAccountsHandler = provider_accounts_handler.ProviderAccountsHandler

type ChainBlock = provider.ChainBlock

var Log = utils.Log

type BlockDaemonGetBlockNumberBody struct {
	JSONRPC    	string `json:"jsonrpc"`
	Method 			string `json:"method"`
	ID 					int `json:"id"`
	Params			[]interface{} `json:"params"`
}

type BlockDaemonCallResponse[T any] struct {
	JSONRPC    	string `json:"jsonrpc"`
	ID 					int `json:"id"`
	Result    	T `json:"result"`
	Error				*BlockDaemonCallResponseError `json:"error"`
}

type BlockDaemonCallResponseError struct {
	Code				int `json:"code"`
	Message			string `json:"message"`
}

type BlockDaemonBlockMeta struct {
	Index int `json:"index"`
	To string `json:"to"`
}

type BlockDaemonBlockEventMeta struct {
	BaseFee int `json:"base_fee"`
	FeeBurned int `json:"fee_burned"`
	GasLimit int `json:"gas_limit"`
	GasPrice int `json:"gas_price"`
	GasUsed int `json:"gas_used"`
}

type BlockDaemonBlockEvent struct {
	ID string `json:"id"`
	TransactionID string `json:"transaction_id"`
	Type string `json:"type"`
	Denomination string `json:"denomination"`
	Source string `json:"source"`
	Meta BlockDaemonBlockEventMeta `json:"meta"`
	Date int `json:"date"`
	Amount float64 `json:"amount"`
	Decimals int `json:"decimals"`
}

type BlockDaemonBlockTransaction struct {
	ID string `json:"id"`
	BlockID string `json:"block_id"`
	Date int `json:"date"`
	Status string `json:"status"`
	NumEvents int `json:"num_events"`
	Meta BlockDaemonBlockMeta `json:"meta"`
	BlockNumber int `json:"block_number"`
	Events []BlockDaemonBlockEvent `json:"events"`
}

type BlockDaemonBlock struct {
	Number int `json:"number"`
	ID string `json:"id"`
	ParentID string `json:"parent_id"`
	Date int `json:"date"`
	Num_Transactions int `json:"num_txs"`
	Transactions []BlockDaemonBlockTransaction `json:"txs"`
}