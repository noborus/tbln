// Package sqlite3 contains SQLite3 specific database dialect of tbln/db.
package sqlite3

import (
	"bytes"
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/noborus/tbln"
	"github.com/noborus/tbln/db"
)

var TestData = `; name: | id | name | age |
; type: | int | text | int |
; TableName: test1
; primarykey: | id |
| 1 | Bob | 19 |
| 2 | Alice | 14 |
`

func SetupSQLiteTest(t *testing.T) *db.TDB {
	conn, err := db.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	r := bytes.NewBufferString(TestData)
	at, err := tbln.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	err = db.WriteTable(conn, at, "", db.ReCreate, db.Normal)
	if err != nil {
		t.Fatal(err)
	}
	err = conn.Commit()
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

func TestReadTableAll(t *testing.T) {
	tdb := SetupSQLiteTest(t)
	type args struct {
		tdb       *db.TDB
		schema    string
		tableName string
	}
	tests := []struct {
		name    string
		args    args
		want    *tbln.TBLN
		wantErr bool
	}{
		{
			name:    "testErr",
			args:    args{tdb: tdb, schema: "", tableName: "test"},
			wantErr: true,
		},
		{
			name:    "test1",
			args:    args{tdb: tdb, schema: "", tableName: "test1"},
			want:    &tbln.TBLN{Rows: [][]string{{"1", "Bob", "19"}, {"2", "Alice", "14"}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.ReadTableAll(tt.args.tdb, tt.args.schema, tt.args.tableName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadTableAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want != nil {
				if !reflect.DeepEqual(got.Rows, tt.want.Rows) {
					t.Errorf("ReadTableAll() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestConstr_GetSchema(t *testing.T) {
	tdb := SetupSQLiteTest(t)
	type args struct {
		conn *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "default schema",
			args:    args{conn: tdb.DB},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		c := &Constr{}
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetSchema(tt.args.conn)
			if (err != nil) != tt.wantErr {
				t.Errorf("Constr.GetSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Constr.GetSchema() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConstr_GetPrimaryKey(t *testing.T) {
	tdb := SetupSQLiteTest(t)
	type args struct {
		conn      *sql.DB
		schema    string
		tableName string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{conn: tdb.DB, schema: "", tableName: "test1"},
			want:    []string{"id"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		c := &Constr{}
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetPrimaryKey(tt.args.conn, tt.args.schema, tt.args.tableName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Constr.GetPrimaryKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Constr.GetPrimaryKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConstr_GetColumnInfo(t *testing.T) {
	tdb := SetupSQLiteTest(t)
	type args struct {
		conn      *sql.DB
		schema    string
		tableName string
	}
	tests := []struct {
		name    string
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{conn: tdb.DB, schema: "", tableName: "test1"},
			want:    []interface{}{"int", "text", "int"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		c := &Constr{}
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetColumnInfo(tt.args.conn, tt.args.schema, tt.args.tableName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Constr.GetColumnInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for n, wantV := range tt.want {
				v := fmt.Sprintf("%s", got["sqlite3_type"][n])
				if v != wantV {
					t.Errorf("Constr.GetColumnInfo() = %v, want %v", v, wantV)
				}
			}
		})
	}
}
