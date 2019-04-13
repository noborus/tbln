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
	Les int
	Src []string
	Dst []string
}

// NewCompare returns a Compare structure
func NewCompare(src, dst Comparer) (*Compare, error) {
	cmp := &Compare{
		src:   src,
		dst:   dst,
		sNext: false,
		dNext: false,
	}
	var err error
	cmp.srcRow, err = src.ReadRow()
	if err != nil {
		return nil, err
	}
	cmp.dstRow, err = dst.ReadRow()
	if err != nil {
		return nil, err
	}
	cmp.pk, err = getPK(src, dst)
	if err != nil {
		return nil, err
	}
	return cmp, nil
}

// ReadDiffRow compares two rows and returns the difference.
func (cmp *Compare) ReadDiffRow() (*DiffRow, error) {
	if cmp.sNext {
		// Ignore errors to continue reading both ends.
		cmp.srcRow, _ = cmp.src.ReadRow()
	}
	if cmp.dNext {
		// Ignore errors to continue reading both ends.
		cmp.dstRow, _ = cmp.dst.ReadRow()
	}

	switch cmp.diffPrimaryKey() {
	case 0:
		cmp.sNext = true
		cmp.dNext = true
		if JoinRow(cmp.srcRow) == JoinRow(cmp.dstRow) {
			return &DiffRow{0, cmp.srcRow, cmp.dstRow}, nil
		}
		return &DiffRow{2, cmp.srcRow, cmp.dstRow}, nil
	case 1:
		cmp.sNext = false
		cmp.dNext = true
		if len(cmp.dstRow) > 0 {
			return &DiffRow{1, nil, cmp.dstRow}, nil
		}
	case -1:
		cmp.sNext = true
		cmp.dNext = false
		if len(cmp.srcRow) > 0 {
			return &DiffRow{-1, cmp.srcRow, nil}, nil
		}
	}
	return nil, io.EOF
}

func (cmp *Compare) PrimaryKeyRow(row []string) []string {
	if row == nil {
		return nil
	}
	pkRow := make([]string, 0, len(cmp.pk))
	for _, pk := range cmp.pk {
		pkRow = append(pkRow, row[pk.pos])
	}
	return pkRow
}

func (cmp *Compare) diffPrimaryKey() int {
	if len(cmp.srcRow) == 0 {
		return 1
	}
	if len(cmp.dstRow) == 0 {
		return -1
	}
	for _, pk := range cmp.pk {
		ret := compareType(pk.typ, cmp.srcRow[pk.pos], cmp.dstRow[pk.pos])
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
