package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/noborus/tbln"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	at, err := tbln.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	prFile, _ := os.Open("test.key")
	defer prFile.Close()

	fi, _ := prFile.Stat()
	size := fi.Size()
	data := make([]byte, size)
	prFile.Read(data)
	privateKey, _ := base64.StdEncoding.DecodeString(string(data))

	at.Sign(privateKey)

	pubFile, _ := os.Open("test.pub")
	defer pubFile.Close()

	fi2, _ := pubFile.Stat()
	size2 := fi2.Size()
	data2 := make([]byte, size2)
	pubFile.Read(data2)
	pubKey, _ := base64.StdEncoding.DecodeString(string(data2))

	if at.VerifySignature(pubKey) {
		fmt.Println("verify")
	} else {
		fmt.Println("error")
		os.Exit(1)
	}
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		log.Fatal(err)
	}
}
