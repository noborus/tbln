package main

import (
	"bytes"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/noborus/tbln"
)

func main() {
	db, err := tbln.DBOpen("postgres", "")
	//	db, err := tbln.DBOpen("mysql", "root:@/noborus")
	if err != nil {
		log.Fatal(err)
	}

	data := `; name: | i||d | name | age |
; pg_type: | int | varchar(40) | int |
; type: | int | text | int |
; TableName: geh1
| 1 | Bob | 19 |
| 2 | Alice | 14 |
| 3 | He||nry | 18 |
`

	r := bytes.NewBufferString(data)
	at, err := tbln.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	at.SetTableName(os.Args[1])
	err = tbln.WriteTable(db, at, true)
	if err != nil {
		log.Fatal(err)
	}
}
