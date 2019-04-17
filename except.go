package tbln

import (
	"io"
)

// ExceptRow reads excepted rows from DiffRow
func (d *DiffRow) ExceptRow() []string {
	switch d.Les {
	case -1:
		return d.Src
	case 2:
		return d.Src
	default:
		return nil
	}
}

// ExceptAll merges two tbln and returns one tbln.
func ExceptAll(t1, t2 Reader) (*Tbln, error) {
	tb := &Tbln{}
	diff, err := NewCompare(t1, t2)
	if err != nil {
		return nil, err
	}
	tb.Definition = t2.GetDefinition()
	tb.Definition.Hashes = nil
	tb.Definition.Signs = nil
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
		row := dd.ExceptRow()
		if row != nil {
			tb.RowNum++
			tb.Rows = append(tb.Rows, row)
		}
	}
}
