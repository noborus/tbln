package tbln

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"regexp"
	"strings"
	"time"
)

// Tbln struct is tbln Definition + Tbln rows.
type Tbln struct {
	*Definition
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
func NewDefinition() *Definition {
	extras := make(map[string]Extra)
	extras["created_at"] = NewExtra(time.Now().Format(time.RFC3339), false)
	hashes := make(map[string]string)
	return &Definition{
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
	if len(row) == 0 {
		return ""
	}
	_, err := b.WriteString("|")
	if err != nil {
		return ""
	}
	for _, column := range row {
		_, err = b.WriteString(" " + escape(column) + " |")
		if err != nil {
			return ""
		}
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
func SplitRow(str string) []string {
	if str[:2] == "| " {
		str = str[2:]
	}
	if str[len(str)-2:] == " |" {
		str = str[0 : len(str)-2]
	}
	rec := strings.Split(str, " | ")
	for i, column := range rec {
		rec[i] = unescape(column)
	}
	return rec
}

// UNESCAPE is unescape || -> |
var UNESCAPE = regexp.MustCompile(`\|(\|+)`)

// unescape vertical bars || -> |
func unescape(str string) string {
	if strings.Contains(str, "|") {
		str = UNESCAPE.ReplaceAllString(str, "$1")
	}
	return str
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
			return columnNum, fmt.Errorf("invalid column num (%d!=%d) %s", columnNum, len(row), row)
		}
	}
	return columnNum, nil
}

// HashType supports the type of hash.
type HashType int

// Types of supported hashes
const (
	SHA256 = iota // import crypto/sha256
	SHA512        // import crypto/sha512
)

func (h HashType) string() string {
	switch h {
	case SHA256:
		return "sha256"
	case SHA512:
		return "sha512"
	default:
		return ""
	}
}

// SumHash is returns the calculated checksum.
func (t *Tbln) SumHash(hashType HashType) (map[string]string, error) {
	if t.Hashes == nil {
		t.Hashes = make(map[string]string)
	}
	var hash hash.Hash
	switch hashType {
	case SHA256:
		hash = sha256.New()
	case SHA512:
		hash = sha512.New()
	default:
		return nil, fmt.Errorf("not support")
	}
	w := NewWriter(hash)
	err := w.writeExtraTarget(t.Definition, true)
	if err != nil {
		return nil, err
	}
	for _, row := range t.Rows {
		err = w.WriteRow(row)
		if err != nil {
			return nil, err
		}
	}
	t.Hashes[hashType.string()] = fmt.Sprintf("%x", hash.Sum(nil))
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
