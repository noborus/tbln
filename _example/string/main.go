package main

import (
	"bytes"
	"log"
	"os"

	"github.com/noborus/tbln"
)

func main() {
	data := `# String example
; name: | id | name | age |
; type: | int | text | int |
; pg_type: | int | varchar(40) | int |
; TableName: geh1
| 1 | Bob | 19 |
| 2 | Alice | 14 |
| 3 | Henry | 19 |
`
	r := bytes.NewBufferString(data)
	at, err := tbln.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		log.Fatal(err)
	}
}
