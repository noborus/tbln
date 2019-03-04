package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/noborus/tbln"
)

func getPubKey(fileName string) []byte {
	pubFile, _ := os.Open(fileName)
	defer pubFile.Close()
	pubf, _ := ioutil.ReadAll(pubFile)
	pubKey, _ := base64.StdEncoding.DecodeString(string(pubf))
	return pubKey
}

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	at, err := tbln.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	// pubKey := getPubKey("test.pub")
	if at.Verify() {
		fmt.Println("verify")
	} else {
		fmt.Println("verify error")
		os.Exit(1)
	}
}
