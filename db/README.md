# tbln/db

[![GoDoc](https://godoc.org/github.com/noborus/tbln/db?status.svg)](https://godoc.org/github.com/noborus/tbln/db)

This library import/export database table.
Also, it can be merged into a table.

## Supported database

## Support DB version.

| DataBase | (since)Version |
|:-----------|:-----------|
| PostgreSQL | 9.5 |
| MySQL | 5.7 |
| SQlite3 |   |

Import requires tbln/db/(database) instead of go SQL database drivers.

```
import (
        "github.com/noborus/tbln"
        "github.com/noborus/tbln/db"
        // MySQL
        _ "github.com/noborus/tbln/db/mysql"
        // PostgreSQL
        _ "github.com/noborus/tbln/db/postgres"
        // SQLite3
        _ "github.com/noborus/tbln/db/sqlite3"
)
```

## Examples

### table export

``` go
package main

import (
        "fmt"
        "log"
        "os"

        "github.com/noborus/tbln"
        "github.com/noborus/tbln/db"
        _ "github.com/noborus/tbln/db/postgres"
)

func main() {
        conn, err := db.Open("postgres", "")
        if err != nil {
                log.Fatal(err)
        }
        at, err := db.ReadTableAll(conn, "", os.Args[1])
        if err != nil {
                log.Fatal(err)
        }
        comment := fmt.Sprintf("DB:%s\tTable:%s", conn.Name, os.Args[1])
        at.Comments = []string{comment}
        err = at.SumHash(tbln.SHA256)
        if err != nil {
                log.Fatal(err)
        }
        err = tbln.WriteAll(os.Stdout, at)
        if err != nil {
                log.Fatal(err)
        }
}
```

### table Import

``` go
package main

import (
	"bytes"
	"log"

	"github.com/noborus/tbln"
	"github.com/noborus/tbln/db"
	_ "github.com/noborus/tbln/db/postgres"
)

func main() {
	conn, err := db.Open("postgres", "")
	if err != nil {
		log.Fatal(err)
	}

	data := `# Simple example
; TableName: simple
; primarykey: | id |
; pg_type: | int | varchar(40) |
; name: | id | name |
; type: | int | text |
| 1 | Bob |
| 2 | Alice |
`

	r := bytes.NewBufferString(data)
	tb, err := tbln.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	err = db.WriteTable(conn, tb, "", db.IfNotExists, db.Normal)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Commit()
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}
```

### table merge

Execute table merge after execution in the table import example.
Same result no matter how many times table merge is executed.
 
``` go
package main

import (
	"bytes"
	"log"
	"os"

	"github.com/noborus/tbln"
	"github.com/noborus/tbln/db"
	_ "github.com/noborus/tbln/db/postgres"
)

var dataMerge = `# Simple example
; TableName: simple
; primarykey: | id |
; pg_type: | int | varchar(40) |
; name: | id | name |
; type: | int | text |
| 1 | Booooob |
| 2 | Alice |
| 3 | Carol |
| 4 | Dave |
`

func main() {
	var err error
	conn, err := db.Open("postgres", "")
	if err != nil {
		log.Fatal(err)
	}

	r := bytes.NewBufferString(dataMerge)
	other, err := tbln.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.MergeTableTbln("", "simple", other, true)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Commit()
	if err != nil {
		log.Fatal(err)
	}

	// Display and confirm
	at, err := db.ReadTableAll(conn, "", "simple")
	if err != nil {
		log.Fatal(err)
	}
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Close()
	if err != nil {
		log.Fatal(err)
	}
}
```