package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
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
	bc.Nodes = NodeSet{make(map[Node]bool)}
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
func (bc *BlockChain) newTransaction(t Transaction) int {
	bc.CurrentTransactions = append(bc.CurrentTransactions, t)
	return len(bc.CurrentTransactions)
}
func (bc BlockChain) lastBlock() Block {
	//Returns the last Block in the chain
	return bc.Chain[len(bc.Chain)-1]
}
func (bc *BlockChain) registerNode(address string) (string, error) {
	uri, err := url.Parse(address)
	if err != nil {
		return "", err
	}
	if uri.Host == "" {
		if destReachable(address) {
			return address, nil
		}
		return "", errors.New("Address is Invalid or Cannot Be Reached")
	}
	if !destReachable(uri.Host) {
		return "", errors.New("Cannot fetch chain from this server")
	}
	bc.Nodes.add(Node{uri.Host})
	return uri.Host, err
}
func (bc BlockChain) validChain(blocks []Block) bool {
	prevBlock := blocks[0]
	i := 1
	log.Println("Validating Chain")
	for i < len(blocks) {
		block := blocks[i]
		//Check that the hash of the block is correct
		if block.PreviousHash != prevBlock.hash() {
			return false
		}
		//Check that the Proof of Work is correct
		if !bc.validProof(prevBlock.Proof, block.Proof) {
			return false
		}
		prevBlock = block
		i++
	}
	return true
}
func (bc *BlockChain) resolveConflicts() bool {
	maxLength := len(bc.Chain)
	var newChain []Block
	for otherNode := range bc.Nodes.set {
		nodeResp, err := http.Get("http://" + otherNode.addr + "/chain")
		if err != nil {
			color.Yellow("Cannot get chain data from: ")
			fmt.Println(otherNode)
			log.Println("Error: " + err.Error())
		} else {
			defer nodeResp.Body.Close()
			resp, err := ioutil.ReadAll(nodeResp.Body)
			if err != nil {
				log.Println("HTTP Read Error: " + err.Error())
			} else {
				var jsonChain []Block
				if err := json.Unmarshal(resp, &jsonChain); err != nil {
					color.Yellow("Corrupted chain data from: ")
					fmt.Println(otherNode)
					log.Println("Error: " + err.Error())
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
