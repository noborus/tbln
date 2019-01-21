package tbln

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

// A Reader reads records from a TBLN-encoded file.
type Reader struct {
	lint   int
	reader *bufio.Reader
	escrep *regexp.Regexp
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		reader: bufio.NewReader(r),
		escrep: regexp.MustCompile(`\|(\|+)`),
	}
}

// Read reads one record from r. The record is a slice of string.
func (r *Reader) Read() ([]string, error) {
	var ret []string
	for {
		line, _, err := r.reader.ReadLine()
		if err != nil {
			return nil, err
		}
		str := string(line)
		if strings.HasPrefix(str, "| ") {
			str = str[2:]
			ret = strings.Split(str, " | ")
			for i, column := range ret {
				if strings.Contains(column, "|") {
					ret[i] = r.escrep.ReplaceAllString(column, "$1")
				}
			}
			return ret, nil
		}
	}
}
