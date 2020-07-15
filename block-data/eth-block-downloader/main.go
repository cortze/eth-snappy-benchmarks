package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	blockids := []int{5671744, 10423672, 10423671, 10423670, 10423669, 10423668, 10423667, 10423666, 10423665, 10423664, 10423663, 10423662, 10423661, 10423660, 10423659}
	blocknames := []string{"block0.json", "block1.json", "block2.json", "block3.json", "block4.json", "block5.json", "block6.json", "block7.json", "block8.json", "block9.json", "block10.json", "block11.json", "block12.json", "block13.json", "block14.json", "block15.json"}
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/a2272e92a733416a8abc2add6eeb6a20")
	if err != nil {
		log.Fatal(err)
	}

	// Get current working directory
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}

	/*
	   header, err := client.HeaderByNumber(context.Background(), nil)
	   if err != nil {
	       log.Fatal(err)
	   }

	   fmt.Println(header.Number.String()) // 5671744
	*/
	var i int
	for i = 0; i < len(blockids); i++ {

		blockNumber := big.NewInt(int64(blockids[i]))
		block, err := client.BlockByNumber(context.Background(), blockNumber)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("block ", block.Number(), "downloaded")
		/*
			b := Block{
				jsonrpc:          *block.Number(),
				id:               *block.Number(),
				difficulty:       *block.Difficulty(),
				gasLimit:         string(block.GasLimit()),
				gasUsed:          string(block.GasUsed()),
				hash:             block.Hash(),
				mixHash:          string(block.ParentHash()),
				nonce:            string(block.Nonce()),
				number:           string(block.Number()),
				parentHash:       block.ParentHash(),
				size:             string(block.Size()),
				stateRoot:        string(block.StateRoot()),
				timestamp:        string(block.TimeStamp()),
				totalDifficulty:  string(block.TotalDifficilty()),
				transactions:     string(block.Transactions()),
				transactionsRoot: string(block.TransactionsRoot()),
				uncles:           string(block.Uncles()),
			}
		*/

		header := block.Header()
		body := block.Body()

		headerj, err := json.Marshal(header)
		if err != nil {
			fmt.Println(err)
			return
		}

		bodyj, err := json.Marshal(body)
		if err != nil {
			fmt.Println(err)
			return
		}

		path := wd + "/data/" + blocknames[i]
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
		if err != nil {
			panic(err)
		}
		_, err = f.Write(headerj)
		if err != nil {
			panic(err)
		}

		f.Seek(0, io.SeekEnd)

		_, err = f.Write(bodyj)
		if err != nil {
			panic(err)
		}

		f.Close()

	}
}
