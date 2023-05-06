package helper

import (
	"regexp"
	"strconv"
	"time"
)

func IntToHex(i int) string {
	return "0x" + strconv.FormatInt(int64(i), 16)
}

func HexToInt(s string) (int64, error) {
	return strconv.ParseInt(regexp.MustCompile("^0x").ReplaceAllString(s, ""), 16, 64)
}

func SetInterval(durationMS int64, fn func()) *time.Ticker {
	// generate ticker and save it
	ticker := time.NewTicker(time.Duration(durationMS * 1000000))
	// run interval handling
	for _ = range (*ticker).C {
		 // trigger function
		 fn()
	}

	return ticker
}