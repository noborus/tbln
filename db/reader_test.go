package db

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/noborus/tbln"
)

func TestTDB_ReadQuery(t *testing.T) {
	type fields struct {
		Name       string
		DB         *sql.DB
		Tx         *sql.Tx
		Style      Style
		Constraint Constraint
	}
	type args struct {
		sql  string
		args []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Reader
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tdb := &TDB{
				Name:       tt.fields.Name,
				DB:         tt.fields.DB,
				Tx:         tt.fields.Tx,
				Style:      tt.fields.Style,
				Constraint: tt.fields.Constraint,
			}
			got, err := tdb.ReadQuery(tt.args.sql, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("TDB.ReadQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TDB.ReadQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReader_ReadRow(t *testing.T) {
	type fields struct {
		Definition *tbln.Definition
		TDB        *TDB
		rows       *sql.Rows
		scanArgs   []interface{}
		values     []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Reader{
				Definition: tt.fields.Definition,
				TDB:        tt.fields.TDB,
				rows:       tt.fields.rows,
				scanArgs:   tt.fields.scanArgs,
				values:     tt.fields.values,
			}
			got, err := tr.ReadRow()
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.ReadRow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reader.ReadRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReadQueryAll(t *testing.T) {
	type args struct {
		tdb   *TDB
		query string
		args  []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    *tbln.Tbln
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadQueryAll(tt.args.tdb, tt.args.query, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadQueryAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadQueryAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReader_readRowsAll(t *testing.T) {
	type fields struct {
		Definition *tbln.Definition
		TDB        *TDB
		rows       *sql.Rows
		scanArgs   []interface{}
		values     []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		want    *tbln.Tbln
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &Reader{
				Definition: tt.fields.Definition,
				TDB:        tt.fields.TDB,
				rows:       tt.fields.rows,
				scanArgs:   tt.fields.scanArgs,
				values:     tt.fields.values,
			}
			got, err := tr.readRowsAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("Reader.readRowsAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Reader.readRowsAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReader_setTableInfo(t *testing.T) {
	type fields struct {
		Definition *tbln.Definition
		TDB        *TDB
		rows       *sql.Rows
		scanArgs   []interface{}
		values     []interface{}
	}
	type args struct {
		constraints map[string][]interface{}
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
			tr := &Reader{
				Definition: tt.fields.Definition,
				TDB:        tt.fields.TDB,
				rows:       tt.fields.rows,
				scanArgs:   tt.fields.scanArgs,
				values:     tt.fields.values,
			}
			tr.setTableInfo(tt.args.constraints)
		})
	}
}

func TestReader_query(t *testing.T) {
	type fields struct {
		Definition *tbln.Definition
		TDB        *TDB
		rows       *sql.Rows
		scanArgs   []interface{}
		values     []interface{}
	}
	type args struct {
		query string
		args  []interface{}
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
			tr := &Reader{
				Definition: tt.fields.Definition,
				TDB:        tt.fields.TDB,
				rows:       tt.fields.rows,
				scanArgs:   tt.fields.scanArgs,
				values:     tt.fields.values,
			}
			if err := tr.query(tt.args.query, tt.args.args...); (err != nil) != tt.wantErr {
				t.Errorf("Reader.query() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReader_setRowInfo(t *testing.T) {
	type fields struct {
		Definition *tbln.Definition
		TDB        *TDB
		rows       *sql.Rows
		scanArgs   []interface{}
		values     []interface{}
	}
	type args struct {
		rows *sql.Rows
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
			tr := &Reader{
				Definition: tt.fields.Definition,
				TDB:        tt.fields.TDB,
				rows:       tt.fields.rows,
				scanArgs:   tt.fields.scanArgs,
				values:     tt.fields.values,
			}
			if err := tr.setRowInfo(tt.args.rows); (err != nil) != tt.wantErr {
				t.Errorf("Reader.setRowInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_toString(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertType(tt.args.dbtype); got != tt.want {
				t.Errorf("convertType() = %v, want %v", got, tt.want)
			}
		})
	}
}
