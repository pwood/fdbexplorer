package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/statusjson"
	"github.com/rivo/tview"
	"sort"
	"strings"
	"sync"
	"time"
)

type SortKey int

const (
	SortIPAddress SortKey = iota
	SortClass
	SortRole
)

type ClusterData struct {
	root statusjson.Root

	processes []statusjson.Process
	views     map[string][]statusjson.Process
	viewFns   map[string]func(statusjson.Process) bool

	sortBy SortKey

	m *sync.RWMutex
}

func (c *ClusterData) Sort(sortKey SortKey) {
	c.m.Lock()
	defer c.m.Unlock()

	c.sortBy = sortKey
	c._filterAndSort()
}

func (c *ClusterData) Update(s statusjson.Root) {
	c.m.Lock()
	defer c.m.Unlock()

	c.root = s
	c.processes = nil

	for _, process := range c.root.Cluster.Processes {
		c.processes = append(c.processes, process)
	}

	c._filterAndSort()
}

func (c *ClusterData) _filterAndSort() {
	switch c.sortBy {
	case SortIPAddress:
		sort.Slice(c.processes, func(i, j int) bool {
			return strings.Compare(c.processes[i].Address, c.processes[j].Address) < 0
		})
	case SortClass:
		sort.Slice(c.processes, func(i, j int) bool {
			return strings.Compare(c.processes[i].Class+c.processes[i].Address, c.processes[j].Class+c.processes[j].Address) < 0
		})
	case SortRole:
		sort.Slice(c.processes, func(i, j int) bool {
			iRole := ""
			if len(c.processes[i].Roles) > 0 {
				iRole = c.processes[i].Roles[0].Role
			}

			jRole := ""
			if len(c.processes[j].Roles) > 0 {
				jRole = c.processes[j].Roles[0].Role
			}

			return strings.Compare(iRole+c.processes[i].Address, jRole+c.processes[j].Address) < 0
		})
	}

	for i, proc := range c.processes {
		c.processes[i] = statusjson.CalculateProcessHealth(proc)
	}

	for n, fn := range c.viewFns {
		var viewProcesses []statusjson.Process

		for _, p := range c.processes {
			if fn(p) {
				viewProcesses = append(viewProcesses, p)
			}
		}

		c.views[n] = viewProcesses
	}
}

func All(_ statusjson.Process) bool {
	return true
}

func RoleMatch(s string) func(statusjson.Process) bool {
	return func(process statusjson.Process) bool {
		for _, r := range process.Roles {
			if r.Role == s {
				return true
			}
		}
		return false
	}
}

func (c *ClusterData) View(name string, fn func(statusjson.Process) bool) *ProcessView {
	c.views[name] = nil
	c.viewFns[name] = fn

	return &ProcessView{
		pd: c,
		n:  name,
	}
}

type ClusterStats struct {
	TxStarted    float64
	TxCommitted  float64
	TxConflicted float64
	TxRejected   float64

	Reads        float64
	Writes       float64
	BytesRead    float64
	BytesWritten float64
}

func (c *ClusterData) Stats() ClusterStats {
	c.m.RLock()
	defer c.m.RUnlock()

	return ClusterStats{
		TxStarted:    c.root.Cluster.Workload.Transactions.Started.Hz,
		TxCommitted:  c.root.Cluster.Workload.Transactions.Committed.Hz,
		TxConflicted: c.root.Cluster.Workload.Transactions.Conflicted.Hz,
		TxRejected:   c.root.Cluster.Workload.Transactions.RejectedForQueuedTooLong.Hz,
		Reads:        c.root.Cluster.Workload.Operations.Reads.Hz,
		Writes:       c.root.Cluster.Workload.Operations.Writes.Hz,
		BytesRead:    c.root.Cluster.Workload.Bytes.Read.Hz,
		BytesWritten: c.root.Cluster.Workload.Bytes.Written.Hz,
	}
}

type ClusterHealth struct {
	Healthy     bool
	Health      string
	MinReplicas int

	RebalanceInFlight int
	RebalanceQueued   int

	RecoveryState       string
	RecoveryDescription string
}

