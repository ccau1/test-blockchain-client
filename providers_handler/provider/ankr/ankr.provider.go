package ankr

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ccau1/test-blockchain-client/provider_accounts_handler"
	ProviderAccountsStrategyTypes "github.com/ccau1/test-blockchain-client/provider_accounts_handler/provider_accounts_strategy"
	"github.com/ccau1/test-blockchain-client/utils/helper"
)

var providerAccountsHandler *ProviderAccountsHandler = &ProviderAccountsHandler{
	// chain accounts handler will fetch accounts based on provider name
	Provider: "ankr",
	// the strategy to use for deciding which account to use for the 
	// coming request
	UseStrategy: &ProviderAccountsStrategyTypes.RequestLimitStrategy{
		// only allow 30 requests
		LimitAmount: 30,
		// every 30 second (30000ms)
		LimitPerInterval: 30000,
	},
}

type AnkrProvider struct {
	
}

func (x *AnkrProvider) SupportedChains() []string {
	return []string {
		"eth",
		"bsc",
	}
}

/*
	health check ping method. Returns error if ping doesn't pass
*/
func (x *AnkrProvider) Ping() error {
	_, err := x.GetLatestBlockNumber("eth")

	return err
}

/*
	get chain type's latest block number
*/
func (x *AnkrProvider) GetLatestBlockNumber(chainType string) (string, error) {
	// Encode the data
	body := &AnkrGetBlockNumberBody{
		JSONRPC: "2.0",
		Method: "eth_blockNumber",
		ID: 1,
	}
	requestBodyByte, _ := json.Marshal(body)
	result, err := callAnkrProvider[string](chainType, requestBodyByte)

	intNum, err := helper.HexToInt(result)

	return fmt.Sprint(intNum), err
}

/*
	get chain type's block by block number
*/
func (x *AnkrProvider) GetByBlockNumber(chainType string, blockNumber string) (ChainBlock, error) {
	// Encode the data
	body := &AnkrGetBlockNumberBody{
		JSONRPC: "2.0",
		Method: "eth_getBlockByNumber",
		ID: 1,
		Params: []interface{}{
			blockNumber,
			true,
		},
	}
	requestBodyByte, _ := json.Marshal(body)
	result, err := callAnkrProvider[ChainBlock](chainType, requestBodyByte)

	if err != nil {
		return ChainBlock{}, err
	}

	// Log.Infof("[GetByBlockNumber] result: %+v", result)

	return result, err
}

func generateUrl(chainType string) (string, error) {
	// get chain account to use
	chainAccount, err := providerAccountsHandler.GetNextAccount(&provider_accounts_handler.GetNextAccountFilter{

	})
	if (err != nil) {
		return "", err
	}
	// format and return url for chain based on:
	// - chain type
	// - chain account id
	return fmt.Sprintf("https://rpc.ankr.com/%s/%s", chainType, chainAccount.ID), nil
}

func callAnkrProvider[Result any](chainType string, body []byte) (Result, error) {
	// call provider to retrieve info
	providerDomain, err := generateUrl(chainType)
	if err != nil {
		return *new(Result), err
	}

	Log.Infof("url: %s", providerDomain)

	res, err := http.Post(
		providerDomain,					// url
		"application/json",			// content-type
		bytes.NewBuffer(body), 	// body (as buffer)
	)
	if err != nil {
		return *new(Result), err
	}

	// get content from response
	defer res.Body.Close()
	resContent, err := ioutil.ReadAll(res.Body)

	if (err != nil) {
		// TODO: set status to down, but also need other handling to allow it be up again
	}

	// Log.Infof("resContent raw: %+v", string(resContent))

	var callResponse AnkrCallResponse[Result]
	err = json.Unmarshal(resContent, &callResponse)

	if err != nil {
		Log.Infof("err: %+v", err)
		return *new(Result), err
	}
	
	// Log.Infof("callResponse: %+v", callResponse)

	if (callResponse.Error != nil) {
		return *new(Result), errors.New(fmt.Sprintf("response error: [%d] %s", callResponse.Error.Code, callResponse.Error.Message))
	}

	return callResponse.Result, nil
}
