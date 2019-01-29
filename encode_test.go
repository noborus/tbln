package tbln

import (
	"bytes"
	"testing"
)

func TestEncode(t *testing.T) {
	data := [][]string{
		{"1", "Bob", "19"},
		{"2", "Alice", "14"},
	}
	dataw := []byte("| 1 | Bob | 19 |\n| 2 | Alice | 14 |\n")

	tbl := NewTbln(2)
	for _, row := range data {
		tbl.AddRow(row)
	}
	b := &bytes.Buffer{}
	err := tbl.Encode(b)
	if err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}
	if bytes.Compare(b.Bytes(), dataw) != 0 {
		t.Errorf("out=[\n%v] ,should=[\n%v]", b, string(dataw))
	}
}

func TestEscapeEncode(t *testing.T) {
	data := [][]string{
		{"a|b", "a | b", "abc ||| def"},
		{"|", " | ", " |"},
	}
	dataw := []byte("| a||b | a || b | abc |||| def |\n| || |  ||  |  || |\n")
	tbl := NewTbln(2)
	for _, row := range data {
		tbl.AddRow(row)
	}
	b := &bytes.Buffer{}
	err := tbl.Encode(b)
	if err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}
	if bytes.Compare(b.Bytes(), dataw) != 0 {
		t.Errorf("out=[\n%v] ,should=[\n%v]", b, string(dataw))
	}
}
