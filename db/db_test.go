package db_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/noborus/tbln"
	"github.com/noborus/tbln/db"
	_ "github.com/noborus/tbln/db/mysql"
	_ "github.com/noborus/tbln/db/postgres"
	_ "github.com/noborus/tbln/db/sqlite3"
)

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

var TestOne = `; name: | one |
; type: | text |
; TableName: testone
| a |
| b |
| c |
`

var TestType = `; name: | a | b | c | d | e | f |
; type: | int | bigint | numeric | timestamp | text | bool |
; TableName: testtype
| 1 | 239239023290 | 3 | 2019-03-14T17:09:14Z | test | true |
| 2 | 99323020392 | 99320 | 1999-03-14T17:09:14Z | t | false |
| 3 | 329 | 329 | 2036-03-14T08:09:14Z | t | true |
`

func createTestTable(t *testing.T, tdb *db.TDB, data string) string {
	t.Helper()
	at := dataTBLN(t, data)
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

func dataTBLN(t *testing.T, data string) *tbln.TBLN {
	t.Helper()
	r := bytes.NewBufferString(data)
	at, err := tbln.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return at
}

func dropTestTable(t *testing.T, tdb *db.TDB, tableName string) {
	t.Helper()
	sql := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	_ = tdb.Begin()
	/*
		if err != nil {
			t.Fatal(err)
		}
	*/
	_, err := tdb.Tx.Exec(sql)
	if err != nil {
		t.Fatal(err)
	}
	err = tdb.Commit()
	if err != nil {
		t.Fatal(err)
	}
}

func fileRead(t *testing.T, fileName string) *tbln.TBLN {
	t.Helper()
	f, err := os.Open(fileName)
	if err != nil {
		t.Fatal(err)
	}
	tb, err := tbln.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	return tb
}

func setHash(tb *tbln.TBLN) *tbln.TBLN {
	tb.AllTargetHash(false)
	tb.ToTargetHash("name", true)
	tb.ToTargetHash("type", true)
	return tb
}

func SetupPostgresTest(t *testing.T) *db.TDB {
	t.Helper()
	conn, err := db.Open("postgres", os.Getenv("SESSION_PG_TEST_DSN"))
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

func SetupSQLite3Test(t *testing.T) *db.TDB {
	t.Helper()
	conn, err := db.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

func SetupMySQLTest(t *testing.T) *db.TDB {
	t.Helper()
	myDsn := os.Getenv("SESSION_MY_TEST_DSN")
	if myDsn == "" {
		myDsn = "root@/" + "test_db"
	}
	conn, err := db.Open("mysql", myDsn)
	if err != nil {
		t.Fatal(err)
	}
	return conn
}

// Test the file targeted by hash only for name and type(sethash(name,type)).
func TestWriteRead(t *testing.T) {
	// db.Debug(true)
	testDataPath := "../testdata"
	tests := []struct {
		name    string
		wantErr bool
	}{
		{name: "abc-n.tbln", wantErr: true},
		{name: "test_type.tbln", wantErr: false},
		{name: "ascii_table.tbln", wantErr: false},
		{name: "eto.tbln", wantErr: false},
	}
	dbDriver := []*db.TDB{
		SetupPostgresTest(t),
		SetupMySQLTest(t),
		SetupSQLite3Test(t),
	}
	for _, tdb := range dbDriver {
		for _, tt := range tests {
			tt := tt
			t.Run(tdb.Name+"_"+tt.name, func(t *testing.T) {
				wt := fileRead(t, filepath.Join(testDataPath, tt.name))
				tableName := wt.TableName()
				wGot := wt.Hashes["sha256"]
				// reset
				_ = tdb.Commit()
				err := tdb.Begin()
				if err != nil {
					t.Fatal(err)
				}
				err = db.WriteTable(tdb, wt, "", db.ReCreate, db.Normal)
				if err != nil {
					t.Fatal(err)
				}
				err = tdb.Commit()
				if err != nil {
					t.Fatal(err)
				}
				rt, err := db.ReadTableAll(tdb, "", tableName)
				if err != nil {
					t.Fatal(err)
				}
				rt = setHash(rt)
				err = rt.SumHash("sha256")
				if err != nil {
					t.Fatal(err)
				}
				rGot := rt.Hashes["sha256"]
				if (string(rGot) != string(wGot)) != tt.wantErr {
					t.Errorf("Hash name = %v, hash = %x, wantHash %x", tt.name, rGot, wGot)
					/*
						tbln.WriteAll(os.Stdout, wt)
						tbln.WriteAll(os.Stdout, rt)
					*/
					return
				}
			})
		}
	}
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
