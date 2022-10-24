package ui

import (
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/data"
	"github.com/pwood/fdbexplorer/data/fdb"
	"github.com/pwood/fdbexplorer/ui/components"
	"github.com/rivo/tview"
	"strings"
	"time"
)

type UpdatableViews func(root fdb.Root)

func New(ch chan data.State) *Main {
	return &Main{ch: ch}
}

type Main struct {
	ch  chan data.State
	app *tview.Application

	updatable []UpdatableViews
	rawJson   []byte

	statusText *tview.TextView
}

func (m *Main) runData() {
	for s := range m.ch {
		var root fdb.Root

		if s.Err == nil {
			s.Err = json.Unmarshal(s.Data, &root)
		}

		text := []string{"[", time.Now().Format("15:04:05"), "] "}

		if s.Err != nil {
			text = append(text, s.Err.Error())
			m.statusText.SetTextColor(tcell.ColorRed)
		} else {
			text = append(text, fmt.Sprintf("Updated in %dms", s.Duration.Milliseconds()))
			m.statusText.SetTextColor(tcell.ColorGreen)
		}

		if s.Interval != 0 {
			text = append(text, fmt.Sprintf(", next in %s.", s.Interval.String()))
		}

		m.statusText.SetText(strings.Join(text, ""))

		if s.Err != nil {
			continue
		}

		m.rawJson = s.Data

		for _, updateFn := range m.updatable {
			updateFn(root)
		}

		m.app.Draw()
	}
}

func (m *Main) Run() {
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

	slideShow := components.NewSlideShow()
	slideShow.Add("Locality", locality)
	slideShow.Add("Usage Overview", usage)
	slideShow.Add("Storage Processes", storage)
	slideShow.Add("Log Processes", logs)

	m.statusText = tview.NewTextView()
	m.statusText.SetTextAlign(tview.AlignRight)
	m.statusText.SetText("")

	bottom := tview.NewFlex()
	bottom.SetBorderPadding(0, 0, 1, 1)
	bottom.AddItem(tview.NewTable().SetContent(&HelpKeys{sorter: sorter}).SetSelectable(false, false), 0, 1, false)
	bottom.AddItem(m.statusText, 0, 1, false)

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

	grid := tview.NewGrid().SetRows(5, 0, 1).SetColumns(0, 0, 0).SetBorders(true)
	grid.AddItem(clusterHealthFlex, 0, 0, 1, 2, 0, 0, false)
	grid.AddItem(clusterWorkloadFlex, 0, 2, 1, 1, 0, 0, false)
	grid.AddItem(slideShow, 1, 0, 1, 3, 0, 0, true)
	grid.AddItem(bottom, 2, 0, 1, 3, 0, 0, false)

	grid.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			slideShow.Prev()
		case tcell.KeyRight:
			slideShow.Next()
		case tcell.KeyF1:
			sorter.Next()
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