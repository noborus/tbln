package tbln

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
)

// Table struct is table Definition + Table rows.
type Table struct {
	Definition
	Hash   map[string]string
	buffer bytes.Buffer
	RowNum int
	Rows   [][]string
}

// NewTable is create table struct.
func NewTable(d Definition) *Table {
	return &Table{
		Definition: d,
	}
}

// AddRows is Add row to Table.
func (t *Table) AddRows(row []string) error {
	var err error
	t.columnNum, err = checkRow(t.columnNum, row)
	if err != nil {
		return err
	}
	t.Rows = append(t.Rows, row)
	t.RowNum++
	return nil
}

func checkRow(columnNum int, row []string) (int, error) {
	if columnNum == 0 {
		columnNum = len(row)
	} else {
		if len(row) != columnNum {
			return columnNum, fmt.Errorf("Error: invalid column num (%d!=%d) %s", columnNum, len(row), row)
		}
	}
	return columnNum, nil
}

// SumHash is returns the calculated checksum.
// Checksum target is exported to buffer and saved.
func (t *Table) SumHash() (map[string]string, error) {
	if t.Hash == nil {
		t.Hash = make(map[string]string)
	}
	writer := bufio.NewWriter(&t.buffer)
	bw := NewWriter(writer)
	err := bw.writeExtraWithHash(t.Definition)
	if err != nil {
		return nil, err
	}
	for _, row := range t.Rows {
		err := bw.WriteRow(row)
		if err != nil {
			return nil, err
		}
	}
	err = writer.Flush()
	if err != nil {
		return nil, err
	}
	sum := sha256.Sum256(t.buffer.Bytes())
	t.Hash["sha256"] = fmt.Sprintf("%x", sum)
	return t.Hash, nil
}
