package views

import (
	"fmt"
	"github.com/pwood/fdbexplorer/output/ui/data/process"
	"github.com/rivo/tview"
)

var helpKeyText = []string{"Sort", "Snapshot", "Interval", "-", "Refresh", "-", "Include", "Exclude"}

type HelpKeys struct {
	tview.TableContentReadOnly

	Sorter   *process.SortControl
	Interval *IntervalControl
	HasEM    bool
}

func (h *HelpKeys) GetCell(_, column int) *tview.TableCell {
	text := ""

	switch column {
	case 0:
		text = fmt.Sprintf("%s (%s)", helpKeyText[column], h.Sorter.SortName())
	case 2:
		text = fmt.Sprintf("%s (%s)", helpKeyText[column], h.Interval.Duration().String())
	case 6, 7:
		if h.HasEM {
			text = helpKeyText[column]
		} else {
			text = "-"
		}
	default:
		text = fmt.Sprintf("%s", helpKeyText[column])
	}

	return tview.NewTableCell(fmt.Sprintf("F%d [black:darkcyan]%s[:-]", column+1, text))
}

func (h *HelpKeys) GetRowCount() int {
	return 1
}

func (h *HelpKeys) GetColumnCount() int {
	return len(helpKeyText)
}
