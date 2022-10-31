package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/data/fdb"
	"github.com/pwood/fdbexplorer/output/ui/components"
	"strings"
	"time"
)

func UpdateProcesses(f func([]fdb.Process)) func(fdb.Root) {
	return func(root fdb.Root) {
		var processes []fdb.Process

		for _, p := range root.Cluster.Processes {
			processes = append(processes, fdb.AnnotateProcessHealth(p))
		}

		f(processes)
	}
}

func All(_ fdb.Process) bool {
	return true
}

func RoleMatch(s string) func(fdb.Process) bool {
	return func(process fdb.Process) bool {
		for _, r := range process.Roles {
			if r.Role == s {
				return true
			}
		}
		return false
	}
}

type ProcessSorter struct {
	i int
}

func (p *ProcessSorter) Next() {
	p.i++
	if p.i > 2 {
		p.i = 0
	}
}

func (p *ProcessSorter) SortName() string {
	switch p.i {
	case 0:
		return "Address"
	case 1:
		return "Role"
	case 2:
		return "Class"
	default:
		return "Unknown"
	}
}

func (p *ProcessSorter) Sort(i fdb.Process, j fdb.Process) bool {
	iKey := i.Address
	jKey := j.Address

	switch p.i {
	case 1:
		iRole := ""
		if len(i.Roles) > 0 {
			iRole = i.Roles[0].Role
		}

		jRole := ""
		if len(j.Roles) > 0 {
			jRole = j.Roles[0].Role
		}

		iKey = iRole + iKey
		jKey = jRole + jKey
	case 2:
		iKey = i.Class + iKey
		jKey = j.Class + jKey
	}

	return strings.Compare(iKey, jKey) < 0
}

func ProcessColour(p fdb.Process) tcell.Color {
	switch p.Health {
	case fdb.HealthCritical:
		return tcell.ColorRed
	case fdb.HealthWarning:
		return tcell.ColorYellow
	case fdb.HealthExcluded:
		return tcell.ColorBlue
	default:
		return tcell.ColorWhite
	}
}

var ColumnIPAddressPort = components.ColumnImpl[fdb.Process]{
	ColName: "IP Address:Port",
	DataFn: func(process fdb.Process) string {
		return process.Address
	},
	ColorFn: ProcessColour,
}

