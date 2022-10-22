package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/data"
	"github.com/pwood/fdbexplorer/data/fdb"
	"github.com/rivo/tview"
	"strconv"
	"sync"
)

func New(ch chan data.State) *View {
	return &View{ch: ch}
}

type View struct {
	ch  chan data.State
	app *tview.Application

	cd *ClusterData
}

func (v *View) runData() {
	for s := range v.ch {
		v.cd.Update(s.ClusterState)
		v.app.Draw()
	}
}

func (v *View) run() {
	v.cd = &ClusterData{m: &sync.RWMutex{}, sortBy: SortIPAddress, views: map[string][]fdb.Process{}, viewFns: map[string]func(fdb.Process) bool{}}
	clusterStatsContent := &ClusterStatsTableContent{cd: v.cd}
	clusterHealthContent := &ClusterHealthTableContent{cd: v.cd}

	pages := tview.NewPages()
	pages.SetBorderPadding(0, 0, 1, 1)

	allView := v.cd.View("all", All)

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
		pv:      v.cd.View("storage", RoleMatch("storage")),
		columns: []ColumnId{ColumnIPAddressPort, ColumnCPUActivity, ColumnRAMUsage, ColumnDiskUsage, ColumnDiskActivity, ColumnKVStorage, ColumnDurabilityRate, ColumnStorageLag, ColumnTotalQueries},
	})
	storage.SetFixed(1, 0)
	storage.SetSelectable(true, false)
	pages.AddPage("2", storage, true, false)

	logs := tview.NewTable().SetContent(&ProcessTableContent{
		pv:      v.cd.View("log", RoleMatch("log")),
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

	clusterInfoFlex := tview.NewFlex()
	clusterInfoFlex.SetDirection(tview.FlexRow)
	clusterInfoFlex.SetBorderPadding(0, 0, 1, 1)
	clusterInfoFlex.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Cluster Info").SetTextColor(tcell.ColorAqua), 1, 1, false)

	clusterHealthFlex := tview.NewFlex()
	clusterHealthFlex.SetDirection(tview.FlexRow)
	clusterHealthFlex.SetBorderPadding(0, 0, 1, 1)
	clusterHealthFlex.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Cluster Health").SetTextColor(tcell.ColorAqua), 1, 1, false)
	clusterHealthFlex.AddItem(tview.NewTable().SetContent(clusterHealthContent).SetSelectable(false, false), 0, 1, false)

	clusterWorkloadFlex := tview.NewFlex()
	clusterWorkloadFlex.SetDirection(tview.FlexRow)
	clusterWorkloadFlex.SetBorderPadding(0, 0, 1, 1)
	clusterWorkloadFlex.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Cluster Workload").SetTextColor(tcell.ColorAqua), 1, 1, false)
	clusterWorkloadFlex.AddItem(tview.NewTable().SetContent(clusterStatsContent).SetSelectable(false, false), 0, 1, false)

	grid := tview.NewGrid().SetRows(5, 1, 0, 1).SetColumns(0, 0, 0).SetBorders(true)
	grid.AddItem(clusterHealthFlex, 0, 0, 1, 2, 0, 0, false)
	grid.AddItem(clusterWorkloadFlex, 0, 2, 1, 1, 0, 0, false)
	grid.AddItem(info, 1, 0, 1, 3, 0, 0, false)
	grid.AddItem(pages, 2, 0, 1, 3, 0, 0, true)
	grid.AddItem(help, 3, 0, 1, 3, 0, 0, false)

	sortIndex := []SortKey{SortRole, SortIPAddress, SortClass}
	sortNow := 0

	v.cd.Sort(sortIndex[sortNow])

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

			v.cd.Sort(sortIndex[sortNow])
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
