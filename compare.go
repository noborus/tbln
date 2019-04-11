package tbln

import (
	"fmt"
	"strconv"
	"strings"
)

type Difference interface {
	Same(row []string) error
	Add(row []string) error
	Mod(srow []string, drow []string) error
	Del(row []string) error
}

type CompareDiff struct {
	Difference
	src    *Reader
	dst    *Reader
	types  []string
	srcRow []string
	dstRow []string
	pkPos  []int
}

func NewCompare(diff Difference, src, dst *Reader) *CompareDiff {
	d := &CompareDiff{
		Difference: diff,
		src:        src,
		dst:        dst,
	}
	d.prepare(src, dst)
	return d
}

func (d *CompareDiff) Compare() error {
	diff := d.Difference
	for {
		switch d.comparePK(d.srcRow, d.dstRow) {
		case 0:
			if JoinRow(d.srcRow) == JoinRow(d.dstRow) {
				diff.Same(d.srcRow)
			} else {
				diff.Mod(d.srcRow, d.dstRow)
			}
			d.srcRow = nextRow(d.src)
			d.dstRow = nextRow(d.dst)
		case 1:
			if len(d.dstRow) > 0 {
				diff.Add(d.dstRow)
				d.dstRow = nextRow(d.dst)
			}
		case -1:
			if len(d.srcRow) > 0 {
				diff.Del(d.srcRow)
				d.srcRow = nextRow(d.src)
			}
		}
		if len(d.srcRow) == 0 && len(d.dstRow) == 0 {
			break
		}
	}
	return nil
}

func nextRow(r *Reader) []string {
	row, err := r.ReadRow()
	if err != nil {
		return []string{}
	}
	return row
}

func (d *CompareDiff) prepare(src, dst *Reader) error {
	var err error
	d.srcRow, err = src.ReadRow()
	if err != nil {
		return err
	}
	d.pkPos, err = getPKey(src)
	if err != nil {
		return err
	}

	d.dstRow, err = dst.ReadRow()
	if err != nil {
		return err
	}
	d.types = src.Types

	return nil
}

func getPKey(tr *Reader) ([]int, error) {
	pk := SplitRow(fmt.Sprintf("%s", tr.ExtraValue("primarykey")))
	if len(pk) == 0 {
		return nil, fmt.Errorf("no primary key")
	}
	pkpos := make([]int, 0)
	for n, v := range tr.Names {
		if pk[0] == v {
			pkpos = append(pkpos, n)
		}
	}
	return pkpos, nil
}

func (d *CompareDiff) comparePK(srcRow, dstRow []string) int {
	if len(srcRow) == 0 {
		return 1
	}
	if len(dstRow) == 0 {
		return -1
	}
	for _, p := range d.pkPos {
		ret := compareType(d.types[p], srcRow[p], dstRow[p])
		if ret != 0 {
			return ret
		}
	}
	return 0
}

func compareType(dtype string, src string, dst string) int {
	switch dtype {
	case "int":
		return compareInt(src, dst)
	case "bigint", "double precision", "numeric":
		return compareFloat(src, dst)
	default:
		return strings.Compare(src, dst)
	}
}

func compareInt(src string, dst string) int {
	var err error
	var s, d int
	if s, err = strconv.Atoi(src); err != nil {
		return strings.Compare(src, dst)
	}
	if d, err = strconv.Atoi(dst); err != nil {
		return strings.Compare(src, dst)
	}
	if s > d {
		return 1
	} else if s < d {
		return -1
	}
	return 0
}

func compareFloat(src string, dst string) int {
	var err error
	var s, d float64
	if s, err = strconv.ParseFloat(src, 64); err != nil {
		return strings.Compare(src, dst)
	}
	if d, err = strconv.ParseFloat(dst, 64); err != nil {
		return strings.Compare(src, dst)
	}
	if s > d {
		return 1
	} else if s < d {
		return -1
	}
	return 0
}
