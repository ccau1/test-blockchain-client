package utils

import (
	"os"
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func InitLogger() {
	// set out to stdout
	Log.Out = os.Stdout
	// set true to displaying calling method
	Log.SetReportCaller(true)
	// if env is not dev, use the logrus formatter
	// for monitoring systems
	if os.Getenv("ENV") != "dev" {
		Log.Formatter = &logrus.JSONFormatter{}
	}
}
