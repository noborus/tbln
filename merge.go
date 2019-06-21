package tbln

import (
	"fmt"
	"io"
)

// MergeMode represents the mode of merge.
type MergeMode int

// Represents the merge mode
const (
	MergeIgnore MergeMode = iota
	MergeUpdate
	MergeDelete
)

// MergeRow reads merged rows from DiffRow
func (d *DiffRow) MergeRow(mode MergeMode) []string {
	switch d.Les {
	case 0:
		return d.Self
	case 1:
		return d.Other
	case -1:
		if mode == MergeDelete {
			return nil
		}
		return d.Self
	case 2:
		if mode == MergeIgnore {
			return d.Self
		}
		return d.Other
	default:
		return nil
	}
}

// MergeAll merges two tbln and returns one tbln.
func MergeAll(t1, t2 Reader, mode MergeMode) (*TBLN, error) {
	tb := &TBLN{}
	diff, err := NewCompare(t1, t2)
	if err != nil {
		return nil, err
	}
	t1Definition := t1.GetDefinition()
	t2Definition := t2.GetDefinition()
	tb.Definition, err = MergeDefinition(t1Definition, t2Definition)
	if err != nil {
		return nil, err
	}
	tb.Rows = make([][]string, 0)
	for {
		dd, err := diff.ReadDiffRow()
		if err != nil {
			if err == io.EOF {
				return tb, nil
			}
			return nil, err
		}
		row := dd.MergeRow(mode)
		if row != nil {
			tb.RowNum++
			tb.Rows = append(tb.Rows, row)
		}
	}
}

// MergeDefinition merges two tbln Definitions.
func MergeDefinition(t1d, t2d *Definition) (*Definition, error) {
	if t1d.columnNum != t2d.columnNum {
		return nil, fmt.Errorf("different column num")
	}
	d := NewDefinition()
	d.Comments = mergeComment(t1d, t2d)
	d.Extras = t1d.Extras
	for k, v := range t2d.Extras {
		d.Extras[k] = v
	}
	return d, nil
}

func mergeComment(t1d, t2d *Definition) []string {
	comments := make([]string, len(t1d.Comments))
	copy(comments, t1d.Comments)

	for _, dc := range t2d.Comments {
		appFlag := true
		for _, sc := range t1d.Comments {
			if dc == sc {
				appFlag = false
				break
			}
		}
		if appFlag {
			comments = append(comments, dc)
		}
	}
	return comments
}
