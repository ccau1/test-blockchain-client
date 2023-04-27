package eth

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func CallPolygonRPC(body []byte) ([]byte, error) {
	// call provider to retrieve info
	providerDomain := "https://polygon-rpc.com/"
	res, err := http.Post(
		providerDomain,					// url
		"application/json",			// content-type
		bytes.NewBuffer(body), 	// body (as buffer)
	)
	if err != nil {
		return nil, err
	}

	// get content from response
	defer res.Body.Close()
	resContent, err := ioutil.ReadAll(res.Body)
	return resContent, nil
}