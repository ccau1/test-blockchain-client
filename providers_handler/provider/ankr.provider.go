package provider

import (
	"fmt"
	"bytes"
	"errors"
	"strings"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"

	"github.com/ccau1/test-blockchain-client/utils"
	"github.com/ccau1/test-blockchain-client/provider_accounts_handler"
	ProviderAccountsStrategyTypes "github.com/ccau1/test-blockchain-client/provider_accounts_handler/provider_accounts_strategy"
)

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

var ankrProviderAccountsHandler *ProviderAccountsHandler = &ProviderAccountsHandler{
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

func (x *AnkrProvider) GetLatestBlockNumber(chainType string) (string, error) {
	// Encode the data
	body := &AnkrGetBlockNumberBody{
		JSONRPC: "2.0",
		Method: "eth_blockNumber",
		ID: 1,
	}
	requestBodyByte, _ := json.Marshal(body)
	result, err := callAnkrProvider[string](chainType, requestBodyByte)

	intNum, err := strconv.ParseInt(strings.ReplaceAll(result, "0x", ""), 16, 64)

	return fmt.Sprint(intNum), err
}

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

	Log.Infof("[GetByBlockNumber] result: %+v", result)

	return result, err
}

func generateUrl(chainType string) (string, error) {
	// get chain account to use
	chainAccount, err := ankrProviderAccountsHandler.GetNextAccount(&provider_accounts_handler.GetNextAccountFilter{

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

	Log.Infof("resContent raw: %+v", string(resContent))

	var callResponse AnkrCallResponse[Result]
	err = json.Unmarshal(resContent, &callResponse)

	if err != nil {
		Log.Infof("err: %+v", err)
		return *new(Result), err
	}
	
	Log.Infof("callResponse: %+v", callResponse)

	if (callResponse.Error != nil) {
		return *new(Result), errors.New(fmt.Sprintf("response error: [%d] %s", callResponse.Error.Code, callResponse.Error.Message))
	}

	return callResponse.Result, nil
}
