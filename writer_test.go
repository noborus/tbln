package tbln

import (
	"bytes"
	"testing"
)

func TestWrite(t *testing.T) {
	data := [][]string{
		{"1", "Bob", "19"},
		{"2", "Alice", "14"},
	}
	dataw := []byte("| 1 | Bob | 19 |\n| 2 | Alice | 14 |\n")
	b := &bytes.Buffer{}
	writer := NewWriter(b)
	for _, row := range data {
		err := writer.Write(row)
		if err != nil {
			t.Errorf("Unexpected error: %s\n", err)
		}
	}
	if bytes.Compare(b.Bytes(), dataw) != 0 {
		t.Errorf("out=[\n%v] ,should=[\n%v]", b, string(dataw))
	}
}
