package tbln

import (
	"fmt"
	"io"
)

// Diff returns diff format
func (d *DiffRow) Diff() string {
	switch d.Les {
	case 0:
		return fmt.Sprintf(" %s", JoinRow(d.Src))
	case 1:
		return fmt.Sprintf("+%s", JoinRow(d.Dst))
	case -1:
		return fmt.Sprintf("-%s", JoinRow(d.Src))
	case 2:
		return fmt.Sprintf("-%s\n+%s", JoinRow(d.Src), JoinRow(d.Dst))
	default:
		return ""
	}
}

// DiffAll Write diff to writer from two readers.
func DiffAll(writer io.Writer, src, dst Reader) error {
	d, err := NewCompare(src, dst)
	if err != nil {
		return err
	}
	for {
		dd, err := d.ReadDiffRow()
		if err != nil {
			break
		}
		fmt.Fprintf(writer, "%s\n", dd.Diff())
	}
	return err
}
