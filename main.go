package main

import (
	"os"
	"fmt"
	"strconv"
	"net/http"
	"github.com/sirupsen/logrus"

	"github.com/ccau1/test-blockchain-client/utils"
)

func init() {
	utils.LoadEnv()
	utils.InitLogger()
}

func main() {
	// init router
	r := GenRoute()

	// gen route doc
	utils.DocGen(r)

	// define host and port
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if (err != nil) {
		port = 3000
	}
	var host = os.Getenv("HOST")

	// start server
	utils.Log.WithFields(logrus.Fields{"host": host, "port": port}).Info("[Server] Starting")
	http.ListenAndServe(fmt.Sprintf("%s:%d", host, port), r)
}