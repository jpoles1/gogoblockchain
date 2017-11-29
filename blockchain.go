package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/fatih/color"
)

//BlockChain is a struct to contain an instance of the blockchain
type BlockChain struct {
	CurrentTransactions []Transaction
	Chain               []Block
	Nodes               NodeSet
}

func (bc BlockChain) start() BlockChain {
	bc.newBlock(100, "1")
	return bc
}
func (bc BlockChain) validProof(lastProof int, proof int) bool {
	guess := strconv.Itoa(lastProof) + strconv.Itoa(proof)
	guessHash := shaHash(guess)
	fmt.Println("Validating...", guessHash)
	return guessHash[:2] == "00"
}
func (bc BlockChain) proofOfWork(lastProof int) int {
	proof := 0
	for bc.validProof(lastProof, proof) == false {
		proof++
	}
	return proof
}
func (bc *BlockChain) newBlock(proof int, prevHash string) Block {
	var newblock = Block{
		Index:        len(bc.Chain) + 1,
		Timestamp:    time.Now().UTC().UnixNano(),
		Transactions: bc.CurrentTransactions,
		Proof:        proof,
		PreviousHash: prevHash,
	}
	bc.CurrentTransactions = nil
	bc.Chain = append(bc.Chain, newblock)
	return newblock
}
func (bc *BlockChain) newTransaction(sender string, recipient string, amount int) int {
	bc.CurrentTransactions = append(bc.CurrentTransactions, Transaction{sender, recipient, amount})
	return len(bc.CurrentTransactions)
}
func (bc BlockChain) lastBlock() Block {
	//Returns the last Block in the chain
	return bc.Chain[len(bc.Chain)-1]
}
func (bc *BlockChain) registerNode(address string) {
	uri, _ := url.Parse(address)
	fmt.Println(uri)
	//bc.Nodes.add(Node{uri.})
}
func (bc BlockChain) validChain(blocks []Block) bool {
	prevBlock := blocks[0]
	i := 1
	fmt.Println("Validating Chain")
	for i < len(blocks) {
		block := blocks[i]
		fmt.Println(prevBlock)
		fmt.Println(block)
		//Check that the hash of the block is correct
		if block.PreviousHash != prevBlock.hash() {
			return false
		}
		//Check that the Proof of Work is correct
		if !bc.validProof(prevBlock.Proof, block.Proof) {
			return false
		}
		fmt.Println("-----------")
		prevBlock = block
		i++
	}
	return true
}
func (bc *BlockChain) resolveConflicts() bool {
	maxLength := len(bc.Chain)
	var newChain []Block
	for otherNode := range bc.Nodes.set {
		fmt.Println(otherNode)
		nodeResp, err := http.Get("http://" + otherNode.addr + "/chain")
		if err != nil || nodeResp.StatusCode != 200 {
			color.Yellow("Cannot get chain data from: ")
			fmt.Println(otherNode)
		} else {
			defer nodeResp.Body.Close()
			resp, err := ioutil.ReadAll(nodeResp.Body)
			if err != nil {
				fmt.Println("HTTP Read Error")
			} else {
				var jsonChain []Block
				if err := json.Unmarshal(resp, &jsonChain); err != nil {
					color.Yellow("Corrupted chain data from: ")
					fmt.Println(otherNode)
					fmt.Println(err)
					fmt.Println("--------------")
				} else if len(jsonChain) > maxLength && bc.validChain(jsonChain) {
					maxLength = len(jsonChain)
					newChain = jsonChain
				}
			}
		}
	}
	if newChain != nil {
		bc.Chain = newChain
		return true
	}
	return false
}
