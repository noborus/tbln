package tbln

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Reader is reader struct.
type Reader struct {
	Table
	Reader *bufio.Reader
}

// NewReader is returns Reader.
func NewReader(reader io.Reader) *Reader {
	return &Reader{
		Table:  Table{Ext: make(map[string]string)},
		Reader: bufio.NewReader(reader),
	}
}

// ReadRow is return one row.
func (tr *Reader) ReadRow() ([]string, error) {
	rec, err := tr.scanLine()
	if err != nil {
		return nil, err
	}
	if (tr.columnNum > 0) && (len(rec) != tr.columnNum) {
		return nil, fmt.Errorf("Error: invalid column num (%d!=%d) %s", tr.columnNum, len(rec), rec)
	}
	return rec, nil
}

func (tr *Reader) scanLine() ([]string, error) {
	for {
		line, _, err := tr.Reader.ReadLine()
		if err != nil {
			return nil, err
		}
		str := string(line)
		if len(str) == 0 {
			return nil, nil
		}
		if len(str) < 3 {
			return nil, fmt.Errorf("Error line: %s", str)
		}
		switch str[:2] {
		case "| ":
			return parseRecord(str), nil
		case "# ":
			tr.Comments = append(tr.Comments, str[2:])
		case "; ":
			extstr := str[2:]
			keypos := strings.Index(extstr, ":")
			if keypos <= 0 {
				return nil, fmt.Errorf("Error: Extra format error %s", extstr)
			}
			key := extstr[:keypos]
			value := extstr[keypos+2:]
			tr.analyzeExt(key, value)
		default:
			return nil, fmt.Errorf("Error: Unsupported line")
		}
	}
}

func (tr *Reader) analyzeExt(key string, value string) error {
	switch key {
	case "name":
		tr.setNames(parseRecord(value))
	case "type":
		tr.setTypes(parseRecord(value))
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
