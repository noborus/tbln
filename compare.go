package tbln

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Comparer is TBLN Compare interface.
type Comparer interface {
	ReadRow() ([]string, error)
	Types() []string
	Names() []string
	PrimaryKey() []string
	GetDefinition() *Definition
}

// Compare is a structure that compares two tbln reads
type Compare struct {
	src    Comparer
	dst    Comparer
	srcRow []string
	dstRow []string
	sNext  bool
	dNext  bool
	pk     []pkey
}

type pkey struct {
	pos  int
	name string
	typ  string
}

// DiffRow contains the difference between two row.
type DiffRow struct {
	les int
	src []string
	dst []string
}

// NewCompare returns a Compare structure
func NewCompare(src, dst Comparer) (*Compare, error) {
	c := &Compare{
		src:   src,
		dst:   dst,
		sNext: false,
		dNext: false,
	}
	var err error
	c.srcRow, err = src.ReadRow()
	if err != nil {
		return nil, err
	}
	c.dstRow, err = dst.ReadRow()
	if err != nil {
		return nil, err
	}
	c.pk, err = getPK(src, dst)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// ReadDiffRow compares two rows and returns the difference.
func (dr *Compare) ReadDiffRow() (*DiffRow, error) {
	if dr.sNext {
		// Ignore errors to continue reading both ends.
		dr.srcRow, _ = dr.src.ReadRow()
	}
	if dr.dNext {
		// Ignore errors to continue reading both ends.
		dr.dstRow, _ = dr.dst.ReadRow()
	}

	switch dr.diffPrimaryKey() {
	case 0:
		dr.sNext = true
		dr.dNext = true
		if JoinRow(dr.srcRow) == JoinRow(dr.dstRow) {
			return &DiffRow{0, dr.srcRow, dr.dstRow}, nil
		}
		return &DiffRow{2, dr.srcRow, dr.dstRow}, nil
	case 1:
		dr.sNext = false
		dr.dNext = true
		if len(dr.dstRow) > 0 {
			return &DiffRow{1, nil, dr.dstRow}, nil
		}
	case -1:
		dr.sNext = true
		dr.dNext = false
		if len(dr.srcRow) > 0 {
			return &DiffRow{-1, dr.srcRow, nil}, nil
		}
	}
	return nil, io.EOF
}

func (dr *Compare) diffPrimaryKey() int {
	if len(dr.srcRow) == 0 {
		return 1
	}
	if len(dr.dstRow) == 0 {
		return -1
	}
	for _, pk := range dr.pk {
		ret := compareType(pk.typ, dr.srcRow[pk.pos], dr.dstRow[pk.pos])
		if ret != 0 {
			return ret
		}
	}
	return 0
}

func getPK(src, dst Comparer) ([]pkey, error) {
	var pos []int
	dPos, err := getPKeyPos(dst)
	if err == nil {
		pos = dPos
	}
	sPos, err := getPKeyPos(src)
	if err == nil {
		pos = sPos
	}
	if len(pos) == 0 {
		return nil, fmt.Errorf("no primary key")
	}
	if len(sPos) != len(dPos) {
		return nil, fmt.Errorf("primary key position")
	}
	pk := make([]pkey, len(pos))
	for i, v := range pos {
		if sPos[i] != dPos[i] {
			return nil, fmt.Errorf("primary key position")
		}
		st := src.Types()
		dt := dst.Types()
		if st[i] != dt[i] {
			return nil, fmt.Errorf("unmatch data type")
		}
		sn := src.Names()
		pk[i] = pkey{v, sn[i], st[i]}
	}
	return pk, nil
}

func getPKeyPos(tr Comparer) ([]int, error) {
	pk := tr.PrimaryKey()
	if len(pk) == 0 {
		return nil, fmt.Errorf("no primary key")
	}
	pkpos := make([]int, 0)
	for _, p := range pk {
		for n, v := range tr.Names() {
			if p == v {
				pkpos = append(pkpos, n)
				break
			}
		}
	}
	return pkpos, nil
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
