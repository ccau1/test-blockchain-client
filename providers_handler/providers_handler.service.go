package providers_handler

import (
	"errors"
	"fmt"

	ankr_provider "github.com/ccau1/test-blockchain-client/providers_handler/provider/ankr"
	block_daemon_provider "github.com/ccau1/test-blockchain-client/providers_handler/provider/block_daemon"
	ProviderStrategies "github.com/ccau1/test-blockchain-client/providers_handler/provider_strategy"
	"golang.org/x/exp/slices"
)

type ProvidersHandler struct {
	providers *[]IProvider
	loaded bool

	providersStrategy ProviderStrategies.IProvidersStrategy
	UseStrategy IProvidersStrategy
}

/*
	providers handler initialization
*/
func (x *ProvidersHandler) Load() *ProvidersHandler {
	// define list of providers
	x.providers = &[]IProvider {
		&ankr_provider.AnkrProvider {
			
		},
		&ankr_provider.AnkrProvider {
			
		},
		&block_daemon_provider.BlockDaemonProvider {
			
		},
	}
	// load provider strategy by either passed in UseStrategy
	// or the default strategy
	if (x.UseStrategy == nil) {
		x.loadProviderStrategy(&DEFAULT_PROVIDERS_STRATEGY{})
	} else {
		x.loadProviderStrategy(x.UseStrategy)
	}
	// return self for chaining
	return x
}

/*
	store and run initialization of the strategy
*/
func (x *ProvidersHandler) loadProviderStrategy(strategy ProviderStrategies.IProvidersStrategy) (*ProvidersHandler) {
	// store strategy into provider instance
	x.providersStrategy = strategy
	// run strategy's load
	x.providersStrategy.Load()
	// return self for chaining
	return x
}

/*
	define a unique key based on the values of filter
*/
func (x *ProvidersHandler) filterKeyGenerator(filter GetNextProviderFilter) string {
	key := ""

	if (filter.ChainType != "") {
		key = fmt.Sprintf("%s_%s", key, filter.ChainType)
	}

	return key
}

/*
	get a provider based on the outcome of the strategy and
	param filters
*/
func (x *ProvidersHandler) GetNextProvider(filter GetNextProviderFilter) (*IProvider, error) {
	// ensure handler is loaded
	x.ensureInitialLoad()

	filteredProviders := *x.providers
	// filter providers based on chain type
	if (filter.ChainType != "") {
		n := 0
    for _, val := range filteredProviders {
        if slices.Contains(val.SupportedChains(), filter.ChainType) {
					filteredProviders[n] = val
            n++
        }
    }

    filteredProviders = filteredProviders[:n]

		if len(filteredProviders) == 0 {
			return nil, errors.New(fmt.Sprintf("no provider available for chain type [%s]", filter.ChainType))
		}
	}
	// get provider from list of providers based on the strategy
	provider, err := x.providersStrategy.GetNextProvider(filteredProviders, &ProviderStrategies.GetNextAccountOptions{
		Key: x.filterKeyGenerator(filter),
	})
	if (err != nil) {
		return nil, errors.New(fmt.Sprintf("[GetNextProvider] GetNextProvider error: %+v", err))
	}
	// return selected provider
	return provider, nil
}

/*
	if not ran initial load, run. Otherwise, skip it
*/
func (x *ProvidersHandler) ensureInitialLoad() {
	if (!x.loaded) {
		// handler is not loaded yet, load first
		x.Load()
		x.loaded = true
	}
}
