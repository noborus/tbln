package main

import (
	"os"

	"log"

	"github.com/noborus/tbln"
	_ "github.com/noborus/tbln/db/mysql"
	_ "github.com/noborus/tbln/db/postgres"
	_ "github.com/noborus/tbln/db/sqlite3"
)

func main() {
	db, err := tbln.DBOpen("postgres", "")
	// db, err := tbln.DBOpen("mysql", "root:@/noborus")
	if err != nil {
		log.Fatal(err)
	}
	at, err := tbln.ReadQueryAll(db, "SELECT * FROM t40 WHERE id=$1", os.Args[1])
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
