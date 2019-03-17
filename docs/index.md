## Overview

TBLN can contain multiple columns like CSV.
Also, it can include checksum and signature inside.

## Products

* [tbln](https://github.com/noborus/tbln) Go library for reading and writing files.
* [tbln/db](https://github.com/noborus/tbln/db) Go library for reading and writing RDBMS.
* [tbln cli tool](https://github.com/noborus/tbln/cmd ) Import/Export TBLN file and RDBMS table.

## Specification

* TBLN contains three types of lines: data, comments, and extras.
* All rows end with a line break.
* The number of columns in all rows must be the same.

### data

```
| column1 | column2 | column3 |
```

* data begins with "\| "(vertical bar + space)  and ends with " \|"(space + vertical bar).
* Multiple columns are separated by " \| "(space + vertical bar + space).
* White space is considered part of a column.
* If "\|" is included in the column, "\|" must be duplicated.
* Otherwise, all values are taken.

```
| -> || , || -> |||
```

### Extras

```
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

#### TableName

TableName is used as a table name when importing into the database.

#### name

name is the column name.
Same as data, it is written in the form of | name1 | name2 | ... |.

#### type

type is the data type of the column.
It is expressed in the form of | int | text | ... |.

#### Hash

Hash is a Hash value of SHA256 or SHA512.
The hash value of Extras below the data and Hash items.
Comments and items above the Extras hash are outside the Hash calculation.
Therefore, the hash value does not change even if you change it.

#### Signature

Signature is a signature of ED25519 format for Hash value.

#### The order of Extras

The order of TBLN is as follows.
1. Comments
2. extras hash not target
3. signature
4. hash value
5. extras hash target
6. data

The target of hash is the line below Hash.
The signature targets the Hash value.

Hash currently supports SHA256 and SHA512.
Signature currently supports ED25519.

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
