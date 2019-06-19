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

func mergeTableRow(dml *dml, d *tbln.DiffRow, shouldDelete bool) *dml {
	switch d.Les {
	case 0:
		return dml
	case 1:
		dml.insert = append(dml.insert, d.Other)
		return dml
	case -1:
		if shouldDelete {
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
func (tdb *TDB) MergeTableTbln(schema string, tableName string, otherTbln *tbln.Tbln, shouldDelete bool) error {
	oRows := otherTbln.Rows
	var rps []RangePrimaryKey
	if !shouldDelete {
		pkPos, err := otherTbln.GetPKeyPos()
		if err == nil && (len(pkPos) > 0) {
			for _, p := range pkPos {
				rp := NewRangePrimaryKey(otherTbln.Names()[p], oRows[0][p], oRows[len(oRows)-1][p])
				rps = append(rps, rp)
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
		dml = mergeTableRow(dml, dd, shouldDelete)
	}
	return tdb.mergeWrite(self.Definition, schema, cmp, dml, shouldDelete)
}

// MergeTable writes all rows to the table.
func (tdb *TDB) MergeTable(schema string, tableName string, other tbln.Reader, shouldDelete bool) error {
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
		dml = mergeTableRow(dml, dd, shouldDelete)
	}
	return tdb.mergeWrite(self.Definition, schema, cmp, dml, shouldDelete)
}

func (tdb *TDB) mergeWrite(definition *tbln.Definition, schema string, cmp *tbln.Compare, dml *dml, shouldDelete bool) error {
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
	if shouldDelete && (len(dml.delete) > 0) {
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
