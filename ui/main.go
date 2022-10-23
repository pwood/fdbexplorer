package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/data"
	"github.com/pwood/fdbexplorer/data/fdb"
	"github.com/pwood/fdbexplorer/ui/components"
	"github.com/rivo/tview"
	"strconv"
)

type UpdatableViews func(root fdb.Root)

func New(ch chan data.State) *Main {
	return &Main{ch: ch}
}

type Main struct {
	ch  chan data.State
	app *tview.Application

	updatable []UpdatableViews
}

func (m *Main) runData() {
	for s := range m.ch {
		for _, updateFn := range m.updatable {
			updateFn(s.ClusterState)
		}

		m.app.Draw()
	}
}

func (m *Main) Run() {
	pages := tview.NewPages()
	pages.SetBorderPadding(0, 0, 1, 1)

	sorter := &ProcessSorter{}

	localityDataContent := components.NewDataTable[fdb.Process](
		[]components.ColumnDef[fdb.Process]{ColumnIPAddressPort, ColumnStatus, ColumnMachine, ColumnLocality, ColumnClass, ColumnRoles, ColumnVersion, ColumnUptime},
		All,
		sorter.Sort)

	usageDataContent := components.NewDataTable[fdb.Process](
		[]components.ColumnDef[fdb.Process]{ColumnIPAddressPort, ColumnRoles, ColumnCPUActivity, ColumnRAMUsage, ColumnNetworkActivity, ColumnDiskUsage, ColumnDiskActivity},
		All,
		sorter.Sort)

	storageDataContent := components.NewDataTable[fdb.Process](
		[]components.ColumnDef[fdb.Process]{ColumnIPAddressPort, ColumnCPUActivity, ColumnRAMUsage, ColumnDiskUsage, ColumnDiskActivity, ColumnKVStorage, ColumnDurabilityRate, ColumnStorageLag, ColumnTotalQueries},
		RoleMatch("storage"),
		sorter.Sort)

	logDataContent := components.NewDataTable[fdb.Process](
		[]components.ColumnDef[fdb.Process]{ColumnIPAddressPort, ColumnCPUActivity, ColumnRAMUsage, ColumnDiskUsage, ColumnDiskActivity, ColumnQueueStorage, ColumnDurabilityRate},
		RoleMatch("log"),
		sorter.Sort)

	clusterHealthContent := components.NewStatsGrid([][]components.ColumnDef[ClusterHealth]{
		{StatClusterHealth, StatRebalanceQueued},
		{StatReplicasRemaining, StatRebalanceInflight},
		{StatRecoveryState, StatEmpty},
		{StatRecoveryDescription, StatEmpty},
	})

	clusterStatsContent := components.NewStatsGrid([][]components.ColumnDef[ClusterStats]{
		{StatTxStarted, StatReads},
		{StatTxCommitted, StatWrites},
		{StatTxConflicted, StatBytesRead},
		{StatTxRejected, StatBytesWritten},
	})

	m.updatable = []UpdatableViews{
		UpdateProcesses(localityDataContent.Update),
		UpdateProcesses(usageDataContent.Update),
		UpdateProcesses(storageDataContent.Update),
		UpdateProcesses(logDataContent.Update),
		UpdateProcessClusterHealth(clusterHealthContent.Update),
		UpdateProcessClusterStats(clusterStatsContent.Update),
	}

	locality := tview.NewTable().SetContent(localityDataContent).SetFixed(1, 0).SetSelectable(true, false)
	usage := tview.NewTable().SetContent(usageDataContent).SetFixed(1, 0).SetSelectable(true, false)
	storage := tview.NewTable().SetContent(storageDataContent).SetFixed(1, 0).SetSelectable(true, false)
	logs := tview.NewTable().SetContent(logDataContent).SetFixed(1, 0).SetSelectable(true, false)

	pages.AddPage("0", locality, true, true)
	pages.AddPage("1", usage, true, false)
	pages.AddPage("2", storage, true, false)
	pages.AddPage("3", logs, true, false)

	pageIndex := []string{"Locality", "Usage Overview", "Storage Processes", "Log Processes"}

	info := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetTextAlign(tview.AlignCenter).
		SetHighlightedFunc(func(added, _, _ []string) {
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

	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			previousSlide()
		case tcell.KeyRight:
			nextSlide()
		case tcell.KeyF1:
			sorter.NextSort()
			localityDataContent.Sort()
			usageDataContent.Sort()
			storageDataContent.Sort()
			logDataContent.Sort()
		case tcell.KeyESC:
			m.app.Stop()
		default:
			return event
		}
		return nil
	})

	m.app = tview.NewApplication().SetRoot(grid, true).SetFocus(locality)

	go m.runData()

	if err := m.app.Run(); err != nil {
		panic(err)
	}
}
