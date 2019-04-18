package tbln

import (
	"fmt"
	"io"
)

// DiffMode represents the mode of diff.
type DiffMode int

// Represents the diff mode
const (
	OnlyAdd DiffMode = iota
	OnlyDiff
	AllDiff
)

func (m DiffMode) String() string {
	switch m {
	case OnlyAdd:
		return "OnlyAdd"
	case OnlyDiff:
		return "OnlyDiff"
	case AllDiff:
		return "AllDiff"
	default:
		return "Unknown"
	}
}

// Diff returns diff format
func (d *DiffRow) Diff(diffMode DiffMode) string {
	switch d.Les {
	case 0:
		if diffMode == AllDiff {
			return fmt.Sprintf(" %s", JoinRow(d.Self))
		}
	case 1:
		if diffMode == OnlyAdd {
			return JoinRow(d.Other)
		}
		return fmt.Sprintf("+%s", JoinRow(d.Other))
	case -1:
		if diffMode == AllDiff || diffMode == OnlyDiff {
			return fmt.Sprintf("-%s", JoinRow(d.Self))
		}
	case 2:
		var str string
		if diffMode == OnlyAdd {
			return JoinRow(d.Other)
		}
		if diffMode == AllDiff || diffMode == OnlyDiff {
			str = fmt.Sprintf("-%s\n", JoinRow(d.Self))
		}
		return str + fmt.Sprintf("+%s", JoinRow(d.Other))
	default:
		return ""
	}
	return ""
}

// DiffAll Write diff to writer from two readers.
func DiffAll(writer io.Writer, t1, t2 Reader, diffMode DiffMode) error {
	d, err := NewCompare(t1, t2)
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
