package tbln

import (
	"reflect"
	"testing"
)

func TestDefinition_SetNames(t *testing.T) {
	type fields struct {
		columnNum int
		tableName string
		comments  []string
		names     []string
		types     []string
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			d := &Definition{
				columnNum: tt.fields.columnNum,
				tableName: tt.fields.tableName,
				Comments:  tt.fields.comments,
				names:     tt.fields.names,
				types:     tt.fields.types,
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
		names     []string
		types     []string
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			d := &Definition{
				columnNum: tt.fields.columnNum,
				names:     tt.fields.names,
				types:     tt.fields.types,
				Extras:    tt.fields.Extras,
			}
			if err := d.SetTypes(tt.args.types); (err != nil) != tt.wantErr {
				t.Errorf("Definition.SetTypes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefinition_SetExtra(t *testing.T) {
	type fields struct {
		columnNum int
		tableName string
		algorithm string
		Comments  []string
		names     []string
		types     []string
		Extras    map[string]Extra
		Hashes    map[string][]byte
		Signs     Signatures
	}
	type args struct {
		keyName string
		value   string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]Extra
		wantv  string
	}{
		{
			name:   "test1",
			fields: fields{Extras: make(map[string]Extra)},
			args:   args{keyName: "test", value: "testValue"},
			want: map[string]Extra{
				"test": {"testValue", false},
			},
			wantv: "testvalue",
		},
		{
			name:   "test1",
			fields: fields{Extras: map[string]Extra{"test": {"testValue", false}}},
			args:   args{keyName: "test"},
			want:   map[string]Extra{},
			wantv:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Definition{
				columnNum: tt.fields.columnNum,
				tableName: tt.fields.tableName,
				algorithm: tt.fields.algorithm,
				Comments:  tt.fields.Comments,
				names:     tt.fields.names,
				types:     tt.fields.types,
				Extras:    tt.fields.Extras,
				Hashes:    tt.fields.Hashes,
				Signs:     tt.fields.Signs,
			}
			d.SetExtra(tt.args.keyName, tt.args.value)
			got := d.Extras
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tbln.SetExtra = %v, want %v", got, tt.want)
			}
			gv := d.ExtraValue(tt.args.keyName)
			if gv == tt.wantv {
				t.Errorf("Tbln.SetExtra = %v, wantv %v", gv, tt.wantv)
			}
		})
	}
}
