package chain_account_handler

import (
	"github.com/ccau1/test-blockchain-client/chain_base/chain_account"
	ChainAccountsStrategyTypes "github.com/ccau1/test-blockchain-client/chain_base/chain_accounts_strategy"
)

type ChainAccount = chain_account.ChainAccount


type ChainAccountsHandler struct {
	loaded bool
	chainList *[]ChainAccount
	chainSelectorStrategy *ChainAccountsStrategyTypes.IChainSelectorStrategy
}

func (x *ChainAccountsHandler) Load() (*ChainAccountsHandler) {
	// load chain account list into handler
	x.LoadChainAccountList()
	// load strategy into handler
	x.LoadChainStrategy(&ChainAccountsStrategyTypes.StrategyLoop{})
	// return self for chaining
	return x
}

func (x *ChainAccountsHandler) LoadChainStrategy(strategy ChainAccountsStrategyTypes.IChainSelectorStrategy) (*ChainAccountsHandler) {
	// load strategy with chain list
	strategy.Load(x.chainList)
	// store strategy in handler
	x.chainSelectorStrategy = &strategy
	// return self for chaining
	return x
}

func (x *ChainAccountsHandler) LoadChainAccountList() (*ChainAccountsHandler) {
	// TODO: need to fetch from DB
	// set list of chain accounts to chainList
	x.chainList = &[]ChainAccount {
		ChainAccount {
			ID: "bc4a5e44c384043c047f9a768e3c6b8d3064fd470b430f8e1a9d265114118e5d",
		},
		ChainAccount {
			ID: "bc4a5e44c384043c047f9a768e3c6b8d3064fd470b430f8e1a9d265114118e5d",
		},
		ChainAccount {
			ID: "bc4a5e44c384043c047f9a768e3c6b8d3064fd470b430f8e1a9d265114118e5d",
		},
		ChainAccount {
			ID: "bc4a5e44c384043c047f9a768e3c6b8d3064fd470b430f8e1a9d265114118e5d",
		},
		ChainAccount {
			ID: "bc4a5e44c384043c047f9a768e3c6b8d3064fd470b430f8e1a9d265114118e5d",
		},
	}
	// if strategy exists, update it with chain list
	if (x.chainSelectorStrategy != nil) {
		(*x.chainSelectorStrategy).Load(x.chainList)
	}
	// return self for chaining
	return x
}

func (x *ChainAccountsHandler) GetNextAccount() (*ChainAccount, error) {
	// ensure handler is loaded
	x.EnsureInitialLoad()
	// determine how to get the next chain account to use
	return (*x.chainSelectorStrategy).GetNextAccount()
}

func (x *ChainAccountsHandler) EnsureInitialLoad() {
	if (!x.loaded) {
		// handler is not loaded yet, load first
		x.Load()
		x.loaded = true
	}
}