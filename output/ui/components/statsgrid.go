package components

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"sync"
)

func NewStatsGrid[D any](grid [][]ColumnDef[D]) *StatsGrid[D] {
	return &StatsGrid[D]{
		grid: grid,

		m: &sync.RWMutex{},
	}
}

type StatsGrid[D any] struct {
	tview.TableContentReadOnly

	grid [][]ColumnDef[D]

	m    *sync.RWMutex
	data D
}

func (dt *StatsGrid[D]) Update(d D) {
	dt.m.Lock()
	dt.data = d
	dt.m.Unlock()
}

func (dt *StatsGrid[D]) GetCell(row, column int) *tview.TableCell {
	actualCol := column / 2
	partCol := column % 2

	cell := dt.grid[row][actualCol]

	if partCol == 0 {
		return tview.NewTableCell(cell.Name()).SetExpansion(1).SetTextColor(tcell.ColorYellow)
	} else {
		return tview.NewTableCell(cell.Data(dt.data)).SetExpansion(1).SetTextColor(cell.Color(dt.data))
	}
}

func (dt *StatsGrid[D]) GetRowCount() int {
	return len(dt.grid)
}

func (dt *StatsGrid[D]) GetColumnCount() int {
	if len(dt.grid) == 0 {
		return 0
	}

	return len(dt.grid[0]) * 2
}
