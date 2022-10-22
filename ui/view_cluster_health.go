package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var healthHeaders = []string{"Healthy", "Replicas Remaining", "Recovery State", "Recovery Description", "Rebalance Queued", "Rebalance In-flight", "", ""}

type ClusterHealthTableContent struct {
	tview.TableContentReadOnly
	cd *ClusterData
}

func (v *ClusterHealthTableContent) GetCell(row, column int) *tview.TableCell {
	ch := v.cd.Health()

	text := ""
	color := tcell.ColorWhite

	switch column {
	case 0:
		return tview.NewTableCell(healthHeaders[row]).SetExpansion(1).SetTextColor(tcell.ColorYellow)
	case 1:
		switch row {
		case 0:
			text = ch.Health
			if ch.Healthy {
				color = tcell.ColorGreen
			} else {
				color = tcell.ColorRed
			}
		case 1:
			text = fmt.Sprintf("%d", ch.MinReplicas)
		case 2:
			text = ch.RecoveryState
		case 3:
			text = ch.RecoveryDescription
		}
	case 2:
		return tview.NewTableCell(healthHeaders[row+4]).SetExpansion(1).SetTextColor(tcell.ColorYellow)
	case 3:
		switch row {
		case 0:
			text = fmt.Sprintf("%0.1f MiB", float64(ch.RebalanceQueued)/1024/1024)
		case 1:
			text = fmt.Sprintf("%0.1f MiB", float64(ch.RebalanceInFlight)/1024/1024)
		case 2:
		case 3:
		}
	}

	return tview.NewTableCell(text).SetExpansion(2).SetTextColor(color)
}

func (v *ClusterHealthTableContent) GetRowCount() int {
	return 4
}

func (v *ClusterHealthTableContent) GetColumnCount() int {
	return 4
}
