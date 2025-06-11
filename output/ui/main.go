package ui

import (
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/data/fdb"
	"github.com/pwood/fdbexplorer/input"
	"github.com/pwood/fdbexplorer/output/ui/components"
	"github.com/pwood/fdbexplorer/output/ui/data/process"
	"github.com/pwood/fdbexplorer/output/ui/views"
	"github.com/rivo/tview"
	"os"
	"strings"
	"time"
)

type UpdatableViews func(dsu process.Update)

func New(ds input.StatusProvider) *Main {
	main := &Main{ds: ds, upCh: make(chan struct{})}

	if em, ok := ds.(input.ExclusionManager); ok {
		main.em = em
	}

	return main
}

type Main struct {
	ds   input.StatusProvider
	em   input.ExclusionManager
	upCh chan struct{}
	app  *tview.Application

	slideShow *components.SlideShow
	sorter    *process.SortControl

	processStore *process.Store
	updatable    []UpdatableViews
	rawJson      []byte

	statusText *tview.TextView
	interval   *views.IntervalControl
}

const (
	StatusInProgress = tcell.ColorYellow
	StatusSuccess    = tcell.ColorGreen
	StatusFailure    = tcell.ColorRed
)

func (m *Main) updateStatus(message string, colour tcell.Color) {
	go m.app.QueueUpdateDraw(func() {
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

	newProcesses := map[string]fdb.Process{}

	for id, p := range root.Cluster.Processes {
		newAddress, tls := strings.CutSuffix(p.Address, ":tls")
		p.Address = newAddress
		p.TLS = tls
		newProcesses[id] = p
	}

	root.Cluster.Processes = newProcesses

	u := process.Update{
		Root: root,
	}

	if em, ok := m.ds.(input.ExclusionManager); ok {
		if excludedProcesses, err := em.ExcludedProcesses(); err != nil {
			m.updateStatus(fmt.Sprintf("Failed to query excluded processes data source: %s", err.Error()), StatusFailure)
			return
		} else {
			u.ExcludedProcesses = excludedProcesses
		}

		if exclusionInProgress, err := em.ExclusionInProgressProcesses(); err != nil {
			m.updateStatus(fmt.Sprintf("Failed to query exclusion in progress data source: %s", err.Error()), StatusFailure)
			return
		} else {
			u.ExclusionInProgress = exclusionInProgress
		}
	}

	m.rawJson = d
	duration := time.Since(start)

	msg := fmt.Sprintf("Updated in %dms, next in %s.", duration.Milliseconds(), m.interval.Duration().String())
	m.updateStatus(msg, StatusSuccess)

	m.app.QueueUpdateDraw(func() {
		for _, uFn := range m.updatable {
			uFn(u)
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
	m.interval = &views.IntervalControl{}

	m.sorter = &process.SortControl{}
	m.processStore = process.NewStore(m.sorter.Sort)

	localityDataContent := components.NewDataTable[process.Process](
		[]components.ColumnDef[process.Process]{views.ColumnSelected, views.ColumnIPAddressPort, views.ColumnTLS, views.ColumnStatus, views.ColumnMachine, views.ColumnLocality, views.ColumnClass, views.ColumnRoles, views.ColumnVersion, views.ColumnUptime})

	m.processStore.AddNotifiable(localityDataContent.Update, views.All)

	usageDataContent := components.NewDataTable[process.Process](
		[]components.ColumnDef[process.Process]{views.ColumnSelected, views.ColumnIPAddressPort, views.ColumnRoles, views.ColumnCPUActivity, views.ColumnRAMUsage, views.ColumnNetworkActivity, views.ColumnDiskUsage, views.ColumnDiskActivity})

	m.processStore.AddNotifiable(usageDataContent.Update, views.All)

	storageDataContent := components.NewDataTable[process.Process](
		[]components.ColumnDef[process.Process]{views.ColumnSelected, views.ColumnIPAddressPort, views.ColumnCPUActivity, views.ColumnRAMUsage, views.ColumnDiskUsage, views.ColumnDiskActivity, views.ColumnKVStorage, views.ColumnStorageDurabilityRate, views.ColumnStorageLag, views.ColumnStorageTotalQueries})

	m.processStore.AddNotifiable(storageDataContent.Update, views.RoleMatch("storage"))

	logDataContent := components.NewDataTable[process.Process](
		[]components.ColumnDef[process.Process]{views.ColumnSelected, views.ColumnIPAddressPort, views.ColumnCPUActivity, views.ColumnRAMUsage, views.ColumnDiskUsage, views.ColumnDiskActivity, views.ColumnLogQueueLength, views.ColumnLogDurabilityRate, views.ColumnLogQueueStorage})

	m.processStore.AddNotifiable(logDataContent.Update, views.RoleMatch("log"))

	clusterHealthContent := components.NewStatsGrid([][]components.ColumnDef[views.ClusterHealth]{
		{views.StatClusterHealth, views.StatRebalanceQueued},
		{views.StatReplicasRemaining, views.StatRebalanceInflight},
		{views.StatRecoveryState, views.StatEmpty},
		{views.StatRecoveryDescription, views.StatDatabaseLocked},
	})

	clusterStatsContent := components.NewStatsGrid([][]components.ColumnDef[views.ClusterStats]{
		{views.StatTxStarted, views.StatReads},
		{views.StatTxCommitted, views.StatWrites},
		{views.StatTxConflicted, views.StatBytesRead},
		{views.StatTxRejected, views.StatBytesWritten},
	})

	backupInstancesContent := components.NewDataTable[fdb.BackupInstance](
		[]components.ColumnDef[fdb.BackupInstance]{views.ColumnBackupInstanceId, views.ColumnBackupInstanceVersion, views.ColumnBackupInstanceConfiguredWorkers, views.ColumnBackupInstanceUsedMemory, views.ColumnBackupInstanceRecentTransfer, views.ColumnBackupInstanceRecentOperations})

	backupTagsContent := components.NewDataTable[fdb.BackupTag](
		[]components.ColumnDef[fdb.BackupTag]{views.ColumnBackupTagId, views.ColumnBackupStatus, views.ColumnBackupRunning, views.ColumnBackupRestorable, views.ColumnBackupSecondsBehind, views.ColumnBackupRestorableVersion, views.ColumnBackupRangeBytes, views.ColumnBackupLogBytes})

	drBackupInstanceColumns := []components.ColumnDef[fdb.DRBackupInstance]{views.ColumnDRBackupInstanceId, views.ColumnDRBackupInstanceLastUpdated, views.ColumnDRBackupInstanceProcessCPU, views.ColumnDRBackupInstanceMemoryUsage, views.ColumnDRBackupInstanceResidentSize, views.ColumnDRBackupInstanceVersion}
	drBackupTagColumns := []components.ColumnDef[fdb.DRBackupTag]{views.ColumnDRBackupTagId, views.ColumnDRBackupTagRunning, views.ColumnDRBackupTagRestorable, views.ColumnDRBackupTagSecondsBehind, views.ColumnDRBackupTagState, views.ColumnDRBackupTagRangeBytes, views.ColumnDRBackupTagLogBytes, views.ColumnDRBackupTagMutationStream}

	drBackupInstancesContent := components.NewDataTable[fdb.DRBackupInstance](drBackupInstanceColumns)
	drBackupTagsContent := components.NewDataTable[fdb.DRBackupTag](drBackupTagColumns)
	drBackupDestInstancesContent := components.NewDataTable[fdb.DRBackupInstance](drBackupInstanceColumns)
	drBackupDestTagsContent := components.NewDataTable[fdb.DRBackupTag](drBackupTagColumns)

	m.updatable = []UpdatableViews{
		m.processStore.Update,
		views.UpdateClusterHealth(clusterHealthContent.Update),
		views.UpdateClusterStats(clusterStatsContent.Update),
		views.UpdateBackupInstances(backupInstancesContent.Update),
		views.UpdateBackupTags(backupTagsContent.Update),
		views.UpdateDrBackupInstances(drBackupInstancesContent.Update),
		views.UpdateDrBackupTags(drBackupTagsContent.Update),
		views.UpdateDrBackupDestInstances(drBackupDestInstancesContent.Update),
		views.UpdateDrBackupDestTags(drBackupDestTagsContent.Update),
	}

	processDataInput := func(table *tview.Table, content *components.DataTable[process.Process]) func(event *tcell.EventKey) *tcell.EventKey {
		return func(event *tcell.EventKey) *tcell.EventKey {
			switch event.Key() {
			case tcell.KeyRune:
				switch event.Rune() {
				case ' ':
					row, _ := table.GetSelection()
					content.Get(row).Metadata.ToggleSelected()
					m.processStore.Sort()
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

	drBackupInstances := tview.NewTable().SetContent(drBackupInstancesContent).SetFixed(1, 0).SetSelectable(false, false)
	drBackupTags := tview.NewTable().SetContent(drBackupTagsContent).SetFixed(1, 0).SetSelectable(false, false)
	drBackupDestInstances := tview.NewTable().SetContent(drBackupDestInstancesContent).SetFixed(1, 0).SetSelectable(false, false)
	drBackupDestTags := tview.NewTable().SetContent(drBackupDestTagsContent).SetFixed(1, 0).SetSelectable(false, false)

	drFlex := tview.NewFlex()
	drFlex.SetDirection(tview.FlexRow)
	drFlex.AddItem(drBackupInstances, 0, 2, false)
	drFlex.AddItem(drBackupTags, 0, 1, false)
	drFlex.SetBorder(true)
	drFlex.SetTitle("'Source' (Local) Cluster")
	drFlex.SetTitleColor(tcell.ColorAqua)

	drDestFelx := tview.NewFlex()
	drDestFelx.SetDirection(tview.FlexRow)
	drDestFelx.AddItem(drBackupDestInstances, 0, 2, false)
	drDestFelx.AddItem(drBackupDestTags, 0, 1, false)
	drDestFelx.SetBorder(true)
	drDestFelx.SetTitle("'Destination' (Remote) Cluster")
	drDestFelx.SetTitleColor(tcell.ColorAqua)

	drContainer := tview.NewFlex()
	drContainer.SetDirection(tview.FlexRow)
	drContainer.AddItem(drFlex, 0, 1, false)
	drContainer.AddItem(drDestFelx, 0, 1, false)

	m.slideShow = components.NewSlideShow()
	m.slideShow.Add("Locality", locality)
	m.slideShow.Add("Usage Overview", usage)
	m.slideShow.Add("Storage Processes", storage)
	m.slideShow.Add("Log Processes", logs)
	m.slideShow.Add("Backups", backupFlex)
	m.slideShow.Add("DR Backups", drContainer)

	m.statusText = tview.NewTextView()
	m.statusText.SetTextAlign(tview.AlignRight)
	m.statusText.SetText("")

	bottom := tview.NewFlex()
	bottom.SetBorderPadding(0, 0, 1, 1)
	bottom.AddItem(tview.NewTable().SetContent(&views.HelpKeys{Sorter: m.sorter, Interval: m.interval, HasEM: m.em != nil}).SetSelectable(false, false), 0, 1, false)
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
	grid.AddItem(m.slideShow, 1, 0, 1, 3, 0, 0, true)
	grid.AddItem(bottom, 2, 0, 1, 3, 0, 0, false)

	grid.SetInputCapture(m.rootAction)

	m.app = tview.NewApplication().SetRoot(grid, true).SetFocus(locality)

	go m.runData()

	if err := m.app.Run(); err != nil {
		panic(err)
	}
}
