package provider_accounts_strategy

import (
	"github.com/ccau1/test-blockchain-client/provider_accounts_handler/provider_account"
)

type ProviderAccount = provider_account.ProviderAccount

type IProviderAccountsStrategy interface {
	Load(chainAccountList *[]ProviderAccount)
	GetNextAccount() (*ProviderAccount, error)
}
