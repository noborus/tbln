package keystore

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/noborus/tbln"
	"github.com/noborus/tbln/cmd/key"
)

// Regist registers public key in keystore
func Regist(keyStore string, keyName string, pubkey []byte) error {
	pubk, err := os.OpenFile(keyStore, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer pubk.Close()

	pkeys, err := tbln.ReadAll(pubk)
	if err != nil {
		return fmt.Errorf("KeyStore read: %s", err)
	}
	err = pkeys.AddRows([]string{keyName, tbln.ED25519, key.PubEncode(pubkey)})
	if err != nil {
		return err
	}
	_, err = pubk.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("KeyStore seek: %s", err)
	}
	err = tbln.WriteAll(pubk, pkeys)
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
