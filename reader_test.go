package tbln

import (
	"bufio"
	"bytes"
	"io"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReader_scanLine(t *testing.T) {
	type fields struct {
		Definition *Definition
		r          *bufio.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name:    "test1",
			fields:  fields{r: bufio.NewReader(bytes.NewBufferString("foo\n"))},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "testRow1",
			fields:  fields{r: bufio.NewReader(bytes.NewBufferString("| 1 |\n"))},
			want:    []string{"1"},
			wantErr: false,
		},
		{
			name:    "testRow2",
			fields:  fields{r: bufio.NewReader(bytes.NewBufferString("| 1 | 2 |\n"))},
			want:    []string{"1", "2"},
			wantErr: false,
		},
		{
			// Do not read after empty line
			name:    "testRowBlank",
			fields:  fields{r: bufio.NewReader(bytes.NewBufferString("| 1 |\n\n| 2 |\n"))},
			want:    []string{"1"},
			wantErr: false,
		},
		{
			name:    "testComment1",
			fields:  fields{r: bufio.NewReader(bytes.NewBufferString("# comment\n"))},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "testComment2",
			fields:  fields{r: bufio.NewReader(bytes.NewBufferString("# comment\n| 1 |\n"))},
			want:    []string{"1"},
			wantErr: false,
		},
		{
			name:    "testExtra1",
			fields:  fields{r: bufio.NewReader(bytes.NewBufferString("; test: test\n"))},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "testExtra2",
			fields:  fields{r: bufio.NewReader(bytes.NewBufferString("# test: test\n| 1 |\n"))},
			want:    []string{"1"},
			wantErr: false,
		},
		{
			name:    "testExtraErr1",
			fields:  fields{r: bufio.NewReader(bytes.NewBufferString("; test\n| 1 |\n"))},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "testExtraErr2",
			fields:  fields{r: bufio.NewReader(bytes.NewBufferString("; Hash: |\n"))},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tr := &FileReader{
				Definition: NewDefinition(),
				r:          tt.fields.r,
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

func TestReader_analyzeExtra(t *testing.T) {
	type fields struct {
		Definition *Definition
		r          *bufio.Reader
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
			tr := &FileReader{
				Definition: tt.fields.Definition,
				r:          tt.fields.r,
			}
			if err := tr.analyzeExtra(tt.args.extstr); (err != nil) != tt.wantErr {
				t.Errorf("Reader.analyzeExtra() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReader_ReadRow(t *testing.T) {
	type fields struct {
		Definition *Definition
		r          *bufio.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name:    "test1",
			fields:  fields{Definition: NewDefinition(), r: bufio.NewReader(bytes.NewBufferString("| a | b |"))},
			want:    []string{"a", "b"},
			wantErr: false,
		},
		{
			name:    "test2",
			fields:  fields{Definition: NewDefinition(), r: bufio.NewReader(bytes.NewBufferString("| a||b | ||b |"))},
			want:    []string{"a|b", "|b"},
			wantErr: false,
		},
		{
			name:    "test3",
			fields:  fields{Definition: &Definition{columnNum: 2}, r: bufio.NewReader(bytes.NewBufferString("| a |\n"))},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &FileReader{
				Definition: tt.fields.Definition,
				r:          tt.fields.r,
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
		want    *TBLN
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{reader: bufio.NewReader(bytes.NewBufferString("| a | b |"))},
			want:    &TBLN{Definition: NewDefinition(), Rows: [][]string{{"a", "b"}}},
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
func TestReadFile(t *testing.T) {
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "testSimple",
			args:    args{reader: openFile(t, filepath.Join("testdata", "simple.tbln"))},
			want:    "simple",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadAll(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !got.Verify() {
				t.Error("verification failure")
			}
			if !reflect.DeepEqual(got.TableName(), tt.want) {
				t.Errorf("ReadFile() = %v, want %v", got.TableName(), tt.want)
			}
		})
	}
}

func Test_unescape(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test1", args: args{str: "||"}, want: "|"},
		{name: "test2", args: args{str: "|||"}, want: "||"},
		{name: "test3", args: args{str: "a||b"}, want: "a|b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := unescape(tt.args.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unescape() = %v, want %v", got, tt.want)
			}
		})
	}
}
