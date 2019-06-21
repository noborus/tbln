package tbln

import "io"

// OwnReader reads records from TBLN.
type OwnReader struct {
	*TBLN
	readRowNum int
}

// NewOwnReader return TBLNReader.
func NewOwnReader(t *TBLN) *OwnReader {
	return &OwnReader{
		TBLN:       t,
		readRowNum: 0,
	}
}

// ReadRow one record (a slice of fields) from TBLN.
func (rr *OwnReader) ReadRow() ([]string, error) {
	if rr.readRowNum >= rr.RowNum {
		return nil, io.EOF
	}
	row := rr.Rows[rr.readRowNum]
	rr.readRowNum++
	return row, nil
}
