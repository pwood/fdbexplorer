package ui

import (
	"fmt"
	"github.com/pwood/fdbexplorer/output/ui/components"
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

func UpdateClusterStats(f func(ClusterStats)) func(DataSourceUpdate) {
	return func(dsu DataSourceUpdate) {
		f(ClusterStats{
			TxStarted:    dsu.root.Cluster.Workload.Transactions.Started.Hz,
			TxCommitted:  dsu.root.Cluster.Workload.Transactions.Committed.Hz,
			TxConflicted: dsu.root.Cluster.Workload.Transactions.Conflicted.Hz,
			TxRejected:   dsu.root.Cluster.Workload.Transactions.RejectedForQueuedTooLong.Hz,
			Reads:        dsu.root.Cluster.Workload.Operations.Reads.Hz,
			Writes:       dsu.root.Cluster.Workload.Operations.Writes.Hz,
			BytesRead:    dsu.root.Cluster.Workload.Bytes.Read.Hz,
			BytesWritten: dsu.root.Cluster.Workload.Bytes.Written.Hz,
		})
	}
}

var StatTxStarted = components.ColumnImpl[ClusterStats]{
	ColName: "Tx Started",
	DataFn: func(cs ClusterStats) string {
		return fmt.Sprintf("%0.1f/s", cs.TxStarted)
	},
}

var StatTxCommitted = components.ColumnImpl[ClusterStats]{
	ColName: "Tx Committed",
	DataFn: func(cs ClusterStats) string {
		return fmt.Sprintf("%0.1f/s", cs.TxCommitted)
	},
}

var StatTxConflicted = components.ColumnImpl[ClusterStats]{
	ColName: "Tx Conflicted",
	DataFn: func(cs ClusterStats) string {
		return fmt.Sprintf("%0.1f/s", cs.TxConflicted)
	},
}

var StatTxRejected = components.ColumnImpl[ClusterStats]{
	ColName: "Tx Rejected",
	DataFn: func(cs ClusterStats) string {
		return fmt.Sprintf("%0.1f/s", cs.TxRejected)
	},
}

var StatReads = components.ColumnImpl[ClusterStats]{
	ColName: "Reads",
	DataFn: func(cs ClusterStats) string {
		return fmt.Sprintf("%0.1f/s", cs.Reads)
	},
}

var StatWrites = components.ColumnImpl[ClusterStats]{
	ColName: "Writes",
	DataFn: func(cs ClusterStats) string {
		return fmt.Sprintf("%0.1f/s", cs.Writes)
	},
}

var StatBytesRead = components.ColumnImpl[ClusterStats]{
	ColName: "Bytes Read",
	DataFn: func(cs ClusterStats) string {
		return convert(cs.BytesRead, 1, "s")
	},
}

var StatBytesWritten = components.ColumnImpl[ClusterStats]{
	ColName: "Bytes Written",
	DataFn: func(cs ClusterStats) string {
		return convert(cs.BytesWritten, 1, "s")
	},
}
