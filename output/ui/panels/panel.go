package panels

import (
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/output/ui/components"
	"github.com/pwood/fdbexplorer/output/ui/data/process"
	"github.com/rivo/tview"
)

type Panel interface {
	Root() tview.Primitive
	Update(process.Update)
}

func handleNodeSelection(table *tview.Table, content *components.DataTable[process.Process], store *process.Store) func(*tcell.EventKey) *tcell.EventKey {
	return func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyRune && event.Rune() == ' ' {
			row, _ := table.GetSelection()
			content.Get(row).Metadata.ToggleSelected()
			store.Sort()
			return nil
		}
		return event
	}
}
