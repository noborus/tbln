package tbln

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

var TestMerge1 = `; name: | id | name | age |
; type: | int | text | int |
; primarykey: | id |
; TableName: test1
| 1 | Bob | 19 |
`

var TestMerge2 = `; name: | id | name | age |
; type: | int | text | int |
; primarykey: | id |
; TableName: test1
| 2 | Alice | 14 |
`

var TestMergeTBLN = `; name: | id | name |
; type: | int | text |
; primarykey: | id |
; TableName: test1
| 1 | Bob |
| 2 | Alice |
`

func TestMergeAll(t *testing.T) {
	type args struct {
		t1 Reader
		t2 Reader
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{t1: NewReader(bytes.NewBufferString(TestMerge1)), t2: NewReader(bytes.NewBufferString(TestMerge2))},
			want:    TestMergeTBLN,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MergeAll(tt.args.t1, tt.args.t2, MergeUpdate)
			if (err != nil) != tt.wantErr {
				t.Errorf("MergeAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			f := bytes.NewBufferString(TestMergeTBLN)
			tb, err := ReadAll(bufio.NewReader(f))
			if err != nil {
				t.Fatalf("NewTBLN error %s", err)
			}
			if !reflect.DeepEqual(got.RowNum, tb.RowNum) {
				t.Errorf("MergeAll() = %v, want %v", got.RowNum, tb.RowNum)
			}
		})
	}
}
