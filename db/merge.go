package db

import (
	"io"

	"github.com/noborus/tbln"
)

type dml struct {
	insert [][]string
	update [][]string
	delete [][]string
}

func mergeTableRow(dml *dml, d *tbln.DiffRow, delete bool) *dml {
	switch d.Les {
	case 0:
		return dml
	case 1:
		dml.insert = append(dml.insert, d.Other)
		return dml
	case -1:
		if delete {
			dml.delete = append(dml.delete, d.Self)
		}
		return dml
	case 2:
		dml.update = append(dml.update, d.Other)
		return dml
	default:
		return dml
	}
}

// MergeTableTbln writes all rows to the table from Tbln.
func (tdb *TDB) MergeTableTbln(schema string, tableName string, otherTbln *tbln.Tbln, delete bool) error {
	orows := otherTbln.Rows
	var rps []RangePrimaryKey
	if !delete {
		pkpos, err := otherTbln.GetPKeyPos()
		if err == nil {
			if len(pkpos) > 0 {
				for _, p := range pkpos {
					rp := NewRangePrimaryKey(otherTbln.Names()[p], orows[0][p], orows[len(orows)-1][p])
					rps = append(rps, rp)
				}
			}
		}
	}
	other := tbln.NewOwnReader(otherTbln)

	self, err := tdb.ReadTable(schema, tableName, rps)
	if err != nil {
		return err
	}
	cmp, err := tbln.NewCompare(self, other)
	if err != nil {
		return err
	}
	dml := &dml{}
	for {
		dd, err := cmp.ReadDiffRow()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		dml = mergeTableRow(dml, dd, delete)
	}
	return tdb.mergeWrite(self.Definition, schema, cmp, dml, delete)
}

// MergeTable writes all rows to the table.
func (tdb *TDB) MergeTable(schema string, tableName string, other tbln.Reader, delete bool) error {
	self, err := tdb.ReadTable(schema, tableName, nil)
	if err != nil {
		return err
	}
	cmp, err := tbln.NewCompare(self, other)
	if err != nil {
		return err
	}
	dml := &dml{}
	for {
		dd, err := cmp.ReadDiffRow()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		dml = mergeTableRow(dml, dd, delete)
	}
	return tdb.mergeWrite(self.Definition, schema, cmp, dml, delete)
}

func (tdb *TDB) mergeWrite(definition *tbln.Definition, schema string, cmp *tbln.Compare, dml *dml, delete bool) error {
	w, err := NewWriter(tdb, definition)
	if err != nil {
		return err
	}
	if schema != "" {
		w.tableFullName = w.quoting(schema) + "." + w.quoting(w.TableName())
	} else {
		w.tableFullName = w.quoting(w.TableName())
	}
	if len(dml.insert) > 0 {
		err = w.insert(dml.insert)
		if err != nil {
			return err
		}
	}
	if len(dml.update) > 0 {
		err = w.update(dml.update, cmp.PK)
		if err != nil {
			return err
		}
	}
	if delete && (len(dml.delete) > 0) {
		err = w.delete(dml.delete, cmp.PK)
		if err != nil {
			return err
		}
	}
	return nil
}

func (w *Writer) insert(insRow [][]string) error {
	err := w.prepareInsert(Normal)
	if err != nil {
		return err
	}
	for _, ins := range insRow {
		err = w.WriteRow(ins)
		if err != nil {
			return err
		}
	}
	return w.stmt.Close()
}

func (w *Writer) update(updRow [][]string, pkeys []tbln.Pkey) error {
	err := w.prepareUpdate(pkeys)
	if err != nil {
		return err
	}
	for _, upd := range updRow {
		err = w.WriteRow(append(upd, tbln.ColumnPrimaryKey(pkeys, upd)...))
		if err != nil {
			return err
		}
	}
	return w.stmt.Close()
}

func (w *Writer) delete(delRow [][]string, pkeys []tbln.Pkey) error {
	err := w.prepareDelete(pkeys)
	if err != nil {
		return err
	}
	for _, del := range delRow {
		err = w.WriteRow(tbln.ColumnPrimaryKey(pkeys, del))
		if err != nil {
			return err
		}
	}
	return w.stmt.Close()
}
