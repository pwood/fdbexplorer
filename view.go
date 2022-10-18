package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/statusjson"
	"github.com/rivo/tview"
	"sync"
)

type View struct {
	ch  chan State
	app *tview.Application

	pd *ProcessData
}

func (v *View) runData() {
	for s := range v.ch {
		v.pd.m.Lock()

		v.pd.processes = []statusjson.Process{}

		for _, process := range s.ClusterState.Cluster.Processes {
			v.pd.processes = append(v.pd.processes, process)
		}

		v.pd.Updated()

		v.pd.m.Unlock()
		v.app.Draw()
	}
}

func (v *View) run() {
	v.pd = &ProcessData{m: &sync.RWMutex{}, sortBy: SortIPAddress}

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	pages := tview.NewPages()

	locality := tview.NewTable().SetContent(&ProcessView{
		pd:      v.pd,
		columns: []ColumnId{ColumnIPAddressPort, ColumnStatus, ColumnMachine, ColumnLocality, ColumnClass, ColumnRoles},
	})
	locality.SetFixed(1, 0)
	locality.SetSelectable(true, false)
	pages.AddPage("locality", locality, true, true)

	usage := tview.NewTable().SetContent(&ProcessView{
		pd:      v.pd,
		columns: []ColumnId{ColumnIPAddressPort, ColumnRoles, ColumnCPUActivity, ColumnRAMUsage, ColumnNetworkActivity, ColumnDiskUsage, ColumnDiskActivity},
	})
	usage.SetFixed(1, 0)
	usage.SetSelectable(true, false)
	pages.AddPage("usage", usage, true, false)

	pageIndex := []string{"locality", "usage"}
	pageNow := 0

	pages.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			pageNow--
			if pageNow < 0 {
				pageNow = len(pageIndex) - 1
			}
		case tcell.KeyRight:
			pageNow++
			if pageNow >= len(pageIndex) {
				pageNow = 0
			}
		default:
			return event
		}

		pages.SwitchToPage(pageIndex[pageNow])
		return nil
	})

	grid := tview.NewGrid().SetRows(5, 0).SetColumns(0, 0, 0).SetBorders(true)
	grid.AddItem(newPrimitive("Header"), 0, 0, 1, 2, 0, 0, false)
	grid.AddItem(newPrimitive("Stats"), 0, 2, 1, 1, 0, 0, false)
	grid.AddItem(pages, 1, 0, 1, 3, 0, 0, true)

	sortIndex := []SortKey{SortIPAddress, SortRole, SortClass}
	sortNow := 0

	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyF1:
			sortNow++
			if sortNow >= len(sortIndex) {
				sortNow = 0
			}

			v.pd.m.Lock()
			v.pd.sortBy = sortIndex[sortNow]
			v.pd.Updated()
			v.pd.m.Unlock()
		default:
			return event
		}

		pages.SwitchToPage(pageIndex[pageNow])
		return nil
	})

	v.app = tview.NewApplication().SetRoot(grid, true).SetFocus(locality)

	go v.runData()

	if err := v.app.Run(); err != nil {
		panic(err)
	}
}
