package keystore

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/noborus/tbln"
)

const TestKeyStore = "../../testdata/keystore.tbln"

func fileCopy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}

func CreateTmpKeyStore(t *testing.T) string {
	tmp, _ := ioutil.TempFile("", "_")
	tmpKeyStore := tmp.Name()
	err := fileCopy(TestKeyStore, tmpKeyStore)
	if err != nil {
		t.Fatal(err)
	}
	return tmpKeyStore
}

func TestRegist(t *testing.T) {
	tmpKeyStore := CreateTmpKeyStore(t)
	defer func() {
		os.Remove(tmpKeyStore)
	}()
	type args struct {
		keyStore string
		keyName  string
		pubkey   []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "testErr",
			args:    args{keyStore: tmpKeyStore + "testErr", keyName: "", pubkey: []byte("")},
			wantErr: true,
		},
		{
			name:    "testErr2",
			args:    args{keyStore: tmpKeyStore, keyName: "", pubkey: []byte("test1")},
			wantErr: true,
		},
		{
			name:    "testRegist",
			args:    args{keyStore: tmpKeyStore, keyName: "test1", pubkey: []byte("test1")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Regist(tt.args.keyStore, tt.args.keyName, tt.args.pubkey); (err != nil) != tt.wantErr {
				t.Errorf("Regist() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_create(t *testing.T) {
	type args struct {
		keyStore string
		keyName  string
		pubkey   []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "testErr",
			args:    args{keyStore: "", keyName: "", pubkey: []byte("")},
			wantErr: true,
		},
		{
			name:    "testErr2",
			args:    args{keyStore: "/tmp/testCreate", keyName: "", pubkey: nil},
			wantErr: true,
		},
		{
			name:    "testCreate",
			args:    args{keyStore: "/tmp/testCreate", keyName: "test1", pubkey: []byte("test1")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := create(tt.args.keyStore, tt.args.keyName, tt.args.pubkey); (err != nil) != tt.wantErr {
				t.Errorf("create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_add(t *testing.T) {
	tmpKeyStore := CreateTmpKeyStore(t)
	defer func() {
		os.Remove(tmpKeyStore)
	}()
	type args struct {
		keyStore string
		keyName  string
		pubkey   []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "testErr",
			args:    args{keyStore: "", keyName: "", pubkey: []byte("")},
			wantErr: true,
		},
		{
			name:    "testErr2",
			args:    args{keyStore: tmpKeyStore, keyName: "", pubkey: nil},
			wantErr: true,
		},
		{
			name:    "testCreate",
			args:    args{keyStore: tmpKeyStore, keyName: "test1", pubkey: []byte("test1")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := add(tt.args.keyStore, tt.args.keyName, tt.args.pubkey); (err != nil) != tt.wantErr {
				t.Errorf("add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSearch(t *testing.T) {
	type args struct {
		keyStore string
		keyName  string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:    "testErr",
			args:    args{keyStore: TestKeyStore + "err", keyName: "err"},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test1",
			args: args{keyStore: TestKeyStore, keyName: "test"},
			want: []byte{0xee, 0x40, 0xc4, 0x2c, 0x05, 0x29, 0x99, 0x1c, 0xbf, 0x64, 0xca, 0x3b, 0x71, 0x33, 0x59, 0x02,
				0x74, 0x9b, 0x1b, 0x33, 0xa5, 0x4b, 0x0c, 0x56, 0x4b, 0x7d, 0x69, 0x95, 0xb9, 0x7d, 0x6c, 0xed},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Search(tt.args.keyStore, tt.args.keyName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				fmt.Printf("%x\n", got)
				t.Errorf("Search() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestList(t *testing.T) {
	type args struct {
		keyStore string
	}
	tests := []struct {
		name    string
		args    args
		want    [][]string
		wantErr bool
	}{
		{
			name:    "testErr",
			args:    args{keyStore: TestKeyStore + "err"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "test1",
			args:    args{keyStore: TestKeyStore},
			want:    [][]string{{"test", "ED25519", "7kDELAUpmRy/ZMo7cTNZAnSbGzOlSwxWS31plbl9bO0="}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := List(tt.args.keyStore)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddKey(t *testing.T) {
	tmpKeyStore := CreateTmpKeyStore(t)
	defer func() {
		os.Remove(tmpKeyStore)
	}()
	type args struct {
		keyStore string
		fileName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "testErr",
			args:    args{keyStore: tmpKeyStore + "err", fileName: ""},
			wantErr: true,
		},
		{
			name:    "test1",
			args:    args{keyStore: tmpKeyStore, fileName: "../../testdata/test.pub"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := AddKey(tt.args.keyStore, tt.args.fileName); (err != nil) != tt.wantErr {
				t.Errorf("AddKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDelKey(t *testing.T) {
	tmpKeyStore := CreateTmpKeyStore(t)
	defer func() {
		os.Remove(tmpKeyStore)
	}()
	type args struct {
		keyStore string
		name     string
		num      int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "testErr",
			args:    args{keyStore: tmpKeyStore + "err", name: "", num: 0},
			wantErr: true,
		},
		{
			name:    "test1",
			args:    args{keyStore: tmpKeyStore, name: "test1", num: 0},
			wantErr: false,
		},
		{
			name:    "test1",
			args:    args{keyStore: tmpKeyStore, name: "test1", num: 100},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DelKey(tt.args.keyStore, tt.args.name, tt.args.num); (err != nil) != tt.wantErr {
				t.Errorf("DelKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_rewriteStore(t *testing.T) {
	type args struct {
		tb *tbln.Tbln
	}
	tests := []struct {
		name    string
		mode    int
		args    args
		wantErr bool
	}{
		{
			name:    "test1Err1",
			mode:    os.O_RDONLY,
			args:    args{tb: tbln.NewTbln()},
			wantErr: true,
		},
		{
			name:    "test1Err2",
			mode:    os.O_APPEND,
			args:    args{tb: tbln.NewTbln()},
			wantErr: true,
		},
		{
			name:    "test1",
			mode:    os.O_RDWR,
			args:    args{tb: tbln.NewTbln()},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tmpKeyStore := CreateTmpKeyStore(t)
		defer func() {
			os.Remove(tmpKeyStore)
		}()
		t.Run(tt.name, func(t *testing.T) {
			file, err := os.OpenFile(tmpKeyStore, tt.mode, 0644)
			if err != nil {
				t.Fatal(err)
			}
			if err := rewriteStore(file, tt.args.tb); (err != nil) != tt.wantErr {
				t.Errorf("rewriteStore() error = %v, wantErr %v", err, tt.wantErr)
			}
			file.Close()
		})
	}
}
