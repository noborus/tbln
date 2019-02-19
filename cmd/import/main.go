package main

import (
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/noborus/tbln"
)

func main() {
	db, err := tbln.DBOpen("postgres", "")
	//db, err := tbln.DBOpen("mysql", "root:@/noborus")
	// db, err := tbln.DBOpen("sqlite3", "sqlite_file")

	if err != nil {
		log.Fatal(err)
	}
	db.Tx, err = db.DB.Begin()
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
	at.SetTableName("language2")
	err = tbln.WriteTable(db, at, true)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
}
