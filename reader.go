package tbln

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

// Reader is TBLN Reader interface.
type Reader interface {
	ReadRow() ([]string, error)
	GetDefinition() *Definition
}

// FileReader reads records from a tbln file.
type FileReader struct {
	*Definition
	r *bufio.Reader
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *FileReader {
	return &FileReader{
		Definition: NewDefinition(),
		r:          bufio.NewReader(r),
	}
}

// ReadRow reads one record (a slice of fields) from tr.
func (tr *FileReader) ReadRow() ([]string, error) {
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

// ReadAll reads all r and returns a tbln struct.
// ReadAll returns when the blank line is reached.
func ReadAll(r io.Reader) (*TBLN, error) {
	tr := NewReader(r)
	at := &TBLN{}
	at.Rows = make([][]string, 0)
	for {
		rec, err := tr.ReadRow()
		if err != nil {
			if err == io.EOF {
				at.Definition = tr.Definition
				at.Hashes = tr.Hashes
				return at, nil
			}
			return nil, err
		}
		// blank line
		if rec == nil {
			// TODO: To be able to process the next block without ending with blank line.
			// continue
			at.Definition = tr.Definition
			at.Hashes = tr.Hashes
			return at, nil
		}
		at.RowNum++
		at.Rows = append(at.Rows, rec)
	}
}

// scanLine reads from tr and returns either one row or a blank line.
// Comments and Extra lines are read until reaching a row or blank line.
func (tr *FileReader) scanLine() ([]string, error) {
	var buf bytes.Buffer
	for {
		buf.Reset()
		for {
			line, isPrefix, err := tr.r.ReadLine()
			if err != nil {
				return nil, err
			}
			buf.Write(line)
			if isPrefix != true {
				break
			}
		}
		str := buf.String()
		switch {
		case strings.HasPrefix(str, "| "):
			return SplitRow(str), nil
		case strings.HasPrefix(str, "#"):
			tr.Comments = append(tr.Comments, strings.TrimSpace(str[1:]))
		case strings.HasPrefix(str, "; "):
			err := tr.analyzeExtra(str)
			if err != nil {
				return nil, err
			}
		case str == "":
			return nil, nil
		default:
			return nil, fmt.Errorf("unsupported line (%s)", str)
		}
	}
}

// Analyze Extra.
// Save the necessary items (name, type, TableName, Hash) in Extra in a variable.
// Save other items in Extras.
func (tr *FileReader) analyzeExtra(extstr string) error {
	extstr = strings.TrimLeft(extstr, "; ")
	keypos := strings.Index(extstr, ":")
	if keypos <= 0 {
		return fmt.Errorf("extra format error %s", extstr)
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
		err := tr.SetHashes(SplitRow(value))
		if err != nil {
			return err
		}
	case "Signature":
		err := tr.SetSignatures(SplitRow(value))
		if err != nil {
			return err
		}
	default:
		tr.SetExtra(key, value)
	}
	return nil
}
