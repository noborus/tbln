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

var TestPrimary = `; name: | id1 | id2 | id3 |
; type: | int | int | int |
; primarykey: | id1 | id2 | id3 |
; TableName: testprim
| 1 | 1 | 1 |
| 2 | 2 | 2 |
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

func TestGetPrimaryKey(t *testing.T) {
	type args struct {
		tdb       *db.TDB
		tableName string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "testSQLite3no",
			args:    args{tdb: SetupSQLite3Test(t), tableName: ""},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "testPostgreSQLno",
			args:    args{tdb: SetupPostgresTest(t), tableName: ""},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "testMySQLno",
			args:    args{tdb: SetupMySQLTest(t), tableName: ""},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "testSQLite3test1",
			args:    args{tdb: SetupSQLite3Test(t), tableName: TestData1},
			want:    []string{"id"},
			wantErr: false,
		},
		{
			name:    "testPostgreSQLtest1",
			args:    args{tdb: SetupPostgresTest(t), tableName: TestData1},
			want:    []string{"id"},
			wantErr: false,
		},
		{
			name:    "testMySQLtest1",
			args:    args{tdb: SetupMySQLTest(t), tableName: TestData1},
			want:    []string{"id"},
			wantErr: false,
		},
		{
			name:    "testSQLite3test2",
			args:    args{tdb: SetupSQLite3Test(t), tableName: TestPrimary},
			want:    []string{"id1", "id2", "id3"},
			wantErr: false,
		},
		{
			name:    "testPostgreSQLtest2",
			args:    args{tdb: SetupPostgresTest(t), tableName: TestPrimary},
			want:    []string{"id1", "id2", "id3"},
			wantErr: false,
		},
		{
			name:    "testMySQLtest2",
			args:    args{tdb: SetupMySQLTest(t), tableName: TestPrimary},
			want:    []string{"id1", "id2", "id3"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		tdb := tt.args.tdb
		tableName := tt.args.tableName
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.tableName != "" {
				tableName = createTestTable(t, tdb, tt.args.tableName)
			}
			got, err := tdb.GetPrimaryKey(tdb.DB, "", tableName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPrimaryKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPrimaryKey() = %v, want %v", got, tt.want)
			}
			if tableName != "" {
				dropTestTable(t, tdb, tableName)
			}
		})
	}
}