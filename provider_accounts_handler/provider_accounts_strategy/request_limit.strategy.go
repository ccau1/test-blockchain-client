package provider_accounts_strategy

import (
	"time"
	"errors"

	"github.com/ccau1/test-blockchain-client/utils"
)

var Log = utils.Log.WithField("class", "RequestLimitStrategy")

type RequestTrackerItem struct {
	// request count
	count int
	// timestamp of interval start time MS
	intervalStart int
}

type RequestLimitStrategy struct {
	providerAccounts *[]ProviderAccount

	// requests amount limited
	LimitAmount int
	// limited every interval (in ms)
	LimitPerInterval int

	// currently used tracker index to run requests
	currentTrackerIdx int
	// request tracker list
	reqTrackerList *[]*RequestTrackerItem
}

func (x *RequestLimitStrategy) Load(newProviderAccountsList *[]ProviderAccount) {
	// if chain account list not provided, throw error
	if (newProviderAccountsList == nil || len(*newProviderAccountsList) == 0) {
		panic(errors.New("chain account list not provided"))
	}
	// store chain list
	x.providerAccounts = newProviderAccountsList
	// instantiate length of reqTrackerList
	newReqTrackerList := make([]*RequestTrackerItem, len(*x.providerAccounts))
	x.reqTrackerList = &newReqTrackerList
}

func (x *RequestLimitStrategy) GetNextAccount() (*ProviderAccount, error) {
	Log := Log.WithField("method", "GetNextAccount")
	// if chain list not available, throw error
	if (x.providerAccounts == nil || len(*x.providerAccounts) == 0) {
		return nil, errors.New("chain account list not loaded")
	}
	chainAccountIdx := x.getUsableProviderAccountIndex()
	Log.Infof("using chain account idx: %d", chainAccountIdx)
	
	(*(*x.reqTrackerList)[chainAccountIdx]).count++

	// return fastest
	return &(*x.providerAccounts)[chainAccountIdx], nil
}

func (x *RequestLimitStrategy) getUsableProviderAccountIndex() int {
	Log := Log.WithField("method", "getUsableProviderAccountIndex")
	// if req tracker list doesn't have this index, instantiate now
	if ((*x.reqTrackerList)[x.currentTrackerIdx] == nil) {
		(*x.reqTrackerList)[x.currentTrackerIdx] = &RequestTrackerItem{
			count: 0,
			intervalStart: x.getCurrentMS(),
		}
	}

	// define current tracker variable
	currentReqTracker := (*x.reqTrackerList)[x.currentTrackerIdx]

	if (x.getCurrentMS() - currentReqTracker.intervalStart > x.LimitPerInterval) {
		Log.Infof("[tracker %d] current tracker interval has passed, reset now", x.currentTrackerIdx)
		// current tracker interval has passed, reset now
		x.cleanReqTracker(currentReqTracker)
		return x.currentTrackerIdx
	} else if ((*currentReqTracker).count < x.LimitAmount) {
		Log.Infof("[tracker %d] still within its limit, continue at this idx [count: %d] [limit: %d]", x.currentTrackerIdx, (*currentReqTracker).count, x.LimitAmount)
		// still within its limit, continue at this idx
		return x.currentTrackerIdx
	} else if ((*x.cleanReqTracker((*x.reqTrackerList)[0])).count == 0) {
		Log.Infof("[tracker %d] first account tracker has reset, go back to first one", x.currentTrackerIdx)
		// first account tracker has reset, go back to first one
		x.currentTrackerIdx = 0
		return 0
	} else if (x.currentTrackerIdx + 1 == len(*x.providerAccounts)) {
		Log.Infof("[tracker %d] already at last account, just use this one", x.currentTrackerIdx)
		// already at last account, just use this one
		return x.currentTrackerIdx
	} else {
		Log.Infof("[tracker %d] current idx reached limit, go to next idx", x.currentTrackerIdx)
		// current idx reached limit, go to next idx
		x.currentTrackerIdx = x.currentTrackerIdx + 1
		return x.getUsableProviderAccountIndex()
	}
}

func (x *RequestLimitStrategy) cleanReqTracker(reqTracker *RequestTrackerItem) *RequestTrackerItem {
	if (x.getCurrentMS() - (*reqTracker).intervalStart > x.LimitPerInterval) {
		// tracker interval has passed, reset it
		(*reqTracker).count = 0
		(*reqTracker).intervalStart = x.getCurrentMS()
	}
	return reqTracker
}

func (x *RequestLimitStrategy) getCurrentMS() int {
	// create a time variable
	now := time.Now()

	// convert to unix time in milliseconds
	unixNano := now.UnixNano()

	return int(unixNano / 1000000)
}