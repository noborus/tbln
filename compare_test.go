package tbln

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

var TestData = `; name: | id | name | age |
; type: | int | text | int |
; primarykey: | id |
; TableName: test1
| 1 | Bob | 19 |
`

var TestData2 = `; name: | id | name | age |
; type: | int | text | int |
; primarykey: | id |
; TableName: test1
| 1 | Alice | 14 |
`

var TestDataSt = `; name: | name |
; type: | text |
; primarykey: | name |
; TableName: test1
| Bob |
`
var TestDataSt2 = `; name: | name |
; type: | text |
; primarykey: | name |
; TableName: test1
| Alice |
`
var TestDataSt3 = `; name: | name |
; type: | text |
; primarykey: | name |
; TableName: test1
| Bob |
| Alice |
`
var TestDataSt4 = `; name: | name |
; type: | text |
; primarykey: | name |
; TableName: test1
`

var TestDataFl = `; name: | fl |
; type: | numeric |
; primarykey: | fl |
; TableName: test1
| 3.14 |
`
var TestDataFl2 = `; name: | fl |
; type: | numeric |
; primarykey: | fl |
; TableName: test1
| 3.1415 |
`

var TestNoPrim = `; name: | id | name | age |
; type: | int | text | int |
; TableName: test1
| 1 | Bob | 19 |
`

var TestErr = `; name: | id | name | age |
; type: | int | text | int |
; primarykey: | name |
; TableName: test1
| 1 | Bob | 19 |
`

func TestNewCompare(t *testing.T) {
	type args struct {
		t1 Reader
		t2 Reader
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Err1",
			args:    args{t1: NewReader(bytes.NewBufferString("")), t2: NewReader(bytes.NewBufferString(""))},
			wantErr: true,
		},
		{
			name:    "Err2",
			args:    args{t1: NewReader(bytes.NewBufferString(TestData)), t2: NewReader(bytes.NewBufferString(""))},
			wantErr: true,
		},
		{
			name:    "testErr1",
			args:    args{t1: NewReader(bytes.NewBufferString(TestErr)), t2: NewReader(bytes.NewBufferString(TestData))},
			wantErr: true,
		},
		{
			name:    "testNo",
			args:    args{t1: NewReader(bytes.NewBufferString(TestNoPrim)), t2: NewReader(bytes.NewBufferString(TestData))},
			wantErr: false,
		},
		{
			name:    "test1",
			args:    args{t1: NewReader(bytes.NewBufferString(TestData)), t2: NewReader(bytes.NewBufferString(TestData2))},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCompare(tt.args.t1, tt.args.t2)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCompare() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCompare_ReadDiffRow(t *testing.T) {
	type fields struct {
		t1 Reader
		t2 Reader
	}
	tests := []struct {
		name    string
		fields  fields
		want    *DiffRow
		wantErr bool
	}{
		{
			name:    "test1",
			fields:  fields{t1: NewReader(bytes.NewBufferString(TestData)), t2: NewReader(bytes.NewBufferString(TestData2))},
			want:    &DiffRow{Les: 2, Self: []string{"1", "Bob", "19"}, Other: []string{"1", "Alice", "14"}},
			wantErr: false,
		},
		{
			name:    "testString1",
			fields:  fields{t1: NewReader(bytes.NewBufferString(TestDataSt)), t2: NewReader(bytes.NewBufferString(TestDataSt2))},
			want:    &DiffRow{Les: 1, Self: nil, Other: []string{"Alice"}},
			wantErr: false,
		},
		{
			name:    "testString2",
			fields:  fields{t1: NewReader(bytes.NewBufferString(TestDataSt)), t2: NewReader(bytes.NewBufferString(TestDataSt3))},
			want:    &DiffRow{Les: 0, Self: []string{"Bob"}, Other: []string{"Bob"}},
			wantErr: false,
		},
		{
			name:    "testString3",
			fields:  fields{t1: NewReader(bytes.NewBufferString(TestDataSt2)), t2: NewReader(bytes.NewBufferString(TestDataSt))},
			want:    &DiffRow{Les: -1, Self: []string{"Alice"}, Other: nil},
			wantErr: false,
		},
		{
			name:    "testString4",
			fields:  fields{t1: NewReader(bytes.NewBufferString(TestDataSt)), t2: NewReader(bytes.NewBufferString(TestDataSt4))},
			want:    &DiffRow{Les: -1, Self: []string{"Bob"}, Other: nil},
			wantErr: false,
		},
		{
			name:    "testFloat",
			fields:  fields{t1: NewReader(bytes.NewBufferString(TestDataFl)), t2: NewReader(bytes.NewBufferString(TestDataFl2))},
			want:    &DiffRow{Les: -1, Self: []string{"3.14"}, Other: nil},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		cmp, err := NewCompare(tt.fields.t1, tt.fields.t2)
		if err != nil {
			if err != io.EOF {
				t.Fatalf("NewCompare error %s", err)
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmp.ReadDiffRow()
			if (err != nil) != tt.wantErr {
				t.Errorf("Compare.ReadDiffRow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Compare.ReadDiffRow() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestCompare_ReadDiffRow2(t *testing.T) {
	type fields struct {
		t1 Reader
		t2 Reader
	}
	tests := []struct {
		name    string
		fields  fields
		want    *DiffRow
		wantErr bool
	}{
		{
			name:    "testString1",
			fields:  fields{t1: NewReader(bytes.NewBufferString(TestDataSt2)), t2: NewReader(bytes.NewBufferString(TestDataSt3))},
			want:    &DiffRow{Les: 1, Self: nil, Other: []string{"Bob"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		cmp, err := NewCompare(tt.fields.t1, tt.fields.t2)
		if err != nil {
			if err != io.EOF {
				t.Fatalf("NewCompare error %s", err)
			}
		}
		t.Run(tt.name, func(t *testing.T) {
			_, err := cmp.ReadDiffRow()
			if (err != nil) != tt.wantErr {
				t.Errorf("Compare.ReadDiffRow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got2, err := cmp.ReadDiffRow()
			if (err != nil) != tt.wantErr {
				t.Errorf("Compare.ReadDiffRow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got2, tt.want) {
				t.Errorf("Compare.ReadDiffRow() = %v, want %v", got2, tt.want)
			}
		})
	}
}

func Test_getPK(t *testing.T) {
	type args struct {
		t1 Reader
		t2 Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []Pkey
		wantErr bool
	}{
		{
			name: "test1",
			args: args{
				t1: NewReader(bytes.NewBufferString(TestData)),
				t2: NewReader(bytes.NewBufferString(TestData2)),
			},
			want:    []Pkey{{Pos: 0, Name: "id", Typ: "int"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		cmp, err := NewCompare(tt.args.t1, tt.args.t2)
		if err != nil {
			t.Fatalf("NewCompare error %s", err)
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := cmp.getPK()
			if (err != nil) != tt.wantErr {
				t.Errorf("getPK() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getPK() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestColumnPrimaryKey(t *testing.T) {
	type args struct {
		pkeys []Pkey
		row   []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test1",
			args: args{pkeys: []Pkey{{Pos: 0, Name: "id", Typ: "int"}}, row: []string{"1", "Bob"}},
			want: []string{"1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ColumnPrimaryKey(tt.args.pkeys, tt.args.row); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ColumnPrimaryKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
