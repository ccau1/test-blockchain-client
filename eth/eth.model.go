package eth

type GetBlockNumberBody struct {
	JSONRPC    	string `json:"jsonrpc"`
	Method 			string `json:"method"`
	ID 					int `json:"id"`
}

type GetBlockByNumberBody struct {
	JSONRPC    	string `json:"jsonrpc"`
	Method 			string `json:"method"`
	Params			[]interface{} `json:"params"`
	ID 					int `json:"id"`
}