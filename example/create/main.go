package main

import (
	"os"

	"github.com/noborus/tbln"
)

func main() {
	tb := tbln.NewTbln()
	tb.SetTableName("sample")
	tb.SetNames([]string{"id", "name"})
	tb.SetTypes([]string{"int", "text"})
	tb.AddRows([]string{"1", "Bob"})
	tb.AddRows([]string{"2", "Alice"})
	tbln.WriteAll(os.Stdout, tb)
}
