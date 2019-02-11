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

// Extra is table definition extra struct.
type Extra struct {
	hashing bool
	value   interface{}
}

// Definition is common table definition struct.
type Definition struct {
	columnNum int
	tableName string
	Comments  []string
	Names     []string
	Types     []string
	Ext       map[string]Extra
}

// NewDefinition is create Definition struct.
func NewDefinition(name string) Definition {
	return Definition{
		tableName: name,
		Ext:       make(map[string]Extra),
	}
}

// SetTableName is set Table Name of the Table.
func (d *Definition) SetTableName(name string) error {
	d.tableName = name
	return nil
}

// SetNames is set Column Name to the Table.
func (d *Definition) SetNames(names []string) error {
	if err := d.setColNum(len(names)); err != nil {
		return err
	}
	d.Names = names
	return nil
}

// SetTypes is set Column Type to Table.
func (d *Definition) SetTypes(types []string) error {
	if err := d.setColNum(len(types)); err != nil {
		return err
	}
	d.Types = types
	return nil
}

func (d *Definition) setColNum(colNum int) error {
	if d.columnNum == 0 {
		d.columnNum = colNum
		return nil
	}
	if colNum != d.columnNum {
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
