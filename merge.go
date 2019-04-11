package tbln

import (
	"fmt"
	"io"
)

type Merge struct {
	*Definition
	w *Writer
}

func NewMerge(iow io.Writer) *Merge {
	return &Merge{
		Definition: NewDefinition(),
		w:          NewWriter(iow),
	}
}

func (m *Merge) Header() error {
	return m.w.WriteDefinition(m.Definition)
}

func (m *Merge) Same(row []string) error {
	return m.w.WriteRow(row)
}

func (m *Merge) Add(row []string) error {
	return m.w.WriteRow(row)
}

func (m *Merge) Mod(srow []string, drow []string) error {
	return m.w.WriteRow(drow)
}

func (m *Merge) Del(row []string) error {
	return m.w.WriteRow(row)
}

func mergeComment(src, dst *Reader) []string {
	comments := make([]string, len(src.Comments))
	copy(comments, src.Comments)

	for _, dc := range dst.Comments {
		appFlag := true
		for _, sc := range src.Comments {
			if dc == sc {
				appFlag = false
				break
			}
		}
		if appFlag {
			comments = append(comments, dc)
		}
	}
	return comments
}

func (m *Merge) MergeDefinition(src, dst *Reader) (*Definition, error) {
	if src.columnNum != dst.columnNum {
		return nil, fmt.Errorf("different column num")
	}
	d := NewDefinition()
	d.Comments = mergeComment(src, dst)
	d.Extras = src.Extras
	for k, v := range dst.Extras {
		d.Extras[k] = v
	}
	return d, nil
}
