package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/data/fdb"
	"github.com/pwood/fdbexplorer/ui/components"
	"strings"
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
		return fmt.Sprintf("%0.1f MiB", float64(h.RebalanceQueued)/Mebibyte)
	},
}

var StatRebalanceInflight = components.ColumnImpl[ClusterHealth]{
	ColName: "Rebalance In-flight",
	DataFn: func(h ClusterHealth) string {
		return fmt.Sprintf("%0.1f MiB", float64(h.RebalanceInFlight)/Mebibyte)
	},
}

var StatEmpty = components.ColumnImpl[ClusterHealth]{
	ColName: "",
	DataFn: func(h ClusterHealth) string {
		return ""
	},
}
