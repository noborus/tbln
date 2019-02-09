package tbln

import "testing"

func TestTable_SetNames(t *testing.T) {
	type fields struct {
		columnNum int
		tableName string
		Names     []string
		Types     []string
		Comments  []string
		Ext       map[string]string
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
		{name: "test1", args: args{names: []string{"a", "b"}}, wantErr: false},
		{name: "test2", fields: fields{columnNum: 3}, args: args{names: []string{"a", "b", "c"}}, wantErr: false},
		{name: "test3", fields: fields{columnNum: 2}, args: args{names: []string{"a", "b", "c"}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := &Definition{
				columnNum: tt.fields.columnNum,
				tableName: tt.fields.tableName,
				Names:     tt.fields.Names,
				Types:     tt.fields.Types,
				Comments:  tt.fields.Comments,
				Ext:       tt.fields.Ext,
			}
			if err := tbl.SetNames(tt.args.names); (err != nil) != tt.wantErr {
				t.Errorf("Table.SetNames() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTable_SetTypes(t *testing.T) {
	type fields struct {
		columnNum int
		tableName string
		Names     []string
		Types     []string
		Comments  []string
		Ext       map[string]string
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
		{name: "test1", args: args{types: []string{"int", "text"}}, wantErr: false},
		{name: "test2", fields: fields{columnNum: 3}, args: args{types: []string{"int", "text", "text"}}, wantErr: false},
		{name: "test3", fields: fields{columnNum: 2}, args: args{types: []string{"int", "text", "text"}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := &Definition{
				columnNum: tt.fields.columnNum,
				tableName: tt.fields.tableName,
				Names:     tt.fields.Names,
				Types:     tt.fields.Types,
				Comments:  tt.fields.Comments,
				Ext:       tt.fields.Ext,
			}
			if err := tbl.SetTypes(tt.args.types); (err != nil) != tt.wantErr {
				t.Errorf("Table.SetTypes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTable_setColNum(t *testing.T) {
	type fields struct {
		columnNum int
		tableName string
		Names     []string
		Types     []string
		Comments  []string
		Ext       map[string]string
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
		{name: "test2", fields: fields{columnNum: 3}, args: args{colNum: 3}, wantErr: false},
		{name: "test3", fields: fields{columnNum: 2}, args: args{colNum: 3}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbl := &Definition{
				columnNum: tt.fields.columnNum,
				tableName: tt.fields.tableName,
				Names:     tt.fields.Names,
				Types:     tt.fields.Types,
				Comments:  tt.fields.Comments,
				Ext:       tt.fields.Ext,
			}
			if err := tbl.setColNum(tt.args.colNum); (err != nil) != tt.wantErr {
				t.Errorf("Table.setColNum() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
