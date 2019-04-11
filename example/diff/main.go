package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/noborus/tbln"
)

var data1 = `# String example
; TableName: geh1
; primarykey: | id |
; pg_type: | int | varchar(40) | int |
; name: | id | name | age |
; type: | int | text | int |
| 1 | Bob | 19 |
| 2 | Alice | 14 |
| 3 | Henry | 19 |
| 5 | hoge | 99 |
| 9 | noborus | 46 |
| 11 | ghoge | 4 |
| 21 | aghoge | 4 |
`

var data2 = `# String example
; TableName: newgeh1
; primarykey: | id |
; pg_type: | int | varchar(40) | int |
; name: | id | name | age |
; type: | int | text | int |
| 1 | Beob | 19 |
| 2 | Alice | 14 |
| 3 | Henry | 19 |
| 4 | noborus | 46 |
| 9 | noborus | 46 |
| 10 | hoge | 46 |
| 11 | ghoge | 5 |
`

func main() {
	r := bytes.NewBufferString(data1)
	src := tbln.NewReader(r)
	r2 := bytes.NewBufferString(data2)
	dst := tbln.NewReader(r2)

	d, err := tbln.NewCompare(src, dst)
	if err != nil {
		log.Fatal(err)
	}
	for {
		dd, err := d.ReadDiffRow()
		if err != nil {
			break
		}
		fmt.Printf("%s\n", dd.Diff())
	}
}
