package tbln

import (
	"bytes"
	"testing"
)

func TestEncode(t *testing.T) {
	data := [][]string{
		{"1", "Bob", "19"},
		{"2", "Alice", "14"},
	}
	dataw := []byte("| 1 | Bob | 19 |\n| 2 | Alice | 14 |\n")

	tbl := NewTbln(2)
	for _, row := range data {
		tbl.AddRow(row)
	}
	b := &bytes.Buffer{}
	err := tbl.Encode(b)
	if err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}
	if bytes.Compare(b.Bytes(), dataw) != 0 {
		t.Errorf("out=[\n%v] ,should=[\n%v]", b, string(dataw))
	}
}

func TestEscapeEncode(t *testing.T) {
	data := [][]string{
		{"a|b", "a | b", "abc ||| def"},
		{"|", " | ", " |"},
	}
	dataw := []byte("| a||b | a || b | abc |||| def |\n| || |  ||  |  || |\n")
	tbl := NewTbln(2)
	for _, row := range data {
		tbl.AddRow(row)
	}
	b := &bytes.Buffer{}
	err := tbl.Encode(b)
	if err != nil {
		t.Errorf("Unexpected error: %s\n", err)
	}
	if bytes.Compare(b.Bytes(), dataw) != 0 {
		t.Errorf("out=[\n%v] ,should=[\n%v]", b, string(dataw))
	}
}

func TestTbln_Encode(t *testing.T) {
	type fields struct {
		Rows     [][]string
		Names    []string
		Types    []string
		Comments []string
		Extra    []string
		cap      int
		colNum   int
	}
	tests := []struct {
		name    string
		fields  fields
		wantW   string
		wantErr bool
	}{
		{
			name:    "test1",
			fields:  fields{Rows: [][]string{{"1", "Bob", "19"}, {"2", "Alice", "14"}}},
			wantW:   "| 1 | Bob | 19 |\n| 2 | Alice | 14 |\n",
			wantErr: false,
		},
		{
			name:    "test2",
			fields:  fields{Comments: []string{"comment"}},
			wantW:   "# comment\n",
			wantErr: false,
		},
		{
			name:    "test3",
			fields:  fields{Names: []string{"id", "name"}, Types: []string{"int", "text"}},
			wantW:   "; name: | id | name |\n; type: | int | text |\n",
			wantErr: false,
		},
		{
			name:    "test4",
			fields:  fields{Rows: [][]string{{"1", "Bob"}, {"2", "Alice"}}, Names: []string{"id", "name"}, Types: []string{"int", "text"}},
			wantW:   "; name: | id | name |\n; type: | int | text |\n| 1 | Bob |\n| 2 | Alice |\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tbln := &Tbln{
				Rows:     tt.fields.Rows,
				Names:    tt.fields.Names,
				Types:    tt.fields.Types,
				Comments: tt.fields.Comments,
				Extra:    tt.fields.Extra,
				cap:      tt.fields.cap,
				colNum:   tt.fields.colNum,
			}
			w := &bytes.Buffer{}
			if err := tbln.Encode(w); (err != nil) != tt.wantErr {
				t.Errorf("Tbln.Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Tbln.Encode() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
