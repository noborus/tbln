package tbln

import (
	"bufio"
	"io"
	"strings"
)

// A Reader reads records from a TBLN-encoded file.
type Reader struct {
	lint   int
	reader *bufio.Reader
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{
		reader: bufio.NewReader(r),
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
			return ret, nil
		}
	}
}
