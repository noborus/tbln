package tbln

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

// A Scanner reads records from a TBLN-encoded file.
type Scanner struct {
	reader  *bufio.Reader
	escrep  *regexp.Regexp
	Comment string   // comment
	Extra   []string // extra data
	Record  []string
	err     error
}

// NewScanner returns a new Reader that reads from r.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		reader: bufio.NewReader(r),
		escrep: regexp.MustCompile(`\|(\|+)`),
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
	body := str[2:]
	switch str[:2] {
	case "| ":
		body = strings.TrimRight(body, " |")
		rec := strings.Split(body, " | ")
		// Unescape vertical bars || -> |
		for i, column := range rec {
			if strings.Contains(column, "|") {
				rec[i] = r.escrep.ReplaceAllString(column, "$1")
			}
		}
		r.Record = rec
		return Record, nil
	case "# ":
		r.Comment = body
		return Comment, nil
	case "; ":
		ext := strings.Split(body, ":")
		if len(ext) != 2 {
			r.err = fmt.Errorf("Error: Extra format error %s", ext)
			return Extra, r.err
		}
		r.Extra = ext
		// r.Extra[ext[0]] = ext[1]
		return Extra, nil
	}
	r.err = fmt.Errorf("Error: Unsupported line")
	return Zero, r.err
}

// Err is Scanner error
func (r *Scanner) Err() error {
	if r.err == io.EOF {
		return nil
	}
	return r.err
}
