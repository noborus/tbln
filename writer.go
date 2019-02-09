package tbln

import (
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
func NewWriter(writer io.Writer, tbl Definition) *Writer {
	return &Writer{
		Definition: tbl,
		Writer:     writer,
	}
}

// WriteDefinition is write table definition.
func (tw *Writer) WriteDefinition() error {
	t := tw.Definition
	err := tw.writeComment(t)
	if err != nil {
		return err
	}
	err = tw.writeExtra(t)
	if err != nil {
		return err
	}
	return nil
}

func (tw *Writer) writeComment(t Definition) error {
	for _, comment := range t.Comments {
		_, err := io.WriteString(tw.Writer, fmt.Sprintf("# %s\n", comment))
		if err != nil {
			return err
		}
	}
	return nil
}

func (tw *Writer) writeExtra(t Definition) error {
	for key, value := range tw.Ext {
		_, err := fmt.Fprintf(tw.Writer, "; %s: %s\n", key, value)
		if err != nil {
			return err
		}
	}
	if len(tw.name) > 0 {
		_, err := fmt.Fprintf(tw.Writer, "; TableName: %s\n", tw.name)
		if err != nil {
			return err
		}
	}
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
