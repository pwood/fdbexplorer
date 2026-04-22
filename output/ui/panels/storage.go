package panels

import (
	"github.com/pwood/fdbexplorer/output/ui/components"
	"github.com/pwood/fdbexplorer/output/ui/data/process"
	"github.com/pwood/fdbexplorer/output/ui/views"
	"github.com/rivo/tview"
)

type StoragePanel struct {
	table   *tview.Table
	content *components.DataTable[process.Process]
}

func NewStorage(store *process.Store) *StoragePanel {
	content := components.NewDataTable[process.Process](
		[]components.ColumnDef[process.Process]{
			views.ColumnSelected, views.ColumnIPAddressPort, views.ColumnCPUActivity,
			views.ColumnRAMUsage, views.ColumnDiskUsage, views.ColumnDiskActivity,
			views.ColumnKVStorage, views.ColumnStorageDurabilityRate,
			views.ColumnStorageLag, views.ColumnStorageTotalQueries,
		})

	store.AddNotifiable(content.Update, views.RoleMatch("storage"))

	table := tview.NewTable().SetContent(content).SetFixed(1, 0).SetSelectable(true, false)
	table.SetInputCapture(handleNodeSelection(table, content, store))

	return &StoragePanel{table: table, content: content}
}

func (p *StoragePanel) Root() tview.Primitive { return p.table }
func (p *StoragePanel) Update(process.Update) {}
