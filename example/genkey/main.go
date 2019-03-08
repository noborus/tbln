package main

import (
	"bufio"
	"encoding/base64"
	"log"
	"os"

	"golang.org/x/crypto/ed25519"
)

func main() {
	pubKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		log.Fatal(err)
	}
	name := os.Args[1]

	privateFile, err := os.Create(name + ".key")
	if err != nil {
		log.Fatal(err)
	}
	buf := bufio.NewWriter(privateFile)
	ps := base64.StdEncoding.EncodeToString(privateKey)
	_, err = buf.WriteString(ps)
	if err != nil {
		log.Fatal(err)
	}
	buf.Flush()
	privateFile.Close()

	pubFile, err := os.Create(name + ".pub")
	if err != nil {
		log.Fatal(err)
	}
	buf2 := bufio.NewWriter(pubFile)
	ps2 := base64.StdEncoding.EncodeToString(pubKey)
	_, err = buf2.WriteString(ps2)
	if err != nil {
		log.Fatal(err)
	}
	buf2.Flush()
	pubFile.Close()
}
