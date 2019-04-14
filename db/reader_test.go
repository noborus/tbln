package db_test

import (
	"reflect"
	"testing"

	"github.com/noborus/tbln"
	"github.com/noborus/tbln/db"
	_ "github.com/noborus/tbln/db/mysql"
	_ "github.com/noborus/tbln/db/postgres"
	_ "github.com/noborus/tbln/db/sqlite3"
)

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
			table:   TestOne,
			want:    wantTbln(t, TestOne),
			wantErr: false,
		},
		{
			name:    "testType",
			args:    args{schema: "", tableName: "testtype", pkey: nil},
			table:   TestType,
			want:    wantTbln(t, TestType),
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
			tt := tt
			tdb := tdb
			var tableName string
			if tt.table != "" {
				tableName = createTestTable(t, tdb, tt.table)
			}
			t.Run(tt.name+"_"+tdb.Name, func(t *testing.T) {
				got, err := db.ReadTableAll(tdb, tt.args.schema, tt.args.tableName)
				if (err != nil) != tt.wantErr {
					t.Errorf("(%s)TDB.ReadTable()error = %v, wantErr %v", tdb.Name, err, tt.wantErr)
					return
				}
				if tt.want != nil {
					if !reflect.DeepEqual(got.Rows, tt.want.Rows) {
						t.Errorf("(%s)ReadTableAll() = %v, want %v", tdb.Name, got.Rows, tt.want.Rows)
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
			tdb := tdb
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

func TestGetTableInfo(t *testing.T) {
	type args struct {
		tdb       *db.TDB
		schema    string
		tableName string
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
			args:    args{tdb: SetupSQLite3Test(t), schema: "", tableName: ""},
			table:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test1",
			args:    args{tdb: SetupSQLite3Test(t), schema: "", tableName: "test1"},
			table:   TestData1,
			want:    wantTbln(t, TestData1),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		var tableName string
		if tt.table != "" {
			tableName = createTestTable(t, tt.args.tdb, tt.table)
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := db.GetTableInfo(tt.args.tdb, tt.args.schema, tableName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTableInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !reflect.DeepEqual(got.Names(), tt.want.Names()) {
					t.Errorf("GetTableInfo() = %v, want %v", got, tt.want)
				}
			}
		})
		if tt.table != "" {
			dropTestTable(t, tt.args.tdb, tableName)
		}
	}
}
