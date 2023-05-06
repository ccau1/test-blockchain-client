package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/ccau1/test-blockchain-client/provider_accounts_handler"
	ProviderAccountsStrategyTypes "github.com/ccau1/test-blockchain-client/provider_accounts_handler/provider_accounts_strategy"
	"github.com/ccau1/test-blockchain-client/utils"
)

var chainAliasMap = map[string]string{
	"eth": "ethereum",
	"btc": "bitcoin",
}

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

var DEFAULT_NETWORK = "mainnet"

var blockDaemonProviderAccountsHandler *ProviderAccountsHandler = &ProviderAccountsHandler{
	// chain accounts handler will fetch accounts based on provider name
	Provider: "block_daemon",
	// the strategy to use for deciding which account to use for the 
	// coming request
	UseStrategy: &ProviderAccountsStrategyTypes.RequestLimitStrategy{
		// only allow 25 requests
		LimitAmount: 25,
		// every 1 second (1000ms)
		LimitPerInterval: 1000,
	},
}

type BlockDaemonProvider struct {
	
}

func (x *BlockDaemonProvider) SupportedChains() []string {
	return []string {
		"eth",
		"btc",
	}
}

func (x *BlockDaemonProvider) GetLatestBlockNumber(chainType string) (string, error) {
	result, err := callBlockDaemon[int](chainType, "sync/block_number")

	if (err != nil) {
		return "", err
	}

	utils.Log.Infof("result: %d", result)

	return fmt.Sprint(result), nil
}

func (x *BlockDaemonProvider) GetByBlockNumber(chainType string, blockNumber string) (ChainBlock, error) {
	result, err := callBlockDaemon[BlockDaemonBlock](chainType, "block/" + blockNumber)

	if (err != nil) {
		return *new(ChainBlock), err
	}

	utils.Log.Infof("result: %+v", result)

	resultChainBlock := ChainBlock{
		Number: intToHex(result.Number),
		ParentHash: result.ParentID,
		GasLimit: intToHex(result.Transactions[0].Events[0].Meta.GasLimit),
		GasUsed: intToHex(result.Transactions[0].Events[0].Meta.GasUsed),
	}

	return resultChainBlock, nil
}

func intToHex(i int) string {
	return "0x" + strconv.FormatInt(int64(i), 16)
}

func callBlockDaemon[Result any](chainType string, method string) (Result, error) {
	// get chain account to use
	chainAccount, err := blockDaemonProviderAccountsHandler.GetNextAccount(&provider_accounts_handler.GetNextAccountFilter{
	})
	if (err != nil) {
		return *new(Result), err
	}
	Log.Infof("selected chainAccount: %+v", chainAccount)

	if (chainAliasMap[chainType] != "") {
		chainType = chainAliasMap[chainType]
	}

	url := fmt.Sprintf("https://svc.blockdaemon.com/universal/v1/%s/%s/%s", chainType, DEFAULT_NETWORK, method)

	Log.Infof("url: %s", url)

	req, err := http.NewRequest("GET", url, nil)

	if (err != nil) {
		return *new(Result), err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("X-API-Key", chainAccount.ID)

	res, err := http.DefaultClient.Do(req)

	if (err != nil) {
		return *new(Result), err
	}

	defer res.Body.Close()
	resContent, err := ioutil.ReadAll(res.Body)

	if (err != nil) {
		return *new(Result), err
	}

	Log.Infof("resContent: %+v\n", string(resContent))

	var callResponse Result
	err = json.Unmarshal(resContent, &callResponse)

	if (err != nil) {
		return *new(Result), err
	}

	return callResponse, err
}
