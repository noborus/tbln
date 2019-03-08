package tbln

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Reader reads records from a tbln file.
type Reader struct {
	*Definition
	r *bufio.Reader
}

// NewReader returns a new Reader that reads from tr.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		Definition: NewDefinition(),
		r:          bufio.NewReader(r),
	}
}

// ReadRow reads one record (a slice of fields) from tr.
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

// ReadAll reads all the remaining records from r.
func ReadAll(r io.Reader) (*Tbln, error) {
	tr := NewReader(r)
	at := &Tbln{}
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
func (tr *Reader) scanLine() ([]string, error) {
	for {
		line, isPrefix, err := tr.r.ReadLine()
		if err != nil {
			return nil, err
		}
		if isPrefix {
			return nil, fmt.Errorf("line too long (%s)", line)
		}
		str := string(line)
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
			return nil, fmt.Errorf("unsupported line (%s)", line)
		}
	}
}

// Analyze Extra.
// Save the necessary items (name, type, TableName, Hash) in Extra in a variable.
// Save other items in Extras.
func (tr *Reader) analyzeExtra(extstr string) error {
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
		if len(tr.Hashes) > 0 {
			tr.Extras[key] = NewExtra(value, true)
		} else {
			tr.Extras[key] = NewExtra(value, false)
		}
	}
	return nil
}
