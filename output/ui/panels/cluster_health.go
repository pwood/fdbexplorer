package panels

import (
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/output/ui/components"
	"github.com/pwood/fdbexplorer/output/ui/data/process"
	"github.com/pwood/fdbexplorer/output/ui/views"
	"github.com/rivo/tview"
)

type ClusterHealthPanel struct {
	flex    *tview.Flex
	content *components.StatsGrid[views.ClusterHealth]
}

func NewClusterHealth() *ClusterHealthPanel {
	content := components.NewStatsGrid([][]components.ColumnDef[views.ClusterHealth]{
		{views.StatClusterHealth, views.StatRebalanceQueued},
		{views.StatReplicasRemaining, views.StatRebalanceInflight},
		{views.StatRecoveryState, views.StatEmpty},
		{views.StatRecoveryDescription, views.StatDatabaseLocked},
	})

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.SetBorderPadding(0, 0, 1, 1)
	flex.AddItem(tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("Cluster Health").SetTextColor(tcell.ColorAqua), 1, 1, false)
	flex.AddItem(tview.NewTable().SetContent(content).SetSelectable(false, false), 0, 1, false)

	return &ClusterHealthPanel{flex: flex, content: content}
}

func (p *ClusterHealthPanel) Root() tview.Primitive { return p.flex }

func (p *ClusterHealthPanel) Update(u process.Update) {
	views.UpdateClusterHealth(p.content.Update)(u)
}
