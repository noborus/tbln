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
	PK     []Pkey
}

// Pkey PrimaryKey information
type Pkey struct {
	Pos  int
	Name string
	Typ  string
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
	cmp.PK, err = getPK(src, dst)
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

// ColumnPrimaryKey return  columns primary key
func (cmp *Compare) ColumnPrimaryKey(row []string) []string {
	if row == nil {
		return nil
	}
	colPK := make([]string, 0, len(cmp.PK))
	for _, pk := range cmp.PK {
		colPK = append(colPK, row[pk.Pos])
	}
	return colPK
}

func (cmp *Compare) diffPrimaryKey() int {
	if len(cmp.srcRow) == 0 {
		return 1
	}
	if len(cmp.dstRow) == 0 {
		return -1
	}
	for _, pk := range cmp.PK {
		ret := compareType(pk.Typ, cmp.srcRow[pk.Pos], cmp.dstRow[pk.Pos])
		if ret != 0 {
			return ret
		}
	}
	return 0
}

func getPK(src, dst Comparer) ([]Pkey, error) {
	sd := src.GetDefinition()
	dd := dst.GetDefinition()
	var pos []int
	dPos, err := dd.GetPKeyPos()
	if err == nil {
		pos = dPos
	}
	sPos, err := sd.GetPKeyPos()
	if err == nil {
		pos = sPos
	}
	if len(pos) == 0 {
		return nil, fmt.Errorf("no primary key")
	}
	if len(sPos) != len(dPos) {
		return nil, fmt.Errorf("primary key position")
	}
	pk := make([]Pkey, len(pos))
	for i, v := range pos {
		if sPos[i] != dPos[i] {
			return nil, fmt.Errorf("primary key position")
		}
		st := sd.Types()
		dt := dd.Types()
		if st[i] != dt[i] {
			return nil, fmt.Errorf("unmatch data type")
		}
		sn := sd.Names()
		pk[i] = Pkey{v, sn[i], st[i]}
	}
	return pk, nil
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
