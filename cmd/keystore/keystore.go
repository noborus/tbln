// Package keystore manages the public key keystore used by tbln command.
//
// The file name ~/.tbln/keystore.tbln is used by default.
//
// example:
//    # TBLN KeyStore
//    ; created_at: 2019-03-23T22:41:55+09:00
//    ; updated_at: 2019-03-28T10:56:21+09:00
//    | test | ED25519 | 7kDELAUpmRy/ZMo7cTNZAnSbGzOlSwxWS31plbl9bO0=
//
// Use Register to register byte public key.
// Use AddKey to register a file format public key.
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

// Register registers public key in keystore
// Register creates a new file if the file not exist.
// Add the public key, if the file already exists.
func Register(keyStore string, keyName string, pubkey []byte) error {
	if _, err := os.Stat(keyStore); os.IsNotExist(err) {
		return create(keyStore, keyName, pubkey)
	}
	return add(keyStore, keyName, pubkey)
}

func create(keyStore string, keyName string, pubkey []byte) error {
	file, err := os.OpenFile(keyStore, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("KeyStore create: %s", err)
	}
	t, err := key.GenTBLNPublic(keyName, pubkey)
	if err != nil {
		return fmt.Errorf("KeyStore generate: %s", err)
	}
	t.Comments = []string{"TBLN KeyStore"}
	err = tbln.WriteAll(file, t)
	if err != nil {
		return fmt.Errorf("KeyStore write: %s", err)
	}
	return nil
}

func createFile(keyStore string, fileName string) error {
	file, err := os.OpenFile(keyStore, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("KeyStore create: %s", err)
	}
	pubkey, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("public key open: %s:%s", fileName, err)
	}
	defer pubkey.Close()

	pt, err := tbln.ReadAll(pubkey)
	if err != nil {
		return fmt.Errorf("public key read: %s:%s", fileName, err)
	}
	return rewriteStore(file, pt)
}

func add(keyStore string, keyName string, pubkey []byte) error {
	if keyName == "" {
		return fmt.Errorf("keyname is required")
	}
	if len(pubkey) == 0 {
		return fmt.Errorf("pubkey is required")
	}
	file, err := os.OpenFile(keyStore, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("KeyStore open: %s", err)
	}
	defer file.Close()

	pKeys, err := tbln.ReadAll(file)
	if err != nil {
		return fmt.Errorf("KeyStore read: %s", err)
	}
	pEnc := base64.StdEncoding.EncodeToString(pubkey)
	err = pKeys.AddRows([]string{keyName, tbln.ED25519, pEnc})
	if err != nil {
		return fmt.Errorf("KeyStore add: %s", err)
	}
	return rewriteStore(file, pKeys)
}

// Search searches public key from keystore
func Search(keyStore string, keyName string) ([]byte, error) {
	pubFile, err := os.Open(keyStore)
	if err != nil {
		return nil, fmt.Errorf("KeyStore open: %s", err)
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

// List returns a list of public keys
func List(keyStore string) ([][]string, error) {
	pubFile, err := os.Open(keyStore)
	if err != nil {
		return nil, fmt.Errorf("KeyStore open: %s", err)
	}
	defer pubFile.Close()
	pt, err := tbln.ReadAll(pubFile)
	if err != nil {
		return nil, fmt.Errorf("key store read: %s", err)
	}
	return pt.Rows, nil
}

// AddKey adds a public key to the KeyStore.
func AddKey(keyStore string, fileName string) error {
	file, err := os.OpenFile(keyStore, os.O_RDWR, 0644)
	if err != nil {
		return createFile(keyStore, fileName)
	}
	defer file.Close()
	storeR := tbln.NewReader(file)

	pubkey, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("public key open: %s:%s", fileName, err)
	}
	defer pubkey.Close()
	keyR := tbln.NewReader(pubkey)

	pKeys, err := tbln.MergeAll(storeR, keyR, tbln.MergeUpdate)
	if err != nil {
		return fmt.Errorf("KeyStore AddKey: %s", err)
	}
	return rewriteStore(file, pKeys)
}

// DelKey deletes the corresponding public key from KeyStore.
// Delete the key that corresponds to num from 0.
// If num is less than 0, delete all keys with the corresponding name.
func DelKey(keyStore string, name string, num int) error {
	file, err := os.OpenFile(keyStore, os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("key store open: %s", err)
	}
	defer file.Close()

	pKeys, err := tbln.ReadAll(file)
	if err != nil {
		return fmt.Errorf("KeyStore read: %s", err)
	}
	var result [][]string
	n := 0
	for _, row := range pKeys.Rows {
		if row[0] == name {
			if num >= 0 && n != num {
				result = append(result, row)
			}
			n++
		} else {
			result = append(result, row)
		}
	}
	pKeys.Rows = result
	return rewriteStore(file, pKeys)
}

func rewriteStore(file *os.File, t *tbln.TBLN) error {
	var err error
	t.Comments = []string{"TBLN KeyStore"}
	t.Extras["updated_at"] = tbln.NewExtra(time.Now().Format(time.RFC3339), false)
	err = file.Truncate(0)
	if err != nil {
		return fmt.Errorf("KeyStore truncate: %s", err)
	}
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("KeyStore seek: %s", err)
	}
	err = tbln.WriteAll(file, t)
	if err != nil {
		return fmt.Errorf("KeyStore write: %s", err)
	}
	return nil
}
