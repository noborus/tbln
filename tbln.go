package tbln

import (
	"fmt"
	"regexp"
)

// Read is TBLN Read interface.
type Read interface {
	ReadRow() ([]string, error)
}

// Write is TBLN Write interface.
type Write interface {
	WriteInfo(Definition) error
	WriteRow([]string) error
}

// ESCREP is escape | -> ||
var ESCREP = regexp.MustCompile(`(\|+)`)

// UNESCREP is unescape || -> |
var UNESCREP = regexp.MustCompile(`\|(\|+)`)

// Definition is common table definition struct.
type Definition struct {
	columnNum int
	name      string
	Names     []string
	Types     []string
	Comments  []string
	Ext       map[string]string
}

// Table struct is table Definition + Table all rows.
type Table struct {
	Definition
	RowNum int
	Rows   [][]string
}

// SetTableName is set Table Name of the Table.
func (t *Definition) SetTableName(name string) error {
	t.name = name
	return nil
}

// setNames is set Column Name to the Table.
func (t *Definition) setNames(names []string) error {
	if err := t.setColNum(len(names)); err != nil {
		return err
	}
	t.Names = names
	return nil
}

// setTypes is set Column Type to Table.
func (t *Definition) setTypes(types []string) error {
	if err := t.setColNum(len(types)); err != nil {
		return err
	}
	t.Types = types
	return nil
}

func (t *Definition) setColNum(colNum int) error {
	if t.columnNum == 0 {
		t.columnNum = colNum
		return nil
	}
	if colNum != t.columnNum {
		return fmt.Errorf("number of columns is different")
	}
	return nil
}
