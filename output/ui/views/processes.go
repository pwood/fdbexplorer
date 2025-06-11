package views

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/data/fdb"
	"github.com/pwood/fdbexplorer/output/ui/components"
	"github.com/pwood/fdbexplorer/output/ui/data/process"
	"strings"
	"time"
)

func All(_ process.Process) bool {
	return true
}

func Selected(p process.Process) bool {
	return p.Metadata.Selected
}

func RoleMatch(s string) func(process.Process) bool {
	return func(process process.Process) bool {
		for _, r := range process.FDBData.Roles {
			if r.Role == s {
				return true
			}
		}
		return false
	}
}

func ProcessColour(p process.Process) tcell.Color {
	switch p.Metadata.Health {
	case process.HealthCritical:
		return tcell.ColorRed
	case process.HealthWarning:
		return tcell.ColorYellow
	case process.HealthExcluded:
		if p.Metadata.ExclusionInProgress {
			return tcell.ColorOlive
		} else {
			return tcell.ColorBlue
		}
	case process.HealthExcludedOnly:
		return tcell.ColorPurple
	default:
		if p.Metadata.Selected {
			return tcell.ColorGreen
		}
		return tcell.ColorWhite
	}
}

var ColumnSelected = components.ColumnImpl[process.Process]{
	ColName: " ",
	DataFn: func(pd process.Process) string {
		if pd.Metadata.Selected {
			return "*"
		}
		return " "
	},
	ColorFn: func(pd process.Process) tcell.Color {
		if pd.Metadata.Selected {
			return tcell.ColorGreen
		}
		return ProcessColour(pd)
	},
}

var ColumnIPAddressPort = components.ColumnImpl[process.Process]{
	ColName: "IP Address:Port",
	DataFn: func(pd process.Process) string {
		return pd.FDBData.Address
	},
	ColorFn: ProcessColour,
}

var ColumnTLS = components.ColumnImpl[process.Process]{
	ColName: "TLS",
	DataFn: func(p process.Process) string {
		if p.FDBData.TLS {
			return "✓"
		}

		return ""
	},
	ColorFn: ProcessColour,
}

var ColumnStatus = components.ColumnImpl[process.Process]{
	ColName: "Status",
	DataFn: func(pd process.Process) string {
		var statuses []string

		if pd.FDBData.Excluded {
			statuses = append(statuses, "Excluded")
		}

		if pd.FDBData.Degraded {
			statuses = append(statuses, "Degraded")
		}

		if pd.FDBData.UnderMaintenance {
			statuses = append(statuses, "Maintenance")
		}

		if len(pd.FDBData.Messages) > 0 {
			statuses = append(statuses, "Message")
		}

		return strings.Join(statuses, " / ")
	},
	ColorFn: ProcessColour,
}

var ColumnMachine = components.ColumnImpl[process.Process]{
	ColName: "Machine",
	DataFn: func(pd process.Process) string {
		return pd.FDBData.Locality[fdb.LocalityMachineID]
	},
	ColorFn: ProcessColour,
}

var ColumnLocality = components.ColumnImpl[process.Process]{
	ColName: "Locality",
	DataFn: func(pd process.Process) string {
		return fmt.Sprintf("%s / %s", pd.FDBData.Locality[fdb.LocalityDataHall], pd.FDBData.Locality[fdb.LocalityDataCenter])
	},
	ColorFn: ProcessColour,
}

var ColumnClass = components.ColumnImpl[process.Process]{
	ColName: "Class",
	DataFn: func(pd process.Process) string {
		return pd.FDBData.Class
	},
	ColorFn: ProcessColour,
}

var ColumnRoles = components.ColumnImpl[process.Process]{
	ColName: "Roles",
	DataFn: func(pd process.Process) string {
		var roles []string

		for _, role := range pd.FDBData.Roles {
			roles = append(roles, role.Role)
		}

		return strings.Join(roles, ", ")
	},
	ColorFn: ProcessColour,
}

var ColumnRAMUsage = components.ColumnImpl[process.Process]{
	ColName: "RAM Usage",
	DataFn: func(pd process.Process) string {
		memUsage := float64(pd.FDBData.Memory.RSSBytes) / float64(pd.FDBData.Memory.AvailableBytes)
		return fmt.Sprintf("%0.1f%% (%s of %s)", memUsage*100, Convert(float64(pd.FDBData.Memory.RSSBytes), 1, None), Convert(float64(pd.FDBData.Memory.AvailableBytes), 1, None))
	},
	ColorFn: ProcessColour,
}

var ColumnDiskUsage = components.ColumnImpl[process.Process]{
	ColName: "Disk Usage",
	DataFn: func(pd process.Process) string {
		usedBytes := pd.FDBData.Disk.TotalBytes - pd.FDBData.Disk.FreeBytes
		diskUsage := float64(usedBytes) / float64(pd.FDBData.Disk.TotalBytes)
		return fmt.Sprintf("%0.1f%% (%s of %s)", diskUsage*100, Convert(float64(usedBytes), 1, None), Convert(float64(pd.FDBData.Disk.TotalBytes), 1, None))
	},
	ColorFn: ProcessColour,
}

