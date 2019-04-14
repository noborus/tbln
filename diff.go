package tbln

import "fmt"

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
