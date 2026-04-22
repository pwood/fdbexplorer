package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/data/fdb"
	"github.com/pwood/fdbexplorer/input"
	"github.com/pwood/fdbexplorer/output/ui/components"
	"github.com/pwood/fdbexplorer/output/ui/data/process"
	"github.com/pwood/fdbexplorer/output/ui/panels"
	"github.com/pwood/fdbexplorer/output/ui/views"
	"github.com/rivo/tview"
)

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
	panels       []panels.Panel
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
		m.processStore.Update(u)
		for _, p := range m.panels {
			p.Update(u)
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

	locality := panels.NewLocality(m.processStore)
	usage := panels.NewUsage(m.processStore)
	storage := panels.NewStorage(m.processStore)
	logs := panels.NewLogs(m.processStore)
	backups := panels.NewBackups()
	drBackups := panels.NewDRBackups()
	clusterHealth := panels.NewClusterHealth()
	clusterWorkload := panels.NewClusterWorkload()

	m.panels = []panels.Panel{backups, drBackups, clusterHealth, clusterWorkload}

	m.slideShow = components.NewSlideShow()
	m.slideShow.Add("Locality", locality.Root())
	m.slideShow.Add("Usage Overview", usage.Root())
	m.slideShow.Add("Storage Processes", storage.Root())
	m.slideShow.Add("Log Processes", logs.Root())
	m.slideShow.Add("Backups", backups.Root())
	m.slideShow.Add("DR Backups", drBackups.Root())

	m.statusText = tview.NewTextView()
	m.statusText.SetTextAlign(tview.AlignRight)
	m.statusText.SetText("")

	bottom := tview.NewFlex()
	bottom.SetBorderPadding(0, 0, 1, 1)
	bottom.AddItem(tview.NewTable().SetContent(&views.HelpKeys{Sorter: m.sorter, Interval: m.interval, HasEM: m.em != nil}).SetSelectable(false, false), 0, 1, false)
	bottom.AddItem(m.statusText, 0, 1, false)

	grid := tview.NewGrid().SetRows(5, 0, 1).SetColumns(0, 0, 0).SetBorders(true)
	grid.AddItem(clusterHealth.Root(), 0, 0, 1, 2, 0, 0, false)
	grid.AddItem(clusterWorkload.Root(), 0, 2, 1, 1, 0, 0, false)
	grid.AddItem(m.slideShow, 1, 0, 1, 3, 0, 0, true)
	grid.AddItem(bottom, 2, 0, 1, 3, 0, 0, false)

	grid.SetInputCapture(m.rootAction)

	m.app = tview.NewApplication().SetRoot(grid, true).SetFocus(locality.Root())

	go m.runData()

	if err := m.app.Run(); err != nil {
		panic(err)
	}
}
