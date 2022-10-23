package ui

import (
	"github.com/pwood/fdbexplorer/data/fdb"
	"strings"
	"sync"
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
