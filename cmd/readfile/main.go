package main

import (
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
	rh := at.Hash["sha256"]
	h, err := at.SumHash()
	if err != nil {
		log.Fatal(err)
	}
	if len(rh) > 0 && rh != h["sha256"] {
		log.Printf("checksum error")
	}
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		log.Fatal(err)
	}
}
