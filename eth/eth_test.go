package eth

import (
	"fmt"
	"regexp"
	"testing"
	"errors"
	"encoding/json"
)

type BlockNumberResponse struct {
	jsonrpc string
	id int
	result string
}

// TestHelloName calls greetings.Hello with a name, checking
// for a valid return value.
func TestCallPolygonRPC_NoError(t *testing.T) {
		body := &GetBlockNumberBody{
			JSONRPC: "2.0",
			Method: "eth_blockNumber",
			ID: 2,
		}
		requestBodyByte, _ := json.Marshal(body)
    name := "2.0"

    want := regexp.MustCompile(name)
    msg, err := CallPolygonRPC(requestBodyByte)

		res := BlockNumberResponse{}
    json.Unmarshal(msg, &res)

		fmt.Printf("%+v", res)

    if !want.MatchString(res.jsonrpc) || err != nil {
        t.Fatalf(`Hello("Gladys") = %q, %v, want match for %#q, nil`, msg, err, want)
    }
}

// TestHelloEmpty calls greetings.Hello with an empty string,
// checking for an error.
func TestHelloEmpty(t *testing.T) {
    msg, err := Hello("")
    if msg != "" || err == nil {
        t.Fatalf(`Hello("") = %q, %v, want "", error`, msg, err)
    }
}

func Hello(str string) (string, error) {
	var err error
	if str == "" {
		err = errors.New("need string")
	}
	return str, err
}