package chain_accounts_strategy

import (
	"github.com/ccau1/test-blockchain-client/chain_base/chain_account"
)

type ChainAccount = chain_account.ChainAccount

type IChainSelectorStrategy interface {
	Load(chainAccountList *[]ChainAccount)
	GetNextAccount() (*ChainAccount, error)
}
