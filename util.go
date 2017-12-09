package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

func trackTime(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s \n", name, elapsed)
}

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
