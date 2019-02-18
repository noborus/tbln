package tbln

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Reader is reader struct.
type Reader struct {
	Definition
	Hash   map[string]string
	Reader *bufio.Reader
}

// NewReader is returns Reader.
func NewReader(reader io.Reader) *Reader {
	return &Reader{
		Definition: NewDefinition(),
		Hash:       make(map[string]string),
		Reader:     bufio.NewReader(reader),
	}
}

// ReadRow is return one row.
func (tr *Reader) ReadRow() ([]string, error) {
	rec, err := tr.scanLine()
	if err != nil || rec == nil {
		return nil, err
	}
	tr.columnNum, err = checkRow(tr.columnNum, rec)
	if err != nil {
		return nil, err
	}
	return rec, nil
}

// Return on one line or blank line.
func (tr *Reader) scanLine() ([]string, error) {
	for {
		line, _, err := tr.Reader.ReadLine()
		if err != nil {
			return nil, err
		}
		str := string(line)
		switch {
		case strings.HasPrefix(str, "| "):
			return parseRecord(str), nil
		case strings.HasPrefix(str, "# "):
			tr.Comments = append(tr.Comments, str[2:])
		case strings.HasPrefix(str, "; "):
			err := tr.analyzeExt(str)
			if err != nil {
				return nil, err
			}
		case str == "":
			return nil, nil
		default:
			return nil, fmt.Errorf("Error: Unsupported line")
		}
	}
}

func (tr *Reader) analyzeExt(extstr string) error {
	extstr = strings.TrimLeft(extstr, "; ")
	keypos := strings.Index(extstr, ":")
	if keypos <= 0 {
		return fmt.Errorf("Error: Extra format error %s", extstr)
	}
	key := extstr[:keypos]
	value := extstr[keypos+2:]
	switch key {
	case "name":
		return tr.SetNames(parseRecord(value))
	case "type":
		return tr.SetTypes(parseRecord(value))
	case "TableName":
		tr.SetTableName(value)
	case "sha256":
		tr.Hash["sha256"] = value
	default:
		if len(tr.Hash) > 0 {
			tr.Ext[key] = Extra{value: value, hashTarget: true}
		} else {
			tr.Ext[key] = Extra{value: value, hashTarget: false}
		}
	}
	return nil
}

func parseRecord(body string) []string {
	if body[:2] == "| " {
		body = body[2:]
	}
	if body[len(body)-2:] == " |" {
		body = body[0 : len(body)-2]
	}
	rec := strings.Split(body, " | ")
	unescape(rec)
	return rec
}

func unescape(rec []string) []string {
	// Unescape vertical bars || -> |
	for i, column := range rec {
		if strings.Contains(column, "|") {
			rec[i] = UNESCREP.ReplaceAllString(column, "$1")
		}
	}
	return rec
}

// ReadAll reads all io.Reader
func ReadAll(reader io.Reader) (*Tbln, error) {
	r := NewReader(reader)
	at := &Tbln{}
	at.Rows = make([][]string, 0)
	var err error
	for {
		rec, err := r.ReadRow()
		if err != nil {
			break
		}
		at.RowNum++
		at.Rows = append(at.Rows, rec)
	}
	if err != io.EOF {
		at.Definition = r.Definition
		at.Hash = r.Hash
		return at, nil
	}
	return nil, err
}
