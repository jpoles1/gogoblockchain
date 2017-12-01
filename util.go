package main

import (
	"crypto/sha256"
	"encoding/base64"
	"net/http"
)

func shaHash(plaintext string) string {
	hasher := sha256.New()
	hasher.Write([]byte(plaintext))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

func destReachable(addr string) bool {
	_, err := http.Get("http://" + addr + "/chain")
	return err == nil
}
