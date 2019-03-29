package db_test

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/noborus/tbln"
	"github.com/noborus/tbln/db"
	_ "github.com/noborus/tbln/db/mysql"
	_ "github.com/noborus/tbln/db/postgres"
	_ "github.com/noborus/tbln/db/sqlite3"
)

var MySQLDBname string
var PostgreSQLDBname string

var TestData1 = `; name: | id | name | age |
; type: | int | text | int |
; primarykey: | id |
; TableName: test1
| 1 | Bob | 19 |
| 2 | Alice | 14 |
`
var TestData2 = `; name: | one |
; type: | text |
; TableName: testone
| a |
| b |
| c |
`

func createTestTable(t *testing.T, tdb *db.TDB, data string) string {
	at := wantTbln(t, data)
	err := tdb.Begin()
	if err != nil {
		t.Fatal(err)
	}
	err = db.WriteTable(tdb, at, "", db.ReCreate, db.Normal)
	if err != nil {
		t.Fatal(err)
	}
	err = tdb.Commit()
	if err != nil {
		t.Fatal(err)
	}
	return at.TableName()
}

func wantTbln(t *testing.T, data string) *tbln.Tbln {
	r := bytes.NewBufferString(data)
	at, err := tbln.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return at
}

func dropTestTable(t *testing.T, tdb *db.TDB, tableName string) {
	sql := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	err := tdb.Begin()
	if err != nil {
		t.Fatal(err)
	}
	_, err = tdb.Tx.Exec(sql)
	if err != nil {
		t.Fatal(err)
	}
	err = tdb.Commit()
	if err != nil {
		t.Fatal(err)
	}
}

func SetupPostgresTest(t *testing.T) *db.TDB {
	// db.Debug(true)
	if PostgreSQLDBname = os.Getenv("TEST_PG_DATABASE"); PostgreSQLDBname == "" {
		PostgreSQLDBname = "test_db"
	}
	dsn := fmt.Sprintf("dbname=%s", PostgreSQLDBname)
	conn, err := db.Open("postgres", dsn)
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

func SetupSQLite3Test(t *testing.T) *db.TDB {
	conn, err := db.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

func SetupMySQLTest(t *testing.T) *db.TDB {
	if MySQLDBname = os.Getenv("TEST_MYSQL_DATABASE"); MySQLDBname == "" {
		MySQLDBname = "test_db"
	}
	conn, err := db.Open("mysql", "root@/"+MySQLDBname)
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

func TestTDB_ReadTableAll(t *testing.T) {
	type args struct {
		schema    string
		tableName string
		pkey      []string
	}
	tests := []struct {
		name    string
		args    args
		table   string
		want    *tbln.Tbln
		wantErr bool
	}{
		{
			name:    "testErr",
			args:    args{schema: "", tableName: "testnotable", pkey: nil},
			table:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test1",
			args:    args{schema: "", tableName: "test1", pkey: nil},
			table:   TestData1,
			want:    wantTbln(t, TestData1),
			wantErr: false,
		},
		{
			name:    "testOne",
			args:    args{schema: "", tableName: "testone", pkey: nil},
			table:   TestData2,
			want:    wantTbln(t, TestData2),
			wantErr: false,
		},
	}
	dbDriver := []*db.TDB{
		SetupPostgresTest(t),
		SetupMySQLTest(t),
		SetupSQLite3Test(t),
	}
	for _, tdb := range dbDriver {
		for _, tt := range tests {
			var tableName string
			if tt.table != "" {
				tableName = createTestTable(t, tdb, tt.table)
			}
			t.Run(tt.name+tdb.Name, func(t *testing.T) {
				got, err := db.ReadTableAll(tdb, tt.args.schema, tt.args.tableName)
				if (err != nil) != tt.wantErr {
					t.Errorf("TDB.ReadTable() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if tt.want != nil {
					if !reflect.DeepEqual(got.Rows, tt.want.Rows) {
						t.Errorf("ReadTableAll() = %v, want %v", got.Rows, tt.want.Rows)
					}
				}
			})
			if tt.table != "" {
				dropTestTable(t, tdb, tableName)
			}
		}
	}
}

func TestTDB_ReadQueryAll(t *testing.T) {
	// db.Debug(true)
	type args struct {
		sql  string
		args []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    [][]string
		wantErr bool
	}{
		{
			name:    "testErr",
			args:    args{sql: "SELECT ERR", args: nil},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "testABC",
			args:    args{sql: "SELECT 'a' as a ,'b' as b ,'c' as c;", args: nil},
			want:    [][]string{{"a", "b", "c"}},
			wantErr: false,
		},
	}
	dbDriver := []*db.TDB{
		SetupPostgresTest(t),
		SetupMySQLTest(t),
		SetupSQLite3Test(t),
	}
	for _, tdb := range dbDriver {
		for _, tt := range tests {
			t.Run(tt.name+tdb.Name, func(t *testing.T) {
				got, err := db.ReadQueryAll(tdb, tt.args.sql, tt.args.args...)
				if (err != nil) != tt.wantErr {
					t.Errorf("TDB.ReadQuery() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if tt.want != nil {
					if !reflect.DeepEqual(got.Rows, tt.want) {
						t.Errorf("ReadTableAll() = %v, want %v", got.Rows, tt.want)
					}
				}
			})
		}
	}
}
