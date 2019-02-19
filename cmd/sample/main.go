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
	data := `; name: | id | name | age |
; type: | int | text | int |
| 1 | Bob | 19 |
| 2 | Alice | 14 |
| 3 | Henry | 18 |
`

	r := bytes.NewBufferString(data)
	tr := tbln.NewReader(r)
	tr.Comments = []string{"comment"}
	var err error
	var info bool
	var tw *tbln.Writer
	for {
		rec, err := tr.ReadRow()
		if err != nil {
			break
		}
		if !info {
			tw = tbln.NewWriter(os.Stdout)
			err = tw.WriteDefinition(tr.Definition)
			if err != nil {
				log.Println(err)
			}
			info = true
		}
		if rec == nil {
			break
		}
		err = tw.WriteRow(rec)
		if err != nil {
			log.Println(err)
		}
	}
	if err != nil {
		log.Println(err)
	}
	db, err := tbln.DBOpen("mysql", "root:@/noborus")
	if err != nil {
		log.Fatal(err)
	}

	var tw2 *tbln.Writer
	td, err := db.ReadTable("test", nil)
	if err != nil {
		log.Fatal(err)
	}
	tw2 = tbln.NewWriter(os.Stdout)
	err = tw2.WriteDefinition(td.Definition)
	if err != nil {
		log.Println(err)
	}
	for {
		rec, err := td.ReadRow()

		if err != nil {
			log.Println(err)
			break
		}
		if rec == nil {
			break
		}
		err = tw2.WriteRow(rec)
		if err != nil {
			log.Println(err)
			break
		}
	}
	if err != nil {
		log.Println(err)
	}

	var tw3 *tbln.DBWriter
	data2 := `; name: | id | name | age |
; type: | int | text | int |
| 1 | Bob | 19 |
| 2 | Alice | 14 |
| 3 | Henry | 18 |
`
	db, err = tbln.DBOpen("postgres", "")
	if err != nil {
		log.Fatal(err)
	}
	r2 := bytes.NewBufferString(data2)
	tr4 := tbln.NewReader(r2)
	info = false
	for {
		rec, err := tr4.ReadRow()
		if err != nil {
			log.Println(err)
			break
		}
		if !info {
			tw3 = tbln.NewDBWriter(db, tr4.Definition, false)
			tw3.Ph = "$"
			tw3.SetTableName("dummy3")
			err = tw3.WriteDefinition()
			if err != nil {
				log.Fatal(err)
			}
			info = true
		}
		if rec == nil {
			break
		}
		err = tw3.WriteRow(rec)
		if err != nil {
			log.Fatal(err)
		}
	}
	if err != nil {
		log.Println(err)
	}
}
