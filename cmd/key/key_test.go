package key

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/noborus/tbln"
)

type MockPasswordReader struct {
}

func (pr MockPasswordReader) ReadPasswordPrompt(prompt string) ([]byte, error) {
	return []byte("test"), nil
}

func TestGenerateKey(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pub, priv, err := GenerateKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(pub) != 32 {
				t.Errorf("GenerateKey() public key length error %d", len(pub))
				return
			}
			if len(priv) != 64 {
				t.Errorf("GenerateKey() private key length error %d", len(priv))
				return
			}
		})
	}
}

func Test_genTBLNPrivate(t *testing.T) {
	type args struct {
		keyName string
	}
	tests := []struct {
		name    string
		pr      MockPasswordReader
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "test1",
			pr:      MockPasswordReader{},
			args:    args{"test1"},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		_, priv, err := GenerateKey()
		if err != nil {
			t.Fatal(err)
		}
		t.Run(tt.name, func(t *testing.T) {
			PR = tt.pr
			got, err := genTBLNPrivate(tt.args.keyName, priv)
			if (err != nil) != tt.wantErr {
				t.Errorf("genTBLNPrivate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.RowNum, tt.want) {
				t.Errorf("genTBLNPrivate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWritePrivateFile(t *testing.T) {
	type args struct {
		fileName   string
		keyName    string
		privateKey []byte
	}
	tests := []struct {
		name    string
		pr      MockPasswordReader
		args    args
		wantErr bool
	}{
		{
			name:    "test1",
			pr:      MockPasswordReader{},
			args:    args{fileName: "test", keyName: "test", privateKey: []byte("test")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tmp, _ := ioutil.TempFile("", "_")
		fileName := tmp.Name()
		defer func() {
			os.Remove(fileName)
		}()
		t.Run(tt.name, func(t *testing.T) {
			if err := WritePrivateFile(fileName, tt.args.keyName, tt.args.privateKey); (err != nil) != tt.wantErr {
				t.Errorf("WritePrivateFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStdInPasswordReader_ReadPasswordPrompt(t *testing.T) {
	type args struct {
		prompt string
	}
	tests := []struct {
		name    string
		pr      MockPasswordReader
		args    args
		wantErr bool
	}{
		{
			name:    "test1",
			pr:      MockPasswordReader{},
			args:    args{"testPrompt"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pr := tt.pr
			_, err := pr.ReadPasswordPrompt(tt.args.prompt)
			if (err != nil) != tt.wantErr {
				t.Errorf("StdInPasswordReader.ReadPasswordPrompt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestWritePublicFile(t *testing.T) {
	type args struct {
		keyName string
		public  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "testErr",
			args:    args{keyName: "", public: []byte("test")},
			wantErr: true,
		},
		{
			name:    "test1",
			args:    args{keyName: "test", public: []byte("test")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tmp, _ := ioutil.TempFile("", "_")
		fileName := tmp.Name()
		defer func() {
			os.Remove(fileName)
		}()
		t.Run(tt.name, func(t *testing.T) {
			if err := WritePublicFile(fileName, tt.args.keyName, tt.args.public); (err != nil) != tt.wantErr {
				t.Errorf("WritePublicFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetPrivateKey(t *testing.T) {
	type args struct {
		privFileName string
		keyName      string
		prompt       bool
	}
	tests := []struct {
		name    string
		pr      MockPasswordReader
		args    args
		wantErr bool
	}{
		{
			name:    "testErr",
			args:    args{privFileName: "", keyName: "test", prompt: false},
			wantErr: true,
		},
		{
			name:    "testErr2",
			args:    args{privFileName: "../../testdata/test.pub", keyName: "test", prompt: false},
			wantErr: true,
		},
		{
			name:    "test1",
			args:    args{privFileName: "../../testdata/test.key", keyName: "test", prompt: false},
			wantErr: false,
		},
		{
			name:    "test1",
			args:    args{privFileName: "../../testdata/test.key", keyName: "test", prompt: true},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		PR = tt.pr
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetPrivateKey(tt.args.privFileName, tt.args.keyName, tt.args.prompt)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPrivateKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGenTBLNPublic(t *testing.T) {
	type args struct {
		keyName string
		pubkey  []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *tbln.Tbln
		wantErr bool
	}{
		{
			name:    "testErr1",
			args:    args{"", []byte("a")},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "testErr2",
			args:    args{"t", []byte("")},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenTBLNPublic(tt.args.keyName, tt.args.pubkey)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenTBLNPublic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenTBLNPublic() = %v, want %v", got, tt.want)
			}
		})
	}
}
