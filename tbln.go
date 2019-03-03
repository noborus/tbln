package tbln

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"log"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/ed25519"
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
	Hashes    map[string][]byte
	Signs     map[string][]byte
}

// NewDefinition is create Definition struct.
func NewDefinition() *Definition {
	extras := make(map[string]Extra)
	extras["created_at"] = NewExtra(time.Now().Format(time.RFC3339), false)
	hashes := make(map[string][]byte)
	signs := make(map[string][]byte)
	return &Definition{
		Extras: extras,
		Hashes: hashes,
		Signs:  signs,
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

func (e *Extra) Value() interface{} {
	return e.value
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
	NOTSUPPORT = iota
	SHA256     // import crypto/sha256
	SHA512     // import crypto/sha512
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

func hashType(hString string) HashType {
	switch hString {
	case "sha256":
		return SHA256
	case "sha512":
		return SHA512
	}
	return NOTSUPPORT
}

// Sign is returns signature for hash.
func (t *Tbln) Sign(privateKey ed25519.PrivateKey) map[string][]byte {
	t.Signs["ED25519"] = ed25519.Sign(privateKey, t.SerialHash())
	return t.Signs
}

// VerifySign validates the signature and hash returns the bool value.
func (t *Tbln) VerifySign(pubKey []byte) bool {
	sign := t.Signs["ED25519"]
	x := ed25519.PublicKey(pubKey)
	if ed25519.Verify(x, t.SerialHash(), sign) {
		return t.Verify()
	}
	return false
}

// Verify validates the hash and returns the bool value.
func (t *Tbln) Verify() bool {
	for key, value := range t.Hashes {
		orig, err := t.hash(hashType(key))
		if err != nil {
			log.Fatal(err)
		}
		if string(orig) != string(value) {
			return false
		}
	}
	return true
}

// SerialHash returns a []byte that serializes Hash's map.
func (d *Definition) SerialHash() []byte {
	hashes := make([]string, 0, len(d.Hashes))
	if val, ok := d.Hashes["sha256"]; ok {
		hashes = append(hashes, "sha256:"+fmt.Sprintf("%x", val))
	}
	if val, ok := d.Hashes["sha512"]; ok {
		hashes = append(hashes, "sha512:"+fmt.Sprintf("%x", val))
	}
	return []byte(JoinRow(hashes))
}

// SumHash calculated checksum.
func (t *Tbln) SumHash(hashType HashType) error {
	h, err := t.hash(hashType)
	if err != nil {
		return err
	}
	t.Hashes[hashType.string()] = h
	return nil
}

// hash is returns the calculated checksum.
func (t *Tbln) hash(hashType HashType) ([]byte, error) {
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
	return hash.Sum(nil), nil
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

// SetSignatures is set signatures.
func (d *Definition) SetSignatures(signs []string) error {
	for _, sign := range signs {
		s := strings.SplitN(sign, ":", 2)
		b, err := hex.DecodeString(s[1])
		if err != nil {
			return err
		}
		d.Signs[s[0]] = b
	}
	return nil
}

// SetHashes is set hashes.
func (d *Definition) SetHashes(hashes []string) error {
	for _, hash := range hashes {
		h := strings.SplitN(hash, ":", 2)
		b, err := hex.DecodeString(h[1])
		if err != nil {
			return err
		}
		d.Hashes[h[0]] = b
	}
	return nil
}

// HashTarget is set as target of hash
func (d *Definition) HashTarget(key string, target bool) {
	if v, ok := d.Extras[key]; ok {
		v.hashTarget = target
		d.Extras[key] = v
	}
}
