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

type GetBlockNumberBody struct {
	JSONRPC    	string `json:"jsonrpc"`
	Method 			string `json:"method"`
	ID 					int `json:"id"`
	Params			[]interface{} `json:"params"`
}

type CallResponse struct {
	JSONRPC    	string `json:"jsonrpc"`
	ID 					int `json:"id"`
	Result    	interface{} `json:"result"`
	Error				*CallResponseError `json:"error"`
}

type CallResponseError struct {
	Code				int `json:"code"`
	Message			string `json:"message"`
}

type AnkrProvider struct {
	
}

type ProviderAccountsHandler = provider_accounts_handler.ProviderAccountsHandler

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

func (x *AnkrProvider) generateUrl(chainType string) (string, error) {
	// get chain account to use
	chainAccount, err := providerAccountsHandler.GetNextAccount()
	if (err != nil) {
		return "", err
	}
	// format and return url for chain based on:
	// - chain type
	// - chain account id
	return fmt.Sprintf("https://rpc.ankr.com/%s/%s", chainType, chainAccount.ID), nil
}

func (x *AnkrProvider) SupportedChains() []string {
	return []string {
		"eth",
		"bsc",
	}
}

func (x *AnkrProvider) GetLatestBlockNumber(chainType string) (string, error) {
	// Encode the data
	body := &GetBlockNumberBody{
		JSONRPC: "2.0",
		Method: "eth_blockNumber",
		ID: 1,
	}
	requestBodyByte, _ := json.Marshal(body)
	result, err := x.call(chainType, requestBodyByte)

	intNum, err := strconv.ParseInt(strings.ReplaceAll(result.(string), "0x", ""), 16, 64)

	return fmt.Sprint(intNum), err
}

func (x *AnkrProvider) GetByBlockNumber(chainType string, blockNumber string) (interface{}, error) {
	// Encode the data
	body := &GetBlockNumberBody{
		JSONRPC: "2.0",
		Method: "eth_getBlockByNumber",
		ID: 1,
		Params: []interface{}{
			blockNumber,
			true,
		},
	}
	requestBodyByte, _ := json.Marshal(body)
	result, err := x.call(chainType, requestBodyByte)

	if err != nil {
		return BlockNumber{}, err
	}

	utils.Log.Infof("[GetByBlockNumber] result: %+v", result)

	return result, err
}

func (x *AnkrProvider) call(chainType string, body []byte) (interface{}, error) {
	// call provider to retrieve info
	providerDomain, err := x.generateUrl(chainType)
	if err != nil {
		return nil, err
	}

	utils.Log.Infof("[call] url: %s", providerDomain)

	res, err := http.Post(
		providerDomain,					// url
		"application/json",			// content-type
		bytes.NewBuffer(body), 	// body (as buffer)
	)
	if err != nil {
		return nil, err
	}

	// get content from response
	defer res.Body.Close()
	resContent, err := ioutil.ReadAll(res.Body)

	utils.Log.Infof("[call] resContent raw: %+v", string(resContent))

	var callResponse CallResponse
	err = json.Unmarshal(resContent, &callResponse)

	if err != nil {
		utils.Log.Infof("[call] err: %+v", err)
		return nil, err
	}
	
	utils.Log.Infof("[call] callResponse: %+v", callResponse)

	if (callResponse.Error != nil) {
		return nil, errors.New(fmt.Sprintf("[call] response error: [%d] %s", callResponse.Error.Code, callResponse.Error.Message))
	}

	if callResponse.Result == nil {
		return nil, errors.New("[call] request result is nil")
	}

	return callResponse.Result, nil
}