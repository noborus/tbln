# TBLN

[![GoDoc](https://godoc.org/github.com/noborus/tbln?status.svg)](https://godoc.org/github.com/noborus/tbln)
[![Build Status](https://travis-ci.org/noborus/tbln.svg?branch=master)](https://travis-ci.org/noborus/tbln)
[![Go Report Card](https://goreportcard.com/badge/github.com/noborus/tbln)](https://goreportcard.com/report/github.com/noborus/tbln)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fnoborus%2Ftbln.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fnoborus%2Ftbln?ref=badge_shield)

[TBLN](https://tbln.dev) is a text format that represents the table.

This repository contains Go library for reading and writing files,
and Go library for reading and writing [RDBMS](db/README.md) tables.

Here is a document about the library.
See [CLI](cmd/README.md) for CLI Tool.

## Features

* TBLN can contain multiple columns like CSV.
* Database tables and import/export are possible.
* It can include the signature with the hash value in the file.
* Merge, sync and diff is possible
    * TBLN file and TBLN file
    * TBLN file and DB Table
    * DB Table and DB Table

Please refer to [TBLN](https://tbln.dev/) for the specification of TBLN.

## Install

```bash
$ go get github.com/noborus/tbln
```

Notes: Requires a version of Go that supports modules. e.g. Go 1.11+

## Example

### Write example

Example of creating TBLN.
(Error check is omitted)

```go
package main

import (
	"os"

	"github.com/noborus/tbln"
)

func main() {
	tb := tbln.NewTBLN()
	tb.SetTableName("sample")
	tb.SetNames([]string{"id", "name"})
	tb.SetTypes([]string{"int", "text"})
	tb.AddRows([]string{"1", "Bob"})
	tb.AddRows([]string{"2", "Alice"})
	tbln.WriteAll(os.Stdout, tb)
}
```

Execution result.

```tbln
; TableName: sample
; created_at: 2019-04-06T01:05:17+09:00
; name: | id | name |
; type: | int | text |
| 1 | Bob |
| 1 | Alice |
```

### Read example

All read to memory by tbln.ReadAll.
Here, the table name is rewritten and output.

```go
package main

import (
	"log"
	"os"

	"github.com/noborus/tbln"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("Requires tbln file")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	at, err := tbln.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	at.SetTableName("newtable")
	err = tbln.WriteAll(os.Stdout, at)
	if err != nil {
		log.Fatal(err)
	}
}
```

### TBLN files can be imported and exported into the database.


Example of importing into PostgreSQL.

```
# \d sample
               Table "public.sample"
 Column |  Type   | Collation | Nullable | Default
--------+---------+-----------+----------+---------
 id     | integer |           |          |
 name   | text    |           |          |
```

Example of exporting it to a TBLN file.

```
; TableName: sample
; character_octet_length: |  | 1073741824 |
; created_at: 2019-04-06T02:03:43+09:00
; is_nullable: | YES | YES |
; numeric_precision: | 32 |  |
; numeric_precision_radix: | 2 |  |
; numeric_scale: | 0 |  |
; postgres_type: | integer | text |
; Signature: | test | ED25519 | 6271909d82000c4f686785cf0f9080971470ad3247b091ca50f6ea12ccc96efde0e1ca77e16723ef0f9d781941dfb92bed094dbf2e4079dd25f5aa9f9f1aab01 |
; Hash: | sha256 | 65f7ce4e15ddc006153fe769b8f328c466cbd1dea4b15aa195ed63daf093668d |
; name: | id | name |
; type: | int | text |
| 1 | Bob |
| 1 | Alice |
```

See [db/README.md](db/README.md) for details.

You can also use tbln cli tool to quickly import and export to a database.

See [cmd/README.md](cmd/README.md) for details.

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fnoborus%2Ftbln.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fnoborus%2Ftbln?ref=badge_large)
