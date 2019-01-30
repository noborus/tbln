package tbln

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Encode is writes TBLN record to writer.
func (tbln *Tbln) Encode(writer io.Writer) error {
	w := bufio.NewWriter(writer)
	tbln.encodeComment(w)
	tbln.encodeExtra(w)
	for _, row := range tbln.Rows {
		err := tbln.encodeRow(w, row)
		if err != nil {
			return err
		}
	}
	return w.Flush()
}

func (tbln *Tbln) encodeComment(w *bufio.Writer) error {
	for _, comment := range tbln.Comments {
		_, err := w.WriteString(fmt.Sprintf("# %s\n", comment))
		if err != nil {
			return err
		}
	}
	return nil
}

func (tbln *Tbln) encodeExtra(w *bufio.Writer) error {
	for _, extra := range tbln.Extra {
		_, err := w.WriteString(fmt.Sprintf("; %s\n", extra))
		if err != nil {
			return err
		}
	}
	if len(tbln.Names) > 0 {
		_, err := w.WriteString("; name: ")
		if err != nil {
			return err
		}
		err = tbln.encodeRow(w, tbln.Names)
	}
	if len(tbln.Types) > 0 {
		_, err := w.WriteString("; type: ")
		if err != nil {
			return err
		}
		err = tbln.encodeRow(w, tbln.Types)
	}
	return nil
}

func (tbln *Tbln) encodeRow(w *bufio.Writer, row []string) error {
	_, err := w.WriteString("|")
	if err != nil {
		return err
	}
	for _, column := range row {
		if strings.Contains(column, "|") {
			column = ESCREP.ReplaceAllString(column, "|$1")
		}
		_, err := w.WriteString(" " + column + " |")
		if err != nil {
			return err
		}
	}
	_, err = w.WriteString("\n")
	if err != nil {
		return err
	}
	return nil
}
