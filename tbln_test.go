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
			tb.ToTargetHash("name", true)
			tb.ToTargetHash("type", true)
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

func TestTbln_calculateHash(t *testing.T) {
	type args struct {
		hashType string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "testErr",
			args:    args{hashType: "err"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "testSHA256",
			args: args{hashType: "sha256"},
			want: []byte{
				0xe3, 0xb0, 0xc4, 0x42, 0x98, 0xfc, 0x1c, 0x14, 0x9a, 0xfb, 0xf4, 0xc8, 0x99, 0x6f, 0xb9, 0x24,
				0x27, 0xae, 0x41, 0xe4, 0x64, 0x9b, 0x93, 0x4c, 0xa4, 0x95, 0x99, 0x1b, 0x78, 0x52, 0xb8, 0x55,
			},
			wantErr: false,
		},
		{
			name: "testSHA512",
			args: args{hashType: "sha512"},
			want: []byte{
				0xcf, 0x83, 0xe1, 0x35, 0x7e, 0xef, 0xb8, 0xbd, 0xf1, 0x54, 0x28, 0x50, 0xd6, 0x6d, 0x80, 0x07,
				0xd6, 0x20, 0xe4, 0x05, 0x0b, 0x57, 0x15, 0xdc, 0x83, 0xf4, 0xa9, 0x21, 0xd3, 0x6c, 0xe9, 0xce,
				0x47, 0xd0, 0xd1, 0x3c, 0x5d, 0x85, 0xf2, 0xb0, 0xff, 0x83, 0x18, 0xd2, 0x87, 0x7e, 0xec, 0x2f,
				0x63, 0xb9, 0x31, 0xbd, 0x47, 0x41, 0x7a, 0x81, 0xa5, 0x38, 0x32, 0x7a, 0xf9, 0x27, 0xda, 0x3e,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := NewTbln()
			got, err := tb.calculateHash(tt.args.hashType)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tbln.calculateHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tbln.calculateHash() = %x, want %x", got, tt.want)
			}
		})
	}
}
