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

func (x *ProvidersHandler) Load() *ProvidersHandler {
	x.providers = &[]IProvider {
		&ankr_provider.AnkrProvider {
			
		},
		&ankr_provider.AnkrProvider {
			
		},
		&block_daemon_provider.BlockDaemonProvider {
			
		},
	}

	if (x.UseStrategy == nil) {
		x.LoadProviderStrategy(&DEFAULT_PROVIDERS_STRATEGY{})
	} else {
		x.LoadProviderStrategy(x.UseStrategy)
	}

	return x
}

func (x *ProvidersHandler) LoadProviderStrategy(strategy ProviderStrategies.IProvidersStrategy) (*ProvidersHandler) {
	x.providersStrategy = strategy

	x.providersStrategy.Load()

	return x
}

func (x *ProvidersHandler) filterKeyGenerator(filter GetNextProviderFilter) string {
	key := ""

	if (filter.ChainType != "") {
		key = fmt.Sprintf("%s_%s", key, filter.ChainType)
	}

	return key
}

func (x *ProvidersHandler) GetNextProvider(filter GetNextProviderFilter) (*IProvider, error) {
	// ensure handler is loaded
	x.EnsureInitialLoad()

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

	provider, err := x.providersStrategy.GetNextProvider(filteredProviders, &ProviderStrategies.GetNextAccountOptions{
		Key: x.filterKeyGenerator(filter),
	})

	if (err != nil) {
		return nil, errors.New(fmt.Sprintf("[GetNextProvider] GetNextProvider error: %+v", err))
	}

	return provider, nil
}

func (x *ProvidersHandler) EnsureInitialLoad() {
	if (!x.loaded) {
		// handler is not loaded yet, load first
		x.Load()
		x.loaded = true
	}
}
