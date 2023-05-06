package providers_handler

import (
	ProviderTypes "github.com/ccau1/test-blockchain-client/providers_handler/provider"
	ProviderStrategies "github.com/ccau1/test-blockchain-client/providers_handler/provider_strategy"
)

type IProvider = ProviderTypes.IProvider
type IProvidersStrategy = ProviderStrategies.IProvidersStrategy

type GetNextProviderFilter struct {
	ProviderStrategies.GetNextAccountOptions
	ChainType string
}

type DEFAULT_PROVIDERS_STRATEGY = ProviderStrategies.FastLinkStrategy
