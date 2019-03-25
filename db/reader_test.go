package db

import (
	"testing"
	"time"
)

func Test_toString(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{nil},
			want: "",
		},
		{
			name: "test2",
			args: args{v: time.Date(2014, 12, 31, 8, 4, 18, 0, time.UTC)},
			want: "2014-12-31T08:04:18Z",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toString(tt.args.v); got != tt.want {
				t.Errorf("toString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertType(t *testing.T) {
	type args struct {
		dbtype string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test1",
			args: args{"text"},
			want: "text",
		},
		{
			name: "test2",
			args: args{"nottype"},
			want: "text",
		},
		{
			name: "test3",
			args: args{"double precision"},
			want: "numeric",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertType(tt.args.dbtype); got != tt.want {
				t.Errorf("convertType() = %v, want %v", got, tt.want)
			}
		})
	}
}
