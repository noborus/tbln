package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/noborus/tbln"

	_ "github.com/noborus/tbln/db/mysql"
	_ "github.com/noborus/tbln/db/postgres"
	_ "github.com/noborus/tbln/db/sqlite3"
)

var data1 = `# String example
; TableName: geh1
; primarykey: | id |
; pg_type: | int | varchar(40) | int |
; name: | id | name | age |
; type: | int | text | int |
| 1 | Bob | 19 |
| 2 | Alice | 14 |
| 3 | Henry | 19 |
| 5 | hoge | 99 |
| 9 | noborus | 46 |
| 11 | ghoge | 4 |
| 21 | aghoge | 4 |
`

var data2 = `# String example
; TableName: newgeh1
; primarykey: | id |
; pg_type: | int | varchar(40) | int |
; name: | id | name | age |
; type: | int | text | int |
| 1 | Beob | 19 |
| 2 | Alice | 14 |
| 3 | Henry | 19 |
| 4 | noborus | 46 |
| 9 | noborus | 46 |
| 10 | hoge | 46 |
| 11 | ghoge | 5 |
`

/*
var data3 = `# Simple Add
; TableName: simple
; primarykey: | id |
; pg_type: | int | varchar(40) |
; name: | id | name |
; type: | int | text |
| 1 | Bob |
| 2 | Alice |
| 3 | Henry |
`
*/

func main() {
	r := bytes.NewBufferString(data1)
	at, err := tbln.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	rr := tbln.NewSelfReader(at)
	/*
		conn, err := db.Open("postgres", "")
		// conn, err := db.Open("mysql", "root:@/noborus")
		// conn, err := db.Open("sqlite3", "sqlite_file")
		if err != nil {
			log.Fatal(err)
		}
		rr, err := conn.ReadTable("", os.Args[1], nil)
		if err != nil {
			log.Fatal(err)
		}
	*/
	r2 := bytes.NewBufferString(data2)
	dst := tbln.NewReader(r2)

	d, err := tbln.NewCompare(rr, dst)
	if err != nil {
		log.Fatal(err)
	}
	for {
		dd, err := d.ReadDiffRow()
		if err != nil {
			break
		}
		fmt.Printf("%s\n", dd.Diff())
	}
}
