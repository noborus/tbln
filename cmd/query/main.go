package main

import (
	"os"

	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/noborus/tbln"
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
