package tbln

import (
	"bufio"
	"testing"
)

func TestWriter_WriteRow(t *testing.T) {
	type fields struct {
		Definition Definition
		Writer     *bufio.Writer
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
				Writer: tt.fields.Writer,
			}
			if err := tw.WriteRow(tt.args.row); (err != nil) != tt.wantErr {
				t.Errorf("Writer.WriteRow() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
