package tbln

import "fmt"

// Diff returns diff format
func (d *DiffTbln) Diff() string {
	switch d.les {
	case 0:
		return fmt.Sprintf(" %s", JoinRow(d.src))
	case 1:
		return fmt.Sprintf("+%s", JoinRow(d.dst))
	case -1:
		return fmt.Sprintf("-%s", JoinRow(d.src))
	case 2:
		return fmt.Sprintf("-%s\n+%s", JoinRow(d.src), JoinRow(d.dst))
	default:
		return ""
	}
}
