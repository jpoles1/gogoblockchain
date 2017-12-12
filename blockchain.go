package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/fatih/color"
)

//BlockChain is a struct to contain an instance of the blockchain
type BlockChain struct {
	CurrentTransactions []Transaction
	Chain               []Block
	Nodes               NodeSet
	nowMining           *sync.Mutex
	minerWaiting        bool
}

func (bc BlockChain) start() BlockChain {
	bc.newBlock(100, "1")
	bc.Nodes = NodeSet{make(map[Node]bool)}
	bc.minerWaiting = false
	bc.nowMining = &sync.Mutex{}
	return bc
}

func (bc BlockChain) proofOfWork(lastProof int) int {
	defer trackTime(time.Now(), "Proof of Work")
	nworkers := 250
	bs := BlockSolver{}
	bs.init(nworkers, lastProof)
	index := 1
	for i := 1; i <= nworkers; i++ {
		bs.proofQueue <- index
		go bs.proofOfWorker(i)
		index++
	}
	proof := <-bs.proofChan
	fmt.Printf("Found solution at proof = %d \n", proof)
	return proof
}

//Only mine if: another miner is not already queued; vote is queued and 5 minutes has passed; 5 or more transactions waiting
//Defer mining if: above and miner is already already working
func (bc *BlockChain) tryToMine() {
	if !bc.minerWaiting {
		fmt.Println("Queueing A Miner")
		bc.minerWaiting = true
		bc.mineBlock()
	} else {
		fmt.Println("Skipping Miner, One Already Queued")
	}
}
func (bc *BlockChain) mineBlock() {
	bc.nowMining.Lock()
	fmt.Println("Starting Miner")
	bc.minerWaiting = false
	lastBlock := bc.lastBlock()
	lastProof := lastBlock.Proof
	proof := bc.proofOfWork(lastProof)
	previousHash := lastBlock.hash()
	bc.newBlock(proof, previousHash)
	bc.nowMining.Unlock()
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
		if !validProof(prevBlock.Proof, block.Proof) {
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
