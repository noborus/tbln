package tbln

import (
	"bufio"
	"io"
	"strings"
)

// A Writer wites records to a tbln encoded file.
type Writer struct {
	writer *bufio.Writer
}

// NewWriter is a new tbln Writer that writes to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: bufio.NewWriter(w),
	}
}

// Writer is writes a single LTSV record to w.
func (w *Writer) Write(record []string) error {
	w.writer.WriteString("| ")
	w.writer.WriteString(strings.Join(record, " | "))
	w.writer.WriteString(" |")
	w.writer.WriteByte('\n')
	return w.writer.Flush()
}
