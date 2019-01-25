package tbln

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

// A Writer wites records to a TBLN encoded file.
type Writer struct {
	writer *bufio.Writer
	escrep *regexp.Regexp
}

// NewWriter is a new TBLN Writer that writes to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: bufio.NewWriter(w),
		escrep: regexp.MustCompile(`(\|+)`),
	}
}

// Writer is writes a single TBLN record to w.
func (w *Writer) Write(record []string) error {
	w.writer.WriteString("|")
	for _, column := range record {
		if strings.Contains(column, "|") {
			column = w.escrep.ReplaceAllString(column, "|$1")
		}
		w.writer.WriteString(" " + column + " |")
	}
	w.writer.WriteString("\n")
	return w.writer.Flush()
}
