package views

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/output/ui/components"
	"github.com/pwood/fdbexplorer/output/ui/data/process"
)

type ClusterHealth struct {
	Healthy     bool
	Health      string
	MinReplicas int

	RebalanceInFlight int
	RebalanceQueued   int

	RecoveryState       string
	RecoveryDescription string

	DatabaseLocked bool
}

func UpdateClusterHealth(f func(ClusterHealth)) func(process.Update) {
	return func(dsu process.Update) {
		f(ClusterHealth{
			Healthy:             dsu.Root.Cluster.Data.State.Health,
			Health:              Titlify(dsu.Root.Cluster.Data.State.Name),
			MinReplicas:         dsu.Root.Cluster.Data.State.MinReplicasRemaining,
			RebalanceQueued:     dsu.Root.Cluster.Data.MovingData.InQueueBytes,
			RebalanceInFlight:   dsu.Root.Cluster.Data.MovingData.InFlightBytes,
			RecoveryState:       Titlify(dsu.Root.Cluster.RecoveryState.Name),
			RecoveryDescription: dsu.Root.Cluster.RecoveryState.Description,
			DatabaseLocked:      dsu.Root.Cluster.DatabaseLockState.Locked,
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
		return Convert(float64(h.RebalanceQueued), 1, None)
	},
}

var StatRebalanceInflight = components.ColumnImpl[ClusterHealth]{
	ColName: "Rebalance In-flight",
	DataFn: func(h ClusterHealth) string {
		return Convert(float64(h.RebalanceInFlight), 1, None)
	},
}

var StatDatabaseLocked = components.ColumnImpl[ClusterHealth]{
	ColName: "Database Locked",
	DataFn: func(h ClusterHealth) string {
		return Boolify(h.DatabaseLocked)
	},
	ColorFn: func(h ClusterHealth) tcell.Color {
		if h.DatabaseLocked {
			return tcell.ColorRed
		} else {
			return tcell.ColorWhite
		}
	},
}

var StatEmpty = components.ColumnImpl[ClusterHealth]{
	ColName: "",
	DataFn: func(h ClusterHealth) string {
		return ""
	},
}
