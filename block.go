package main

import (
	"fmt"
)

//Block is a struct to contain a blockchain transaction object
type Block struct {
	Index        int
	Timestamp    int64
	Transactions []Transaction
	Proof        int
	PreviousHash string
}

func (b Block) hash() string {
	blockstr := fmt.Sprint(b)
	return shaHash(blockstr)
}
