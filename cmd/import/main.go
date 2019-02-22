package main

import (
	"log"
	"os"

	"github.com/noborus/tbln"
	"github.com/noborus/tbln/db"

	_ "github.com/noborus/tbln/db/mysql"
	_ "github.com/noborus/tbln/db/postgres"
	_ "github.com/noborus/tbln/db/sqlite3"
)

func main() {
	// conn, err := db.Open("postgres", "")
	// conn, err := db.Open("mysql", "root:@/noborus")
	conn, err := db.Open("sqlite3", "sqlite_file")

	if err != nil {
		log.Fatal(err)
	}
	conn.Tx, err = conn.DB.Begin()
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	at, err := tbln.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	err = db.WriteTable(conn, at, true)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
