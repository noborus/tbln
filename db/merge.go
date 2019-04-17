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
		dml.insert = append(dml.insert, d.Dst)
		return dml
	case -1:
		if delete {
			dml.delete = append(dml.delete, d.Src)
		}
		return dml
	case 2:
		dml.update = append(dml.update, d.Dst)
		return dml
	default:
		return dml
	}
}

// MergeTable writes all rows to the table.
func (tdb *TDB) MergeTable(schema string, tableName string, dst tbln.Reader, delete bool) error {
	src, err := tdb.ReadTable(schema, tableName, nil)
	if err != nil {
		return err
	}
	cmp, err := tbln.NewCompare(src, dst)
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

	w, err := NewWriter(tdb, src.Definition)
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
		err = w.update(dml.update, cmp)
		if err != nil {
			return err
		}
	}
	if delete && (len(dml.delete) > 0) {
		err = w.delete(dml.delete, cmp)
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

func (w *Writer) update(updRow [][]string, cmp *tbln.Compare) error {
	err := w.prepareUpdate(cmp.PK)
	if err != nil {
		return err
	}
	for _, upd := range updRow {
		err = w.WriteRow(append(upd, tbln.ColumnPrimaryKey(cmp.PK, upd)...))
		if err != nil {
			return err
		}
	}
	return w.stmt.Close()
}

func (w *Writer) delete(delRow [][]string, cmp *tbln.Compare) error {
	err := w.prepareDelete(cmp.PK)
	if err != nil {
		return err
	}
	for _, del := range delRow {
		err = w.WriteRow(tbln.ColumnPrimaryKey(cmp.PK, del))
		if err != nil {
			return err
		}
	}
	return w.stmt.Close()
}
