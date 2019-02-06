package tbln

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Writer is writer struct.
type Writer struct {
	Definition
	Writer *bufio.Writer
}

// NewWriter is Writer
func NewWriter(writer io.Writer, tbl Definition) *Writer {
	return &Writer{
		Definition: tbl,
		Writer:     bufio.NewWriter(writer),
	}
}

// WriteInfo is output table definition.
func (tw *Writer) WriteInfo() error {
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
		_, err := tw.Writer.WriteString(fmt.Sprintf("# %s\n", comment))
		if err != nil {
			return err
		}
	}
	return nil
}

func (tw *Writer) writeExtra(t Definition) error {
	for key, value := range tw.Ext {
		_, err := tw.Writer.WriteString(fmt.Sprintf("; %s: %s\n", key, value))
		if err != nil {
			return err
		}
	}
	if len(tw.Names) > 0 {
		_, err := tw.Writer.WriteString("; name: ")
		if err != nil {
			return err
		}
		err = tw.WriteRow(tw.Names)
		if err != nil {
			return err
		}
	}
	if len(tw.Types) > 0 {
		_, err := tw.Writer.WriteString("; type: ")
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
	_, err := tw.Writer.WriteString("|")
	if err != nil {
		return err
	}
	for _, column := range row {
		if strings.Contains(column, "|") {
			column = ESCREP.ReplaceAllString(column, "|$1")
		}
		_, err := tw.Writer.WriteString(" " + column + " |")
		if err != nil {
			return err
		}
	}
	_, err = tw.Writer.WriteString("\n")
	if err != nil {
		return err
	}
	return nil
}
