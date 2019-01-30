package tbln

import (
	"fmt"
	"io"
)

// Decode is read TBLN record from io.Reader.
func (tbln *Tbln) Decode(reader io.Reader) error {
	scan := NewScanner(reader)
	for {
		t, err := scan.Scan()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		switch t {
		case Record:
			if tbln.colNum <= 0 {
				tbln.colNum = len(scan.Record)
			}
			if tbln.colNum == len(scan.Record) {
				tbln.AddRow(scan.Record)
			} else {
				return fmt.Errorf("number of column is invalid %d != %v", tbln.colNum, scan.Record)
			}
		case Comment:
			tbln.Comments = append(tbln.Comments, scan.Comment)
		case Extra:
			tbln.analyzeExt(scan.Extra)
		}
	}
}

func (tbln *Tbln) analyzeExt(ext []string) error {
	switch ext[0] {
	case "name":
		tbln.SetNames(parseRecord(ext[1]))
	case "type":
		tbln.SetTypes(parseRecord(ext[1]))
	default:
		tbln.Extra = append(tbln.Extra, ext[0]+" : "+ext[1])
	}
	return nil
}
