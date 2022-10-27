package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/data/fdb"
	"github.com/pwood/fdbexplorer/ui/components"
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

const Kibibyte float64 = 1024
const Mebibyte float64 = Kibibyte * 1024
const Gibibyte float64 = Mebibyte * 1024
const Tebibyte float64 = Gibibyte * 1024

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

		used := float64(process.Memory.RSSBytes) / Gibibyte
		available := float64(process.Memory.AvailableBytes) / Gibibyte

		return fmt.Sprintf("%0.1f%% (%0.1f of %0.1f GiB)", memUsage*100, used, available)
	},
	ColorFn: ProcessColour,
}

var ColumnDiskUsage = components.ColumnImpl[fdb.Process]{
	ColName: "Disk Usage",
	DataFn: func(process fdb.Process) string {
		usedBytes := process.Disk.TotalBytes - process.Disk.FreeBytes
		diskUsage := float64(usedBytes) / float64(process.Disk.TotalBytes)

		used := float64(usedBytes) / Gibibyte
		available := float64(process.Disk.TotalBytes) / Gibibyte

		return fmt.Sprintf("%0.1f%% (%0.1f of %0.1f GiB)", diskUsage*100, used, available)
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
		used := process.Roles[idx].KVUsedBytes / Gibibyte
		return fmt.Sprintf("%0.1f GiB", used)
	},
	ColorFn: ProcessColour,
}

var ColumnLogQueueStorage = components.ColumnImpl[fdb.Process]{
	ColName: "Queue Storage",
	DataFn: func(process fdb.Process) string {
		idx := findRole(process.Roles, "log")
		used := process.Roles[idx].QueueUsedBytes / Mebibyte
		return fmt.Sprintf("%0.1f MiB", used)
	},
	ColorFn: ProcessColour,
}

var ColumnLogQueueLength = components.ColumnImpl[fdb.Process]{
	ColName: "Queue Length",
	DataFn: func(process fdb.Process) string {
		idx := findRole(process.Roles, "log")
		length := (process.Roles[idx].InputBytes.Counter - process.Roles[idx].DurableBytes.Counter) / Mebibyte
		return fmt.Sprintf("%0.1f MiB", length)
	},
	ColorFn: ProcessColour,
}

var ColumnStorageDurabilityRate = components.ColumnImpl[fdb.Process]{
	ColName: "Input / Durable Rate",
	DataFn: func(process fdb.Process) string {
		idx := findRole(process.Roles, "storage")
		input := process.Roles[idx].InputBytes.Hz / Mebibyte
		durable := process.Roles[idx].DurableBytes.Hz / Mebibyte
		return fmt.Sprintf("%0.1f MiB/s / %0.1f MiB/s", input, durable)
	},
	ColorFn: ProcessColour,
}

var ColumnLogDurabilityRate = components.ColumnImpl[fdb.Process]{
	ColName: "Input / Durable Rate",
	DataFn: func(process fdb.Process) string {
		idx := findRole(process.Roles, "log")
		input := process.Roles[idx].InputBytes.Hz / Mebibyte
		durable := process.Roles[idx].DurableBytes.Hz / Mebibyte
		return fmt.Sprintf("%0.1f MiB/s / %0.1f MiB/s", input, durable)
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

func convert(v float64, dp int, perSecond bool) string {
	f := fmt.Sprintf("%%0.%df %%s%%s", dp)
	suffix := ""
	if perSecond {
		suffix = "/s"
	}

	if v < Kibibyte {
		return fmt.Sprintf(f, v, "B", suffix)
	} else if v < Mebibyte {
		return fmt.Sprintf(f, v/Kibibyte, "KiB", suffix)
	} else if v < Gibibyte {
		return fmt.Sprintf(f, v/Mebibyte, "MiB", suffix)
	} else if v < Tebibyte {
		return fmt.Sprintf(f, v/Gibibyte, "GiB", suffix)
	} else {
		return fmt.Sprintf(f, v/Tebibyte, "TiB", suffix)
	}
}
