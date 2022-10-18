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

type ProcessData struct {
	processes []statusjson.Process
	views     map[string][]statusjson.Process
	viewFns   map[string]func(statusjson.Process) bool

	sortBy SortKey

	m *sync.RWMutex
}

func (p *ProcessData) Sort(sortKey SortKey) {
	p.m.Lock()
	defer p.m.Unlock()

	p.sortBy = sortKey
	p._filterAndSort()
}

func (p *ProcessData) Update(processes []statusjson.Process) {
	p.m.Lock()
	defer p.m.Unlock()

	p.processes = processes
	p._filterAndSort()
}

func (p *ProcessData) _filterAndSort() {
	switch p.sortBy {
	case SortIPAddress:
		sort.Slice(p.processes, func(i, j int) bool {
			return strings.Compare(p.processes[i].Address, p.processes[j].Address) < 0
		})
	case SortClass:
		sort.Slice(p.processes, func(i, j int) bool {
			return strings.Compare(p.processes[i].Class+p.processes[i].Address, p.processes[j].Class+p.processes[j].Address) < 0
		})
	case SortRole:
		sort.Slice(p.processes, func(i, j int) bool {
			iRole := ""
			if len(p.processes[i].Roles) > 0 {
				iRole = p.processes[i].Roles[0].Role
			}

			jRole := ""
			if len(p.processes[j].Roles) > 0 {
				jRole = p.processes[j].Roles[0].Role
			}

			return strings.Compare(iRole+p.processes[i].Address, jRole+p.processes[j].Address) < 0
		})
	}

	for n, fn := range p.viewFns {
		var viewProcesses []statusjson.Process

		for _, p := range p.processes {
			if fn(p) {
				viewProcesses = append(viewProcesses, p)
			}
		}

		p.views[n] = viewProcesses
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

func (p *ProcessData) View(name string, fn func(statusjson.Process) bool) *ProcessView {
	p.views[name] = nil
	p.viewFns[name] = fn

	return &ProcessView{
		pd: p,
		n:  name,
	}
}

type ProcessView struct {
	pd *ProcessData
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
		return tview.NewTableCell(columns[cid].DataFn(v.pv.Get(row - 1)))
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
			if process.Excluded {
				return "Excluded"
			} else {
				return ""
			}
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

			used := process.Memory.UsedBytes / 1024 / 1024
			available := process.Memory.AvailableBytes / 1024 / 1024

			return fmt.Sprintf("%0.1f%% (%d MiB of %d MiB)", memUsage*100, used, available)
		},
	},
	ColumnDiskUsage: {
		Name: "Disk Usage",
		DataFn: func(process statusjson.Process) string {
			usedBytes := process.Disk.TotalBytes - process.Disk.FreeBytes
			diskUsage := float64(usedBytes) / float64(process.Disk.TotalBytes)

			used := float64(usedBytes) / 1024 / 1024 / 1024
			available := float64(process.Disk.TotalBytes) / 1024 / 1024 / 1024

			return fmt.Sprintf("%0.1f%% (%0.1f GiB of %0.1f GiB)", diskUsage*100, used, available)
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
