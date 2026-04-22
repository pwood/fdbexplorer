package panels

import (
	"github.com/pwood/fdbexplorer/output/ui/components"
	"github.com/pwood/fdbexplorer/output/ui/data/process"
	"github.com/pwood/fdbexplorer/output/ui/views"
	"github.com/rivo/tview"
)

type LocalityPanel struct {
	table   *tview.Table
	content *components.DataTable[process.Process]
}

func NewLocality(store *process.Store) *LocalityPanel {
	content := components.NewDataTable[process.Process](
		[]components.ColumnDef[process.Process]{
			views.ColumnSelected, views.ColumnIPAddressPort, views.ColumnTLS, views.ColumnStatus,
			views.ColumnMachine, views.ColumnLocality, views.ColumnClass, views.ColumnRoles,
			views.ColumnVersion, views.ColumnUptime,
		})

	store.AddNotifiable(content.Update, views.All)

	table := tview.NewTable().SetContent(content).SetFixed(1, 0).SetSelectable(true, false)
	table.SetInputCapture(handleNodeSelection(table, content, store))

	return &LocalityPanel{table: table, content: content}
}

func (p *LocalityPanel) Root() tview.Primitive { return p.table }
func (p *LocalityPanel) Update(process.Update) {}
