package main

import (
	"fmt"
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
	// conn, err := db.Open("mysql", "root:@/noborus")
	// conn, err := db.Open("sqlite3", "sqlite_file")
	if err != nil {
		log.Fatal(err)
	}
	at, err := db.ReadTableAll(conn, os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	comment := fmt.Sprintf("DB:%s\tTable:%s", conn.Name, os.Args[1])
	at.Comments = []string{comment}
	_, err = at.SumHash(tbln.SHA256)
	if err != nil {
		log.Fatal(err)
	}
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		log.Fatal(err)
	}
}
