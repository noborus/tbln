package tbln

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// HashType supports the type of hash.
type HashType int

// Types of supported hashes
const (
	SHA256 = iota // import crypto/sha256
	SHA512        // import crypto/sha512
)

// Tbln struct is tbln Definition + Tbln rows.
type Tbln struct {
	Definition
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
	Extras    map[string]Extra
	Hashes    map[string]string
}

// NewDefinition is create Definition struct.
func NewDefinition() Definition {
	extras := make(map[string]Extra)
	extras["created_at"] = NewExtra(time.Now().Format(time.RFC3339), false)
	hashes := make(map[string]string)
	return Definition{
		Extras: extras,
		Hashes: hashes,
	}
}

// Extra is table definition extra struct.
type Extra struct {
	value      interface{}
	hashTarget bool
}

// NewExtra is return new extra struct.
func NewExtra(value interface{}, hashTarget bool) Extra {
	return Extra{
		value:      value,
		hashTarget: hashTarget,
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

// JoinRow makes a Row array a character string.
func JoinRow(row []string) string {
	var b strings.Builder
	b.WriteString("|")
	for _, column := range row {
		b.WriteString(" " + escape(column) + " |")
	}
	return b.String()
}

// ESCAPE is escape | -> ||
var ESCAPE = regexp.MustCompile(`(\|+)`)

func escape(str string) string {
	if strings.Contains(str, "|") {
		str = ESCAPE.ReplaceAllString(str, "|$1")
	}
	return str
}

// SplitRow divides a character string into a row array.
func SplitRow(body string) []string {
	if body[:2] == "| " {
		body = body[2:]
	}
	if body[len(body)-2:] == " |" {
		body = body[0 : len(body)-2]
	}
	rec := strings.Split(body, " | ")
	return unescape(rec)
}

// UNESCAPE is unescape || -> |
var UNESCAPE = regexp.MustCompile(`\|(\|+)`)

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

func checkRow(ColumnNum int, row []string) (int, error) {
	if ColumnNum == 0 {
		ColumnNum = len(row)
	} else {
		if len(row) != ColumnNum {
			return ColumnNum, fmt.Errorf("Error: invalid column num (%d!=%d) %s", ColumnNum, len(row), row)
		}
	}
	return ColumnNum, nil
}

// SumHash is returns the calculated checksum.
// Checksum target is exported to buffer and saved.
func (t *Tbln) SumHash(hashType HashType) (map[string]string, error) {
	if t.Hashes == nil {
		t.Hashes = make(map[string]string)
	}
	writer := bufio.NewWriter(&t.buffer)
	bw := NewWriter(writer)
	err := bw.writeExtraTarget(t.Definition, true)
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
	switch hashType {
	case SHA256:
		sum := sha256.Sum256(t.buffer.Bytes())
		t.Hashes["sha256"] = fmt.Sprintf("%x", sum)
	case SHA512:
		sum := sha512.Sum512(t.buffer.Bytes())
		t.Hashes["sha512"] = fmt.Sprintf("%x", sum)
	default:
		return nil, fmt.Errorf("not support")
	}
	return t.Hashes, nil
}

// TableName returns the Table Name.
func (d *Definition) TableName() string {
	return d.tableName
}

// SetTableName is set Table Name of the Table.
func (d *Definition) SetTableName(name string) {
	d.tableName = name
	if len(name) != 0 {
		d.Extras["TableName"] = NewExtra(name, false)
	} else {
		delete(d.Extras, "TableName")
	}
}

// SetNames is set Column Name to the Table.
func (d *Definition) SetNames(names []string) error {
	if err := d.setColNum(len(names)); err != nil {
		return err
	}
	d.Names = names
	if names != nil {
		d.Extras["name"] = NewExtra(JoinRow(names), true)
	} else {
		delete(d.Extras, "name")
	}
	return nil
}

// SetTypes is set Column Type to Table.
func (d *Definition) SetTypes(types []string) error {
	if types == nil {
		return nil
	}
	if err := d.setColNum(len(types)); err != nil {
		return err
	}
	d.Types = types
	if types != nil {
		d.Extras["type"] = NewExtra(JoinRow(types), true)
	} else {
		delete(d.Extras, "type")
	}
	return nil
}

// ColumnNum returns the number of columns.
func (d *Definition) ColumnNum() int {
	return d.columnNum
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

// SetHashes is set hashes.
func (d *Definition) SetHashes(hashes []string) {
	for _, hash := range hashes {
		h := strings.SplitN(hash, ":", 2)
		d.Hashes[h[0]] = h[1]
	}
}

// HashTarget is set as target of hash
func (d *Definition) HashTarget(key string, target bool) {
	if v, ok := d.Extras[key]; ok {
		v.hashTarget = target
		d.Extras[key] = v
	}
}
