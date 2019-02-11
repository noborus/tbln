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
	Definition
	Writer io.Writer
}

// NewWriter is Writer
func NewWriter(writer io.Writer, definition Definition) *Writer {
	return &Writer{
		Definition: definition,
		Writer:     writer,
	}
}

// WriteDefinition is write table definition.
func (tw *Writer) WriteDefinition() error {
	err := tw.writeComment()
	if err != nil {
		return err
	}
	err = tw.writeExtra()
	if err != nil {
		return err
	}
	return nil
}

func (tw *Writer) writeComment() error {
	for _, comment := range tw.Comments {
		_, err := io.WriteString(tw.Writer, fmt.Sprintf("# %s\n", comment))
		if err != nil {
			return err
		}
	}
	return nil
}

func (tw *Writer) writeExtra() error {
	err := tw.writeExtraWithOutHash()
	if err != nil {
		return err
	}
	return tw.writeExtraWithHash()
}

func (tw *Writer) writeExtraWithOutHash() error {
	if len(tw.tableName) > 0 {
		_, err := fmt.Fprintf(tw.Writer, "; TableName: %s\n", tw.tableName)
		if err != nil {
			return err
		}
	}
	for key, extra := range tw.Ext {
		if extra.hashing {
			continue
		}
		_, err := fmt.Fprintf(tw.Writer, "; %s: %s\n", key, extra.value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (tw *Writer) writeExtraWithHash() error {
	if len(tw.Names) > 0 {
		_, err := io.WriteString(tw.Writer, "; name: ")
		if err != nil {
			return err
		}
		err = tw.WriteRow(tw.Names)
		if err != nil {
			return err
		}
	}
	if len(tw.Types) > 0 {
		_, err := io.WriteString(tw.Writer, "; type: ")
		if err != nil {
			return err
		}
		err = tw.WriteRow(tw.Types)
		if err != nil {
			return err
		}
	}
	for key, extra := range tw.Ext {
		if !extra.hashing {
			continue
		}
		_, err := fmt.Fprintf(tw.Writer, "; %s: %s\n", key, extra.value)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteRow is write one row.
func (tw *Writer) WriteRow(row []string) error {
	_, err := io.WriteString(tw.Writer, "|")
	if err != nil {
		return err
	}
	for _, column := range row {
		if strings.Contains(column, "|") {
			column = ESCREP.ReplaceAllString(column, "|$1")
		}
		_, err := io.WriteString(tw.Writer, " "+column+" |")
		if err != nil {
			return err
		}
	}
	_, err = io.WriteString(tw.Writer, "\n")
	if err != nil {
		return err
	}
	return nil
}

// WriteAll write all table.
func WriteAll(writer io.Writer, table *Table) error {
	tw := NewWriter(writer, table.Definition)
	err := tw.WriteDefinition()
	if err != nil {
		return err
	}
	for _, row := range table.Rows {
		err = tw.WriteRow(row)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteHashAll write hash ...
func WriteHashAll(writer io.Writer, table *Table) error {
	tw := NewWriter(writer, table.Definition)
	err := tw.writeExtraWithOutHash()
	if err != nil {
		return err
	}
	tmp := tw.Writer

	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	tw.Writer = w
	err = tw.writeExtraWithHash()
	if err != nil {
		return err
	}
	for _, row := range table.Rows {
		err = tw.WriteRow(row)
		if err != nil {
			return err
		}
	}
	w.Flush()
	tw.Writer = tmp
	sum := sha256.Sum256(b.Bytes())
	_, err = fmt.Fprintf(tw.Writer, "; sha256: %x\n", sum)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(tw.Writer, "%s\n", b.String())
	if err != nil {
		return err
	}
	return nil
}
