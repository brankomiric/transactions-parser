package parser

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/brankomiric/transactions-parser/utils"
)

var defaultRpc = "https://mainnet.infura.io/v3/c9657d3c5621495c9f6b60c3913df958"

type Parser interface {
	// last parsed block
	GetCurrentBlock() int
	// add address to observer
	Subscribe(address string) bool
	// list of inbound or outbound transactions for an address
	GetTransactions(address string) []Transaction
}

type EthTxParser struct {
	rpcNodeUrl  string
	storage     storage
	lastBlock   string
	subscribers map[string]bool
}

func New() *EthTxParser {
	p := &EthTxParser{}
	rpcNode := os.Getenv("RPC_URL")
	if rpcNode == "" {
		log.Println("No RPC_URL set, using default RPC")
		p.rpcNodeUrl = defaultRpc
	}
	p.storage.records = make(map[string][]Transaction)
	p.subscribers = make(map[string]bool)
	go p.worker()
	return p
}

// Returns latest block number as int
func (p *EthTxParser) GetCurrentBlock() int {
	hex := p.getCurrentBlockHex()
	value, err := strconv.ParseInt(hex, 0, 64)
	if err != nil {
		log.Println(err.Error())
	}
	return int(value)
}

func (p *EthTxParser) Subscribe(address string) bool {
	p.subscribers[address] = true
	return true
}

func (p *EthTxParser) GetTransactions(address string) []Transaction {
	if p.subscribers[address] {
		return p.storage.Get(address)
	}
	log.Printf("%s not subscribed", address)
	return nil
}

func (s *EthTxParser) worker() {
	tickCh := time.NewTicker(time.Second * 5)
	for v := range tickCh.C {
		s.sync(v)
	}
}

// Returns latest block hex
func (p *EthTxParser) getCurrentBlockHex() string {
	rpcReq := utils.RpcRequest{
		Url: p.rpcNodeUrl,
	}
	body := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_blockNumber",
		"id":      1,
	}
	resp, err := rpcReq.QueryRpc(body)
	if err != nil {
		log.Println(err.Error())
	}
	var respBody GetBlockNumberResp
	err = json.Unmarshal(resp, &respBody)
	if err != nil {
		log.Println(err.Error())
	}
	return respBody.Result
}

func (p *EthTxParser) getListOfTransactionsForBlock(blockHex string) ([]Transaction, error) {
	rpcReq := utils.RpcRequest{
		Url: p.rpcNodeUrl,
	}
	body := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "eth_getBlockByNumber",
		"params":  []interface{}{blockHex, true},
		"id":      1,
	}
	resp, err := rpcReq.QueryRpc(body)
	if err != nil {
		return nil, err
	}
	var respBody GetBlockByNumberResp
	err = json.Unmarshal(resp, &respBody)
	if err != nil {
		return nil, err
	}
	return respBody.Result.Transactions, nil
}

func (p *EthTxParser) sync(t time.Time) {
	log.Printf("Sync invoked at: %s", t.Format("2006-01-02 15:04:05"))
	lastBlockHex := p.getCurrentBlockHex()
	if lastBlockHex == p.lastBlock {
		return
	}
	p.lastBlock = lastBlockHex
	txs, err := p.getListOfTransactionsForBlock(p.lastBlock)
	if err != nil {
		log.Println("Failed to sync storage")
		return
	}
	for _, tx := range txs {
		p.storage.Add(tx.From, tx)
		p.storage.Add(tx.To, tx)
	}
	log.Println("Successful sync to storage")
}
