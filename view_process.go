package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/statusjson"
	"github.com/rivo/tview"
	"sort"
	"strings"
	"sync"
)

type SortKey int

const (
	SortIPAddress SortKey = iota
	SortClass
	SortRole
)

type ProcessData struct {
	tview.TableContentReadOnly

	processes []statusjson.Process
	sortBy    SortKey

	m *sync.RWMutex
}

func (p *ProcessData) Updated() {
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
}

type ProcessView struct {
	tview.TableContentReadOnly
	pd *ProcessData

	columns []ColumnId
}

func (v *ProcessView) GetCell(row, column int) *tview.TableCell {
	v.pd.m.RLock()
	defer v.pd.m.RUnlock()

	cid := v.columns[column]

	if row == 0 {
		return tview.NewTableCell(columns[cid].Name).SetExpansion(1).SetTextColor(tcell.ColorAqua).SetSelectable(false)
	} else {
		return tview.NewTableCell(columns[cid].DataFn(v.pd.processes[row-1]))
	}
}

func (v *ProcessView) GetRowCount() int {
	v.pd.m.RLock()
	defer v.pd.m.RUnlock()

	return len(v.pd.processes) + 1
}

func (v *ProcessView) GetColumnCount() int {
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
			return fmt.Sprintf("%0.1f RPS / %0.1f WPS", process.Disk.Reads.Hz, process.Disk.Writes.Hz)
		},
	},
	ColumnNetworkActivity: {
		Name: "Network Activity",
		DataFn: func(process statusjson.Process) string {
			return fmt.Sprintf("%0.1f Mbps / %0.1f Mbps", process.Network.MegabitsSent.Hz, process.Network.MegabitsReceived.Hz)
		},
	},
}
