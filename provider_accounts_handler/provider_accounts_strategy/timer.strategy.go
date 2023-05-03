package provider_accounts_strategy

import (
	"time"
	"errors"
	"math"
	"math/rand"
)

/**

DISCLAIMER:: 	Wrong implementation, timer should be used for different endpoints (per provider) instead
							of applying to different accounts of the same endpoint. This will be removed when handling
							different provider endpoints has been implemented. Current level is for handling multiple
							accounts to one provider endpoint

**/

// the fastest time idx can only stay the king for
// this amount of MS before being contested again
var FASTEST_TIME_STREAK_LIMIT_MS = 1000 * 60 * 10

// the percentage change a user will try out another
// chain account for time instead of using fastest
// 0 - 100 as percent
var VOLUNTEER_TEST_TIME_PERCENTAGE = 30

type TimerStrategy struct {
	providerAccounts *[]ProviderAccount
	timerMap *[]int

	fastestIdx int
	fastestTime int
	lastUpdatedFastest int
}

func (x *TimerStrategy) Load(newProviderAccountsList *[]ProviderAccount) {
	// if chain account list not provided, throw error
	if (newProviderAccountsList == nil || len(*newProviderAccountsList) == 0) {
		panic(errors.New("chain account list not provided"))
	}
	// store chain list
	x.providerAccounts = newProviderAccountsList
	// set timerMap slice size based on chain account list size
	timerMap := make([]int, len(*x.providerAccounts))
	x.timerMap = &timerMap
	// trackers for fastest
	x.fastestIdx = -1
	x.fastestTime = -1
	x.lastUpdatedFastest = -1
}

func (x *TimerStrategy) GetNextAccount() (*ProviderAccount, error) {
	// if chain list not available, throw error
	if (x.providerAccounts == nil || len(*x.providerAccounts) == 0) {
		return nil, errors.New("chain account list not loaded")
	}

	// if random number is less than VOLUNTEER_TEST_TIME_PERCENTAGE,
	// do a random speed test
	if (x.fastestIdx == -1 || rand.Intn(100) <= VOLUNTEER_TEST_TIME_PERCENTAGE) {
		// try another account. We'll randomize idx to try
		tryoutIdx := int(math.Floor((rand.Float64() - 0.01) * float64(len(*x.providerAccounts))))
		// run speed test at idx
		x.runSpeedTest(tryoutIdx)
	}
	// return fastest
	return &(*x.providerAccounts)[x.fastestIdx], nil
}

func (x *TimerStrategy) runSpeedTest(idx int) {
	// get chain account by idx
	providerAcc := &(*x.providerAccounts)[idx]
	// test out speed
	(*x.timerMap)[idx] = x.getAccountSpeed(providerAcc)
	// set response time and update fastest
	x.setResponseTime(idx, (*x.timerMap)[idx])
}

func (x *TimerStrategy) getAccountSpeed(providerAccount *ProviderAccount) int {
	// TODO: how to test time? how to get url to ping?
	return 1
}

func (x *TimerStrategy) setResponseTime(idx int, time int) {
	(*x.timerMap)[idx] = time

	if (x.fastestTime == -1 || time < x.fastestTime || x.getCurrentMS() - x.lastUpdatedFastest > FASTEST_TIME_STREAK_LIMIT_MS) {
		// if no time has been set as fastest OR time beats current fastest OR fastest time has stayed king too long,
		// set current one as fastest
		x.fastestTime = time
		x.fastestIdx = idx
		x.lastUpdatedFastest = x.getCurrentMS()
	}
}

func (x *TimerStrategy) getCurrentMS() int {
		// create a time variable
		now := time.Now()

		// convert to unix time in milliseconds
		unixNano := now.UnixNano()

		return int(unixNano / 1000000)
}