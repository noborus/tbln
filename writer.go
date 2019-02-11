package tbln

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"strings"
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

// WriteDefinition is write table definition.
func (w *Writer) WriteDefinition(d Definition) error {
	err := w.writeComment(d)
	if err != nil {
		return err
	}
	err = w.writeExtra(d)
	if err != nil {
		return err
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
	err := w.writeExtraWithOutHash(d)
	if err != nil {
		return err
	}
	return w.writeExtraWithHash(d)
}

func (w *Writer) writeExtraWithOutHash(d Definition) error {
	if len(d.tableName) > 0 {
		_, err := fmt.Fprintf(w.Writer, "; TableName: %s\n", d.tableName)
		if err != nil {
			return err
		}
	}
	for key, extra := range d.Ext {
		if extra.hashing {
			continue
		}
		_, err := fmt.Fprintf(w.Writer, "; %s: %s\n", key, extra.value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) writeExtraWithHash(d Definition) error {
	if len(d.Names) > 0 {
		_, err := io.WriteString(w.Writer, "; name: ")
		if err != nil {
			return err
		}
		err = w.WriteRow(d.Names)
		if err != nil {
			return err
		}
	}
	if len(d.Types) > 0 {
		_, err := io.WriteString(w.Writer, "; type: ")
		if err != nil {
			return err
		}
		err = w.WriteRow(d.Types)
		if err != nil {
			return err
		}
	}
	for key, extra := range d.Ext {
		if !extra.hashing {
			continue
		}
		_, err := fmt.Fprintf(w.Writer, "; %s: %s\n", key, extra.value)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteRow is write one row.
func (w *Writer) WriteRow(row []string) error {
	_, err := io.WriteString(w.Writer, "|")
	if err != nil {
		return err
	}
	for _, column := range row {
		if strings.Contains(column, "|") {
			column = ESCREP.ReplaceAllString(column, "|$1")
		}
		_, err := io.WriteString(w.Writer, " "+column+" |")
		if err != nil {
			return err
		}
	}
	_, err = io.WriteString(w.Writer, "\n")
	if err != nil {
		return err
	}
	return nil
}

// WriteAll write all table.
func WriteAll(writer io.Writer, table *Table) error {
	w := NewWriter(writer)
	err := w.WriteDefinition(table.Definition)
	if err != nil {
		return err
	}
	for _, row := range table.Rows {
		err = w.WriteRow(row)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteHashAll write hash ...
func WriteHashAll(writer io.Writer, table *Table) error {
	w := NewWriter(writer)
	err := w.writeExtraWithOutHash(table.Definition)
	if err != nil {
		return err
	}
	b, err := targetBytes(table)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(b)
	_, err = fmt.Fprintf(w.Writer, "; sha256: %x\n", sum)
	if err != nil {
		return err
	}
	_, err = w.Writer.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func targetBytes(table *Table) ([]byte, error) {
	var b bytes.Buffer
	writer := bufio.NewWriter(&b)
	bw := NewWriter(writer)
	err := bw.writeExtraWithHash(table.Definition)
	if err != nil {
		return nil, err
	}
	for _, row := range table.Rows {
		err = bw.WriteRow(row)
		if err != nil {
			return nil, err
		}
	}
	writer.Flush()
	return b.Bytes(), nil
}