var ColumnCPUActivity = components.ColumnImpl[process.Process]{
	ColName: "CPU Activity",
	DataFn: func(pd process.Process) string {
		return fmt.Sprintf("%0.1f%%", pd.FDBData.CPU.UsageCores*100)
	},
	ColorFn: ProcessColour,
}

var ColumnDiskActivity = components.ColumnImpl[process.Process]{
	ColName: "Disk Activity",
	DataFn: func(pd process.Process) string {
		busy := pd.FDBData.Disk.Busy * 100
		return fmt.Sprintf("%0.1f RPS / %0.1f WPS / %0.1f%%", pd.FDBData.Disk.Reads.Hz, pd.FDBData.Disk.Writes.Hz, busy)
	},
	ColorFn: ProcessColour,
}

var ColumnNetworkActivity = components.ColumnImpl[process.Process]{
	ColName: "Network Activity",
	DataFn: func(pd process.Process) string {
		return fmt.Sprintf("%0.1f Mbps / %0.1f Mbps", pd.FDBData.Network.MegabitsSent.Hz, pd.FDBData.Network.MegabitsReceived.Hz)
	},
	ColorFn: ProcessColour,
}

var ColumnVersion = components.ColumnImpl[process.Process]{
	ColName: "Version",
	DataFn: func(pd process.Process) string {
		return pd.FDBData.Version
	},
	ColorFn: ProcessColour,
}

var ColumnUptime = components.ColumnImpl[process.Process]{
	ColName: "Uptime",
	DataFn: func(process process.Process) string {
		return (time.Duration(process.FDBData.Uptime) * time.Second).String()
	},
	ColorFn: ProcessColour,
}

var ColumnKVStorage = components.ColumnImpl[process.Process]{
	ColName: "KV Storage",
	DataFn: func(pd process.Process) string {
		idx := findRole(pd.FDBData.Roles, "storage")
		return Convert(pd.FDBData.Roles[idx].KVUsedBytes, 1, None)
	},
	ColorFn: ProcessColour,
}

var ColumnLogQueueStorage = components.ColumnImpl[process.Process]{
	ColName: "Queue Storage",
	DataFn: func(pd process.Process) string {
		idx := findRole(pd.FDBData.Roles, "log")
		return Convert(pd.FDBData.Roles[idx].QueueUsedBytes, 1, None)
	},
	ColorFn: ProcessColour,
}

var ColumnLogQueueLength = components.ColumnImpl[process.Process]{
	ColName: "Queue Length",
	DataFn: func(pd process.Process) string {
		idx := findRole(pd.FDBData.Roles, "log")
		length := pd.FDBData.Roles[idx].InputBytes.Counter - pd.FDBData.Roles[idx].DurableBytes.Counter
		return Convert(length, 1, None)
	},
	ColorFn: ProcessColour,
}

var ColumnStorageDurabilityRate = components.ColumnImpl[process.Process]{
	ColName: "Input / Durable Rate",
	DataFn: func(pd process.Process) string {
		idx := findRole(pd.FDBData.Roles, "storage")
		return fmt.Sprintf("%s / %s", Convert(pd.FDBData.Roles[idx].InputBytes.Hz, 1, "s"), Convert(pd.FDBData.Roles[idx].DurableBytes.Hz, 1, "s"))
	},
	ColorFn: ProcessColour,
}

var ColumnLogDurabilityRate = components.ColumnImpl[process.Process]{
	ColName: "Input / Durable Rate",
	DataFn: func(pd process.Process) string {
		idx := findRole(pd.FDBData.Roles, "log")
		return fmt.Sprintf("%s / %s", Convert(pd.FDBData.Roles[idx].InputBytes.Hz, 1, "s"), Convert(pd.FDBData.Roles[idx].DurableBytes.Hz, 1, "s"))
	},
	ColorFn: ProcessColour,
}

var ColumnStorageLag = components.ColumnImpl[process.Process]{
	ColName: "Data / Durability Lag",
	DataFn: func(pd process.Process) string {
		idx := findRole(pd.FDBData.Roles, "storage")
		return fmt.Sprintf("%0.1fs / %0.1fs", pd.FDBData.Roles[idx].DataLag.Seconds, pd.FDBData.Roles[idx].DurabilityLag.Seconds)
	},
	ColorFn: ProcessColour,
}

var ColumnStorageTotalQueries = components.ColumnImpl[process.Process]{
	ColName: "Queries",
	DataFn: func(pd process.Process) string {
		idx := findRole(pd.FDBData.Roles, "storage")
		return fmt.Sprintf("%0.1f/s", pd.FDBData.Roles[idx].TotalQueries.Hz)
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
