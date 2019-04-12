package tbln

import (
	"fmt"
	"io"
)

// MergeRow reads merged rows from DiffRow
func (d *DiffRow) MergeRow() []string {
	switch d.les {
	case 0:
		return d.src
	case 1:
		return d.dst
	case -1:
		return d.src
	case 2:
		return d.dst
	default:
		return []string{}
	}
}

// MergeAll merges two tbln readers and returns one tbln.
func MergeAll(src, dst Read) (*Tbln, error) {
	tb := &Tbln{}
	diff, err := NewCompare(src, dst)
	if err != nil {
		return nil, err
	}
	sDefinition := src.GetDefinition()
	dDefinition := dst.GetDefinition()
	tb.Definition, err = MergeDefinition(sDefinition, dDefinition)
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
		tb.RowNum++
		tb.Rows = append(tb.Rows, dd.MergeRow())
	}
}

// MergeDefinition merges two tbln Definitions.
func MergeDefinition(src, dst *Definition) (*Definition, error) {
	if src.columnNum != dst.columnNum {
		return nil, fmt.Errorf("different column num")
	}
	d := NewDefinition()
	d.Comments = mergeComment(src, dst)
	d.Extras = src.Extras
	for k, v := range dst.Extras {
		d.Extras[k] = v
	}
	return d, nil
}

func mergeComment(src, dst *Definition) []string {
	comments := make([]string, len(src.Comments))
	copy(comments, src.Comments)

	for _, dc := range dst.Comments {
		appFlag := true
		for _, sc := range src.Comments {
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
