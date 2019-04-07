package main

import (
	"log"
	"os"

	"github.com/noborus/tbln"
)

func main() {
	var err error
	tb := tbln.NewTbln()
	tb.SetTableName("sample")
	// SetNames sets column names
	err = tb.SetNames([]string{"id", "name"})
	if err != nil {
		log.Fatal(err)
	}
	// SetTypes sets the column type
	err = tb.SetTypes([]string{"int", "text"})
	if err != nil {
		log.Fatal(err)
	}
	// Add a row.
	// The number of columns should be the same
	// number of columns in Names and Types.
	err = tb.AddRows([]string{"1", "Bob"})
	if err != nil {
		log.Fatal(err)
	}
	err = tb.AddRows([]string{"2", "Alice"})
	if err != nil {
		log.Fatal(err)
	}
	tbln.WriteAll(os.Stdout, tb)
}
