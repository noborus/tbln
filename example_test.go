package tbln_test

import (
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
