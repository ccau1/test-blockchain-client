package chain_accounts_strategy

import (
	"log"
	"errors"
)

type StrategyLoop struct {
	chainList *[]ChainAccount
	currentPos int
}

func (x *StrategyLoop) Load(chainAccountList *[]ChainAccount) {
	// store chain list
	x.chainList = chainAccountList
	// set default position to 0
	x.currentPos = 0
}

func (x *StrategyLoop) GetNextAccount() (*ChainAccount, error) {
	// if chain list not available, throw error
	if (x.chainList == nil || len(*x.chainList) == 0) {
		return nil, errors.New("chain account list not loaded")
	}
	// set new current position
	if (x.currentPos + 1 >= len(*x.chainList)) {
		log.Println("In 1")
		// reached end, reset to 0
		x.currentPos = 0
	} else {
		log.Println("In 2")
		// more to add, add now
		x.currentPos = x.currentPos + 1
	}
	log.Printf("current pos: %d\n", x.currentPos)
	// return chain account with no errors
	return &(*x.chainList)[x.currentPos], nil
}