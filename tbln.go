package tbln

import (
	"fmt"
	"regexp"
)

// TBLN is Read/Writer interface.
type TBLN interface {
	ReadRow() ([]string, error)
	WriteInfo(Table) error
	WriteRow([]string) error
}

// ESCREP is escape | -> ||
var ESCREP = regexp.MustCompile(`(\|+)`)

// UNESCREP is unescape || -> |
var UNESCREP = regexp.MustCompile(`\|(\|+)`)

// Table is common TBLN struct.
type Table struct {
	columnNum int
	name      string
	Names     []string
	Types     []string
	Comments  []string
	Ext       map[string]string
}

// SetTableName is set Table Name of the Table.
func (t *Table) SetTableName(name string) error {
	t.name = name
	return nil
}

// setNames is set Column Name to the Table.
func (t *Table) setNames(names []string) error {
	if err := t.setColNum(len(names)); err != nil {
		return err
	}
	t.Names = names
	return nil
}

// setTypes is set Column Type to Table.
func (t *Table) setTypes(types []string) error {
	if err := t.setColNum(len(types)); err != nil {
		return err
	}
	t.Types = types
	return nil
}

func (t *Table) setColNum(colNum int) error {
	if t.columnNum == 0 {
		t.columnNum = colNum
		return nil
	}
	if colNum != t.columnNum {
		return fmt.Errorf("number of columns is different")
	}
	return nil
}
