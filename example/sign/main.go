package main

import (
	"encoding/base64"
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
	err = at.SumHash(tbln.SHA256)
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

	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		log.Fatal(err)
	}
}
