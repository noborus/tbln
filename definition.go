package tbln

import (
	"encoding/hex"
	"fmt"
	"time"
)

// Definition is common table definition struct.
type Definition struct {
	columnNum int
	tableName string
	algorithm string
	Comments  []string
	Names     []string
	Types     []string
	Extras    map[string]Extra
	Hashes    map[string][]byte
	Signs     Signatures
}

// NewDefinition is create Definition struct.
func NewDefinition() *Definition {
	extras := make(map[string]Extra)
	extras["created_at"] = NewExtra(time.Now().Format(time.RFC3339), false)
	hashes := make(map[string][]byte)
	signs := make(Signatures)
	return &Definition{
		algorithm: ED25519,
		Extras:    extras,
		Hashes:    hashes,
		Signs:     signs,
	}
}

// Extra is table definition extra struct.
type Extra struct {
	value      interface{}
	hashTarget bool
}

// NewExtra is return new extra struct.
func NewExtra(value interface{}, hashTarget bool) Extra {
	return Extra{
		value:      value,
		hashTarget: hashTarget,
	}
}

// Value is return extra value.
func (e *Extra) Value() interface{} {
	return e.value
}

// ExtraValue is return extra value.
func (d *Definition) ExtraValue(ekey string) interface{} {
	ext := d.Extras[ekey]
	return ext.Value()
}

// Signature algorithm
const (
	ED25519 = "ED25519"
)

// Signatures is a map of signature name and signature.
type Signatures map[string]Signature

// Signature struct stores a signature, a name, and an algorithm.
type Signature struct {
	sign      []byte
	algorithm string
}

// TableName returns the Table Name.
func (d *Definition) TableName() string {
	return d.tableName
}

// SetExtra is set Extra of the Definition.
func (d *Definition) SetExtra(keyName string, value string) {
	target := false
	if len(value) != 0 {
		if len(d.Hashes) > 0 {
			target = true
		}
		d.Extras[keyName] = NewExtra(value, target)
	} else {
		delete(d.Extras, keyName)
	}
}

// SetTableName is set Table Name of the Definition..
func (d *Definition) SetTableName(name string) {
	d.tableName = name
	d.SetExtra("TableName", name)
}

// SetNames is set Column Name to the Table.
func (d *Definition) SetNames(names []string) error {
	if names != nil {
		if err := d.setColNum(len(names)); err != nil {
			return err
		}
	}
	d.Names = names
	d.SetExtra("name", JoinRow(names))
	return nil
}

// SetTypes is set Column Type to Table.
func (d *Definition) SetTypes(types []string) error {
	if types != nil {
		if err := d.setColNum(len(types)); err != nil {
			return err
		}
	}
	d.Types = types
	d.SetExtra("type", JoinRow(types))
	return nil
}

// ColumnNum returns the number of columns.
func (d *Definition) ColumnNum() int {
	return d.columnNum
}

func (d *Definition) setColNum(colNum int) error {
	if d.columnNum == 0 {
		d.columnNum = colNum
		return nil
	}
	if colNum != d.columnNum {
		return fmt.Errorf("number of columns is different")
	}
	return nil
}

// ToTargetHash is set as target of hash
func (d *Definition) ToTargetHash(key string, target bool) {
	if v, ok := d.Extras[key]; ok {
		v.hashTarget = target
		d.Extras[key] = v
	}
}

// AllTargetHash is set all target of hash
func (d *Definition) AllTargetHash(target bool) {
	for k, v := range d.Extras {
		v.hashTarget = target
		d.Extras[k] = v
	}
}

// SerializeHash returns a []byte that serializes Hash's map.
func (d *Definition) SerializeHash() []byte {
	hashes := make([]string, 0, len(d.Hashes))
	if val, ok := d.Hashes["sha256"]; ok {
		hashes = append(hashes, "sha256:"+fmt.Sprintf("%x", val))
	}
	if val, ok := d.Hashes["sha512"]; ok {
		hashes = append(hashes, "sha512:"+fmt.Sprintf("%x", val))
	}
	return []byte(JoinRow(hashes))
}

// SetSignatures is set signatures.
func (d *Definition) SetSignatures(sign []string) error {
	if len(sign) != 3 {
		return fmt.Errorf("not analyze signature")
	}
	b, err := hex.DecodeString(sign[2])
	if err != nil {
		return err
	}
	if sign[1] != ED25519 {
		return fmt.Errorf("not support algotithm: %s", sign[1])
	}
	d.Signs[sign[0]] = Signature{sign: b, algorithm: sign[1]}
	return nil
}

// SetHashes is set hashes.
func (d *Definition) SetHashes(hashes []string) error {
	if len(hashes) != 2 {
		return fmt.Errorf("not analyze hashes")
	}
	b, err := hex.DecodeString(hashes[1])
	if err != nil {
		return err
	}
	d.Hashes[hashes[0]] = b
	return nil
}
