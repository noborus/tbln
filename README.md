# TBLN

[![GoDoc](https://godoc.org/github.com/noborus/tbln?status.svg)](https://godoc.org/github.com/noborus/tbln)
[![Build Status](https://travis-ci.org/noborus/tbln.svg?branch=master)](https://travis-ci.org/noborus/tbln)
[![Go Report Card](https://goreportcard.com/badge/github.com/noborus/tbln)](https://goreportcard.com/report/github.com/noborus/tbln)
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fnoborus%2Ftbln.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fnoborus%2Ftbln?ref=badge_shield)

[TBLN](https://tbln.dev) is a text format that represents the table.

This repository contains Go library for reading and writing files,
and Go library for reading and writing [RDBMS](db/README.md) tables.
There is also a [CLI](cmd/README.md) tool that uses them.

## Features

* TBLN can contain multiple columns like CSV.
* Database tables and import/export are possible.
* It can include the signature with the hash value in the file.

Please refer to [TBLN](https://tbln.dev/) for the specification of TBLN.

## TBLN file example

Data only.

```
| 1 | Bob |
| 2 | Alice |
```

Add comment, column name and data type.

```
# comment
; name: | id | name |
; type: | int | text |
| 1 | Bob |
| 2 | Alice |
```

Database export and signature.

```
; TableName: simple
; character_octet_length: |  | 1073741824 |
; created_at: 2019-03-12T15:41:42+09:00
; is_nullable: | YES | YES |
; numeric_precision: | 32 |  |
; numeric_precision_radix: | 2 |  |
; numeric_scale: | 0 |  |
; postgres_type: | integer | text |
; Signature: | test | ED25519 | dfe0077a4baa689dec15365642de8d736b30b678fc4b6725acf25cd760528ed365dc18855a11fc4473ca0a2d36499819de95caba3ac44937ac7c04465e7af901 |
; Hash: | sha256 | 3191722649a6388498c435e411cb6534b740d9b3a5c7ac281dd824b4ba78e968 |
; name: | id | name |
; type: | int | text |
| 1 | Bob |
| 2 | Alice |
```

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fnoborus%2Ftbln.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fnoborus%2Ftbln?ref=badge_large)
