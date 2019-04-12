package tbln

import "io"

// SelfReader reads records from Tbln.
type SelfReader struct {
	*Tbln
	readRowNum int
}

// NewSelfReader return TblnReader.
func NewSelfReader(t *Tbln) *SelfReader {
	return &SelfReader{
		Tbln:       t,
		readRowNum: 0,
	}
}

// ReadRow one record (a slice of fields) from Tbln.
func (rr *SelfReader) ReadRow() ([]string, error) {
	if rr.readRowNum >= rr.RowNum {
		return nil, io.EOF
	}
	row := rr.Rows[rr.readRowNum]
	rr.readRowNum++
	return row, nil
}
