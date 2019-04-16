package tbln

import (
	"fmt"
	"io"
)

// DiffMode represents the mode of diff.
type DiffMode int

// Represents the diff mode
const (
	OnlyAdd = iota
	OnlyDiff
	AllDiff
)

// Diff returns diff format
func (d *DiffRow) Diff(diffMode DiffMode) string {
	switch d.Les {
	case 0:
		if diffMode == AllDiff {
			return fmt.Sprintf(" %s", JoinRow(d.Src))
		}
	case 1:
		if diffMode == OnlyAdd {
			return fmt.Sprintf("%s", JoinRow(d.Dst))
		}
		return fmt.Sprintf("+%s", JoinRow(d.Dst))
	case -1:
		if diffMode == AllDiff || diffMode == OnlyDiff {
			return fmt.Sprintf("-%s", JoinRow(d.Src))
		}
	case 2:
		var str string
		if diffMode == AllDiff || diffMode == OnlyDiff {
			str = fmt.Sprintf("-%s\n", JoinRow(d.Src))
		}
		if diffMode == OnlyAdd {
			return fmt.Sprintf("%s", JoinRow(d.Dst))
		}
		return str + fmt.Sprintf("+%s", JoinRow(d.Dst))
	default:
		return ""
	}
	return ""
}

// DiffAll Write diff to writer from two readers.
func DiffAll(writer io.Writer, src, dst Reader, diffMode DiffMode) error {
	d, err := NewCompare(src, dst)
	if err != nil {
		return err
	}
	for {
		dd, err := d.ReadDiffRow()
		if err != nil {
			break
		}
		dString := dd.Diff(diffMode)
		if dString != "" {
			fmt.Fprintf(writer, "%s\n", dString)
		}
	}
	return err
}
