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
	value      interface{}
	hashTarget bool
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
func NewDefinition() Definition {
	return Definition{
		Ext: make(map[string]Extra),
	}
}

// SetTableName is set Table Name of the Table.
func (d *Definition) SetTableName(name string) {
	d.tableName = name
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
