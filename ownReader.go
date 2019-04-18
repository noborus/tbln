package tbln

import "io"

// OwnReader reads records from Tbln.
type OwnReader struct {
	*Tbln
	readRowNum int
}

// NewOwnReader return TblnReader.
func NewOwnReader(t *Tbln) *OwnReader {
	return &OwnReader{
		Tbln:       t,
		readRowNum: 0,
	}
}

// ReadRow one record (a slice of fields) from Tbln.
func (rr *OwnReader) ReadRow() ([]string, error) {
	if rr.readRowNum >= rr.RowNum {
		return nil, io.EOF
	}
	row := rr.Rows[rr.readRowNum]
	rr.readRowNum++
	return row, nil
}
