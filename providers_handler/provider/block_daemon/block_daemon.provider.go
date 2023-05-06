package provider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ccau1/test-blockchain-client/provider_accounts_handler"
	ProviderAccountsStrategyTypes "github.com/ccau1/test-blockchain-client/provider_accounts_handler/provider_accounts_strategy"
	"github.com/ccau1/test-blockchain-client/utils"
	"github.com/ccau1/test-blockchain-client/utils/helper"
)

var DEFAULT_NETWORK = "mainnet"

var providerAccountsHandler *ProviderAccountsHandler = &ProviderAccountsHandler{
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

var chainAliasMap = map[string]string{
	"eth": "ethereum",
	"btc": "bitcoin",
}

type BlockDaemonProvider struct {
	
}

func (x *BlockDaemonProvider) SupportedChains() []string {
	return []string {
		"eth",
		"btc",
	}
}

func (x *BlockDaemonProvider) Ping() error {
	_, err := x.GetLatestBlockNumber("eth")

	return err
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

	// utils.Log.Infof("result: %+v", result)

	resultChainBlock := ChainBlock{
		Number: helper.IntToHex(result.Number),
		ParentHash: result.ParentID,
		GasLimit: helper.IntToHex(result.Transactions[0].Events[0].Meta.GasLimit),
		GasUsed: helper.IntToHex(result.Transactions[0].Events[0].Meta.GasUsed),
	}

	return resultChainBlock, nil
}

func callBlockDaemon[Result any](chainType string, method string) (Result, error) {
	// get chain account to use
	chainAccount, err := providerAccountsHandler.GetNextAccount(&provider_accounts_handler.GetNextAccountFilter{
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

	// Log.Infof("resContent: %+v\n", string(resContent))

	var callResponse Result
	err = json.Unmarshal(resContent, &callResponse)

	if (err != nil) {
		return *new(Result), err
	}

	return callResponse, err
}
