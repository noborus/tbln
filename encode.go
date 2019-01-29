package tbln

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Encode is writes TBLN record to w.
func (tbln *Tbln) Encode(w io.Writer) error {
	writer := bufio.NewWriter(w)
	tbln.encodeComment(writer)
	tbln.encodeExtra(writer)
	for _, row := range tbln.Rows {
		err := tbln.encodeRow(writer, row)
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func (tbln *Tbln) encodeComment(writer *bufio.Writer) error {
	for _, comment := range tbln.Comments {
		_, err := writer.WriteString(fmt.Sprintf("# %s\n", comment))
		if err != nil {
			return err
		}
	}
	return nil
}

func (tbln *Tbln) encodeExtra(writer *bufio.Writer) error {
	for _, extra := range tbln.Extra {
		_, err := writer.WriteString(fmt.Sprintf("; %s\n", extra))
		if err != nil {
			return err
		}
	}
	if len(tbln.Names) > 0 {
		_, err := writer.WriteString("; name: ")
		if err != nil {
			return err
		}
		err = tbln.encodeRow(writer, tbln.Names)
	}
	if len(tbln.Types) > 0 {
		_, err := writer.WriteString("; type: ")
		if err != nil {
			return err
		}
		err = tbln.encodeRow(writer, tbln.Types)
	}
	return nil
}

func (tbln *Tbln) encodeRow(writer *bufio.Writer, row []string) error {
	_, err := writer.WriteString("|")
	if err != nil {
		return err
	}
	for _, column := range row {
		if strings.Contains(column, "|") {
			column = tbln.escrep.ReplaceAllString(column, "|$1")
		}
		_, err := writer.WriteString(" " + column + " |")
		if err != nil {
			return err
		}
	}
	_, err = writer.WriteString("\n")
	if err != nil {
		return err
	}
	return nil
}
