package utils

import (
	"strings"
	"os"
	"fmt"
	"path"
	"runtime"
	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func InitLogger() {
	// set out to stdout
	Log.Out = os.Stdout
	// set true to displaying calling method
	Log.SetReportCaller(true)
	Log.Formatter = &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			function := path.Base(f.Function)[strings.LastIndex(path.Base(f.Function), ".") + 1:]
			return fmt.Sprintf("%s() \033[0;31m::\033[0m", function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
	}
	// if env is not dev, use the json formatter
	// for monitoring systems
	if os.Getenv("ENV") != "dev" {
		Log.Formatter = &logrus.JSONFormatter{}
	}
}
