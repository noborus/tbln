package db

import (
	"testing"
)

func TestTDB_quoting(t *testing.T) {
	type fields struct {
		Style Style
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{name: "test1", fields: fields{Style: Style{Quote: `"`}}, args: args{`test`}, want: `test`},
		{name: "test2", fields: fields{Style: Style{Quote: `"`}}, args: args{`Test`}, want: `"Test"`},
		{name: "test3", fields: fields{Style: Style{Quote: `"`}}, args: args{`t'est`}, want: `"t'est"`},
		{name: "test4", fields: fields{Style: Style{Quote: `"`}}, args: args{`t"est`}, want: `"t""est"`},
		{name: "test5", fields: fields{Style: Style{Quote: `"`}}, args: args{`test123`}, want: `test123`},
		{name: "test6", fields: fields{Style: Style{Quote: `"`}}, args: args{`test;SELECT 1`}, want: `"test;SELECT 1"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TDB := &TDB{
				Style: tt.fields.Style,
			}
			if got := TDB.quoting(tt.args.name); got != tt.want {
				t.Errorf("TDB.quoting() = %v, want %v", got, tt.want)
			}
		})
	}
}
