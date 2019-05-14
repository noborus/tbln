package db

import (
	"reflect"
	"testing"

	"github.com/noborus/tbln"
)

func Test_mergeTableRow(t *testing.T) {
	type args struct {
		dml    *dml
		d      *tbln.DiffRow
		delete bool
	}
	tests := []struct {
		name string
		args args
		want *dml
	}{
		{
			name: "test1",
			args: args{
				dml:    &dml{},
				d:      &tbln.DiffRow{Les: 0, Self: []string{"1"}, Other: []string{"1"}},
				delete: false,
			},
			want: &dml{},
		},
		{
			name: "test2",
			args: args{
				dml:    &dml{},
				d:      &tbln.DiffRow{Les: 1, Self: []string{"1"}, Other: []string{"2"}},
				delete: false,
			},
			want: &dml{insert: [][]string{{"2"}}},
		},
		{
			name: "test3",
			args: args{
				dml:    &dml{},
				d:      &tbln.DiffRow{Les: -1, Self: []string{"2"}, Other: []string{"1"}},
				delete: false,
			},
			want: &dml{},
		},
		{
			name: "test4",
			args: args{
				dml:    &dml{},
				d:      &tbln.DiffRow{Les: -1, Self: []string{"2"}, Other: []string{"1"}},
				delete: true,
			},
			want: &dml{delete: [][]string{{"2"}}},
		},
		{
			name: "test5",
			args: args{
				dml:    &dml{},
				d:      &tbln.DiffRow{Les: 2, Self: []string{"1", "2"}, Other: []string{"1", "3"}},
				delete: true,
			},
			want: &dml{update: [][]string{{"1", "3"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeTableRow(tt.args.dml, tt.args.d, tt.args.delete); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeTableRow() = %v, want %v", got, tt.want)
			}
		})
	}
}
