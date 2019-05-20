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

	prFile, err := os.Open("test.key")
	if err != nil {
		log.Fatal(err)
	}
	defer prFile.Close()
	fi, _ := prFile.Stat()
	size := fi.Size()
	data := make([]byte, size)
	_, err = prFile.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	privateKey, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		log.Fatal(err)
	}
	_, err = at.Sign("test", privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		log.Fatal(err)
	}
}
