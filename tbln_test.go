package tbln

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTbln_AddRows(t *testing.T) {
	type fields struct {
		Definition *Definition
		Hash       map[string]string
		RowNum     int
		Rows       [][]string
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
		{
			name:    "test1",
			fields:  fields{Definition: &Definition{columnNum: 1}},
			args:    args{row: []string{"1"}},
			wantErr: false,
		},
		{
			name:    "test2",
			fields:  fields{Definition: &Definition{columnNum: 2}},
			args:    args{row: []string{"1"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &Tbln{
				Definition: tt.fields.Definition,
				RowNum:     tt.fields.RowNum,
				Rows:       tt.fields.Rows,
			}
			if err := tb.AddRows(tt.args.row); (err != nil) != tt.wantErr {
				t.Errorf("Tbln.AddRows() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkRow(t *testing.T) {
	type args struct {
		ColumnNum int
		row       []string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{row: []string{"1"}, ColumnNum: 1},
			want:    1,
			wantErr: false,
		},
		{
			name:    "test2",
			args:    args{row: []string{"1", "2"}, ColumnNum: 0},
			want:    2,
			wantErr: false,
		},
		{
			name:    "test3",
			args:    args{row: []string{"1"}, ColumnNum: 2},
			want:    2,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkRow(tt.args.ColumnNum, tt.args.row)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkRow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTbln_SumHash(t *testing.T) {
	type fields struct {
		Definition *Definition
		Hash       map[string]string
		RowNum     int
		Rows       [][]string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name:    "testBlank",
			fields:  fields{Definition: NewDefinition()},
			want:    "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			wantErr: false,
		},
		{
			name:    "testOneRow",
			fields:  fields{Definition: NewDefinition(), Rows: [][]string{{"1"}}, RowNum: 1},
			want:    "3c5c7b4b1fcd47206cbc619c23b59a27b97f730abca146d67157d2d9df8ca9dc",
			wantErr: false,
		},
		{
			name:    "testNames",
			fields:  fields{Definition: &Definition{Extras: make(map[string]Extra), Hashes: make(map[string][]byte), Names: []string{"id"}, columnNum: 1}, Rows: [][]string{{"1"}}, RowNum: 1},
			want:    "e5ce5f72c836840efdbcbf7639075966944253ef438a305761d888158a6b22a8",
			wantErr: false,
		},
		{
			name:    "testFullRow",
			fields:  fields{Definition: &Definition{Extras: make(map[string]Extra), Hashes: make(map[string][]byte), Names: []string{"id", "name"}, Types: []string{"int", "text"}, columnNum: 2}, Rows: [][]string{{"1", "test"}}, RowNum: 1},
			want:    "fcc150288d592d5c0cf13eed4b1054f6fadbfd2c48cde10954b44d6b7fc42623",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &Tbln{
				Definition: tt.fields.Definition,
				RowNum:     tt.fields.RowNum,
				Rows:       tt.fields.Rows,
			}
			if err := tb.SetNames(tb.Names); err != nil {
				t.Errorf("Tbln.SetNames error = %v", err)
				return

			}
			if err := tb.SetTypes(tb.Types); err != nil {
				t.Errorf("Tbln.SetTypes error = %v", err)
				return
			}
			if err := tb.SumHash(SHA256); (err != nil) != tt.wantErr {
				t.Errorf("Tbln.SumHash(SHA256) error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := fmt.Sprintf("%x", tb.Hashes["sha256"])
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tbln.SumHash(SHA256) = %v, want %v", got, tt.want)
			}
		})
	}
}
