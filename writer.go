package tbln

import (
	"fmt"
	"io"
	"sort"
)

// Writer writes records to a tbln encoded file.
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
func WriteAll(writer io.Writer, tbln *Tbln) error {
	w := NewWriter(writer)
	err := w.WriteDefinition(tbln.Definition)
	if err != nil {
		return err
	}
	for _, row := range tbln.Rows {
		err = w.WriteRow(row)
		if err != nil {
			return err
		}
	}
	return nil
}

// WriteDefinition writes Definition (comment and extra) to w.
func (w *Writer) WriteDefinition(d *Definition) error {
	err := w.writeComment(d)
	if err != nil {
		return err
	}
	return w.writeExtra(d)
}

// writeComment writes comment to w.
func (w *Writer) writeComment(d *Definition) error {
	for _, comment := range d.Comments {
		_, err := io.WriteString(w.Writer, fmt.Sprintf("# %s\n", comment))
		if err != nil {
			return err
		}
	}
	return nil
}

// writeExtra writes extra to w.
func (w *Writer) writeExtra(d *Definition) error {
	err := w.writeExtraTarget(d, false)
	if err != nil {
		return err
	}
	if len(d.Signs) > 0 {
		err = w.writeSigns(d)
		if err != nil {
			return err
		}
	}
	if len(d.Hashes) > 0 {
		err = w.writeHashes(d)
		if err != nil {
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
		_, err := fmt.Fprintf(w.Writer, "; %s: %s\n", k, d.Extras[k].value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) writeSigns(d *Definition) error {
	signs := make([]string, 0, len(d.Signs))
	for k, v := range d.Signs {
		signs = append(signs, k+":"+fmt.Sprintf("%x", v))
	}
	_, err := fmt.Fprintf(w.Writer, "; Signature: %s\n", JoinRow(signs))
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) writeHashes(d *Definition) error {
	hash := d.SerialHash()
	_, err := fmt.Fprintf(w.Writer, "; Hash: %s\n", string(hash))
	if err != nil {
		return err
	}
	return nil
}
