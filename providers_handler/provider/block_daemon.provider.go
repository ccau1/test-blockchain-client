package provider

import (
	"github.com/ccau1/test-blockchain-client/provider_accounts_handler"
	ProviderAccountsStrategyTypes "github.com/ccau1/test-blockchain-client/provider_accounts_handler/provider_accounts_strategy"
)

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

type BlockDaemonProviderAccountsHandler = provider_accounts_handler.ProviderAccountsHandler

var blockDaemonProviderAccountsHandler *BlockDaemonProviderAccountsHandler = &BlockDaemonProviderAccountsHandler{
	// chain accounts handler will fetch accounts based on provider name
	Provider: "block_daemon",
	// the strategy to use for deciding which account to use for the 
	// coming request
	UseStrategy: &ProviderAccountsStrategyTypes.RequestLimitStrategy{
		// only allow 30 requests
		LimitAmount: 30,
		// every 30 second (30000ms)
		LimitPerInterval: 30000,
	},
}

type BlockDaemonProvider struct {
	
}

func (x *BlockDaemonProvider) SupportedChains() []string {
	return []string { }
}

func (x *BlockDaemonProvider) GetLatestBlockNumber(chainType string) (string, error) {
	return "", nil
}

func (x *BlockDaemonProvider) GetByBlockNumber(chainType string, blockNumber string) (ChainBlock, error) {
	return ChainBlock{}, nil
}
