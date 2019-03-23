package keystore

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/noborus/tbln"
	"github.com/noborus/tbln/cmd/key"
)

// Regist registers public key in keystore
func Regist(keyStore string, keyName string, pubkey []byte) error {
	if _, err := os.Stat(keyStore); os.IsNotExist(err) {
		return newKeyStore(keyStore, keyName, pubkey)
	}
	return addKeyStore(keyStore, keyName, pubkey)
}

func newKeyStore(keyStore string, keyName string, pubkey []byte) error {
	file, err := os.OpenFile(keyStore, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	t, err := key.GeneratePublic(keyName, pubkey)
	if err != nil {
		return err
	}
	t.Comments = []string{fmt.Sprintf("TBLN KeyStore")}
	err = tbln.WriteAll(file, t)
	if err != nil {
		return fmt.Errorf("KeyStore write: %s", err)
	}
	return nil
}

func addKeyStore(keyStore string, keyName string, pubkey []byte) error {
	file, err := os.OpenFile(keyStore, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	pkeys, err := tbln.ReadAll(file)
	if err != nil {
		return fmt.Errorf("KeyStore read: %s", err)
	}
	pkeys.Extras["updated_at"] = tbln.NewExtra(time.Now().Format(time.RFC3339), false)
	pEnc := base64.StdEncoding.EncodeToString(pubkey)
	err = pkeys.AddRows([]string{keyName, tbln.ED25519, pEnc})
	if err != nil {
		return err
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("KeyStore seek: %s", err)
	}
	err = tbln.WriteAll(file, pkeys)
	if err != nil {
		return fmt.Errorf("KeyStore write: %s", err)
	}
	return nil
}

// Search searches public key from keystore
func Search(keyStore string, keyName string) ([]byte, error) {
	pubFile, err := os.Open(keyStore)
	if err != nil {
		return nil, err
	}
	defer pubFile.Close()
	pt, err := tbln.ReadAll(pubFile)
	if err != nil {
		return nil, fmt.Errorf("key store read: %s", err)
	}
	for _, row := range pt.Rows {
		if row[0] == keyName {
			return base64.StdEncoding.DecodeString(row[2])
		}
	}
	return nil, fmt.Errorf("public key not found: %s", keyName)
}
