# transactions-parser

Simple lib for accumulating blockchain data. This is only for show, as it stores everything in memory. Switch to some persistent storage if intending to use for real.

Set *RPC_URL* env var if you whish to use your own node, otherwise it will default to Ethereum Mainnet free Infura Node.

### Usage Example
```
import (
	txp "github.com/brankomiric/transactions-parser/parser"
)

func main() {
    parser := txp.New()

    blockNumber := parser.GetCurrentBlock()

    parser.Subscribe("<eth_address>")

	txs := parser.GetTransactions("<eth_address>")
	log.Printf("Tx count %d\n", len(txs1))
	log.Printf("Tx0: from %s, to %s, gas %s, input %s\n", txs[0].From, txs[0].To, txs[0].Gas, txs[0].Input)
}
```