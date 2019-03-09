# tbln

Tbln is a text format that represents the table.

This repository has Go library and CLI tool which can read/write file and RDBMS.

## Features

Tbln can contain multiple fields like csv.
Also, it can include checksum and signature inside.

## Specification

* Tbln contains three types of lines: data, comments, and extras.
* All rows end with a new line(LF).
* The number of fields in all rows must be the same.

### data

```
| fields1 | fields2 | fields3 |
```

* data begins with "| "(vertical bar + space)  and ends with " |"(space + vertical bar).
* Multiple fields are separated by " | ".
* White space is considered part of a field.
* If "|" is included in the field, "|" must be duplicated.
* Otherwise, all values are taken.

```
| -> || , || -> |||
```

### Comments

```
# Comments
```

* Comments begin with "# ".
* Comments are not interpreted.

### Extras

```
; ItemName: Value
````

* Extras begin with ";". Extras can be interpreted as a header.
* Extras is basically written in the item name: value.
* Extras has item names that are interpreted in some special ways.

## TBLN format example

Simple and only data.

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

```
# comment
; name: | id | name |
; type: | int | text |
| 1 | B||o||b |
| 2 | Alice |
```
