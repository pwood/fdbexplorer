package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/statusjson"
	"github.com/rivo/tview"
	"strconv"
	"sync"
)

type View struct {
	ch  chan State
	app *tview.Application

	pd *ProcessData
}

func (v *View) runData() {
	for s := range v.ch {
		var processes []statusjson.Process

		for _, process := range s.ClusterState.Cluster.Processes {
			processes = append(processes, process)
		}

		v.pd.Update(processes)
		v.app.Draw()
	}
}

func (v *View) run() {
	v.pd = &ProcessData{m: &sync.RWMutex{}, sortBy: SortIPAddress, views: map[string][]statusjson.Process{}, viewFns: map[string]func(statusjson.Process) bool{}}

	pages := tview.NewPages()

	allView := v.pd.View("all", All)

	locality := tview.NewTable().SetContent(&ProcessTableContent{
		pv:      allView,
		columns: []ColumnId{ColumnIPAddressPort, ColumnStatus, ColumnMachine, ColumnLocality, ColumnClass, ColumnRoles, ColumnVersion, ColumnUptime},
	})
	locality.SetFixed(1, 0)
	locality.SetSelectable(true, false)
	pages.AddPage("0", locality, true, true)

	usage := tview.NewTable().SetContent(&ProcessTableContent{
		pv:      allView,
		columns: []ColumnId{ColumnIPAddressPort, ColumnRoles, ColumnCPUActivity, ColumnRAMUsage, ColumnNetworkActivity, ColumnDiskUsage, ColumnDiskActivity},
	})
	usage.SetFixed(1, 0)
	usage.SetSelectable(true, false)
	pages.AddPage("1", usage, true, false)

	storage := tview.NewTable().SetContent(&ProcessTableContent{
		pv:      v.pd.View("storage", RoleMatch("storage")),
		columns: []ColumnId{ColumnIPAddressPort, ColumnCPUActivity, ColumnRAMUsage, ColumnDiskUsage, ColumnDiskActivity, ColumnKVStorage, ColumnDurabilityRate, ColumnStorageLag, ColumnTotalQueries},
	})
	storage.SetFixed(1, 0)
	storage.SetSelectable(true, false)
	pages.AddPage("2", storage, true, false)

	logs := tview.NewTable().SetContent(&ProcessTableContent{
		pv:      v.pd.View("log", RoleMatch("log")),
		columns: []ColumnId{ColumnIPAddressPort, ColumnCPUActivity, ColumnRAMUsage, ColumnDiskUsage, ColumnDiskActivity, ColumnQueueStorage, ColumnDurabilityRate},
	})
	logs.SetFixed(1, 0)
	logs.SetSelectable(true, false)
	pages.AddPage("3", logs, true, false)

	pageIndex := []string{"Locality", "Usage Overview", "Storage Processes", "Log Processes"}

	info := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetTextAlign(tview.AlignCenter).
		SetHighlightedFunc(func(added, removed, remaining []string) {
			pages.SwitchToPage(added[0])
		})

	// Create the pages for all slides.
	previousSlide := func() {
		slide, _ := strconv.Atoi(info.GetHighlights()[0])
		slide = (slide - 1 + len(pageIndex)) % len(pageIndex)
		info.Highlight(strconv.Itoa(slide)).ScrollToHighlight()
	}
	nextSlide := func() {
		slide, _ := strconv.Atoi(info.GetHighlights()[0])
		slide = (slide + 1) % len(pageIndex)
		info.Highlight(strconv.Itoa(slide)).ScrollToHighlight()
	}

	for index, title := range pageIndex {
		_, _ = fmt.Fprintf(info, `%d ["%d"][yellow]%s[white][""]  `, index+1, index, title)
	}

	info.Highlight("0")

	help := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetTextAlign(tview.AlignLeft)

	_, _ = fmt.Fprintf(help, ` F1 [black:darkcyan]Sort[:-] `)

	grid := tview.NewGrid().SetRows(5, 1, 0, 1).SetColumns(0, 0, 0).SetBorders(true)
	grid.AddItem(tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("Header"), 0, 0, 1, 2, 0, 0, false)
	grid.AddItem(tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText("Stats"), 0, 2, 1, 1, 0, 0, false)
	grid.AddItem(info, 1, 0, 1, 3, 0, 0, false)
	grid.AddItem(pages, 2, 0, 1, 3, 0, 0, true)
	grid.AddItem(help, 3, 0, 1, 3, 0, 0, false)

	sortIndex := []SortKey{SortIPAddress, SortRole, SortClass}
	sortNow := 0

	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			previousSlide()
		case tcell.KeyRight:
			nextSlide()
		case tcell.KeyF1:
			sortNow++
			if sortNow >= len(sortIndex) {
				sortNow = 0
			}

			v.pd.Sort(sortIndex[sortNow])
		default:
			return event
		}
		return nil
	})

	v.app = tview.NewApplication().SetRoot(grid, true).SetFocus(locality)

	go v.runData()

	if err := v.app.Run(); err != nil {
		panic(err)
	}
}
