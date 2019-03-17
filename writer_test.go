package tbln

import (
	"bytes"
	"testing"
)

func TestWriter_WriteRow(t *testing.T) {
	type args struct {
		row []string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{[]string{"a", "b"}},
			want:    "| a | b |\n",
			wantErr: false,
		},
		{
			name:    "test2",
			args:    args{[]string{"a|b", "b"}},
			want:    "| a||b | b |\n",
			wantErr: false,
		},
		{
			name:    "test3",
			args:    args{[]string{"a||b", "b"}},
			want:    "| a|||b | b |\n",
			wantErr: false,
		},
		{
			name:    "test4",
			args:    args{[]string{"", ""}},
			want:    "|  |  |\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			w := &Writer{
				Writer: buf,
			}
			if err := w.WriteRow(tt.args.row); (err != nil) != tt.wantErr {
				t.Errorf("Writer.WriteRow() error = %v, wantErr %v", err, tt.wantErr)
			}
			if buf.String() != tt.want {
				t.Errorf("Writer.WriteRow() = [%v], want [%v]", buf, tt.want)
			}
		})
	}
}

func TestWriteAll(t *testing.T) {
	type args struct {
		tbln *Tbln
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
		wantErr    bool
	}{
		{
			name:       "test1",
			args:       args{tbln: NewTbln()},
			wantWriter: "",
			wantErr:    false,
		},
		{
			name:       "test2",
			args:       args{tbln: &Tbln{Definition: NewDefinition(), Rows: [][]string{{"a", "b"}}}},
			wantWriter: "| a | b |\n",
			wantErr:    false,
		},
		{
			name:       "test3",
			args:       args{tbln: &Tbln{Definition: NewDefinition(), Rows: [][]string{{"a", "b"}, {"c", "d"}}}},
			wantWriter: "| a | b |\n| c | d |\n",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			at := tt.args.tbln
			delete(at.Extras, "created_at")
			writer := &bytes.Buffer{}
			if err := WriteAll(writer, tt.args.tbln); (err != nil) != tt.wantErr {
				t.Errorf("WriteAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("WriteAll() = [%v], want [%v]", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestWriter_writeComment(t *testing.T) {
	type args struct {
		d *Definition
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
		wantErr    bool
	}{
		{
			name:       "test1",
			args:       args{d: &Definition{}},
			wantWriter: "",
			wantErr:    false,
		},
		{
			name:       "test2",
			args:       args{d: &Definition{Comments: []string{"comment"}}},
			wantWriter: "# comment\n",
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			w := &Writer{
				Writer: buf,
			}
			if err := w.writeComment(tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("Writer.writeComment() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotWriter := buf.String(); gotWriter != tt.wantWriter {
				t.Errorf("WriteAll() = [%v], want [%v]", gotWriter, tt.wantWriter)
			}
		})
	}
}
