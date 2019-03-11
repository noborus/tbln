# TBLN

TBLN is a text format that represents the table.

This repository has Go library and [CLI](cmd/README.md) tool which can read/write file and [RDBMS](db/README.md).

## Features

TBLN can contain multiple columns like CSV.
Also, it can include checksum and signature inside.

## Specification

* TBLN contains three types of lines: data, comments, and extras.
* All rows end with a line break.
* The number of columns in all rows must be the same.

### data

```
| column1 | column2 | column3 |
```

* data begins with "| "(vertical bar + space)  and ends with " |"(space + vertical bar).
* Multiple columns are separated by " | "(space + vertical bar + space).
* White space is considered part of a column.
* If "|" is included in the column, "|" must be duplicated.
* Otherwise, all values are taken.

```
| -> || , || -> |||
```

### Extras

```s
; ItemName: Value
````

* Extras begin with "; ". Extras can be interpreted as a header.
* Extras can be used to indicate the column name and column data type.
* Extras is basically written in the item name: value.
* Extras has a predefined item name.

#### Predefined item name in extras.

| Item name | detail |
|:----------|:--------|
| TableName | table name |
| name      | column name |
| type      | column type |
| Hash      | data and extras checksum hash |
| Signature | signature for hash |

#### The order of Extras

The order of TBLN is as follows.
1. Comments
2. extras hash not target
3. signature
4. hash value
5. extras hash target
6. data

![simple-tbln](https://user-images.githubusercontent.com/2296563/54079389-0ba63580-431f-11e9-8c21-2ce39aeee4e3.png)

The target of hash is the line below Hash.
The signature targets the Hash value.

Hash currently supports SHA 256 and SHA 512.
Signature currently supports ED 25519.

### Comments

```
# Comments
```

* Comments begin with "#".
* Comments are not interpreted.

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

Database export

```
; TableName: simple
; character_octet_length: |  | 1073741824 |
; created_at: 2019-03-10T09:56:06+09:00
; is_nullable: | YES | YES |
; numeric_precision: | 32 |  |
; numeric_precision_radix: | 2 |  |
; numeric_scale: | 0 |  |
; postgres_type: | integer | text |
; Signature: | ED25519:6a693478f1aa5ff91847e2b4e7ec633358dc2f9561454791289cb1c644ef070e37e089a37e8b11324f32f12c439fd2bd6c802144ebf2686df04811455573dd05 |
; Hash: | sha256:3191722649a6388498c435e411cb6534b740d9b3a5c7ac281dd824b4ba78e968 |
; name: | id | name |
; type: | int | text |
| 1 | Bob |
| 2 | Alice |
```
