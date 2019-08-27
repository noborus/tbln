package tbln_test

import (
	"bytes"
	"log"
	"os"
	"strings"

	"github.com/noborus/tbln"
)

func Example() {
	in := `; TableName: simple
; name: | id | name |
; type: | int | text |
| 1 | Bob |
| 2 | Alice |
`
	at, err := tbln.ReadAll(strings.NewReader(in))
	if err != nil {
		log.Fatal(err)
	}
	at.SetTableName("newtable")
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	//; TableName: newtable
	//; name: | id | name |
	//; type: | int | text |
	//| 1 | Bob |
	//| 2 | Alice |
}

func ExampleTBLN() {
	var err error
	tb := tbln.NewTBLN()
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
	err = tbln.WriteAll(os.Stdout, tb)
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	//; TableName: sample
	//; name: | id | name |
	//; type: | int | text |
	//| 1 | Bob |
	//| 2 | Alice |
}

func ExampleDiffAll() {
	TestDiff1 := `; name: | id | name | age |
; type: | int | text | int |
; primarykey: | id |
; TableName: test1
| 1 | Bob | 19 |
`

	TestDiff2 := `; name: | id | name | age |
; type: | int | text | int |
; primarykey: | id |
; TableName: test1
| 1 | Bob | 19 |
| 2 | Alice | 14 |
`
	err := tbln.DiffAll(os.Stdout,
		tbln.NewReader(bytes.NewBufferString(TestDiff1)),
		tbln.NewReader(bytes.NewBufferString(TestDiff2)),
		tbln.AllDiff)
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	//| 1 | Bob | 19 |
	//+| 2 | Alice | 14 |
}
