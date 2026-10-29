package panels

import (
	"slices"

	"github.com/pwood/fdbexplorer/output/ui/components"
	"github.com/pwood/fdbexplorer/output/ui/data/process"
	"github.com/pwood/fdbexplorer/output/ui/views"
	"github.com/rivo/tview"
)

type LogsPanel struct {
	table   *tview.Table
	content *components.DataTable[process.Process]
}

func NewLogs(store *process.Store, localities []components.ColumnDef[process.Process]) *LogsPanel {
	content := components.NewDataTable[process.Process](
		slices.Concat(
			[]components.ColumnDef[process.Process]{
				views.ColumnSelected, views.ColumnIPAddressPort,
			},
			localities,
			[]components.ColumnDef[process.Process]{
				views.ColumnCPUActivity, views.ColumnRAMUsage, views.ColumnDiskUsage, views.ColumnDiskActivity,
				views.ColumnLogQueueLength, views.ColumnLogDurabilityRate, views.ColumnLogQueueStorage,
			},
		))

	store.AddNotifiable(content.Update, views.RoleMatch("log"))

	table := tview.NewTable().SetContent(content).SetFixed(1, 0).SetSelectable(true, false)
	table.SetInputCapture(handleNodeSelection(table, content, store))

	return &LogsPanel{table: table, content: content}
}

func (p *LogsPanel) Root() tview.Primitive { return p.table }
func (p *LogsPanel) Update(process.Update) {}
