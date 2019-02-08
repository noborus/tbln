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
	WriteDefinition(Definition) error
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

// NewDefinition is create Definition struct.
func NewDefinition(name string) Definition {
	return Definition{
		name: name,
	}
}

// SetTableName is set Table Name of the Table.
func (t *Definition) SetTableName(name string) error {
	t.name = name
	return nil
}

// SetNames is set Column Name to the Table.
func (t *Definition) SetNames(names []string) error {
	if err := t.setColNum(len(names)); err != nil {
		return err
	}
	t.Names = names
	return nil
}

// SetTypes is set Column Type to Table.
func (t *Definition) SetTypes(types []string) error {
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

// Table struct is table Definition + Table rows.
type Table struct {
	Definition
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
