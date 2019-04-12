// Package tbln reads and writes tbln files.
//
// tbln can contain multiple fields like csv.
// Also, it can include checksum and signature inside.
// Tbln contains three types of lines: data, comments, and extras.
// All rows end with a new line(LF).
// The number of fields in all rows must be the same.
//
// Data begins with "| " and ends with " |". Multiple fields are separated by " | ".
// (Note that the delimiter contains spaces.)
//	| fields1 | fields2 | fields3 |
// White space is considered part of a field.
// If "|" is included in the field, "|" must be duplicated.
// 	| -> || , || -> |||
//
// Comments begin with "# ". Comments are not interpreted.
//	# Comments
//
// Extras begin with ";". extras can be interpreted as a header.
// Extras is basically written in the item name: value.
//	; ItemName: Value
// Extras is optional.
// Extras has item names that are interpreted in some special ways.
package tbln

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"regexp"
	"strings"

	"golang.org/x/crypto/ed25519"
)

// Tbln struct is Tbln Definition + Tbln rows.
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

// Read is TBLN Read interface.
type Read interface {
	ReadRow() ([]string, error)
	Types() []string
	Names() []string
	PrimaryKey() []string
	GetDefinition() *Definition
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
	if len(str) < 4 {
		return nil
	}
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
		return len(row), nil
	}
	if len(row) != columnNum {
		return columnNum, fmt.Errorf("invalid column num (%d!=%d) %s", columnNum, len(row), row)
	}
	return columnNum, nil
}

// Types of supported hashes
const (
	SHA256 = "sha256" // import crypto/sha256
	SHA512 = "sha512" // import crypto/sha512
)

// SumHash calculated checksum.
func (t *Tbln) SumHash(hashType string) error {
	h, err := t.calculateHash(hashType)
	if err != nil {
		return err
	}
	if string(t.Hashes[hashType]) != string(h) {
		t.Hashes[hashType] = h
		// Delete signature
		t.Signs = make(map[string]Signature)
	}
	return nil
}

// calculateHash is returns the calculated checksum.
func (t *Tbln) calculateHash(hashType string) ([]byte, error) {
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

// Sign is returns signature for hash.
func (t *Tbln) Sign(name string, pkey []byte) (map[string]Signature, error) {
	if t == nil || t.Definition == nil {
		return nil, fmt.Errorf("no algorithm")
	}
	if t.algorithm != ED25519 {
		return nil, fmt.Errorf("unsupported algorithm")
	}
	if len(pkey) != ed25519.PrivateKeySize {
		return nil, fmt.Errorf("bad private key length")
	}
	t.Signs[name] = Signature{sign: ed25519.Sign(pkey, t.SerializeHash()), algorithm: ED25519}
	return t.Signs, nil
}

// VerifySignature returns the boolean value of the signature verification and Verify().
func (t *Tbln) VerifySignature(name string, pubkey []byte) bool {
	if t == nil || t.Definition == nil || len(t.Signs) == 0 {
		return false
	}
	if len(pubkey) != ed25519.PublicKeySize {
		return false
	}
	s := t.Signs[name]
	if s.algorithm == ED25519 {
		x := ed25519.PublicKey(pubkey)
		if ed25519.Verify(x, t.SerializeHash(), s.sign) {
			return t.Verify()
		}
	}
	return false
}

// Verify returns the boolean value of the hash verification.
func (t *Tbln) Verify() bool {
	if len(t.Hashes) == 0 {
		return false
	}
	for name, old := range t.Hashes {
		new, err := t.calculateHash(name)
		if err != nil {
			return false
		}
		if string(old) != string(new) {
			return false
		}
	}
	return true
}
