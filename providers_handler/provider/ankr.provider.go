package provider

import (
	"fmt"
	"bytes"
	"net/http"
	"io/ioutil"

	"github.com/ccau1/test-blockchain-client/provider_accounts_handler"
	ProviderAccountsStrategyTypes "github.com/ccau1/test-blockchain-client/provider_accounts_handler/provider_accounts_strategy"
)

type AnkrProvider struct {
	
}

type ProviderAccountsHandler = provider_accounts_handler.ProviderAccountsHandler

var chainAccountsHandler *ProviderAccountsHandler = &ProviderAccountsHandler{
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
	chainAccount, err := chainAccountsHandler.GetNextAccount()
	if (err != nil) {
		return "", err
	}
	// format and return url for chain based on:
	// - chain type
	// - chain account id
	return fmt.Sprintf("https://rpc.ankr.com/%s/%s", chainType, chainAccount.ID), nil
}


func (x *AnkrProvider) Call(chainType string, body []byte) ([]byte, error) {
	// call provider to retrieve info
	providerDomain, err := x.generateUrl(chainType)
	if err != nil {
		return nil, err
	}
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
	return resContent, nil
}

func (x *AnkrProvider) CallFactory(chainType string) (func([]byte) ([]byte, error)) {
	// return a function that only needs to take in body parameter
	return func (body []byte) ([]byte, error) {
		return x.Call(chainType, body)
	}
}