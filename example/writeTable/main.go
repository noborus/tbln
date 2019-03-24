package main

import (
	"bytes"
	"log"
	"os"

	"github.com/noborus/tbln"
	"github.com/noborus/tbln/db"

	_ "github.com/noborus/tbln/db/mysql"
	_ "github.com/noborus/tbln/db/postgres"
	_ "github.com/noborus/tbln/db/sqlite3"
)

func main() {
	conn, err := db.Open("postgres", "")
	//	conn, err := db.Open("mysql", "root:@/noborus")
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
	err = db.WriteTable(conn, at, "", db.IfNotExists, db.Normal)
	if err != nil {
		log.Fatal(err)
	}
}
