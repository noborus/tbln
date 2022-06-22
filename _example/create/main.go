package main

import (
	"log"
	"os"

	"github.com/noborus/tbln"
)

func main() {
	tb := tbln.NewTBLN()
	tb.SetTableName("sample")
	// SetNames sets column names
	if err := tb.SetNames([]string{"id", "name"}); err != nil {
		log.Fatal(err)
	}
	// SetTypes sets the column type
	if err := tb.SetTypes([]string{"int", "text"}); err != nil {
		log.Fatal(err)
	}
	// Add a row.
	// The number of columns should be the same
	// number of columns in Names and Types.
	if err := tb.AddRows([]string{"1", "Bob"}); err != nil {
		log.Fatal(err)
	}
	if err := tb.AddRows([]string{"2", "Alice"}); err != nil {
		log.Fatal(err)
	}
	if err := tbln.WriteAll(os.Stdout, tb); err != nil {
		log.Fatal(err)
	}
}
