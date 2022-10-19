package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var statsHeaders = []string{"Tx Started", "Tx Committed", "Tx Conflicted", "Tx Rejected", "Reads", "Writes", "Bytes Read", "Bytes Written"}

type ClusterStatsTableContent struct {
	tview.TableContentReadOnly
	cd *ClusterData
}

func (v *ClusterStatsTableContent) GetCell(row, column int) *tview.TableCell {
	cs := v.cd.Stats()

	text := ""

	switch column {
	case 0:
		return tview.NewTableCell(statsHeaders[row]).SetExpansion(1).SetTextColor(tcell.ColorYellow)
	case 1:
		switch row {
		case 0:
			text = fmt.Sprintf("%0.1f/s", cs.TxStarted)
		case 1:
			text = fmt.Sprintf("%0.1f/s", cs.TxCommitted)
		case 2:
			text = fmt.Sprintf("%0.1f/s", cs.TxConflicted)
		case 3:
			text = fmt.Sprintf("%0.1f/s", cs.TxRejected)
		}
	case 2:
		return tview.NewTableCell(statsHeaders[row+4]).SetExpansion(1).SetTextColor(tcell.ColorYellow)
	case 3:
		switch row {
		case 0:
			text = fmt.Sprintf("%0.1f/s", cs.Reads)
		case 1:
			text = fmt.Sprintf("%0.1f/s", cs.Writes)
		case 2:
			text = fmt.Sprintf("%0.1f MiB/s", cs.BytesRead/1024/1024)
		case 3:
			text = fmt.Sprintf("%0.1f MiB/s", cs.BytesWritten/1024/1024)
		}
	}

	return tview.NewTableCell(text).SetExpansion(2)
}

func (v *ClusterStatsTableContent) GetRowCount() int {
	return 4
}

func (v *ClusterStatsTableContent) GetColumnCount() int {
	return 4
}
