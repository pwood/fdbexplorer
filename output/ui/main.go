package ui

import (
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/data/fdb"
	"github.com/pwood/fdbexplorer/input"
	"github.com/pwood/fdbexplorer/output/ui/components"
	"github.com/pwood/fdbexplorer/output/ui/data"
	"github.com/pwood/fdbexplorer/output/ui/views"
	"github.com/rivo/tview"
	"os"
	"strings"
	"time"
)

type UpdatableViews func(dsu data.DataSourceUpdate)

func New(ds input.StatusProvider) *Main {
	return &Main{ds: ds, upCh: make(chan struct{})}
}

type Main struct {
	ds   input.StatusProvider
	upCh chan struct{}
	app  *tview.Application

	metadataStore *data.MetadataStore
	updatable     []UpdatableViews
	rawJson       []byte

	statusText *tview.TextView
	interval   *IntervalControl
}

const (
	StatusInProgress = tcell.ColorYellow
	StatusSuccess    = tcell.ColorGreen
	StatusFailure    = tcell.ColorRed
)

func (m *Main) updateStatus(message string, colour tcell.Color) {
	m.app.QueueUpdate(func() {
		text := []string{"[", time.Now().Format("15:04:05"), "] ", message}
		m.statusText.SetText(strings.Join(text, "")).SetTextColor(colour)
	})
}

func (m *Main) runData() {
	m.updateFromDS()

	for {
		select {
		case <-time.After(m.interval.Duration()):
		case <-m.upCh:
		}

		m.updateFromDS()
	}
}

func (m *Main) updateFromDS() {
	m.updateStatus("Updating data...", StatusInProgress)
	start := time.Now()

	d, err := m.ds.Status()
	if err != nil {
		m.updateStatus(fmt.Sprintf("Failed to query Root data source: %s", err.Error()), StatusFailure)
		return
	}

	var root fdb.Root
	if err := json.Unmarshal(d, &root); err != nil {
		m.updateStatus(fmt.Sprintf("Failed to unmarshal data: %s", err.Error()), StatusFailure)
		return
	}

	dsu := data.DataSourceUpdate{
		Root: root,
	}

	if em, ok := m.ds.(input.ExclusionManager); ok {
		if excludedProcesses, err := em.ExcludedProcesses(); err != nil {
			m.updateStatus(fmt.Sprintf("Failed to query excluded processes data source: %s", err.Error()), StatusFailure)
			return
		} else {
			dsu.ExcludedProcesses = excludedProcesses
		}

		if exclusionInProgress, err := em.ExclusionInProgressProcesses(); err != nil {
			m.updateStatus(fmt.Sprintf("Failed to query exclusion in progress data source: %s", err.Error()), StatusFailure)
			return
		} else {
			dsu.ExclusionInProgress = exclusionInProgress
		}
	}

	m.rawJson = d
	duration := time.Since(start)

	msg := fmt.Sprintf("Updated in %dms, next in %s.", duration.Milliseconds(), m.interval.Duration().String())
	m.updateStatus(msg, StatusSuccess)

	m.app.QueueUpdateDraw(func() {
		for _, updateFn := range m.updatable {
			updateFn(dsu)
		}
	})
}

func (m *Main) snapshotData() (string, error) {
	fileName := fmt.Sprintf("fdbexplorer-status-snapshot-%d.json", time.Now().Unix())

	f, err := os.Create(fileName)
	defer func() {
		_ = f.Close()
	}()

	if err != nil {
		return "", fmt.Errorf("open: %w", err)
	}

	if n, err := f.Write(m.rawJson); err != nil {
		return "", fmt.Errorf("write: %w", err)
	} else if n != len(m.rawJson) {
		return "", fmt.Errorf("write: only %d of %d bytes written", n, len(m.rawJson))
	}

	return fileName, nil
}

