package chain_base

import (
	"fmt"
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/ccau1/test-blockchain-client/chain_base/chain_account_handler"
)

var chainAccountsHandler *chain_account_handler.ChainAccountsHandler = &chain_account_handler.ChainAccountsHandler{}

func getChainUrl(chainType string) (string, error) {
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

func CallPolygonRPC(chainType string, body []byte) ([]byte, error) {
	// call provider to retrieve info
	providerDomain, err := getChainUrl(chainType)
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

func CallPolygonRPCFactory(chainType string) (func([]byte) ([]byte, error)) {
	// return a function that only needs to take in body parameter
	return func (body []byte) ([]byte, error) {
		return CallPolygonRPC(chainType, body)
	}
}