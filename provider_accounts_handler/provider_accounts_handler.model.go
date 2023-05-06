package provider_accounts_handler

import (
	"github.com/ccau1/test-blockchain-client/provider_accounts_handler/provider_account"
	ProviderAccountsStrategyTypes "github.com/ccau1/test-blockchain-client/provider_accounts_handler/provider_accounts_strategy"
	"github.com/ccau1/test-blockchain-client/utils"
)

var Log = utils.Log

type ProviderAccount = provider_account.ProviderAccount
type IProviderAccountsStrategy = ProviderAccountsStrategyTypes.IProviderAccountsStrategy

type GetNextAccountFilter struct {
	
}
