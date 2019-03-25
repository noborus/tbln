package db

import (
	"testing"
)

func TestTDB_quoting(t *testing.T) {
	st := Style{Quote: `"`}
	type args struct {
		obj string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "test1", args: args{`test`}, want: `test`},
		{name: "test2", args: args{`Test`}, want: `"Test"`},
		{name: "test3", args: args{`t'est`}, want: `"t'est"`},
		{name: "test4", args: args{`t"est`}, want: `"t""est"`},
		{name: "test5", args: args{`test123`}, want: `test123`},
		{name: "test6", args: args{`test;SELECT 1`}, want: `"test;SELECT 1"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			TDB := &TDB{
				Style: st,
			}
			if got := TDB.quoting(tt.args.obj); got != tt.want {
				t.Errorf("TDB.quoting() = %v, want %v", got, tt.want)
			}
		})
	}
}
