package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/ui/components"
)

var StatClusterHealth = components.ColumnImpl[ClusterHealth]{
	ColName: "Healthy",
	DataFn: func(h ClusterHealth) string {
		return h.Health
	},
	ColorFn: func(h ClusterHealth) tcell.Color {
		if h.Healthy {
			return tcell.ColorGreen
		} else {
			return tcell.ColorRed
		}
	},
}

var StatReplicasRemaining = components.ColumnImpl[ClusterHealth]{
	ColName: "Replicas Remaining",
	DataFn: func(h ClusterHealth) string {
		return fmt.Sprintf("%d", h.MinReplicas)
	},
}

var StatRecoveryState = components.ColumnImpl[ClusterHealth]{
	ColName: "Recovery State",
	DataFn: func(h ClusterHealth) string {
		return h.RecoveryState
	},
}

var StatRecoveryDescription = components.ColumnImpl[ClusterHealth]{
	ColName: "Recovery Description",
	DataFn: func(h ClusterHealth) string {
		return h.RecoveryDescription
	},
}

var StatRebalanceQueued = components.ColumnImpl[ClusterHealth]{
	ColName: "Rebalance Queued",
	DataFn: func(h ClusterHealth) string {
		return fmt.Sprintf("%0.1f MiB", float64(h.RebalanceQueued)/Mibibyte)
	},
}

var StatRebalanceInflight = components.ColumnImpl[ClusterHealth]{
	ColName: "Rebalance In-flight",
	DataFn: func(h ClusterHealth) string {
		return fmt.Sprintf("%0.1f MiB", float64(h.RebalanceInFlight)/Mibibyte)
	},
}

var StatEmpty = components.ColumnImpl[ClusterHealth]{
	ColName: "",
	DataFn: func(h ClusterHealth) string {
		return ""
	},
}
