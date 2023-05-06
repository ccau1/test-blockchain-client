package provider_strategy

import (
	"reflect"
	"time"

	"github.com/ccau1/test-blockchain-client/utils"
	"github.com/ccau1/test-blockchain-client/utils/helper"
)

var Log = utils.Log

// only check every 1 min or more
var CHECK_RESPONSE_TIME_INTERVAL int64 = 60 * 1000

type ResponseTimeItem struct {
	ms int64
	provider *IProvider
}

type FastLinkStrategy struct {
	providersResponseTime map[string]*ResponseTimeItem
	lastCheckedResponseTime int64
	responseTimeTicker *time.Ticker
}

func (x *FastLinkStrategy) GetNextProvider(providers []IProvider, options *GetNextAccountOptions) (*IProvider, error) {
	Log.Infof("start selecting provider: %s\n", reflect.TypeOf(providers[0]))

	x.addProvidersToScheduledSpeedTest(providers)

	/*
		selecting criteria:

		1. response time
		2. up/down status
	*/

	return &providers[0], nil
}

func (x *FastLinkStrategy) Load() {
	// instantiate map
	x.providersResponseTime = make(map[string]*ResponseTimeItem)
	// set interval speed tests
	go x.LoadIntervalSpeedTest(CHECK_RESPONSE_TIME_INTERVAL)
}

func (x *FastLinkStrategy) LoadIntervalSpeedTest(intervalMS int64) {
	// if ticker already exists, stop previous ticker
	if (x.responseTimeTicker != nil) {
		(*x.responseTimeTicker).Stop()
	}

	// run initial function call
	x.intervalSpeedTestJob()
	// generate ticker and save it
	x.responseTimeTicker = helper.SetInterval(intervalMS, x.intervalSpeedTestJob)
}

func (x *FastLinkStrategy) addProvidersToScheduledSpeedTest(providers []IProvider) {
	for _, provider := range providers {
		providerName := reflect.TypeOf(provider).String()
		if (x.providersResponseTime[providerName] == nil) {
			// not added yet. Add now
			x.providersResponseTime[providerName] = &ResponseTimeItem{
				ms: 100000,
				provider: &provider,
			}
		}
	}
}

func (x *FastLinkStrategy) intervalSpeedTestJob() {
	if (time.Now().UnixMilli() - x.lastCheckedResponseTime < CHECK_RESPONSE_TIME_INTERVAL) {
		// not passed check time interval yet, skip
		return
	}

	Log.Infof("speed test job: START")

	// update last checked time right away so others won't do it as well
	// others as in other instances of this service also running checks
	// and response time is stored in memCache
	x.lastCheckedResponseTime = time.Now().UnixMilli()

	var startTime, endTime int64

	for key, responseTimeItem := range x.providersResponseTime {
		startTime = time.Now().UnixMilli()
		err := (*responseTimeItem.provider).Ping()

		if (err != nil) {
			// TODO: handle error. Store provider as down somewhere?
			// is it based on provider account?
		}

		endTime = time.Now().UnixMilli()

		responseTimeItem.ms = endTime - startTime
		Log.Infof("speed test result [%s]: %+v", key, responseTimeItem.ms)
	}

	Log.Infof("speed test job: END")
}
