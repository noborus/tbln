package tbln

import (
	"fmt"
	"regexp"
)

// Tbln represents a table with rows and columns.
type Tbln struct {
	Rows     [][]string
	Names    []string
	Types    []string
	Comments []string
	Extra    []string
	cap      int
	colNum   int
}

// ESCREP is escape | -> ||
var ESCREP = regexp.MustCompile(`(\|+)`)

// UNESCREP is unescape || -> |
var UNESCREP = regexp.MustCompile(`\|(\|+)`)

// NewTbln makes a model of TBLN.
func NewTbln(rowNum int) *Tbln {
	return &Tbln{
		Rows:     make([][]string, 0, rowNum),
		Names:    make([]string, 0),
		Types:    make([]string, 0),
		Comments: []string{},
		Extra:    make([]string, 0),
		cap:      rowNum,
		colNum:   0,
	}
}

// AddRow is Add a row to Tbln.
func (tbln *Tbln) AddRow(row []string) error {
	if err := tbln.setColNum(len(row)); err != nil {
		return err
	}
	if len(tbln.Rows) >= tbln.cap {
		return fmt.Errorf("over capacity [%d]", tbln.cap)
	}
	tbln.Rows = append(tbln.Rows, row)
	return nil
}

func (tbln *Tbln) setColNum(colNum int) error {
	if tbln.colNum == 0 {
		tbln.colNum = colNum
		return nil
	}
	if colNum != tbln.colNum {
		return fmt.Errorf("number of columns is different")
	}
	return nil
}

// SetNames is Set Column Name to Tbln.
func (tbln *Tbln) SetNames(names []string) error {
	if err := tbln.setColNum(len(names)); err != nil {
		return err
	}
	tbln.Names = names
	return nil
}

// SetTypes is Set Column Type to Tbln.
func (tbln *Tbln) SetTypes(types []string) error {
	if err := tbln.setColNum(len(types)); err != nil {
		return err
	}
	tbln.Types = types
	return nil
}

// AddComment is Add a comment to Tbln.
func (tbln *Tbln) AddComment(comment string) error {
	tbln.Comments = append(tbln.Comments, comment)
	return nil
}

// AddExtra is Add a extra to Tbln.
func (tbln *Tbln) AddExtra(extra string) error {
	tbln.Extra = append(tbln.Extra, extra)
	return nil
}
