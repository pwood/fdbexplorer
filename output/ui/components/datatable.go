package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"sync"
)

type ColumnDef[D any] interface {
	Name() string
	Data(D) string
	Color(D) tcell.Color
}

type ColumnImpl[D any] struct {
	ColName string
	DataFn  func(D) string
	ColorFn func(D) tcell.Color
}

func (c ColumnImpl[D]) Name() string {
	return c.ColName
}

func (c ColumnImpl[D]) Data(d D) string {
	return c.DataFn(d)
}

func (c ColumnImpl[D]) Color(d D) tcell.Color {
	if c.ColorFn == nil {
		return tcell.ColorWhite
	}
	return c.ColorFn(d)
}

func NewDataTable[D any](columns []ColumnDef[D]) *DataTable[D] {
	return &DataTable[D]{
		columns: columns,

		m: &sync.RWMutex{},
	}
}

type DataTable[D any] struct {
	tview.TableContentReadOnly

	columns []ColumnDef[D]

	m    *sync.RWMutex
	data []D
}

func (dt *DataTable[D]) Update(d []D) {
	dt.m.Lock()
	dt.data = d
	dt.m.Unlock()
}

func (dt *DataTable[D]) Get(row int) *D {
	dt.m.RLock()
	defer dt.m.RUnlock()

	return &dt.data[row-1]
}

func (dt *DataTable[D]) GetCell(row, column int) *tview.TableCell {
	col := dt.columns[column]

	if row == 0 {
		cell := tview.NewTableCell(col.Name()).SetTextColor(tcell.ColorAqua).SetSelectable(false)

		if len(col.Name()) > 1 {
			cell.SetExpansion(1)
		}

		return cell
	} else {
		dt.m.RLock()
		defer dt.m.RUnlock()

		di := dt.data[row-1]
		return tview.NewTableCell(col.Data(di)).SetTextColor(col.Color(di))
	}
}

func (dt *DataTable[D]) GetRowCount() int {
	dt.m.RLock()
	defer dt.m.RUnlock()

	return len(dt.data) + 1
}

func (dt *DataTable[D]) GetColumnCount() int {
	return len(dt.columns)
}
