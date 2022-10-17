package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"sync"
)

type ProcessData struct {
	tview.TableContentReadOnly

	processes []Process

	m *sync.RWMutex
}

type ProcessView struct {
	tview.TableContentReadOnly
	pd *ProcessData

	columns []string
	dataFn  func(Process, int) *tview.TableCell
}

func (v *ProcessView) GetCell(row, column int) *tview.TableCell {
	v.pd.m.RLock()
	defer v.pd.m.RUnlock()

	if row == 0 {
		return tview.NewTableCell(v.columns[column]).SetExpansion(1).SetTextColor(tcell.ColorAqua).SetSelectable(false)
	} else {
		return v.dataFn(v.pd.processes[row-1], column)
	}
}

func (v *ProcessView) GetRowCount() int {
	v.pd.m.RLock()
	defer v.pd.m.RUnlock()

	return len(v.pd.processes) + 1
}

func (v *ProcessView) GetColumnCount() int {
	return len(v.columns)
}
