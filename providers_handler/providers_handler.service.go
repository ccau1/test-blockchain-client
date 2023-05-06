package providers_handler

import (
	"errors"
	"fmt"

	ProviderTypes "github.com/ccau1/test-blockchain-client/providers_handler/provider"
	ProviderStrategies "github.com/ccau1/test-blockchain-client/providers_handler/provider_strategy"
)

type ProvidersHandler struct {
	providers *[]IProvider
	loaded bool

	providersStrategy ProviderStrategies.IProvidersStrategy
	UseStrategy IProvidersStrategy
}

func (x *ProvidersHandler) Load() *ProvidersHandler {
	x.providers = &[]IProvider {
		&ProviderTypes.AnkrProvider {
			
		},
		&ProviderTypes.AnkrProvider {
			
		},
		&ProviderTypes.BlockDaemonProvider {
			
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
        if contains(val.SupportedChains(), filter.ChainType) {
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

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}