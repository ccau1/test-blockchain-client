package provider_accounts_strategy

import (
	"log"
	"errors"
)

type LoopStrategy struct {
	providerAccounts *[]ProviderAccount
	currentPos int
}

func (x *LoopStrategy) Load(newProviderAccountsList *[]ProviderAccount) {
	// store chain list
	x.providerAccounts = newProviderAccountsList
	// set default position to 0
	x.currentPos = 0
}

func (x *LoopStrategy) GetNextAccount() (*ProviderAccount, error) {
	// if chain list not available, throw error
	if (x.providerAccounts == nil || len(*x.providerAccounts) == 0) {
		return nil, errors.New("chain account list not loaded")
	}
	// set new current position
	if (x.currentPos + 1 >= len(*x.providerAccounts)) {
		// reached end, reset to 0
		x.currentPos = 0
	} else {
		// more to add, add now
		x.currentPos = x.currentPos + 1
	}
	log.Printf("chain account fetch pos: %d\n", x.currentPos)
	// return chain account with no errors
	return &(*x.providerAccounts)[x.currentPos], nil
}