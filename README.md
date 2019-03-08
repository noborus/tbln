# tbln

tbln is a text format that represents the table.

This repository has Go library and CLI tool which can read/write file and RDBMS.

## Specification

1. One line ends with a line feed.
2. Value consists of multiple columns.
3. Lines begin with "| " (vertical bar + space) and end with " |" (space + vertical bar).
4. Values are separated by " | "  (space + vertical bar + space).
5. If you want to include "|" in the value, it is "||". Increase | after | one by one.
6. Otherwise, all values are taken.

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
