// Package mysql contains MySQL specific database dialect of tbln/db.
package mysql

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/noborus/tbln"
	"github.com/noborus/tbln/db"
)

var MySQLDBname string

var TestData = `; name: | id | name | age |
; type: | int | text | int |
; primarykey: | id |
; TableName: test1_mysql
| 1 | Bob | 19 |
| 2 | Alice | 14 |
`

func SetupMySQLTest(t *testing.T) *db.TDB {
	if MySQLDBname = os.Getenv("TEST_MYSQL_DATABASE"); MySQLDBname == "" {
		MySQLDBname = "test_db"
	}
	conn, err := db.Open("mysql", "root@/"+MySQLDBname)
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

func dummyDBConn(t *testing.T) *sql.DB {
	conn := SetupMySQLTest(t).DB
	conn.Close()
	return conn
}

func TestReadTableAll(t *testing.T) {
	tdb := SetupMySQLTest(t)
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
			args:    args{tdb: tdb, schema: "", tableName: "notablename"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test1",
			args:    args{tdb: tdb, schema: "", tableName: "test1_mysql"},
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
	tdb := SetupMySQLTest(t)
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
			want:    MySQLDBname,
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
	tdb := SetupMySQLTest(t)
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
			name:    "testErr",
			args:    args{conn: dummyDBConn(t), schema: "", tableName: "test1_mysql"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test1",
			args:    args{conn: tdb.DB, schema: "", tableName: "test1_mysql"},
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
	tdb := SetupMySQLTest(t)
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
			name:    "testErr",
			args:    args{conn: dummyDBConn(t), schema: "", tableName: "test1_mysql"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test1",
			args:    args{conn: tdb.DB, schema: "", tableName: "test1_mysql"},
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
				v := fmt.Sprintf("%s", got["mysql_type"][n])
				if v != wantV {
					t.Errorf("Constr.GetColumnInfo() = %v, want %v", v, wantV)
				}
			}
		})
	}
}
