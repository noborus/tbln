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
	Reader *bufio.Reader
}

// NewReader is returns Reader.
func NewReader(reader io.Reader) *Reader {
	return &Reader{
		Definition: Definition{Ext: make(map[string]string)},
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
			tr.analyzeExt(str)
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
		tr.SetNames(parseRecord(value))
	case "type":
		tr.SetTypes(parseRecord(value))
	default:
		tr.Ext[key] = value
	}
	return nil
}

func parseRecord(body string) []string {
	body = strings.TrimLeft(body, "| ")
	body = strings.TrimRight(body, " |")
	rec := strings.Split(body, " | ")
	unescape(rec)
	return rec
}

func unescape(rec []string) {
	// Unescape vertical bars || -> |
	for i, column := range rec {
		if strings.Contains(column, "|") {
			rec[i] = UNESCREP.ReplaceAllString(column, "$1")
		}
	}
}

// ReadAll reads all io.Reader
func ReadAll(reader io.Reader) (*Table, error) {
	r := NewReader(reader)
	at := &Table{}
	at.Rows = make([][]string, 0)
	var err error
	var info bool
	for {
		rec, err := r.ReadRow()
		if err != nil {
			break
		}
		if !info {
			at.Definition = r.Definition
		}
		at.RowNum++
		at.Rows = append(at.Rows, rec)
	}
	if err != io.EOF {
		return at, nil
	}
	return nil, err
}
