package tbln

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Compare is a structure that compares two tbln reads
type Compare struct {
	t1     Reader
	t2     Reader
	t1Row  []string
	t2Row  []string
	t1Next bool
	t2Next bool

	PK []Pkey
}

// Pkey PrimaryKey information
type Pkey struct {
	Pos  int
	Name string
	Typ  string
}

// DiffRow contains the difference between two row.
type DiffRow struct {
	Les   int
	Self  []string
	Other []string
}

// NewCompare returns a Reader interface
func NewCompare(t1, t2 Reader) (*Compare, error) {
	cmp := &Compare{
		t1:     t1,
		t2:     t2,
		t1Next: false,
		t2Next: false,
	}
	var err error
	cmp.t1Row, err = t1.ReadRow()
	if err != nil {
		if err != io.EOF {
			return nil, err
		}
	}
	cmp.t2Row, err = t2.ReadRow()
	if err != nil {
		if err != io.EOF {
			return nil, err
		}
	}
	cmp.PK, err = cmp.getPK()
	if err != nil {
		return nil, err
	}
	return cmp, nil
}

// ReadDiffRow compares two rows and returns the difference.
func (cmp *Compare) ReadDiffRow() (*DiffRow, error) {
	var err error
	if cmp.t1Next {
		cmp.t1Row, err = cmp.t1.ReadRow()
		// Ignore errors to continue reading both ends.
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
		}
	}
	if cmp.t2Next {
		cmp.t2Row, err = cmp.t2.ReadRow()
		// Ignore errors to continue reading both ends.
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
		}
	}

	switch cmp.diffPrimaryKey() {
	case 0:
		cmp.t1Next = true
		cmp.t2Next = true
		if JoinRow(cmp.t1Row) == JoinRow(cmp.t2Row) {
			return &DiffRow{0, cmp.t1Row, cmp.t2Row}, nil
		}
		return &DiffRow{2, cmp.t1Row, cmp.t2Row}, nil
	case 1:
		cmp.t1Next = false
		cmp.t2Next = true
		if len(cmp.t2Row) > 0 {
			return &DiffRow{1, nil, cmp.t2Row}, nil
		}
	case -1:
		cmp.t1Next = true
		cmp.t2Next = false
		if len(cmp.t1Row) > 0 {
			return &DiffRow{-1, cmp.t1Row, nil}, nil
		}
	}
	return nil, io.EOF
}

func (cmp *Compare) diffPrimaryKey() int {
	if len(cmp.t1Row) == 0 {
		return 1
	}
	if len(cmp.t2Row) == 0 {
		return -1
	}
	for _, pk := range cmp.PK {
		ret := compareType(pk.Typ, cmp.t1Row[pk.Pos], cmp.t2Row[pk.Pos])
		if ret != 0 {
			return ret
		}
	}
	return 0
}

// ColumnPrimaryKey return  columns primary key
func ColumnPrimaryKey(pKeys []Pkey, row []string) []string {
	if row == nil {
		return nil
	}
	colPK := make([]string, 0, len(pKeys))
	for _, pk := range pKeys {
		colPK = append(colPK, row[pk.Pos])
	}
	return colPK
}

func (cmp *Compare) getPK() ([]Pkey, error) {
	t1d := cmp.t1.GetDefinition()
	t2d := cmp.t2.GetDefinition()
	var pos []int
	t2Pos, err := t2d.GetPKeyPos()
	if err == nil && len(t2Pos) > 0 {
		pos = t2Pos
	}
	t1Pos, err := t1d.GetPKeyPos()
	if err == nil && len(t1Pos) > 0 {
		pos = t1Pos
	}
	// If both do not have a primary key,
	// use all columns as primary keys.
	if len(pos) == 0 {
		t1Pos = make([]int, t1d.ColumnNum())
		for i := 0; i < t1d.ColumnNum(); i++ {
			t1Pos[i] = i
		}
		t2Pos = make([]int, t2d.ColumnNum())
		for i := 0; i < t2d.ColumnNum(); i++ {
			t2Pos[i] = i
		}
		pos = t1Pos
	}

	if len(t1Pos) > 0 && len(t2Pos) > 0 {
		if len(t1Pos) != len(t2Pos) {
			return nil, fmt.Errorf("primary key position")
		}
		for i := range pos {
			if t1Pos[i] != t2Pos[i] {
				return nil, fmt.Errorf("primary key position")
			}
		}
	}
	pk := make([]Pkey, len(pos))
	t1t := t1d.Types()
	t2t := t2d.Types()
	if (len(t1t) < len(pos)) || (len(t1t) != len(t2t)) {
		return nil, fmt.Errorf("mismatch data type: %d:%d", len(t1t), len(t2t))
	}
	for i, v := range pos {
		if t1t[i] != t2t[i] {
			return nil, fmt.Errorf("mismatch data type: %s:%s", t1t[i], t2t[i])
		}
		n := t1d.Names()
		pk[i] = Pkey{v, n[i], t1t[i]}
	}
	return pk, nil
}

func compareType(dType string, t1 string, t2 string) int {
	switch dType {
	case "int":
		return compareInt(t1, t2)
	case "bigint", "double precision", "numeric":
		return compareFloat(t1, t2)
	default:
		return strings.Compare(t1, t2)
	}
}

func compareInt(t1 string, t2 string) int {
	var err error
	var a, b int
	if a, err = strconv.Atoi(t1); err != nil {
		return strings.Compare(t1, t2)
	}
	if b, err = strconv.Atoi(t2); err != nil {
		return strings.Compare(t1, t2)
	}
	if a > b {
		return 1
	} else if a < b {
		return -1
	}
	return 0
}

func compareFloat(t1 string, t2 string) int {
	var err error
	var a, b float64
	if a, err = strconv.ParseFloat(t1, 64); err != nil {
		return strings.Compare(t1, t2)
	}
	if b, err = strconv.ParseFloat(t2, 64); err != nil {
		return strings.Compare(t1, t2)
	}
	if a > b {
		return 1
	} else if a < b {
		return -1
	}
	return 0
}
