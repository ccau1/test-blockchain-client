package provider_accounts_handler

import (
	"github.com/ccau1/test-blockchain-client/utils"
	"github.com/ccau1/test-blockchain-client/provider_accounts_handler/provider_account"
	ProviderAccountsStrategyTypes "github.com/ccau1/test-blockchain-client/provider_accounts_handler/provider_accounts_strategy"
)

var Log = utils.Log.WithField("class", "ProviderAccountsHandler")

type ProviderAccount = provider_account.ProviderAccount
type IProviderAccountsStrategy = ProviderAccountsStrategyTypes.IProviderAccountsStrategy

type DEFAULT_PROVIDER_ACCOUNTS_STRATEGY = ProviderAccountsStrategyTypes.LoopStrategy


type ProviderAccountsHandler struct {
	// whether initial load has been ran
	loaded bool
	// provider accounts list
	providerAccounts *[]ProviderAccount
	// strategy being used
	providerAccountsStrategy *IProviderAccountsStrategy
	// user input provider name
	Provider string
	// user input strategy to use
	UseStrategy IProviderAccountsStrategy
}

func (x *ProviderAccountsHandler) Load() (*ProviderAccountsHandler) {
	// load chain account list into handler
	x.LoadProviderAccountList()
	// load strategy into handler
	if (x.UseStrategy == nil) {
		x.LoadChainStrategy(&DEFAULT_PROVIDER_ACCOUNTS_STRATEGY{})
	} else {
		x.LoadChainStrategy(x.UseStrategy)
	}
	// return self for chaining
	return x
}

func (x *ProviderAccountsHandler) LoadChainStrategy(strategy ProviderAccountsStrategyTypes.IProviderAccountsStrategy) (*ProviderAccountsHandler) {
	// load strategy with chain list
	strategy.Load(x.providerAccounts)
	// store strategy in handler
	x.providerAccountsStrategy = &strategy
	// return self for chaining
	return x
}

func (x *ProviderAccountsHandler) LoadProviderAccountList() (*ProviderAccountsHandler) {
	Log := Log.WithField("method", "LoadProviderAccountList")
	// TODO: need to fetch from DB based on x.Provider
	Log.Infof("fetching DB for chain provider accounts with provider type: %s", x.Provider)
	// set list of chain accounts to providerAccounts
	x.providerAccounts = &[]ProviderAccount {
		ProviderAccount {
			ID: "bc4a5e44c384043c047f9a768e3c6b8d3064fd470b430f8e1a9d265114118e5d",
			ProviderType: "ankr",
		},
		ProviderAccount {
			ID: "bc4a5e44c384043c047f9a768e3c6b8d3064fd470b430f8e1a9d265114118e5d",
			ProviderType: "ankr",
		},
		ProviderAccount {
			ID: "bc4a5e44c384043c047f9a768e3c6b8d3064fd470b430f8e1a9d265114118e5d",
			ProviderType: "ankr",
		},
		ProviderAccount {
			ID: "bc4a5e44c384043c047f9a768e3c6b8d3064fd470b430f8e1a9d265114118e5d",
			ProviderType: "ankr",
		},
		ProviderAccount {
			ID: "bc4a5e44c384043c047f9a768e3c6b8d3064fd470b430f8e1a9d265114118e5d",
			ProviderType: "ankr",
		},
	}
	// if strategy exists, update it with chain list
	if (x.providerAccountsStrategy != nil) {
		(*x.providerAccountsStrategy).Load(x.providerAccounts)
	}
	// return self for chaining
	return x
}

func (x *ProviderAccountsHandler) GetNextAccount() (*ProviderAccount, error) {
	// ensure handler is loaded
	x.EnsureInitialLoad()
	// determine how to get the next chain account to use
	return (*x.providerAccountsStrategy).GetNextAccount()
}

func (x *ProviderAccountsHandler) EnsureInitialLoad() {
	if (!x.loaded) {
		// handler is not loaded yet, load first
		x.Load()
		x.loaded = true
	}
}