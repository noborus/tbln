package tbln

import (
	"bytes"
	"io"
	"testing"
)

func TestNewScanner(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"test1", args{r: bytes.NewBufferString("")}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewScanner(tt.args.r); len(got.Comment) != tt.want {
				t.Errorf("NewScanner() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScanner_Scan(t *testing.T) {
	tests := []struct {
		name    string
		scanner *Scanner
		want    ScanType
		wantErr bool
	}{
		{"test1", NewScanner(bytes.NewBufferString("# comment")), Comment, false},
		{"test2", NewScanner(bytes.NewBufferString("; name: |a|")), Extra, false},
		{"test3", NewScanner(bytes.NewBufferString("| a |")), Record, false},
		{"test4", NewScanner(bytes.NewBufferString("\n")), Zero, false},
		{"test5", NewScanner(bytes.NewBufferString("")), Zero, true},
		{"test6", NewScanner(bytes.NewBufferString("a")), Zero, true},
		{"test7", NewScanner(bytes.NewBufferString(";;")), Zero, true},
		{"test8", NewScanner(bytes.NewBufferString("; a:b:c")), Extra, true},
		{"test9", NewScanner(bytes.NewBufferString("| a||b |")), Record, false},
		{"test10", NewScanner(bytes.NewBufferString("a b")), Zero, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.scanner
			got, err := r.Scan()
			if (err != nil) != tt.wantErr {
				t.Errorf("Scanner.Scan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Scanner.Scan() = %v, want %v", got, tt.want)
			}
		})
	}
}
