package tbln

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Tbln struct is tbln Definition + Tbln rows.
type Tbln struct {
	Definition
	Hash   map[string]string
	buffer bytes.Buffer
	RowNum int
	Rows   [][]string
}

// NewTbln is create tbln struct.
func NewTbln() *Tbln {
	return &Tbln{
		Definition: NewDefinition(),
	}
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
	ext := make(map[string]Extra)
	ext["created_at"] = NewExtra(time.Now().Format(time.RFC3339))
	return Definition{
		Ext: ext,
	}
}

// Read is TBLN Read interface.
type Read interface {
	ReadRow() ([]string, error)
}

// Write is TBLN Write interface.
type Write interface {
	WriteDefinition(Definition) error
	WriteRow([]string) error
}

// ESCAPE is escape | -> ||
var ESCAPE = regexp.MustCompile(`(\|+)`)

// UNESCAPE is unescape || -> |
var UNESCAPE = regexp.MustCompile(`\|(\|+)`)

func joinRow(row []string) string {
	var b strings.Builder
	b.WriteString("|")
	for _, column := range row {
		if strings.Contains(column, "|") {
			column = ESCAPE.ReplaceAllString(column, "|$1")
		}
		b.WriteString(" " + column + " |")
	}
	return b.String()
}

func splitRow(body string) []string {
	if body[:2] == "| " {
		body = body[2:]
	}
	if body[len(body)-2:] == " |" {
		body = body[0 : len(body)-2]
	}
	rec := strings.Split(body, " | ")
	return unescape(rec)
}

// unescape vertical bars || -> |
func unescape(rec []string) []string {
	for i, column := range rec {
		if strings.Contains(column, "|") {
			rec[i] = UNESCAPE.ReplaceAllString(column, "$1")
		}
	}
	return rec
}

// AddRows is Add row to Table.
func (t *Tbln) AddRows(row []string) error {
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
func (t *Tbln) SumHash() (map[string]string, error) {
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

// Extra is table definition extra struct.
type Extra struct {
	value      interface{}
	hashTarget bool
}

// NewExtra is return new extra struct.
func NewExtra(value interface{}) Extra {
	return Extra{
		value:      value,
		hashTarget: false,
	}
}

// HashTarget is set as target of hash
func (d *Definition) HashTarget(key string, target bool) {
	if v, ok := d.Ext[key]; ok {
		v.hashTarget = target
		d.Ext[key] = v
	}
}
