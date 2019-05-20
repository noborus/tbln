package main

import (
	"log"
	"os"

	"github.com/noborus/tbln"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("Requires tbln file")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	at, err := tbln.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	at.SetTableName("newtable")
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		log.Fatal(err)
	}
}
