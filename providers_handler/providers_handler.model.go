package providers_handler

import (
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

func (x *ProvidersHandler) GetNextProvider(filter GetNextProviderFilter) *IProvider {
	// ensure handler is loaded
	x.EnsureInitialLoad()
	// TODO: implement strategies to determine which provider to return here
	// determine how to get the next provider to use
	return &(*x.providers)[0]
}

func (x *ProvidersHandler) EnsureInitialLoad() {
	if (!x.loaded) {
		// handler is not loaded yet, load first
		x.Load()
		x.loaded = true
	}
}