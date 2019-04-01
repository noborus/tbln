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

var WriteTestData = `; name: | id | name | age |
; type: | int | text | int |
; primarykey: | id |
; TableName: test_write
| 1 | Bob | 19 |
| 2 | Alice | 14 |
`

func SetupTbln(t *testing.T) *tbln.Tbln {
	r := bytes.NewBufferString(WriteTestData)
	at, err := tbln.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return at
}

func TestWriteTable(t *testing.T) {
	type args struct {
		tdb    *db.TDB
		tbln   *tbln.Tbln
		schema string
		cmode  db.CreateMode
		imode  db.InsertMode
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "testPostgreSQL",
			args:    args{tdb: SetupPostgresTest(t), tbln: SetupTbln(t), schema: "", cmode: db.ReCreate, imode: db.Normal},
			wantErr: false,
		},
		{
			name:    "testSQLite3",
			args:    args{tdb: SetupSQLite3Test(t), tbln: SetupTbln(t), schema: "", cmode: db.ReCreate, imode: db.Normal},
			wantErr: false,
		},
		{
			name:    "testMySQL",
			args:    args{tdb: SetupMySQLTest(t), tbln: SetupTbln(t), schema: "", cmode: db.ReCreate, imode: db.Normal},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.tdb.Begin()
			if err != nil {
				t.Fatal(err)
			}
			if err = db.WriteTable(tt.args.tdb, tt.args.tbln, tt.args.schema, tt.args.cmode, tt.args.imode); (err != nil) != tt.wantErr {
				t.Errorf("WriteTable() error = %v, wantErr %v", err, tt.wantErr)
			}
			err = tt.args.tdb.Commit()
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
