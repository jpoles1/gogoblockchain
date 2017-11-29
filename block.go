package main

import (
	"crypto/sha256"
	"encoding/base64"
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

func shaHash(plaintext string) string {
	hasher := sha256.New()
	hasher.Write([]byte(plaintext))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}
func (b Block) hash() string {
	blockstr := fmt.Sprint(b)
	return shaHash(blockstr)
}
