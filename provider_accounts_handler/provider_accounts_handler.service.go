package provider_accounts_handler

import (
	ProviderAccountsStrategyTypes "github.com/ccau1/test-blockchain-client/provider_accounts_handler/provider_accounts_strategy"
)

type DEFAULT_PROVIDER_ACCOUNTS_STRATEGY = ProviderAccountsStrategyTypes.LoopStrategy

/*
	Central control for managing a list of accounts for providers
*/
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

/*
	initialize provider accounts handler
*/
func (x *ProviderAccountsHandler) Load() (*ProviderAccountsHandler) {
	// load chain account list into handler
	x.loadProviderAccountList()
	// load strategy into handler
	if (x.UseStrategy == nil) {
		x.loadProviderAccountsStrategy(&DEFAULT_PROVIDER_ACCOUNTS_STRATEGY{})
	} else {
		x.loadProviderAccountsStrategy(x.UseStrategy)
	}
	// return self for chaining
	return x
}

/*
	load provider accounts strategy for handling next account selection
*/
func (x *ProviderAccountsHandler) loadProviderAccountsStrategy(strategy ProviderAccountsStrategyTypes.IProviderAccountsStrategy) (*ProviderAccountsHandler) {
	// load strategy with chain list
	strategy.Load(x.providerAccounts)
	// store strategy in handler
	x.providerAccountsStrategy = &strategy
	// return self for chaining
	return x
}

/*
	load provider account list to run strategy against
*/
func (x *ProviderAccountsHandler) loadProviderAccountList() (*ProviderAccountsHandler) {
	Log.Infof("fetching DB for chain provider accounts with provider type: %s", x.Provider)
	// TODO: need to fetch from DB based on x.Provider
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
		ProviderAccount {
			ID: "xKAwowrWTIy6hh18fBlZDXEtJf2luTN66xbkTuZEHbGK6Cog",
			ProviderType: "block_daemon",
		},
	}
	// filter by provider
	if (x.Provider != "") {
		n := 0
    for _, val := range *x.providerAccounts {
        if (val.ProviderType == x.Provider) {
					(*x.providerAccounts)[n] = val
            n++
        }
    }

		newList := (*x.providerAccounts)[:n]

    x.providerAccounts = &newList
	}
	// if strategy exists, update it with chain list
	if (x.providerAccountsStrategy != nil) {
		(*x.providerAccountsStrategy).Load(x.providerAccounts)
	}
	// return self for chaining
	return x
}

/*
	get an account from the list of accounts based on strategy
*/
func (x *ProviderAccountsHandler) GetNextAccount(filter *GetNextAccountFilter) (*ProviderAccount, error) {
	// ensure handler is loaded
	x.ensureInitialLoad()
	// determine how to get the next chain account to use
	return (*x.providerAccountsStrategy).GetNextAccount()
}

/*
	if not initialized, load now. Otherwise, skip
*/
func (x *ProviderAccountsHandler) ensureInitialLoad() {
	if (!x.loaded) {
		// handler is not loaded yet, load first
		x.Load()
		x.loaded = true
	}
}