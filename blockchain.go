package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strconv"
	"time"
)

//BlockChain is a struct to contain an instance of the blockchain
type BlockChain struct {
	CurrentTransactions []Transaction
	Chain               []Block
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
func shaHash(plaintext string) string {
	hasher := sha256.New()
	hasher.Write([]byte(plaintext))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}
