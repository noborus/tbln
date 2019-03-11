package tbln

import (
	"reflect"
	"testing"
)

func TestNewDefinition(t *testing.T) {
	tests := []struct {
		name string
		want *Definition
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDefinition(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDefinition() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewExtra(t *testing.T) {
	type args struct {
		value      interface{}
		hashTarget bool
	}
	tests := []struct {
		name string
		args args
		want Extra
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewExtra(tt.args.value, tt.args.hashTarget); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewExtra() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtra_Value(t *testing.T) {
	type fields struct {
		value      interface{}
		hashTarget bool
	}
	tests := []struct {
		name   string
		fields fields
		want   interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Extra{
				value:      tt.fields.value,
				hashTarget: tt.fields.hashTarget,
			}
			if got := e.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Extra.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefinition_TableName(t *testing.T) {
	type fields struct {
		columnNum int
		tableName string
		Comments  []string
		Names     []string
		Types     []string
		Extras    map[string]Extra
		Hashes    map[string][]byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Definition{
				columnNum: tt.fields.columnNum,
				tableName: tt.fields.tableName,
				Comments:  tt.fields.Comments,
				Names:     tt.fields.Names,
				Types:     tt.fields.Types,
				Extras:    tt.fields.Extras,
				Hashes:    tt.fields.Hashes,
			}
			if got := d.TableName(); got != tt.want {
				t.Errorf("Definition.TableName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefinition_SetTableName(t *testing.T) {
	type fields struct {
		columnNum int
		tableName string
		Comments  []string
		Names     []string
		Types     []string
		Extras    map[string]Extra
		Hashes    map[string][]byte
		Signs     map[string][]byte
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Definition{
				columnNum: tt.fields.columnNum,
				tableName: tt.fields.tableName,
				Comments:  tt.fields.Comments,
				Names:     tt.fields.Names,
				Types:     tt.fields.Types,
				Extras:    tt.fields.Extras,
				Hashes:    tt.fields.Hashes,
			}
			d.SetTableName(tt.args.name)
		})
	}
}

func TestDefinition_SetNames(t *testing.T) {
	type fields struct {
		columnNum int
		tableName string
		comments  []string
		Names     []string
		Types     []string
		Extras    map[string]Extra
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
		{name: "test1", fields: fields{Extras: make(map[string]Extra)}, args: args{names: []string{"a", "b"}}, wantErr: false},
		{name: "test2", fields: fields{Extras: make(map[string]Extra), columnNum: 3}, args: args{names: []string{"a", "b", "c"}}, wantErr: false},
		{name: "test3", fields: fields{Extras: make(map[string]Extra), columnNum: 2}, args: args{names: []string{"a", "b", "c"}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Definition{
				columnNum: tt.fields.columnNum,
				tableName: tt.fields.tableName,
				Comments:  tt.fields.comments,
				Names:     tt.fields.Names,
				Types:     tt.fields.Types,
				Extras:    tt.fields.Extras,
			}
			if err := d.SetNames(tt.args.names); (err != nil) != tt.wantErr {
				t.Errorf("Definition.SetNames() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefinition_SetTypes(t *testing.T) {
	type fields struct {
		columnNum int
		tableName string
		Comments  []string
		Names     []string
		Types     []string
		Extras    map[string]Extra
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
		{name: "test1", fields: fields{Extras: make(map[string]Extra)}, args: args{types: []string{"int", "text"}}, wantErr: false},
		{name: "test2", fields: fields{Extras: make(map[string]Extra), columnNum: 3}, args: args{types: []string{"int", "text", "text"}}, wantErr: false},
		{name: "test3", fields: fields{Extras: make(map[string]Extra), columnNum: 2}, args: args{types: []string{"int", "text", "text"}}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Definition{
				columnNum: tt.fields.columnNum,
				tableName: tt.fields.tableName,
				Comments:  tt.fields.Comments,
				Names:     tt.fields.Names,
				Types:     tt.fields.Types,
				Extras:    tt.fields.Extras,
			}
			if err := d.SetTypes(tt.args.types); (err != nil) != tt.wantErr {
				t.Errorf("Definition.SetTypes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefinition_ColumnNum(t *testing.T) {
	type fields struct {
		columnNum int
		tableName string
		Comments  []string
		Names     []string
		Types     []string
		Extras    map[string]Extra
		Hashes    map[string][]byte
		Signs     map[string][]byte
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Definition{
				columnNum: tt.fields.columnNum,
				tableName: tt.fields.tableName,
				Comments:  tt.fields.Comments,
				Names:     tt.fields.Names,
				Types:     tt.fields.Types,
				Extras:    tt.fields.Extras,
				Hashes:    tt.fields.Hashes,
			}
			if got := d.ColumnNum(); got != tt.want {
				t.Errorf("Definition.ColumnNum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefinition_setColNum(t *testing.T) {
	type fields struct {
		columnNum int
		tableName string
		Comments  []string
		Names     []string
		Types     []string
		Extras    map[string]Extra
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
			d := &Definition{
				columnNum: tt.fields.columnNum,
				tableName: tt.fields.tableName,
				Comments:  tt.fields.Comments,
				Names:     tt.fields.Names,
				Types:     tt.fields.Types,
				Extras:    tt.fields.Extras,
			}
			if err := d.setColNum(tt.args.colNum); (err != nil) != tt.wantErr {
				t.Errorf("Definition.setColNum() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefinition_ToTargetHash(t *testing.T) {
	type fields struct {
		columnNum int
		tableName string
		Comments  []string
		Names     []string
		Types     []string
		Extras    map[string]Extra
		Hashes    map[string][]byte
		Signs     map[string][]byte
	}
	type args struct {
		key    string
		target bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Definition{
				columnNum: tt.fields.columnNum,
				tableName: tt.fields.tableName,
				Comments:  tt.fields.Comments,
				Names:     tt.fields.Names,
				Types:     tt.fields.Types,
				Extras:    tt.fields.Extras,
				Hashes:    tt.fields.Hashes,
			}
			d.ToTargetHash(tt.args.key, tt.args.target)
		})
	}
}

func TestDefinition_SerializeHash(t *testing.T) {
	type fields struct {
		columnNum int
		tableName string
		Comments  []string
		Names     []string
		Types     []string
		Extras    map[string]Extra
		Hashes    map[string][]byte
		Signs     map[string][]byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Definition{
				columnNum: tt.fields.columnNum,
				tableName: tt.fields.tableName,
				Comments:  tt.fields.Comments,
				Names:     tt.fields.Names,
				Types:     tt.fields.Types,
				Extras:    tt.fields.Extras,
				Hashes:    tt.fields.Hashes,
			}
			if got := d.SerializeHash(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Definition.SerializeHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefinition_SetSignatures(t *testing.T) {
	type fields struct {
		columnNum int
		tableName string
		Comments  []string
		Names     []string
		Types     []string
		Extras    map[string]Extra
		Hashes    map[string][]byte
		Signs     map[string][]byte
	}
	type args struct {
		signs []string
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
			d := &Definition{
				columnNum: tt.fields.columnNum,
				tableName: tt.fields.tableName,
				Comments:  tt.fields.Comments,
				Names:     tt.fields.Names,
				Types:     tt.fields.Types,
				Extras:    tt.fields.Extras,
				Hashes:    tt.fields.Hashes,
			}
			if err := d.SetSignatures(tt.args.signs); (err != nil) != tt.wantErr {
				t.Errorf("Definition.SetSignatures() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefinition_SetHashes(t *testing.T) {
	type fields struct {
		columnNum int
		tableName string
		Comments  []string
		Names     []string
		Types     []string
		Extras    map[string]Extra
		Hashes    map[string][]byte
		Signs     map[string][]byte
	}
	type args struct {
		hashes []string
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
			d := &Definition{
				columnNum: tt.fields.columnNum,
				tableName: tt.fields.tableName,
				Comments:  tt.fields.Comments,
				Names:     tt.fields.Names,
				Types:     tt.fields.Types,
				Extras:    tt.fields.Extras,
				Hashes:    tt.fields.Hashes,
			}
			if err := d.SetHashes(tt.args.hashes); (err != nil) != tt.wantErr {
				t.Errorf("Definition.SetHashes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