var ColumnStatus = components.ColumnImpl[fdb.Process]{
	ColName: "Status",
	DataFn: func(process fdb.Process) string {
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
	ColorFn: ProcessColour,
}

var ColumnMachine = components.ColumnImpl[fdb.Process]{
	ColName: "Machine",
	DataFn: func(process fdb.Process) string {
		return process.Locality.MachineID
	},
	ColorFn: ProcessColour,
}

var ColumnLocality = components.ColumnImpl[fdb.Process]{
	ColName: "Locality",
	DataFn: func(process fdb.Process) string {
		return fmt.Sprintf("%s / %s", process.Locality.DataHall, process.Locality.DCID)
	},
	ColorFn: ProcessColour,
}

var ColumnClass = components.ColumnImpl[fdb.Process]{
	ColName: "Class",
	DataFn: func(process fdb.Process) string {
		return process.Class
	},
	ColorFn: ProcessColour,
}

var ColumnRoles = components.ColumnImpl[fdb.Process]{
	ColName: "Roles",
	DataFn: func(process fdb.Process) string {
		var roles []string

		for _, role := range process.Roles {
			roles = append(roles, role.Role)
		}

		return strings.Join(roles, ", ")
	},
	ColorFn: ProcessColour,
}

var ColumnRAMUsage = components.ColumnImpl[fdb.Process]{
	ColName: "RAM Usage",
	DataFn: func(process fdb.Process) string {
		memUsage := float64(process.Memory.RSSBytes) / float64(process.Memory.AvailableBytes)
		return fmt.Sprintf("%0.1f%% (%s of %s)", memUsage*100, convert(float64(process.Memory.RSSBytes), 1, None), convert(float64(process.Memory.AvailableBytes), 1, None))
	},
	ColorFn: ProcessColour,
}

var ColumnDiskUsage = components.ColumnImpl[fdb.Process]{
	ColName: "Disk Usage",
	DataFn: func(process fdb.Process) string {
		usedBytes := process.Disk.TotalBytes - process.Disk.FreeBytes
		diskUsage := float64(usedBytes) / float64(process.Disk.TotalBytes)
		return fmt.Sprintf("%0.1f%% (%s of %s)", diskUsage*100, convert(float64(usedBytes), 1, None), convert(float64(process.Disk.TotalBytes), 1, None))
	},
	ColorFn: ProcessColour,
}

var ColumnCPUActivity = components.ColumnImpl[fdb.Process]{
	ColName: "CPU Activity",
	DataFn: func(process fdb.Process) string {
		return fmt.Sprintf("%0.1f%%", process.CPU.UsageCores*100)
	},
	ColorFn: ProcessColour,
}

var ColumnDiskActivity = components.ColumnImpl[fdb.Process]{
	ColName: "Disk Activity",
	DataFn: func(process fdb.Process) string {
		busy := process.Disk.Busy * 100
		return fmt.Sprintf("%0.1f RPS / %0.1f WPS / %0.1f%%", process.Disk.Reads.Hz, process.Disk.Writes.Hz, busy)
	},
	ColorFn: ProcessColour,
}

var ColumnNetworkActivity = components.ColumnImpl[fdb.Process]{
	ColName: "Network Activity",
	DataFn: func(process fdb.Process) string {
		return fmt.Sprintf("%0.1f Mbps / %0.1f Mbps", process.Network.MegabitsSent.Hz, process.Network.MegabitsReceived.Hz)
	},
	ColorFn: ProcessColour,
}

var ColumnVersion = components.ColumnImpl[fdb.Process]{
	ColName: "Version",
	DataFn: func(process fdb.Process) string {
		return process.Version
	},
	ColorFn: ProcessColour,
}

var ColumnUptime = components.ColumnImpl[fdb.Process]{
	ColName: "Uptime",
	DataFn: func(process fdb.Process) string {
		return (time.Duration(process.Uptime) * time.Second).String()
	},
	ColorFn: ProcessColour,
}

var ColumnKVStorage = components.ColumnImpl[fdb.Process]{
	ColName: "KV Storage",
	DataFn: func(process fdb.Process) string {
		idx := findRole(process.Roles, "storage")
		return convert(process.Roles[idx].KVUsedBytes, 1, None)
	},
	ColorFn: ProcessColour,
}

var ColumnLogQueueStorage = components.ColumnImpl[fdb.Process]{
	ColName: "Queue Storage",
	DataFn: func(process fdb.Process) string {
		idx := findRole(process.Roles, "log")
		return convert(process.Roles[idx].QueueUsedBytes, 1, None)
	},
	ColorFn: ProcessColour,
}

var ColumnLogQueueLength = components.ColumnImpl[fdb.Process]{
	ColName: "Queue Length",
	DataFn: func(process fdb.Process) string {
		idx := findRole(process.Roles, "log")
		length := process.Roles[idx].InputBytes.Counter - process.Roles[idx].DurableBytes.Counter
		return convert(length, 1, None)
	},
	ColorFn: ProcessColour,
}

var ColumnStorageDurabilityRate = components.ColumnImpl[fdb.Process]{
	ColName: "Input / Durable Rate",
	DataFn: func(process fdb.Process) string {
		idx := findRole(process.Roles, "storage")
		return fmt.Sprintf("%s / %s", convert(process.Roles[idx].InputBytes.Hz, 1, "s"), convert(process.Roles[idx].DurableBytes.Hz, 1, "s"))
	},
	ColorFn: ProcessColour,
}

var ColumnLogDurabilityRate = components.ColumnImpl[fdb.Process]{
	ColName: "Input / Durable Rate",
	DataFn: func(process fdb.Process) string {
		idx := findRole(process.Roles, "log")
		return fmt.Sprintf("%s / %s", convert(process.Roles[idx].InputBytes.Hz, 1, "s"), convert(process.Roles[idx].DurableBytes.Hz, 1, "s"))
	},
	ColorFn: ProcessColour,
}

var ColumnStorageLag = components.ColumnImpl[fdb.Process]{
	ColName: "Data / Durability Lag",
	DataFn: func(process fdb.Process) string {
		idx := findRole(process.Roles, "storage")
		return fmt.Sprintf("%0.1fs / %0.1fs", process.Roles[idx].DataLag.Seconds, process.Roles[idx].DurabilityLag.Seconds)
	},
	ColorFn: ProcessColour,
}

var ColumnStorageTotalQueries = components.ColumnImpl[fdb.Process]{
	ColName: "Queries",
	DataFn: func(process fdb.Process) string {
		idx := findRole(process.Roles, "storage")
		return fmt.Sprintf("%0.1f/s", process.Roles[idx].TotalQueries.Hz)
	},
	ColorFn: ProcessColour,
}

func findRole(roles []fdb.Role, role string) int {
	for i, straw := range roles {
		if straw.Role == role {
			return i
		}
	}

	return 0
}
