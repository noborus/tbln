package main

import (
	"log"
	"os"

	"github.com/noborus/tbln"
	// _ "github.com/noborus/tbln/db/mysql"
	// _ "github.com/noborus/tbln/db/postgres"
	// _ "github.com/noborus/tbln/db/sqlite3"
	_ "github.com/lib/pq"
)

func main() {
	db, err := tbln.DBOpen("postgres", "")
	// db, err := tbln.DBOpen("mysql", "root:@/noborus")
	// db, err := tbln.DBOpen("sqlite3", "sqlite_file")
	if err != nil {
		log.Fatal(err)
	}
	at, err := tbln.ReadTableAll(db, os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	_, err = at.SumHash()
	if err != nil {
		log.Fatal(err)
	}
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		log.Fatal(err)
	}
}
