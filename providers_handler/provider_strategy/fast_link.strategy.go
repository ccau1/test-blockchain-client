package provider_strategy

import (
	"github.com/ccau1/test-blockchain-client/utils"
)

var Log = utils.Log

type FastLinkStrategy struct {

}

func (x *FastLinkStrategy) GetNextProvider(providers []IProvider, options *GetNextAccountOptions) (*IProvider, error) {
	Log.Infof("start selecting provider\n")

	/*
		selecting criteria:

		1. response time
		2. up/down status
	*/

	return &providers[0], nil
}