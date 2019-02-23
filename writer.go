package tbln

import (
	"fmt"
	"io"
)

// Writer is writer struct.
type Writer struct {
	Writer io.Writer
}

// NewWriter is Writer
func NewWriter(writer io.Writer) *Writer {
	return &Writer{
		Writer: writer,
	}
}

// WriteDefinition is write tbln definition.
func (w *Writer) WriteDefinition(d Definition) error {
	err := w.writeComment(d)
	if err != nil {
		return err
	}
	return w.writeExtra(d)
}

// WriteRow is write one row.
func (w *Writer) WriteRow(row []string) error {
	_, err := io.WriteString(w.Writer, JoinRow(row)+"\n")
	return err
}

// WriteAll write tbln.
func WriteAll(writer io.Writer, tbln *Tbln) error {
	w := NewWriter(writer)
	err := w.writeExtraTarget(tbln.Definition, false)
	if err != nil {
		return err
	}
	if tbln.buffer.Len() > 0 {
		for n, v := range tbln.Hash {
			_, err := fmt.Fprintf(w.Writer, "; %s: %s\n", n, v)
			if err != nil {
				return err
			}
		}
		_, err := w.Writer.Write(tbln.buffer.Bytes())
		if err != nil {
			return err
		}
	} else {
		err := w.writeExtraTarget(tbln.Definition, true)
		if err != nil {
			return err
		}
		for _, row := range tbln.Rows {
			err = w.WriteRow(row)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (w *Writer) writeComment(d Definition) error {
	for _, comment := range d.Comments {
		_, err := io.WriteString(w.Writer, fmt.Sprintf("# %s\n", comment))
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) writeExtra(d Definition) error {
	err := w.writeExtraTarget(d, false)
	if err != nil {
		return err
	}
	return w.writeExtraTarget(d, true)
}

func (w *Writer) writeExtraTarget(d Definition, targetFlag bool) error {
	ext := sortExtra(d.Ext)
	for _, entry := range ext {
		if entry.v.hashTarget != targetFlag {
			continue
		}
		_, err := fmt.Fprintf(w.Writer, "; %s: %s\n", entry.n, entry.v.value)
		if err != nil {
			return err
		}
	}
	return nil
}
