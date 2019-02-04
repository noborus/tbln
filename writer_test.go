package tbln

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

func TestNewWriter(t *testing.T) {
	type args struct {
		tbl Table
	}
	tests := []struct {
		name       string
		args       args
		want       *Writer
		wantWriter string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			if got := NewWriter(writer, tt.args.tbl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWriter() = %v, want %v", got, tt.want)
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("NewWriter() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}

func TestWriter_WriteInfo(t *testing.T) {
	type fields struct {
		Table  Table
		Writer *bufio.Writer
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tw := &Writer{
				Table:  tt.fields.Table,
				Writer: tt.fields.Writer,
			}
			if err := tw.WriteInfo(); (err != nil) != tt.wantErr {
				t.Errorf("Writer.WriteInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWriter_writeComment(t *testing.T) {
	type fields struct {
		Table  Table
		Writer *bufio.Writer
	}
	type args struct {
		t Table
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
			tw := &Writer{
				Table:  tt.fields.Table,
				Writer: tt.fields.Writer,
			}
			if err := tw.writeComment(tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("Writer.writeComment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWriter_writeExtra(t *testing.T) {
	type fields struct {
		Table  Table
		Writer *bufio.Writer
	}
	type args struct {
		t Table
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
			tw := &Writer{
				Table:  tt.fields.Table,
				Writer: tt.fields.Writer,
			}
			if err := tw.writeExtra(tt.args.t); (err != nil) != tt.wantErr {
				t.Errorf("Writer.writeExtra() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestWriter_WriteRow(t *testing.T) {
	type fields struct {
		Table  Table
		Writer *bufio.Writer
	}
	type args struct {
		row []string
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
			tw := &Writer{
				Table:  tt.fields.Table,
				Writer: tt.fields.Writer,
			}
			if err := tw.WriteRow(tt.args.row); (err != nil) != tt.wantErr {
				t.Errorf("Writer.WriteRow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
