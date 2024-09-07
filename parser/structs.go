package parser

type GetBlockNumberResp struct {
	RpcVersion string `json:"jsonrpc"`
	Result     string `json:"result"`
	Id         int    `json:"id"`
}

type GetBlockByNumberResp struct {
	RpcVersion string                 `json:"jsonrpc"`
	Result     GetBlockByNumberResult `json:"result"`
	Id         int                    `json:"id"`
}

type GetBlockByNumberResult struct {
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	ChainId          string `json:"chainId"`
	From             string `json:"from"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	Hash             string `json:"hash"`
	Input            string `json:"input"`
	Nonce            string `json:"nonce"`
	R                string `json:"r"`
	S                string `json:"s"`
	To               string `json:"to"`
	TransactionIndex string `json:"transactionIndex"`
	Type             string `json:"type"`
	Value            string `json:"value"`
}
