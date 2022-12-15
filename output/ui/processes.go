package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/data/fdb"
	"github.com/pwood/fdbexplorer/output/ui/components"
	"strings"
	"time"
)

type ProcessData struct {
	Process  fdb.Process
	Metadata *ProcessMetadata
}

type Health int

const (
	HealthCritical Health = iota
	HealthWarning
	HealthNormal
	HealthExcluded
)

type ProcessMetadata struct {
	Health   Health
	Selected bool
}

func (p *ProcessMetadata) ToggleSelected() {
	p.Selected = !p.Selected
}

func (p *ProcessMetadata) Update(proc fdb.Process) {
	p.Health = HealthNormal

	if proc.Excluded || proc.UnderMaintenance {
		p.Health = HealthExcluded
	}

	if len(proc.Messages) > 0 {
		p.Health = HealthWarning
	}

	if proc.Degraded {
		p.Health = HealthCritical
	}
}

type metadataStore struct {
	metadata map[string]*ProcessMetadata
}

func (m *metadataStore) Update(f func([]ProcessData)) func(fdb.Root) {
	return func(root fdb.Root) {
		var processes []ProcessData

		for _, proc := range root.Cluster.Processes {
			meta := m.findOrCreateMetadata(proc.Locality[fdb.LocalityProcessID])
			meta.Update(proc)
			processes = append(processes, ProcessData{Process: proc, Metadata: meta})
		}

		f(processes)
	}
}

func (m *metadataStore) ClearSelected() {
	for _, pm := range m.metadata {
		pm.Selected = false
	}
}

func (m *metadataStore) findOrCreateMetadata(id string) *ProcessMetadata {
	pd, ok := m.metadata[id]

	if !ok {
		pd = &ProcessMetadata{}
		m.metadata[id] = pd
	}

	return pd
}

func All(_ ProcessData) bool {
	return true
}

func RoleMatch(s string) func(ProcessData) bool {
	return func(process ProcessData) bool {
		for _, r := range process.Process.Roles {
			if r.Role == s {
				return true
			}
		}
		return false
	}
}

func ProcessColour(p ProcessData) tcell.Color {
	switch p.Metadata.Health {
	case HealthCritical:
		return tcell.ColorRed
	case HealthWarning:
		return tcell.ColorYellow
	case HealthExcluded:
		return tcell.ColorBlue
	default:
		if p.Metadata.Selected {
			return tcell.ColorGreen
		}
		return tcell.ColorWhite
	}
}

var ColumnSelected = components.ColumnImpl[ProcessData]{
	ColName: " ",
	DataFn: func(pd ProcessData) string {
		if pd.Metadata.Selected {
			return "*"
		}
		return " "
	},
	ColorFn: func(pd ProcessData) tcell.Color {
		if pd.Metadata.Selected {
			return tcell.ColorGreen
		}
		return ProcessColour(pd)
	},
}

var ColumnIPAddressPort = components.ColumnImpl[ProcessData]{
	ColName: "IP Address:Port",
	DataFn: func(pd ProcessData) string {
		return pd.Process.Address
	},
	ColorFn: ProcessColour,
}

var ColumnStatus = components.ColumnImpl[ProcessData]{
	ColName: "Status",
	DataFn: func(pd ProcessData) string {
		var statuses []string

		if pd.Process.Excluded {
			statuses = append(statuses, "Excluded")
		}

		if pd.Process.Degraded {
			statuses = append(statuses, "Degraded")
		}

		if pd.Process.UnderMaintenance {
			statuses = append(statuses, "Maintenance")
		}

		if len(pd.Process.Messages) > 0 {
			statuses = append(statuses, "Message")
		}

		return strings.Join(statuses, " / ")
	},
	ColorFn: ProcessColour,
}

var ColumnMachine = components.ColumnImpl[ProcessData]{
	ColName: "Machine",
	DataFn: func(pd ProcessData) string {
		return pd.Process.Locality[fdb.LocalityMachineID]
	},
	ColorFn: ProcessColour,
}

var ColumnLocality = components.ColumnImpl[ProcessData]{
	ColName: "Locality",
	DataFn: func(pd ProcessData) string {
		return fmt.Sprintf("%s / %s", pd.Process.Locality[fdb.LocalityDataHall], pd.Process.Locality[fdb.LocalityDataCenter])
	},
	ColorFn: ProcessColour,
}

var ColumnClass = components.ColumnImpl[ProcessData]{
	ColName: "Class",
	DataFn: func(pd ProcessData) string {
		return pd.Process.Class
	},
	ColorFn: ProcessColour,
}

var ColumnRoles = components.ColumnImpl[ProcessData]{
	ColName: "Roles",
	DataFn: func(pd ProcessData) string {
		var roles []string

		for _, role := range pd.Process.Roles {
			roles = append(roles, role.Role)
		}

		return strings.Join(roles, ", ")
	},
	ColorFn: ProcessColour,
}

var ColumnRAMUsage = components.ColumnImpl[ProcessData]{
	ColName: "RAM Usage",
	DataFn: func(pd ProcessData) string {
		memUsage := float64(pd.Process.Memory.RSSBytes) / float64(pd.Process.Memory.AvailableBytes)
		return fmt.Sprintf("%0.1f%% (%s of %s)", memUsage*100, convert(float64(pd.Process.Memory.RSSBytes), 1, None), convert(float64(pd.Process.Memory.AvailableBytes), 1, None))
	},
	ColorFn: ProcessColour,
}

