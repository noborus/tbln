package tbln

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

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
			fields:  fields{Definition: Definition{columnNum: 2}, Reader: bufio.NewReader(bytes.NewBufferString("| a |\n"))},
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
		// TODO: Add test cases.
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

func Test_parseRecord(t *testing.T) {
	type args struct {
		body string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseRecord(tt.args.body); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseRecord() = %v, want %v", got, tt.want)
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
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unescape(tt.args.rec)
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
		// TODO: Add test cases.
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
