package tbln

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// A Scanner reads records from a TBLN-encoded file.
type Scanner struct {
	reader  *bufio.Reader
	Comment string   // comment
	Extra   []string // extra data
	Record  []string
	err     error
}

// NewScanner returns a new Reader that reads from r.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		reader: bufio.NewReader(r),
		Record: make([]string, 0),
	}
}

// ScanType is Type of line to return
type ScanType int

// Scan Type
const (
	Zero = iota
	Record
	Comment
	Extra
)

// Scan reads one record from r.
func (r *Scanner) Scan() (ScanType, error) {
	line, _, err := r.reader.ReadLine()
	if err != nil {
		r.err = err
		return Zero, err
	}
	str := string(line)
	if len(str) == 0 {
		return Zero, nil
	}
	if len(str) < 3 {
		r.err = fmt.Errorf("Error line: %s", str)
		return Zero, r.err
	}
	switch str[:2] {
	case "| ":
		r.Record = parseRecord(str)
		return Record, nil
	case "# ":
		r.Comment = str[2:]
		return Comment, nil
	case "; ":
		ext := strings.Split(str[2:], ":")
		if len(ext) != 2 {
			r.err = fmt.Errorf("Error: Extra format error %s", ext)
			return Extra, r.err
		}
		ext[1] = strings.TrimLeft(ext[1], " ")
		r.Extra = ext
		return Extra, nil
	}
	r.err = fmt.Errorf("Error: Unsupported line")
	return Zero, r.err
}

func parseRecord(body string) []string {
	body = strings.TrimLeft(body, "| ")
	body = strings.TrimRight(body, " |")
	rec := strings.Split(body, " | ")
	unescape(rec)
	return rec
}

func unescape(rec []string) {
	// Unescape vertical bars || -> |
	for i, column := range rec {
		if strings.Contains(column, "|") {
			rec[i] = UNESCREP.ReplaceAllString(column, "$1")
		}
	}
}

// Err is Scanner error
func (r *Scanner) Err() error {
	if r.err == io.EOF {
		return nil
	}
	return r.err
}
