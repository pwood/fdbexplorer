package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/data/fdb"
	"github.com/pwood/fdbexplorer/ui/components"
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

// Deprecated
type ClusterData struct {
	root fdb.Root

	m *sync.RWMutex
}

func (c *ClusterData) Update(s fdb.Root) {
	c.m.Lock()
	defer c.m.Unlock()

	c.root = s
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

const Mibibyte float64 = 1024 * 1024
const Gibibyte float64 = Mibibyte * 1024

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
}

var ColumnMachine = components.ColumnImpl[fdb.Process]{
	ColName: "Machine",
	DataFn: func(process fdb.Process) string {
		return process.Locality.MachineID
	},
}

var ColumnLocality = components.ColumnImpl[fdb.Process]{
	ColName: "Locality",
	DataFn: func(process fdb.Process) string {
		return fmt.Sprintf("%s / %s", process.Locality.DataHall, process.Locality.DCID)
	},
}

var ColumnClass = components.ColumnImpl[fdb.Process]{
	ColName: "Class",
	DataFn: func(process fdb.Process) string {
		return process.Class
	},
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
}

var ColumnRAMUsage = components.ColumnImpl[fdb.Process]{
	ColName: "RAM Usage",
	DataFn: func(process fdb.Process) string {
		memUsage := float64(process.Memory.UsedBytes) / float64(process.Memory.AvailableBytes)

		used := float64(process.Memory.UsedBytes) / Gibibyte
		available := float64(process.Memory.AvailableBytes) / Gibibyte

		return fmt.Sprintf("%0.1f%% (%0.1f of %0.1f GiB)", memUsage*100, used, available)
	},
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
}
var ColumnCPUActivity = components.ColumnImpl[fdb.Process]{
	ColName: "CPU Activity",
	DataFn: func(process fdb.Process) string {
		return fmt.Sprintf("%0.1f%%", process.CPU.UsageCores*100)
	},
}
var ColumnDiskActivity = components.ColumnImpl[fdb.Process]{
	ColName: "Disk Activity",
	DataFn: func(process fdb.Process) string {
		busy := process.Disk.Busy * 100
		return fmt.Sprintf("%0.1f RPS / %0.1f WPS / %0.1f%%", process.Disk.Reads.Hz, process.Disk.Writes.Hz, busy)
	},
}
var ColumnNetworkActivity = components.ColumnImpl[fdb.Process]{
	ColName: "Network Activity",
	DataFn: func(process fdb.Process) string {
		return fmt.Sprintf("%0.1f Mbps / %0.1f Mbps", process.Network.MegabitsSent.Hz, process.Network.MegabitsReceived.Hz)
	},
}
var ColumnVersion = components.ColumnImpl[fdb.Process]{
	ColName: "Version",
	DataFn: func(process fdb.Process) string {
		return process.Version
	},
}

var ColumnUptime = components.ColumnImpl[fdb.Process]{
	ColName: "Uptime",
	DataFn: func(process fdb.Process) string {
		return (time.Duration(process.Uptime) * time.Second).String()
	},
}

var ColumnKVStorage = components.ColumnImpl[fdb.Process]{
	ColName: "KV Storage",
	DataFn: func(process fdb.Process) string {
		used := process.Roles[0].KVUsedBytes / Gibibyte
		return fmt.Sprintf("%0.1f GiB", used)
	},
}

var ColumnQueueStorage = components.ColumnImpl[fdb.Process]{
	ColName: "Queue Storage",
	DataFn: func(process fdb.Process) string {
		used := process.Roles[0].QueueUsedBytes / Mibibyte
		return fmt.Sprintf("%0.1f MiB", used)
	},
}

var ColumnDurabilityRate = components.ColumnImpl[fdb.Process]{
	ColName: "Input / Durable Rate",
	DataFn: func(process fdb.Process) string {
		input := process.Roles[0].InputBytes.Hz / Mibibyte
		durable := process.Roles[0].DurableBytes.Hz / Mibibyte
		return fmt.Sprintf("%0.1f MiB/s / %0.1f MiB/s", input, durable)
	},
}

var ColumnStorageLag = components.ColumnImpl[fdb.Process]{
	ColName: "Data / Durability Lag",
	DataFn: func(process fdb.Process) string {
		return fmt.Sprintf("%0.1fs / %0.1fs", process.Roles[0].DataLag.Seconds, process.Roles[0].DurabilityLag.Seconds)
	},
}

var ColumnTotalQueries = components.ColumnImpl[fdb.Process]{
	ColName: "Queries",
	DataFn: func(process fdb.Process) string {
		return fmt.Sprintf("%0.1f/s", process.Roles[0].TotalQueries.Hz)
	},
}
