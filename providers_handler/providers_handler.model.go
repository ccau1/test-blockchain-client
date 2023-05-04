package providers_handler

import (
	"fmt"
	"errors"
	ProviderTypes "github.com/ccau1/test-blockchain-client/providers_handler/provider"
)

type IProvider = ProviderTypes.IProvider

type GetNextProviderFilter struct {
	ChainType string
}

type ProvidersHandler struct {
	providers *[]IProvider
	loaded bool
}

func (x *ProvidersHandler) Load() *ProvidersHandler {
	// set different provides available fetched from DB
	x.providers = &[]IProvider {
		&ProviderTypes.AnkrProvider {
			
		},
		&ProviderTypes.AnkrProvider {
			
		},
		&ProviderTypes.AnkrProvider {
			
		},
	}

	return x
}

func (x *ProvidersHandler) GetNextProvider(filter GetNextProviderFilter) (*IProvider, error) {
	// ensure handler is loaded
	x.EnsureInitialLoad()
	filteredProviders := *x.providers
	n := 0
	if (filter.ChainType != "") {
		n = 0
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

	// TODO: implement strategies to determine which provider to return here
	// determine how to get the next provider to use
	return &(filteredProviders)[0], nil
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