package provider_strategy

import (
	"github.com/ccau1/test-blockchain-client/providers_handler/provider"
)

type GetNextAccountOptions struct {
	Key string `json:"key"`
}

type IProvider = provider.IProvider

type IProvidersStrategy interface {
	Load()
	GetNextProvider(providers []IProvider, options *GetNextAccountOptions) (*IProvider, error)
}