func (m *Main) Run() {
	sorter := &SortControl{}
	m.interval = &IntervalControl{}

	_, haveEM := m.ds.(input.ExclusionManager)

	localityDataContent := components.NewDataTable[views.ProcessData](
		[]components.ColumnDef[views.ProcessData]{views.ColumnSelected, views.ColumnIPAddressPort, views.ColumnStatus, views.ColumnMachine, views.ColumnLocality, views.ColumnClass, views.ColumnRoles, views.ColumnVersion, views.ColumnUptime},
		views.All,
		sorter.Sort)

	usageDataContent := components.NewDataTable[views.ProcessData](
		[]components.ColumnDef[views.ProcessData]{views.ColumnSelected, views.ColumnIPAddressPort, views.ColumnRoles, views.ColumnCPUActivity, views.ColumnRAMUsage, views.ColumnNetworkActivity, views.ColumnDiskUsage, views.ColumnDiskActivity},
		views.All,
		sorter.Sort)

	storageDataContent := components.NewDataTable[views.ProcessData](
		[]components.ColumnDef[views.ProcessData]{views.ColumnSelected, views.ColumnIPAddressPort, views.ColumnCPUActivity, views.ColumnRAMUsage, views.ColumnDiskUsage, views.ColumnDiskActivity, views.ColumnKVStorage, views.ColumnStorageDurabilityRate, views.ColumnStorageLag, views.ColumnStorageTotalQueries},
		views.RoleMatch("storage"),
		sorter.Sort)

	logDataContent := components.NewDataTable[views.ProcessData](
		[]components.ColumnDef[views.ProcessData]{views.ColumnSelected, views.ColumnIPAddressPort, views.ColumnCPUActivity, views.ColumnRAMUsage, views.ColumnDiskUsage, views.ColumnDiskActivity, views.ColumnLogQueueLength, views.ColumnLogDurabilityRate, views.ColumnLogQueueStorage},
		views.RoleMatch("log"),
		sorter.Sort)

	clusterHealthContent := components.NewStatsGrid([][]components.ColumnDef[views.ClusterHealth]{
		{views.StatClusterHealth, views.StatRebalanceQueued},
		{views.StatReplicasRemaining, views.StatRebalanceInflight},
		{views.StatRecoveryState, views.StatEmpty},
		{views.StatRecoveryDescription, views.StatEmpty},
	})

	clusterStatsContent := components.NewStatsGrid([][]components.ColumnDef[views.ClusterStats]{
		{views.StatTxStarted, views.StatReads},
		{views.StatTxCommitted, views.StatWrites},
		{views.StatTxConflicted, views.StatBytesRead},
		{views.StatTxRejected, views.StatBytesWritten},
	})

	backupInstancesContent := components.NewDataTable[fdb.BackupInstance](
		[]components.ColumnDef[fdb.BackupInstance]{views.ColumnBackupInstanceId, views.ColumnBackupInstanceVersion, views.ColumnBackupInstanceConfiguredWorkers, views.ColumnBackupInstanceUsedMemory, views.ColumnBackupInstanceRecentTransfer, views.ColumnBackupInstanceRecentOperations},
		func(_ fdb.BackupInstance) bool { return true },
		func(i fdb.BackupInstance, j fdb.BackupInstance) bool { return strings.Compare(i.Id, j.Id) < 0 })

	backupTagsContent := components.NewDataTable[fdb.BackupTag](
		[]components.ColumnDef[fdb.BackupTag]{views.ColumnBackupTagId, views.ColumnBackupStatus, views.ColumnBackupRunning, views.ColumnBackupRestorable, views.ColumnBackupSecondsBehind, views.ColumnBackupRestorableVersion, views.ColumnBackupRangeBytes, views.ColumnBackupLogBytes},
		func(_ fdb.BackupTag) bool { return true },
		func(i fdb.BackupTag, j fdb.BackupTag) bool { return strings.Compare(i.Id, j.Id) < 0 })

	m.metadataStore = &data.MetadataStore{Metadata: map[string]*views.ProcessMetadata{}}

	m.updatable = []UpdatableViews{
		m.metadataStore.Update(localityDataContent.Update),
		m.metadataStore.Update(usageDataContent.Update),
		m.metadataStore.Update(storageDataContent.Update),
		m.metadataStore.Update(logDataContent.Update),
		views.UpdateClusterHealth(clusterHealthContent.Update),
		views.UpdateClusterStats(clusterStatsContent.Update),
		views.UpdateBackupInstances(backupInstancesContent.Update),
		views.UpdateBackupTags(backupTagsContent.Update),
	}

	sortAll := func() {
		localityDataContent.Sort()
		usageDataContent.Sort()
		storageDataContent.Sort()
		logDataContent.Sort()
	}

	processDataInput := func(table *tview.Table, content *components.DataTable[views.ProcessData]) func(event *tcell.EventKey) *tcell.EventKey {
		return func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyRune:
				switch event.Rune() {
				case ' ':
					row, _ := table.GetSelection()
					content.Get(row).Metadata.ToggleSelected()
					sortAll()
					return nil
				}
			}

			return event
		}
	}

	locality := tview.NewTable().SetContent(localityDataContent).SetFixed(1, 0).SetSelectable(true, false)
	locality.SetInputCapture(processDataInput(locality, localityDataContent))
	usage := tview.NewTable().SetContent(usageDataContent).SetFixed(1, 0).SetSelectable(true, false)
	usage.SetInputCapture(processDataInput(usage, usageDataContent))
	storage := tview.NewTable().SetContent(storageDataContent).SetFixed(1, 0).SetSelectable(true, false)
	storage.SetInputCapture(processDataInput(storage, storageDataContent))
	logs := tview.NewTable().SetContent(logDataContent).SetFixed(1, 0).SetSelectable(true, false)
	logs.SetInputCapture(processDataInput(logs, logDataContent))

	backupInstances := tview.NewTable().SetContent(backupInstancesContent).SetFixed(1, 0).SetSelectable(false, false)
	backupTags := tview.NewTable().SetContent(backupTagsContent).SetFixed(1, 0).SetSelectable(false, false)

	backupFlex := tview.NewFlex()
	backupFlex.SetDirection(tview.FlexRow)
	backupFlex.AddItem(backupInstances, 0, 1, false)
	backupFlex.AddItem(backupTags, 0, 1, false)

	slideShow := components.NewSlideShow()
	slideShow.Add("Locality", locality)
	slideShow.Add("Usage Overview", usage)
	slideShow.Add("Storage Processes", storage)
	slideShow.Add("Log Processes", logs)
	slideShow.Add("Backups", backupFlex)

	m.statusText = tview.NewTextView()
	m.statusText.SetTextAlign(tview.AlignRight)
	m.statusText.SetText("")

	bottom := tview.NewFlex()
	bottom.SetBorderPadding(0, 0, 1, 1)
	bottom.AddItem(tview.NewTable().SetContent(&HelpKeys{sorter: sorter, interval: m.interval, haveEM: haveEM}).SetSelectable(false, false), 0, 1, false)
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
			sortAll()
		case tcell.KeyF2:
			if filename, err := m.snapshotData(); err != nil {
				m.updateStatus(fmt.Sprintf("Failed to write snapshot: %s", err.Error()), StatusFailure)
			} else {
				m.updateStatus(fmt.Sprintf("Snapshot written: %s", filename), StatusSuccess)
			}
		case tcell.KeyF3:
			m.interval.Next()
		case tcell.KeyF5:
			m.upCh <- struct{}{}
		case tcell.KeyF7:
			if haveEM {

			}
		case tcell.KeyF8:
			if haveEM {

			}
		case tcell.KeyESC:
			m.app.Stop()
		case tcell.KeyCtrlL:
			go m.app.Draw()
		case tcell.KeyRune:
			switch event.Rune() {
			case '\\':
				m.metadataStore.ClearSelected()
				sortAll()
			default:
				return event
			}
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
