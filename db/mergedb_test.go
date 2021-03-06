package db_test

import (
	"bytes"
	"testing"

	"github.com/noborus/tbln"
	"github.com/noborus/tbln/db"
	_ "github.com/noborus/tbln/db/mysql"
	_ "github.com/noborus/tbln/db/postgres"
	_ "github.com/noborus/tbln/db/sqlite3"
)

var TestMerge1 = `; name: | id | name | age |
; type: | int | text | int |
; primarykey: | id |
; TableName: test1
| 3 | Carol | 24 |
`

var TestMerge2 = `; name: | id | name | age |
; type: | int | text | int |
; primarykey: | id |
; TableName: test1
| 1 | Boob | 10 |
`

func TestTDB_MergeTableTBLN(t *testing.T) {
	type args struct {
		schema       string
		tableName    string
		shouldDelete bool
	}
	tests := []struct {
		name     string
		dataTBLN string
		args     args
		wantErr  bool
	}{
		{
			name:     "test1",
			dataTBLN: TestMerge1,
			args: args{
				schema:       "",
				tableName:    TestData1,
				shouldDelete: false,
			},
		},
		{
			name:     "test2",
			dataTBLN: TestMerge2,
			args: args{
				schema:       "",
				tableName:    TestData1,
				shouldDelete: false,
			},
		},
		{
			name:     "testDelete",
			dataTBLN: TestMerge2,
			args: args{
				schema:       "",
				tableName:    TestData1,
				shouldDelete: true,
			},
		},
	}
	dbDriver := []*db.TDB{
		SetupPostgresTest(t),
		SetupMySQLTest(t),
		SetupSQLite3Test(t),
	}
	for _, tdb := range dbDriver {
		for _, tt := range tests {
			tt := tt
			tdb := tdb
			var tableName string
			if tt.args.tableName != "" {
				tableName = createTestTable(t, tdb, tt.args.tableName)
			}
			t.Run(tt.name, func(t *testing.T) {
				mergeTBLN := dataTBLN(t, tt.dataTBLN)
				if err := tdb.MergeTableTBLN(tt.args.schema, tableName, mergeTBLN, tt.args.shouldDelete); (err != nil) != tt.wantErr {
					t.Errorf("TDB.MergeTableTBLN() %s error = %v, wantErr %v", tdb.Name, err, tt.wantErr)
				}
			})
			if tableName != "" {
				dropTestTable(t, tdb, tableName)
			}
		}
	}
}

func TestTDB_MergeTable(t *testing.T) {
	type args struct {
		schema       string
		tableName    string
		shouldDelete bool
	}
	tests := []struct {
		name     string
		dataTBLN string
		args     args
		wantErr  bool
	}{
		{
			name:     "test1",
			dataTBLN: TestMerge1,
			args: args{
				schema:       "",
				tableName:    TestData1,
				shouldDelete: false,
			},
		},
		{
			name:     "test2",
			dataTBLN: TestMerge2,
			args: args{
				schema:       "",
				tableName:    TestData1,
				shouldDelete: false,
			},
		},
		{
			name:     "testDelete",
			dataTBLN: TestMerge2,
			args: args{
				schema:       "",
				tableName:    TestData1,
				shouldDelete: true,
			},
		},
	}
	dbDriver := []*db.TDB{
		SetupPostgresTest(t),
		SetupMySQLTest(t),
		SetupSQLite3Test(t),
	}
	for _, tdb := range dbDriver {
		for _, tt := range tests {
			tt := tt
			tdb := tdb
			var tableName string
			if tt.args.tableName != "" {
				tableName = createTestTable(t, tdb, tt.args.tableName)
			}
			t.Run(tt.name, func(t *testing.T) {
				r := bytes.NewBufferString(tt.dataTBLN)
				other := tbln.NewReader(r)
				if err := tdb.MergeTable(tt.args.schema, tableName, other, tt.args.shouldDelete); (err != nil) != tt.wantErr {
					t.Errorf("TDB.MergeTable() error = %v, wantErr %v", err, tt.wantErr)
				}
			})
			if tableName != "" {
				dropTestTable(t, tdb, tableName)
			}
		}
	}
}
