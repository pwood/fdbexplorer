package panels

import (
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/output/ui/components"
	"github.com/pwood/fdbexplorer/output/ui/data/process"
	"github.com/pwood/fdbexplorer/output/ui/views"
	"github.com/rivo/tview"
)

type ClusterWorkloadPanel struct {
	flex    *tview.Flex
	content *components.StatsGrid[views.ClusterStats]
}

func NewClusterWorkload() *ClusterWorkloadPanel {
	content := components.NewStatsGrid([][]components.ColumnDef[views.ClusterStats]{
		{views.StatTxStarted, views.StatReads},
		{views.StatTxCommitted, views.StatWrites},
		{views.StatTxConflicted, views.StatBytesRead},
		{views.StatTxRejected, views.StatBytesWritten},
	})

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.SetBorderPadding(0, 0, 1, 1)
	flex.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Cluster Workload").SetTextColor(tcell.ColorAqua), 1, 1, false)
	flex.AddItem(tview.NewTable().SetContent(content).SetSelectable(false, false), 0, 1, false)

	return &ClusterWorkloadPanel{flex: flex, content: content}
}

func (p *ClusterWorkloadPanel) Root() tview.Primitive { return p.flex }

func (p *ClusterWorkloadPanel) Update(u process.Update) {
	views.UpdateClusterStats(p.content.Update)(u)
}
