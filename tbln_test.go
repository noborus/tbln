package tbln

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func openFile(t *testing.T, fileName string) *os.File {
	file, err := os.Open(fileName)
	if err != nil {
		t.Error(err)
	}
	return file
}

func decodeHashHelper(d string) []byte {
	b, err := hex.DecodeString(d)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func decode64Helper(d string) []byte {
	b, err := base64.StdEncoding.DecodeString(d)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func TestTBLN_AddRows(t *testing.T) {
	type fields struct {
		Definition *Definition
		Hash       map[string]string
		RowNum     int
		Rows       [][]string
	}
	type args struct {
		row []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "test1",
			fields:  fields{Definition: &Definition{columnNum: 1}},
			args:    args{row: []string{"1"}},
			wantErr: false,
		},
		{
			name:    "test2",
			fields:  fields{Definition: &Definition{columnNum: 2}},
			args:    args{row: []string{"1"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		tb := &TBLN{
			Definition: tt.fields.Definition,
			RowNum:     tt.fields.RowNum,
			Rows:       tt.fields.Rows,
		}
		t.Run(tt.name, func(t *testing.T) {
			if err := tb.AddRows(tt.args.row); (err != nil) != tt.wantErr {
				t.Errorf("TBLN.AddRows() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkRow(t *testing.T) {
	type args struct {
		ColumnNum int
		row       []string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "test1",
			args:    args{row: []string{"1"}, ColumnNum: 1},
			want:    1,
			wantErr: false,
		},
		{
			name:    "test2",
			args:    args{row: []string{"1", "2"}, ColumnNum: 0},
			want:    2,
			wantErr: false,
		},
		{
			name:    "test3",
			args:    args{row: []string{"1"}, ColumnNum: 2},
			want:    2,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkRow(tt.args.ColumnNum, tt.args.row)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkRow() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkRow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTBLN_SumHash(t *testing.T) {
	type fields struct {
		Definition *Definition
		Hash       map[string]string
		RowNum     int
		Rows       [][]string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name:    "testBlank",
			fields:  fields{Definition: NewDefinition()},
			want:    "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
			wantErr: false,
		},
		{
			name:    "testOneRow",
			fields:  fields{Definition: NewDefinition(), Rows: [][]string{{"1"}}, RowNum: 1},
			want:    "3c5c7b4b1fcd47206cbc619c23b59a27b97f730abca146d67157d2d9df8ca9dc",
			wantErr: false,
		},
		{
			name:    "testNames",
			fields:  fields{Definition: &Definition{Extras: make(map[string]Extra), Hashes: make(map[string][]byte), names: []string{"id"}, columnNum: 1}, Rows: [][]string{{"1"}}, RowNum: 1},
			want:    "e5ce5f72c836840efdbcbf7639075966944253ef438a305761d888158a6b22a8",
			wantErr: false,
		},
		{
			name:    "testFullRow",
			fields:  fields{Definition: &Definition{Extras: make(map[string]Extra), Hashes: make(map[string][]byte), names: []string{"id", "name"}, types: []string{"int", "text"}, columnNum: 2}, Rows: [][]string{{"1", "test"}}, RowNum: 1},
			want:    "fcc150288d592d5c0cf13eed4b1054f6fadbfd2c48cde10954b44d6b7fc42623",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tb := &TBLN{
				Definition: tt.fields.Definition,
				RowNum:     tt.fields.RowNum,
				Rows:       tt.fields.Rows,
			}
			if err := tb.SetNames(tb.names); err != nil {
				t.Errorf("TBLN.SetNames error = %v", err)
				return

			}
			if err := tb.SetTypes(tb.types); err != nil {
				t.Errorf("TBLN.SetTypes error = %v", err)
				return
			}
			tb.ToTargetHash("name", true)
			tb.ToTargetHash("type", true)
			if err := tb.SumHash(SHA256); (err != nil) != tt.wantErr {
				t.Errorf("TBLN.SumHash(SHA256) error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			got := fmt.Sprintf("%x", tb.Hashes["sha256"])
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TBLN.SumHash(SHA256) = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTBLN_calculateHash(t *testing.T) {
	type args struct {
		hashType string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "testErr",
			args:    args{hashType: "err"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "testSHA256",
			args: args{hashType: "sha256"},
			want: []byte{
				0xe3, 0xb0, 0xc4, 0x42, 0x98, 0xfc, 0x1c, 0x14, 0x9a, 0xfb, 0xf4, 0xc8, 0x99, 0x6f, 0xb9, 0x24,
				0x27, 0xae, 0x41, 0xe4, 0x64, 0x9b, 0x93, 0x4c, 0xa4, 0x95, 0x99, 0x1b, 0x78, 0x52, 0xb8, 0x55,
			},
			wantErr: false,
		},
		{
			name: "testSHA512",
			args: args{hashType: "sha512"},
			want: []byte{
				0xcf, 0x83, 0xe1, 0x35, 0x7e, 0xef, 0xb8, 0xbd, 0xf1, 0x54, 0x28, 0x50, 0xd6, 0x6d, 0x80, 0x07,
				0xd6, 0x20, 0xe4, 0x05, 0x0b, 0x57, 0x15, 0xdc, 0x83, 0xf4, 0xa9, 0x21, 0xd3, 0x6c, 0xe9, 0xce,
				0x47, 0xd0, 0xd1, 0x3c, 0x5d, 0x85, 0xf2, 0xb0, 0xff, 0x83, 0x18, 0xd2, 0x87, 0x7e, 0xec, 0x2f,
				0x63, 0xb9, 0x31, 0xbd, 0x47, 0x41, 0x7a, 0x81, 0xa5, 0x38, 0x32, 0x7a, 0xf9, 0x27, 0xda, 0x3e,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := NewTBLN()
			got, err := tb.calculateHash(tt.args.hashType)
			if (err != nil) != tt.wantErr {
				t.Errorf("TBLN.calculateHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TBLN.calculateHash() = %x, want %x", got, tt.want)
			}
		})
	}
}

func TestTBLN_Sign(t *testing.T) {
	type args struct {
		name string
		pkey []byte
	}
	tests := []struct {
		name     string
		fileName string
		args     args
		want     map[string]Signature
		wantErr  bool
	}{
		{
			name:     "testErr1",
			fileName: "",
			args: args{
				name: "test",
				pkey: []byte(""),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:     "testErr2",
			fileName: filepath.Join("testdata", "abc.tbln"),
			args: args{
				name: "test",
				pkey: []byte(""),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:     "test1",
			fileName: filepath.Join("testdata", "abc.tbln"),
			args: args{
				name: "test",
				pkey: decodeHashHelper("a490a355aca9be30a72710e81a05ef0d8fea0877e7031bc7cd66970f7f9e3537ee40c42c0529991cbf64ca3b71335902749b1b33a54b0c564b7d6995b97d6ced"),
			},
			want: map[string]Signature{
				"test": {
					sign:      decodeHashHelper("b289138915aaa0510962bac8acd253e753dfb7e41d220fcba4a83be16a08282349dc2c5198ef30c0c45e73b9859a80916ff758ad483814353d23ad55e8681205"),
					algorithm: ED25519,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var tb *TBLN
			if tt.fileName == "" {
				tb = &TBLN{}
			} else {
				var err error
				f := openFile(t, tt.fileName)
				tb, err = ReadAll(f)
				if err != nil {
					t.Errorf("TBLN file open %s", err)
				}
			}
			got, err := tb.Sign(tt.args.name, tt.args.pkey)
			if (err != nil) != tt.wantErr {
				t.Errorf("TBLN.Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TBLN.Sign() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTBLN_VerifySignature(t *testing.T) {
	type args struct {
		name   string
		pubkey []byte
	}
	tests := []struct {
		name     string
		fileName string
		args     args
		want     bool
	}{
		{
			name:     "testABCfalse",
			fileName: filepath.Join("testdata", "abc-s.tbln"),
			args: args{
				name:   "test",
				pubkey: []byte(""),
			},
			want: false,
		},
		{
			name:     "testABCtrue",
			fileName: filepath.Join("testdata", "abc-s.tbln"),
			args: args{
				name:   "test",
				pubkey: decode64Helper("7kDELAUpmRy/ZMo7cTNZAnSbGzOlSwxWS31plbl9bO0="),
			},
			want: true,
		},
		{
			name:     "testABnoSign",
			fileName: filepath.Join("testdata", "abc.tbln"),
			args: args{
				name:   "test",
				pubkey: decode64Helper("7kDELAUpmRy/ZMo7cTNZAnSbGzOlSwxWS31plbl9bO0="),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		f := openFile(t, tt.fileName)
		tb, err := ReadAll(f)
		if err != nil {
			t.Errorf("TBLN file open %s", err)
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := tb.VerifySignature(tt.args.name, tt.args.pubkey); got != tt.want {
				t.Errorf("TBLN.VerifySignature() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTBLN_Verify(t *testing.T) {
	tests := []struct {
		name       string
		Definition *Definition
		want       bool
	}{
		{
			name:       "testNot",
			Definition: NewDefinition(),
			want:       false,
		},
		{
			name:       "testErr",
			Definition: &Definition{Hashes: map[string][]byte{"sha256": []byte("testErr")}},
			want:       false,
		},
		{
			name:       "testErr2",
			Definition: &Definition{Hashes: map[string][]byte{"sha255": []byte("testErr")}},
			want:       false,
		},
		{
			name:       "test1",
			Definition: &Definition{Hashes: map[string][]byte{"sha256": decodeHashHelper("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")}},
			want:       true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tb := &TBLN{
				Definition: tt.Definition,
			}
			if got := tb.Verify(); got != tt.want {
				t.Errorf("TBLN.Verify() = %v, want %v", got, tt.want)
			}
		})
	}
}