func (c *ClusterData) Health() ClusterHealth {
	c.m.RLock()
	defer c.m.RUnlock()

	return ClusterHealth{
		Healthy:             c.root.Cluster.Data.State.Health,
		Health:              strings.Title(strings.Replace(c.root.Cluster.Data.State.Name, "_", " ", -1)),
		MinReplicas:         c.root.Cluster.Data.State.MinReplicasRemaining,
		RebalanceQueued:     c.root.Cluster.Data.MovingData.InQueueBytes,
		RebalanceInFlight:   c.root.Cluster.Data.MovingData.InFlightBytes,
		RecoveryState:       strings.Title(strings.Replace(c.root.Cluster.RecoveryState.Name, "_", " ", -1)),
		RecoveryDescription: c.root.Cluster.RecoveryState.Description,
	}
}

type ProcessView struct {
	pd *ClusterData
	n  string
}

func (p *ProcessView) Count() int {
	p.pd.m.RLock()
	defer p.pd.m.RUnlock()

	return len(p.pd.views[p.n])
}

func (p *ProcessView) Get(i int) statusjson.Process {
	p.pd.m.RLock()
	defer p.pd.m.RUnlock()

	return p.pd.views[p.n][i]
}

type ProcessTableContent struct {
	tview.TableContentReadOnly
	pv *ProcessView

	columns []ColumnId
}

func (v *ProcessTableContent) GetCell(row, column int) *tview.TableCell {
	cid := v.columns[column]

	if row == 0 {
		return tview.NewTableCell(columns[cid].Name).SetExpansion(1).SetTextColor(tcell.ColorAqua).SetSelectable(false)
	} else {
		process := v.pv.Get(row - 1)
		color := tcell.ColorWhite

		switch process.Health {
		case statusjson.HealthCritical:
			color = tcell.ColorRed
		case statusjson.HealthWarning:
			color = tcell.ColorYellow
		case statusjson.HealthExcluded:
			color = tcell.ColorBlue
		}

		return tview.NewTableCell(columns[cid].DataFn(process)).SetTextColor(color)
	}
}

func (v *ProcessTableContent) GetRowCount() int {
	return v.pv.Count() + 1
}

func (v *ProcessTableContent) GetColumnCount() int {
	return len(v.columns)
}

type ColumnId int

const (
	ColumnIPAddressPort ColumnId = iota
	ColumnStatus
	ColumnMachine
	ColumnLocality
	ColumnClass
	ColumnRoles
	ColumnCPUActivity
	ColumnRAMUsage
	ColumnDiskUsage
	ColumnDiskActivity
	ColumnNetworkActivity
	ColumnVersion
	ColumnUptime
	ColumnKVStorage
	ColumnQueueStorage
	ColumnDurabilityRate
	ColumnStorageLag
	ColumnTotalQueries
)

type ColumnDef struct {
	Name   string
	DataFn func(statusjson.Process) string
}

