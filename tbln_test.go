package tbln

import (
	"bytes"
	"reflect"
	"testing"
)

func TestTbln_AddRows(t *testing.T) {
	type fields struct {
		Definition Definition
		Hash       map[string]string
		buffer     bytes.Buffer
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
			fields:  fields{Definition: Definition{ColumnNum: 1}},
			args:    args{row: []string{"1"}},
			wantErr: false,
		},
		{
			name:    "test2",
			fields:  fields{Definition: Definition{ColumnNum: 2}},
			args:    args{row: []string{"1"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &Tbln{
				Definition: tt.fields.Definition,
				buffer:     tt.fields.buffer,
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
		Definition Definition
		Hash       map[string]string
		buffer     bytes.Buffer
		RowNum     int
		Rows       [][]string
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]string
		wantErr bool
	}{
		{
			name:    "testBlank",
			fields:  fields{},
			want:    map[string]string{"sha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
			wantErr: false,
		},
		{
			name:    "testOneRow",
			fields:  fields{Rows: [][]string{{"1"}}, RowNum: 1},
			want:    map[string]string{"sha256": "3c5c7b4b1fcd47206cbc619c23b59a27b97f730abca146d67157d2d9df8ca9dc"},
			wantErr: false,
		},
		{
			name:    "testNames",
			fields:  fields{Definition: Definition{Ext: make(map[string]Extra), Names: []string{"id"}, ColumnNum: 1}, Rows: [][]string{{"1"}}, RowNum: 1},
			want:    map[string]string{"sha256": "e5ce5f72c836840efdbcbf7639075966944253ef438a305761d888158a6b22a8"},
			wantErr: false,
		},
		{
			name:    "testFullRow",
			fields:  fields{Definition: Definition{Ext: make(map[string]Extra), Names: []string{"id", "name"}, Types: []string{"int", "text"}, ColumnNum: 2}, Rows: [][]string{{"1", "test"}}, RowNum: 1},
			want:    map[string]string{"sha256": "fcc150288d592d5c0cf13eed4b1054f6fadbfd2c48cde10954b44d6b7fc42623"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &Tbln{
				Definition: tt.fields.Definition,
				buffer:     tt.fields.buffer,
				RowNum:     tt.fields.RowNum,
				Rows:       tt.fields.Rows,
			}
			tb.SetNames(tb.Names)
			tb.SetTypes(tb.Types)
			got, err := tb.SumHash(SHA256)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tbln.SumHash(SHA256) error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tbln.SumHash(SHA256) = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefinition_SetNames(t *testing.T) {
	type fields struct {
		ColumnNum int
		TableName string
		Comments  []string
		Names     []string
		Types     []string
		Ext       map[string]Extra
	}
	type args struct {
		names []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "test1", fields: fields{Ext: make(map[string]Extra)}, args: args{names: []string{"a", "b"}}, wantErr: false},
		{name: "test2", fields: fields{Ext: make(map[string]Extra), ColumnNum: 3}, args: args{names: []string{"a", "b", "c"}}, wantErr: false},
		{name: "test3", fields: fields{Ext: make(map[string]Extra), ColumnNum: 2}, args: args{names: []string{"a", "b", "c"}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Definition{
				ColumnNum: tt.fields.ColumnNum,
				TableName: tt.fields.TableName,
				Comments:  tt.fields.Comments,
				Names:     tt.fields.Names,
				Types:     tt.fields.Types,
				Ext:       tt.fields.Ext,
			}
			if err := d.SetNames(tt.args.names); (err != nil) != tt.wantErr {
				t.Errorf("Definition.SetNames() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefinition_SetTypes(t *testing.T) {
	type fields struct {
		ColumnNum int
		TableName string
		Comments  []string
		Names     []string
		Types     []string
		Ext       map[string]Extra
	}
	type args struct {
		types []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "test1", fields: fields{Ext: make(map[string]Extra)}, args: args{types: []string{"int", "text"}}, wantErr: false},
		{name: "test2", fields: fields{Ext: make(map[string]Extra), ColumnNum: 3}, args: args{types: []string{"int", "text", "text"}}, wantErr: false},
		{name: "test3", fields: fields{Ext: make(map[string]Extra), ColumnNum: 2}, args: args{types: []string{"int", "text", "text"}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Definition{
				ColumnNum: tt.fields.ColumnNum,
				TableName: tt.fields.TableName,
				Comments:  tt.fields.Comments,
				Names:     tt.fields.Names,
				Types:     tt.fields.Types,
				Ext:       tt.fields.Ext,
			}
			if err := d.SetTypes(tt.args.types); (err != nil) != tt.wantErr {
				t.Errorf("Definition.SetTypes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefinition_setColNum(t *testing.T) {
	type fields struct {
		ColumnNum int
		TableName string
		Comments  []string
		Names     []string
		Types     []string
		Ext       map[string]Extra
	}
	type args struct {
		colNum int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{name: "test1", args: args{colNum: 2}, wantErr: false},
		{name: "test2", fields: fields{ColumnNum: 3}, args: args{colNum: 3}, wantErr: false},
		{name: "test3", fields: fields{ColumnNum: 2}, args: args{colNum: 3}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Definition{
				ColumnNum: tt.fields.ColumnNum,
				TableName: tt.fields.TableName,
				Comments:  tt.fields.Comments,
				Names:     tt.fields.Names,
				Types:     tt.fields.Types,
				Ext:       tt.fields.Ext,
			}
			if err := d.setColNum(tt.args.colNum); (err != nil) != tt.wantErr {
				t.Errorf("Definition.setColNum() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
