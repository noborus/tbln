package tbln

import (
	"io"
)

// ExceptRow reads excepted rows from DiffRow
func (d *DiffRow) ExceptRow() []string {
	switch d.Les {
	case -1:
		return d.Self
	case 2:
		return d.Self
	default:
		return nil
	}
}

// ExceptAll merges two tbln and returns one tbln.
func ExceptAll(t1, t2 Reader) (*TBLN, error) {
	diff, err := NewCompare(t1, t2)
	if err != nil {
		return nil, err
	}
	tb := NewTBLN()
	tb.Definition = t2.GetDefinition()
	tb.Definition.Hashes = make(map[string][]byte)
	tb.Definition.Signs = make(Signatures)
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
