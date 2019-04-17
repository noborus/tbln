package tbln

import (
	"bytes"
	"testing"
)

var TestDiff1 = `; name: | id | name | age |
; type: | int | text | int |
; primarykey: | id |
; TableName: test1
| 1 | Bob | 19 |
`

var TestDiff2 = `; name: | id | name | age |
; type: | int | text | int |
; primarykey: | id |
; TableName: test1
| 1 | Bob | 19 |
| 2 | Alice | 14 |
`

var TestDiff3 = `; name: | id | name | age |
; type: | int | text | int |
; primarykey: | id |
; TableName: test1
| 1 | Henry | 19 |
| 2 | Alice | 14 |
`

var TestAllDiff = ` | 1 | Bob | 19 |
+| 2 | Alice | 14 |
`
var TestAllDiff2 = ` | 1 | Bob | 19 |
-| 2 | Alice | 14 |
`
var TestAllDiff3 = `-| 1 | Bob | 19 |
+| 1 | Henry | 19 |
 | 2 | Alice | 14 |
`

func TestDiffAll(t *testing.T) {
	type args struct {
		t1       Reader
		t2       Reader
		diffMode DiffMode
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
		wantErr    bool
	}{
		{
			name:       "test1",
			args:       args{t1: NewReader(bytes.NewBufferString(TestDiff1)), t2: NewReader(bytes.NewBufferString(TestDiff2)), diffMode: AllDiff},
			wantWriter: TestAllDiff,
			wantErr:    false,
		},
		{
			name:       "test2",
			args:       args{t1: NewReader(bytes.NewBufferString(TestDiff2)), t2: NewReader(bytes.NewBufferString(TestDiff1)), diffMode: AllDiff},
			wantWriter: TestAllDiff2,
			wantErr:    false,
		},
		{
			name:       "test3",
			args:       args{t1: NewReader(bytes.NewBufferString(TestDiff2)), t2: NewReader(bytes.NewBufferString(TestDiff3)), diffMode: AllDiff},
			wantWriter: TestAllDiff3,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			if err := DiffAll(writer, tt.args.t1, tt.args.t2, tt.args.diffMode); (err != nil) != tt.wantErr {
				t.Errorf("DiffAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWriter := writer.String(); gotWriter != tt.wantWriter {
				t.Errorf("DiffAll() = %v, want %v", gotWriter, tt.wantWriter)
			}
		})
	}
}
