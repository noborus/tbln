package tbln

import (
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"log"
	"regexp"
	"strings"

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

// SumHash calculated checksum.
func (t *Tbln) SumHash(hashType HashType) error {
	h, err := t.calculateHash(hashType)
	if err != nil {
		return err
	}
	t.Hashes[hashType.string()] = h
	return nil
}

// calculateHash is returns the calculated checksum.
func (t *Tbln) calculateHash(hashType HashType) ([]byte, error) {
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
func (t *Tbln) Sign(privateKey ed25519.PrivateKey) map[string][]byte {
	t.Signs["ED25519"] = ed25519.Sign(privateKey, t.SerializeHash())
	return t.Signs
}

// VerifySignature returns the boolean value of the signature verification and Verify().
func (t *Tbln) VerifySignature(pubKey []byte) bool {
	sign := t.Signs["ED25519"]
	x := ed25519.PublicKey(pubKey)
	if ed25519.Verify(x, t.SerializeHash(), sign) {
		return t.Verify()
	}
	return false
}

// Verify returns the boolean value of the hash varification.
func (t *Tbln) Verify() bool {
	for name, old := range t.Hashes {
		new, err := t.calculateHash(hashType(name))
		if err != nil {
			log.Fatal(err)
		}
		if string(old) != string(new) {
			return false
		}
	}
	return true
}