var columns = map[ColumnId]ColumnDef{
	ColumnIPAddressPort: {
		Name: "IP Address:Port",
		DataFn: func(process statusjson.Process) string {
			return process.Address
		},
	},
	ColumnStatus: {
		Name: "Status",
		DataFn: func(process statusjson.Process) string {
			var statuses []string

			if process.Excluded {
				statuses = append(statuses, "Excluded")
			}

			if process.Degraded {
				statuses = append(statuses, "Degraded")
			}

			if process.UnderMaintenance {
				statuses = append(statuses, "Maintenance")
			}

			if len(process.Messages) > 0 {
				statuses = append(statuses, "Message")
			}

			return strings.Join(statuses, " / ")
		},
	},
	ColumnMachine: {
		Name: "Machine",
		DataFn: func(process statusjson.Process) string {
			return process.Locality.MachineID
		},
	},
	ColumnLocality: {
		Name: "Locality",
		DataFn: func(process statusjson.Process) string {
			return fmt.Sprintf("%s / %s", process.Locality.DataHall, process.Locality.DCID)
		},
	},
	ColumnClass: {
		Name: "Class",
		DataFn: func(process statusjson.Process) string {
			return process.Class
		},
	},
	ColumnRoles: {
		Name: "Roles",
		DataFn: func(process statusjson.Process) string {
			var roles []string

			for _, role := range process.Roles {
				roles = append(roles, role.Role)
			}

			return strings.Join(roles, ", ")
		},
	},
	ColumnRAMUsage: {
		Name: "RAM Usage",
		DataFn: func(process statusjson.Process) string {
			memUsage := float64(process.Memory.UsedBytes) / float64(process.Memory.AvailableBytes)

			used := float64(process.Memory.UsedBytes) / 1024 / 1024 / 1024
			available := float64(process.Memory.AvailableBytes) / 1024 / 1024 / 1024

			return fmt.Sprintf("%0.1f%% (%0.1f of %0.1f MiB)", memUsage*100, used, available)
		},
	},
	ColumnDiskUsage: {
		Name: "Disk Usage",
		DataFn: func(process statusjson.Process) string {
			usedBytes := process.Disk.TotalBytes - process.Disk.FreeBytes
			diskUsage := float64(usedBytes) / float64(process.Disk.TotalBytes)

			used := float64(usedBytes) / 1024 / 1024 / 1024
			available := float64(process.Disk.TotalBytes) / 1024 / 1024 / 1024

			return fmt.Sprintf("%0.1f%% (%0.1f of %0.1f GiB)", diskUsage*100, used, available)
		},
	},
	ColumnCPUActivity: {
		Name: "CPU Activity",
		DataFn: func(process statusjson.Process) string {
			return fmt.Sprintf("%0.1f%%", process.CPU.UsageCores*100)
		},
	},
	ColumnDiskActivity: {
		Name: "Disk Activity",
		DataFn: func(process statusjson.Process) string {
			busy := process.Disk.Busy * 100
			return fmt.Sprintf("%0.1f RPS / %0.1f WPS / %0.1f%%", process.Disk.Reads.Hz, process.Disk.Writes.Hz, busy)
		},
	},
	ColumnNetworkActivity: {
		Name: "Network Activity",
		DataFn: func(process statusjson.Process) string {
			return fmt.Sprintf("%0.1f Mbps / %0.1f Mbps", process.Network.MegabitsSent.Hz, process.Network.MegabitsReceived.Hz)
		},
	},
	ColumnVersion: {
		Name: "Version",
		DataFn: func(process statusjson.Process) string {
			return process.Version
		},
	},
	ColumnUptime: {
		Name: "Uptime",
		DataFn: func(process statusjson.Process) string {
			return (time.Duration(process.Uptime) * time.Second).String()
		},
	},
	ColumnKVStorage: {
		Name: "KV Storage",
		DataFn: func(process statusjson.Process) string {
			used := float64(process.Roles[0].KVUsedBytes) / 1024 / 1024
			return fmt.Sprintf("%0.1f MiB", used)
		},
	},
	ColumnQueueStorage: {
		Name: "Queue Storage",
		DataFn: func(process statusjson.Process) string {
			used := float64(process.Roles[0].QueueUsedBytes) / 1024 / 1024
			return fmt.Sprintf("%0.1f MiB", used)
		},
	},
	ColumnDurabilityRate: {
		Name: "Input / Durable Rate",
		DataFn: func(process statusjson.Process) string {
			input := float64(process.Roles[0].InputBytes.Hz) / 1024 / 1024
			durable := float64(process.Roles[0].DurableBytes.Hz) / 1024 / 1024

			return fmt.Sprintf("%0.1f MiB/s / %0.1f MiB/s", input, durable)
		},
	},
	ColumnStorageLag: {
		Name: "Data / Durability Lag",
		DataFn: func(process statusjson.Process) string {
			return fmt.Sprintf("%0.1fs / %0.1fs", process.Roles[0].DataLag.Seconds, process.Roles[0].DurabilityLag.Seconds)
		},
	},
	ColumnTotalQueries: {
		Name: "Queries",
		DataFn: func(process statusjson.Process) string {
			return fmt.Sprintf("%0.1f/s", process.Roles[0].TotalQueries.Hz)
		},
	},
}
