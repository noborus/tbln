package tbln

import (
	"fmt"
	"io"
	"sort"
)

// Writer writes records to a TBLN encoded file.
type Writer struct {
	Writer io.Writer
}

// NewWriter returns a new Writer that writes to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		Writer: w,
	}
}

// WriteRow writes a single tbln record to w along with any necessary escaping.
// A record is a slice of strings with each string being one field.
func (w *Writer) WriteRow(row []string) error {
	_, err := io.WriteString(w.Writer, JoinRow(row)+"\n")
	return err
}

// WriteAll writes multiple tbln records to w using Write.
func WriteAll(writer io.Writer, tbln *TBLN) error {
	w := NewWriter(writer)
	if err := w.WriteDefinition(tbln.Definition); err != nil {
		return err
	}
	for _, row := range tbln.Rows {
		if err := w.WriteRow(row); err != nil {
			return err
		}
	}
	return nil
}

// WriteDefinition writes Definition (comment and extra) to w.
func (w *Writer) WriteDefinition(d *Definition) error {
	if err := w.writeComment(d); err != nil {
		return err
	}
	return w.writeExtra(d)
}

// writeComment writes comment to w.
func (w *Writer) writeComment(d *Definition) error {
	for _, comment := range d.Comments {
		if _, err := io.WriteString(w.Writer, fmt.Sprintf("# %s\n", comment)); err != nil {
			return err
		}
	}
	return nil
}

// writeExtra writes extra to w.
func (w *Writer) writeExtra(d *Definition) error {
	if err := w.writeExtraTarget(d, false); err != nil {
		return err
	}
	if len(d.Signs) > 0 {
		if err := w.writeSigns(d); err != nil {
			return err
		}
	}
	if len(d.Hashes) > 0 {
		if err := w.writeHashes(d); err != nil {
			return err
		}
	}
	return w.writeExtraTarget(d, true)
}

func (w *Writer) writeExtraTarget(d *Definition, targetFlag bool) error {
	keys := make([]string, 0, len(d.Extras))
	for k := range d.Extras {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if d.Extras[k].hashTarget != targetFlag {
			continue
		}
		if _, err := fmt.Fprintf(w.Writer, "; %s: %s\n", k, d.Extras[k].value); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) writeSigns(d *Definition) error {
	for k, v := range d.Signs {
		signs := make([]string, 0, 3)
		signs = append(signs, k)
		signs = append(signs, v.algorithm)
		signs = append(signs, fmt.Sprintf("%x", v.sign))
		if _, err := fmt.Fprintf(w.Writer, "; Signature: %s\n", JoinRow(signs)); err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) writeHashes(d *Definition) error {
	for k, v := range d.Hashes {
		hashes := make([]string, 0, 2)
		hashes = append(hashes, k)
		hashes = append(hashes, fmt.Sprintf("%x", v))
		if _, err := fmt.Fprintf(w.Writer, "; Hash: %s\n", JoinRow(hashes)); err != nil {
			return err
		}
	}
	return nil
}
