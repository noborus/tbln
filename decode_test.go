package tbln

import (
	"bytes"
	"io"
	"testing"
)

func TestTbln_Decode(t *testing.T) {
	data := `
# comment
; name: | id | name | age |
; type: | int | text | int |
| 1 | Bob | 19 |
| 2 | Alice | 14 |
| 3 | Henry | 18 |
`
	type fields struct {
		Rows     [][]string
		Names    []string
		Types    []string
		Comments []string
		Extra    []string
		cap      int
		colNum   int
	}
	type args struct {
		reader io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{bytes.NewBufferString(data)},
			wantErr: false,
		},
		{
			name:    "test2",
			args:    args{bytes.NewBufferString("||")},
			wantErr: true,
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
			err := tbln.Decode(tt.args.reader)
			if (err != nil) != tt.wantErr {
				t.Log(tbln)
				t.Errorf("Tbln.Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTbln_analyzeExt(t *testing.T) {
	type fields struct {
		Rows     [][]string
		Names    []string
		Types    []string
		Comments []string
		Extra    []string
		cap      int
		colNum   int
	}
	type args struct {
		ext []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
			if err := tbln.analyzeExt(tt.args.ext); (err != nil) != tt.wantErr {
				t.Errorf("Tbln.analyzeExt() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_unescape(t *testing.T) {
	type args struct {
		rec []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "test1", args: args{[]string{"||"}}, want: []string{"|"}},
		{name: "test1", args: args{[]string{"|||"}}, want: []string{"||"}},
		{name: "test1", args: args{[]string{"ab||c"}}, want: []string{"ab|c"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			unescape(tt.args.rec)
			if tt.args.rec[0] != tt.want[0] {
				t.Errorf("Tbln.unescape() error = %v, want %v", tt.args.rec, tt.want)
			}
		})
	}
}
