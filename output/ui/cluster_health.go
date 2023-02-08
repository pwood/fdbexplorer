package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/output/ui/components"
)

type ClusterHealth struct {
	Healthy     bool
	Health      string
	MinReplicas int

	RebalanceInFlight int
	RebalanceQueued   int

	RecoveryState       string
	RecoveryDescription string
}

func UpdateClusterHealth(f func(ClusterHealth)) func(DataSourceUpdate) {
	return func(dsu DataSourceUpdate) {
		f(ClusterHealth{
			Healthy:             dsu.root.Cluster.Data.State.Health,
			Health:              titlify(dsu.root.Cluster.Data.State.Name),
			MinReplicas:         dsu.root.Cluster.Data.State.MinReplicasRemaining,
			RebalanceQueued:     dsu.root.Cluster.Data.MovingData.InQueueBytes,
			RebalanceInFlight:   dsu.root.Cluster.Data.MovingData.InFlightBytes,
			RecoveryState:       titlify(dsu.root.Cluster.RecoveryState.Name),
			RecoveryDescription: dsu.root.Cluster.RecoveryState.Description,
		})
	}
}

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
		return convert(float64(h.RebalanceQueued), 1, None)
	},
}

var StatRebalanceInflight = components.ColumnImpl[ClusterHealth]{
	ColName: "Rebalance In-flight",
	DataFn: func(h ClusterHealth) string {
		return convert(float64(h.RebalanceInFlight), 1, None)
	},
}

var StatEmpty = components.ColumnImpl[ClusterHealth]{
	ColName: "",
	DataFn: func(h ClusterHealth) string {
		return ""
	},
}
