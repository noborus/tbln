package tbln

import (
	"bufio"
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestReader_scanLine(t *testing.T) {
	type fields struct {
		Definition Definition
		Reader     *bufio.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name:    "test1",
			fields:  fields{Reader: bufio.NewReader(bytes.NewBufferString("foo\n"))},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test2",
			fields:  fields{Reader: bufio.NewReader(bytes.NewBufferString("| 1 |\n"))},
			want:    []string{"1"},
			wantErr: false,
		},
		{
			name:    "test3",
			fields:  fields{Reader: bufio.NewReader(bytes.NewBufferString("| 1 | 2 |\n"))},
			want:    []string{"1", "2"},
			wantErr: false,
		},
		{
			name:    "test4",
			fields:  fields{Reader: bufio.NewReader(bytes.NewBufferString("# comment\n"))},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test5",
			fields:  fields{Reader: bufio.NewReader(bytes.NewBufferString("# comment\n| 1 |\n"))},
			want:    []string{"1"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Reader{
				Definition: tt.fields.Definition,
				Reader:     tt.fields.Reader,
			}
			got, err := tr.scanLine()
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.scanLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reader.scanLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_SplitRow(t *testing.T) {
	type args struct {
		body string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test1",
			args: args{body: "| 1 | 2 |"},
			want: []string{"1", "2"},
		},
		{
			name: "test2",
			args: args{body: "| 1 | || |"},
			want: []string{"1", "|"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitRow(tt.args.body); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReader_analyzeExt(t *testing.T) {
	type fields struct {
		Definition Definition
		Reader     *bufio.Reader
	}
	type args struct {
		extstr string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "test1",
			fields:  fields{Definition: NewDefinition()},
			args:    args{extstr: "; name: | id | name |"},
			wantErr: false,
		},
		{
			name:    "test2",
			fields:  fields{Definition: NewDefinition()},
			args:    args{extstr: "; type"},
			wantErr: true,
		},
		{
			name:    "test3",
			fields:  fields{Definition: NewDefinition()},
			args:    args{extstr: "; type: | int |"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Reader{
				Definition: tt.fields.Definition,
				Reader:     tt.fields.Reader,
			}
			if err := tr.analyzeExt(tt.args.extstr); (err != nil) != tt.wantErr {
				t.Errorf("Reader.analyzeExt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReader_ReadRow(t *testing.T) {
	type fields struct {
		Definition Definition
		Reader     *bufio.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name:    "test1",
			fields:  fields{Reader: bufio.NewReader(bytes.NewBufferString("| a | b |"))},
			want:    []string{"a", "b"},
			wantErr: false,
		},
		{
			name:    "test2",
			fields:  fields{Reader: bufio.NewReader(bytes.NewBufferString("| a||b | ||b |"))},
			want:    []string{"a|b", "|b"},
			wantErr: false,
		},
		{
			name:    "test3",
			fields:  fields{Definition: Definition{ColumnNum: 2}, Reader: bufio.NewReader(bytes.NewBufferString("| a |\n"))},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Reader{
				Definition: tt.fields.Definition,
				Reader:     tt.fields.Reader,
			}
			got, err := tr.ReadRow()
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.ReadRow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reader.ReadRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadAll(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    *Tbln
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{reader: bufio.NewReader(bytes.NewBufferString("| a | b |"))},
			want:    &Tbln{Definition: NewDefinition(), Rows: [][]string{{"a", "b"}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadAll(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Rows, tt.want.Rows) {
				t.Errorf("ReadAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unescape(t *testing.T) {
	type args struct {
		rec []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "test1", args: args{[]string{"||"}}, want: []string{"|"}},
		{name: "test2", args: args{[]string{"|||"}}, want: []string{"||"}},
		{name: "test3", args: args{[]string{"a||b"}}, want: []string{"a|b"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := unescape(tt.args.rec); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unescape() = %v, want %v", got, tt.want)
			}
		})
	}
}