var ColumnDiskUsage = components.ColumnImpl[ProcessData]{
	ColName: "Disk Usage",
	DataFn: func(pd ProcessData) string {
		usedBytes := pd.Process.Disk.TotalBytes - pd.Process.Disk.FreeBytes
		diskUsage := float64(usedBytes) / float64(pd.Process.Disk.TotalBytes)
		return fmt.Sprintf("%0.1f%% (%s of %s)", diskUsage*100, convert(float64(usedBytes), 1, None), convert(float64(pd.Process.Disk.TotalBytes), 1, None))
	},
	ColorFn: ProcessColour,
}

var ColumnCPUActivity = components.ColumnImpl[ProcessData]{
	ColName: "CPU Activity",
	DataFn: func(pd ProcessData) string {
		return fmt.Sprintf("%0.1f%%", pd.Process.CPU.UsageCores*100)
	},
	ColorFn: ProcessColour,
}

var ColumnDiskActivity = components.ColumnImpl[ProcessData]{
	ColName: "Disk Activity",
	DataFn: func(pd ProcessData) string {
		busy := pd.Process.Disk.Busy * 100
		return fmt.Sprintf("%0.1f RPS / %0.1f WPS / %0.1f%%", pd.Process.Disk.Reads.Hz, pd.Process.Disk.Writes.Hz, busy)
	},
	ColorFn: ProcessColour,
}

var ColumnNetworkActivity = components.ColumnImpl[ProcessData]{
	ColName: "Network Activity",
	DataFn: func(pd ProcessData) string {
		return fmt.Sprintf("%0.1f Mbps / %0.1f Mbps", pd.Process.Network.MegabitsSent.Hz, pd.Process.Network.MegabitsReceived.Hz)
	},
	ColorFn: ProcessColour,
}

var ColumnVersion = components.ColumnImpl[ProcessData]{
	ColName: "Version",
	DataFn: func(pd ProcessData) string {
		return pd.Process.Version
	},
	ColorFn: ProcessColour,
}

var ColumnUptime = components.ColumnImpl[ProcessData]{
	ColName: "Uptime",
	DataFn: func(process ProcessData) string {
		return (time.Duration(process.Process.Uptime) * time.Second).String()
	},
	ColorFn: ProcessColour,
}

var ColumnKVStorage = components.ColumnImpl[ProcessData]{
	ColName: "KV Storage",
	DataFn: func(pd ProcessData) string {
		idx := findRole(pd.Process.Roles, "storage")
		return convert(pd.Process.Roles[idx].KVUsedBytes, 1, None)
	},
	ColorFn: ProcessColour,
}

var ColumnLogQueueStorage = components.ColumnImpl[ProcessData]{
	ColName: "Queue Storage",
	DataFn: func(pd ProcessData) string {
		idx := findRole(pd.Process.Roles, "log")
		return convert(pd.Process.Roles[idx].QueueUsedBytes, 1, None)
	},
	ColorFn: ProcessColour,
}

var ColumnLogQueueLength = components.ColumnImpl[ProcessData]{
	ColName: "Queue Length",
	DataFn: func(pd ProcessData) string {
		idx := findRole(pd.Process.Roles, "log")
		length := pd.Process.Roles[idx].InputBytes.Counter - pd.Process.Roles[idx].DurableBytes.Counter
		return convert(length, 1, None)
	},
	ColorFn: ProcessColour,
}

var ColumnStorageDurabilityRate = components.ColumnImpl[ProcessData]{
	ColName: "Input / Durable Rate",
	DataFn: func(pd ProcessData) string {
		idx := findRole(pd.Process.Roles, "storage")
		return fmt.Sprintf("%s / %s", convert(pd.Process.Roles[idx].InputBytes.Hz, 1, "s"), convert(pd.Process.Roles[idx].DurableBytes.Hz, 1, "s"))
	},
	ColorFn: ProcessColour,
}

var ColumnLogDurabilityRate = components.ColumnImpl[ProcessData]{
	ColName: "Input / Durable Rate",
	DataFn: func(pd ProcessData) string {
		idx := findRole(pd.Process.Roles, "log")
		return fmt.Sprintf("%s / %s", convert(pd.Process.Roles[idx].InputBytes.Hz, 1, "s"), convert(pd.Process.Roles[idx].DurableBytes.Hz, 1, "s"))
	},
	ColorFn: ProcessColour,
}

var ColumnStorageLag = components.ColumnImpl[ProcessData]{
	ColName: "Data / Durability Lag",
	DataFn: func(pd ProcessData) string {
		idx := findRole(pd.Process.Roles, "storage")
		return fmt.Sprintf("%0.1fs / %0.1fs", pd.Process.Roles[idx].DataLag.Seconds, pd.Process.Roles[idx].DurabilityLag.Seconds)
	},
	ColorFn: ProcessColour,
}

var ColumnStorageTotalQueries = components.ColumnImpl[ProcessData]{
	ColName: "Queries",
	DataFn: func(pd ProcessData) string {
		idx := findRole(pd.Process.Roles, "storage")
		return fmt.Sprintf("%0.1f/s", pd.Process.Roles[idx].TotalQueries.Hz)
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
