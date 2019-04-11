package tbln

import "io"

type Diff struct {
	*Definition
	w io.Writer
}

func NewDiff(w io.Writer) *Diff {
	return &Diff{
		Definition: NewDefinition(),
		w:          w,
	}
}

func (d Diff) Same(row []string) error {
	_, err := io.WriteString(d.w, " "+JoinRow(row)+"\n")
	return err
}

func (d Diff) Add(row []string) error {
	_, err := io.WriteString(d.w, "+"+JoinRow(row)+"\n")
	return err
}

func (d Diff) Mod(srow []string, drow []string) error {
	_, err := io.WriteString(d.w, "-"+JoinRow(srow)+"\n")
	if err != nil {
		return err
	}
	_, err = io.WriteString(d.w, "+"+JoinRow(drow)+"\n")
	return err
}

func (d Diff) Del(row []string) error {
	_, err := io.WriteString(d.w, "-"+JoinRow(row)+"\n")
	return err
}
