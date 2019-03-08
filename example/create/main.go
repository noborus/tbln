package main

import (
	"log"
	"os"

	"github.com/noborus/tbln"
)

func main() {
	t := tbln.NewTbln()
	t.SetTableName("dummy")
	err := t.AddRows([]string{"yes", "no"})
	if err != nil {
		log.Fatal(err)
	}
	err = t.SetNames([]string{"c1", "c2"})
	if err != nil {
		log.Fatal(err)
	}
	err = tbln.WriteAll(os.Stdout, t)
	if err != nil {
		log.Fatal(err)
	}

}
