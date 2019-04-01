package db_test

import (
	"bytes"
	"fmt"
	"os"
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

var WriteTestData2 = `; name: | id | name | age |
; type: | int | text | int |
; primarykey: | id |
| 1 | Bob | 19 |
| 2 | Alice | 14 |
`

func SetupTbln2(t *testing.T) *tbln.Tbln {
	r := bytes.NewBufferString(WriteTestData2)
	at, err := tbln.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	return at
}

func fileRead(t *testing.T, fileName string) *tbln.Tbln {
	f, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
		t.Fatal(err)
	}
	tb, err := tbln.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}
	return tb
}
func TestWriteTableFile(t *testing.T) {
	type args struct {
		tdb    *db.TDB
		schema string
		cmode  db.CreateMode
		imode  db.InsertMode
	}
	tests := []struct {
		name     string
		fileName string
		args     args
		wantErr  bool
	}{
		{
			name:     "testErr",
			fileName: "../testdata/simple-err.tbln",
			args: args{
				tdb:    SetupSQLite3Test(t),
				schema: "",
				cmode:  db.ReCreate,
				imode:  db.Normal,
			},
			wantErr: true,
		},
		{
			name:     "simple",
			fileName: "../testdata/simple.tbln",
			args: args{
				schema: "",
				cmode:  db.ReCreate,
				imode:  db.Normal,
			},
			wantErr: false,
		},
		{
			name:     "simpleIgnore",
			fileName: "../testdata/simple.tbln",
			args: args{
				schema: "",
				cmode:  db.ReCreate,
				imode:  db.OrIgnore,
			},
			wantErr: false,
		},
		{
			name:     "simplemulti",
			fileName: "../testdata/simplemulti.tbln",
			args: args{
				schema: "",
				cmode:  db.ReCreate,
				imode:  db.OrIgnore,
			},
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
			tb := fileRead(t, tt.fileName)
			t.Run(tt.name+tdb.Name, func(t *testing.T) {
				err := tdb.Begin()
				if err != nil {
					t.Fatal(err)
				}
				if err = db.WriteTable(tdb, tb, tt.args.schema, tt.args.cmode, tt.args.imode); (err != nil) != tt.wantErr {
					t.Errorf("WriteTable() error = %v, wantErr %v", err, tt.wantErr)
				}
				err = tdb.Commit()
				if err != nil {
					t.Fatal(err)
				}
			})
		}
	}
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
			name:    "testErr",
			args:    args{tdb: SetupPostgresTest(t), tbln: SetupTbln2(t), schema: "", cmode: db.ReCreate, imode: db.Normal},
			wantErr: true,
		},
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

func createDefinition1(t *testing.T) *tbln.Definition {
	at := tbln.NewDefinition()
	err := at.SetNames([]string{"a", "b"})
	if err != nil {
		t.Fatal(err)
	}
	err = at.SetTypes([]string{"int", "int"})
	if err != nil {
		t.Fatal(err)
	}
	return at
}

func createDefinition2(t *testing.T) *tbln.Definition {
	at := tbln.NewDefinition()
	err := at.SetNames([]string{"a", "b"})
	if err != nil {
		t.Fatal(err)
	}
	return at
}
func createDefinition3(t *testing.T) *tbln.Definition {
	at := tbln.NewDefinition()
	err := at.SetTypes([]string{"int", "int"})
	if err != nil {
		t.Fatal(err)
	}
	return at
}

func TestWriter_WriteDefinition(t *testing.T) {
	type fields struct {
		TDB        *db.TDB
		Definition *tbln.Definition
		ReplaceLN  bool
	}
	type args struct {
		cmode db.CreateMode
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "testErr1",
			fields: fields{
				TDB:        SetupPostgresTest(t),
				Definition: &tbln.Definition{},
				ReplaceLN:  true,
			},
			args:    args{cmode: db.CreateOnly},
			wantErr: true,
		},
		{
			name: "testErr2",
			fields: fields{
				TDB:        SetupPostgresTest(t),
				Definition: tbln.NewDefinition(),
				ReplaceLN:  true,
			},
			args:    args{cmode: db.CreateOnly},
			wantErr: true,
		},
		{
			name: "testNoCreate1",
			fields: fields{
				TDB:        SetupPostgresTest(t),
				Definition: createDefinition1(t),
				ReplaceLN:  true,
			},
			args:    args{cmode: db.NotCreate},
			wantErr: false,
		},
		{
			name: "testNoCreate2",
			fields: fields{
				TDB:        SetupPostgresTest(t),
				Definition: createDefinition2(t),
				ReplaceLN:  true,
			},
			args:    args{cmode: db.NotCreate},
			wantErr: false,
		},
		{
			name: "testNoCreate3",
			fields: fields{
				TDB:        SetupPostgresTest(t),
				Definition: createDefinition3(t),
				ReplaceLN:  true,
			},
			args:    args{cmode: db.NotCreate},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &db.Writer{
				TDB:        tt.fields.TDB,
				Definition: tt.fields.Definition,
				ReplaceLN:  tt.fields.ReplaceLN,
			}
			if err := w.WriteDefinition(tt.args.cmode); (err != nil) != tt.wantErr {
				t.Errorf("Writer.WriteDefinition() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
