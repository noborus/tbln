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

var data1 = `# String example
; TableName: geh1
; primarykey: | id |
; pg_type: | int | varchar(40) |
; name: | id | name |
; type: | int | text |
| 1 | Bob |
| 2 | Alice |
| 3 | Hanry |
| 4 | hohog |
`

func main() {
	var err error
	conn, err := db.Open("postgres", "")
	if err != nil {
		log.Fatal(err)
	}

	r := bytes.NewBufferString(data1)
	dr := tbln.NewReader(r)
	err = conn.MergeTable("", "simple", dr, true)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Commit()
	if err != nil {
		log.Fatal(err)
	}

	at, err := db.ReadTableAll(conn, "", "simple")
	if err != nil {
		log.Fatal(err)
	}
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		log.Fatal(err)
	}
}
