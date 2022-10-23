package ui

import (
	"github.com/pwood/fdbexplorer/data/fdb"
	"strings"
)

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

func UpdateProcessClusterStats(f func(ClusterStats)) func(fdb.Root) {
	return func(root fdb.Root) {
		f(ClusterStats{
			TxStarted:    root.Cluster.Workload.Transactions.Started.Hz,
			TxCommitted:  root.Cluster.Workload.Transactions.Committed.Hz,
			TxConflicted: root.Cluster.Workload.Transactions.Conflicted.Hz,
			TxRejected:   root.Cluster.Workload.Transactions.RejectedForQueuedTooLong.Hz,
			Reads:        root.Cluster.Workload.Operations.Reads.Hz,
			Writes:       root.Cluster.Workload.Operations.Writes.Hz,
			BytesRead:    root.Cluster.Workload.Bytes.Read.Hz,
			BytesWritten: root.Cluster.Workload.Bytes.Written.Hz,
		})
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

func UpdateProcessClusterHealth(f func(ClusterHealth)) func(fdb.Root) {
	return func(root fdb.Root) {
		f(ClusterHealth{
			Healthy:             root.Cluster.Data.State.Health,
			Health:              strings.Title(strings.Replace(root.Cluster.Data.State.Name, "_", " ", -1)),
			MinReplicas:         root.Cluster.Data.State.MinReplicasRemaining,
			RebalanceQueued:     root.Cluster.Data.MovingData.InQueueBytes,
			RebalanceInFlight:   root.Cluster.Data.MovingData.InFlightBytes,
			RecoveryState:       strings.Title(strings.Replace(root.Cluster.RecoveryState.Name, "_", " ", -1)),
			RecoveryDescription: root.Cluster.RecoveryState.Description,
		})
	}
}
