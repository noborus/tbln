# tbln/db

[![GoDoc](https://godoc.org/github.com/noborus/tbln/db?status.svg)](https://godoc.org/github.com/noborus/tbln/db)

This library import/export database table.

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

        data := `; name: | i||d | name | age |
; pg_type: | int | varchar(40) | int |
; type: | int | text | int |
; TableName: geh1
| 1 | Bob | 19 |
| 2 | Alice | 14 |
| 3 | He||nry | 18 |
`

        r := bytes.NewBufferString(data)
        at, err := tbln.ReadAll(r)
        if err != nil {
                log.Fatal(err)
        }
        at.SetTableName(os.Args[1])
        err = db.WriteTable(conn, at, "", db.IfNotExists, db.Normal)
        if err != nil {
                log.Fatal(err)
        }
}
```
