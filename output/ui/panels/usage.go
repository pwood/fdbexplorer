package panels

import (
	"slices"

	"github.com/pwood/fdbexplorer/output/ui/components"
	"github.com/pwood/fdbexplorer/output/ui/data/process"
	"github.com/pwood/fdbexplorer/output/ui/views"
	"github.com/rivo/tview"
)

type UsagePanel struct {
	table   *tview.Table
	content *components.DataTable[process.Process]
}

func NewUsage(store *process.Store, localities []components.ColumnDef[process.Process]) *UsagePanel {
	content := components.NewDataTable[process.Process](
		slices.Concat(
			[]components.ColumnDef[process.Process]{
				views.ColumnSelected, views.ColumnIPAddressPort,
			},
			localities,
			[]components.ColumnDef[process.Process]{
				views.ColumnRoles, views.ColumnCPUActivity, views.ColumnRAMUsage, views.ColumnNetworkActivity,
				views.ColumnDiskUsage, views.ColumnDiskActivity,
			},
		))

	store.AddNotifiable(content.Update, views.All)

	table := tview.NewTable().SetContent(content).SetFixed(1, 0).SetSelectable(true, false)
	table.SetInputCapture(handleNodeSelection(table, content, store))

	return &UsagePanel{table: table, content: content}
}

func (p *UsagePanel) Root() tview.Primitive { return p.table }
func (p *UsagePanel) Update(process.Update) {}
