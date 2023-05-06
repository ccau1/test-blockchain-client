package provider

import (
	"github.com/ccau1/test-blockchain-client/provider_accounts_handler"
)

type ProviderAccountsHandler = provider_accounts_handler.ProviderAccountsHandler

type ChainBlockTransaction struct {
	BlockHash string `json:"blockHash"`
	BlockNumber string `json:"blockNumber"`
	From string `json:"from"`
	Gas string `json:"gas"`
	GasPrice string `json:"gasPrice"`
	Hash string `json:"hash"`
	Input string `json:"input"`
	Nonce string `json:"nonce"`
	To string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Value string `json:"value"`
	Type string `json:"type"`
	ChainId string `json:"chainId"`
	V string `json:"v"`
	R string `json:"r"`
	S string `json:"s"`
}

type ChainBlock struct {
	Difficulty string `json:"difficulty"`
	ExtraData string `json:"extraData"`
	GasLimit string `json:"gasLimit"`
	GasUsed string `json:"gasUsed"`
	Hash string `json:"hash"`
	LogsBloom string `json:"logsBloom"`
	Miner string `json:"miner"`
	MixHash string `json:"mixHash"`
	Nonce string `json:"nonce"`
	Number string `json:"number"`
	ParentHash string `json:"parentHash"`
	ReceiptsRoot string `json:"receiptsRoot"`
	Sha3Uncles string `json:"sha3Uncles"`
	Size string `json:"size"`
	StateRoot string `json:"stateRoot"`
	Timestamp string `json:"timestamp"`
	TotalDifficulty string `json:"totalDifficulty"`
	Transactions []ChainBlockTransaction `json:"transactions"`
	TransactionsRoot string `json:"transactionsRoot"`
	Uncles []string `json:"uncles"`
}

type IProvider interface {
	// define a list of chains this provider supports
	SupportedChains() []string
	// simple ping function to test connection
	Ping() error
	// chain action methods for fetching/updating which every provider should provide
	GetLatestBlockNumber(chainType string) (string, error)
	GetByBlockNumber(chainType string, blockNumber string) (ChainBlock, error)
}