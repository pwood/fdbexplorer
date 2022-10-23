package ui

import (
	"fmt"
	"github.com/pwood/fdbexplorer/ui/components"
)

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
		return fmt.Sprintf("%0.1f MiB/s", cs.BytesRead/Mibibyte)
	},
}

var StatBytesWritten = components.ColumnImpl[ClusterStats]{
	ColName: "Bytes Written",
	DataFn: func(cs ClusterStats) string {
		return fmt.Sprintf("%0.1f MiB/s", cs.BytesWritten/Mibibyte)
	},
}
