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
		Definition: NewDefinition(),
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

// ReadAll reads all io.Reader
func ReadAll(reader io.Reader) (*Tbln, error) {
	r := NewReader(reader)
	at := &Tbln{}
	at.Rows = make([][]string, 0)
	for {
		rec, err := r.ReadRow()
		if err != nil {
			if err == io.EOF {
				at.Definition = r.Definition
				at.Hash = r.Hash
				return at, nil
			}
			return nil, err
		}
		at.RowNum++
		at.Rows = append(at.Rows, rec)
	}
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
			return SplitRow(str), nil
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
		return tr.SetNames(SplitRow(value))
	case "type":
		return tr.SetTypes(SplitRow(value))
	case "TableName":
		tr.SetTableName(value)
	case "Hash":
		tr.SetHashes(SplitRow(value))
	default:
		if len(tr.Hash) > 0 {
			tr.Ext[key] = Extra{value: value, hashTarget: true}
		} else {
			tr.Ext[key] = Extra{value: value, hashTarget: false}
		}
	}
	return nil
}
