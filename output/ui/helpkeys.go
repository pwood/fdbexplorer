package ui

import (
	"fmt"
	"github.com/rivo/tview"
)

var helpKeyText = []string{"Sort", "Snapshot", "Interval", "-", "Refresh", "-", "Include", "Exclude"}

type HelpKeys struct {
	tview.TableContentReadOnly

	sorter   *SortControl
	interval *IntervalControl
	haveEM   bool
}

func (h *HelpKeys) GetCell(_, column int) *tview.TableCell {
	text := ""

	switch column {
	case 0:
		text = fmt.Sprintf("%s (%s)", helpKeyText[column], h.sorter.SortName())
	case 2:
		text = fmt.Sprintf("%s (%s)", helpKeyText[column], h.interval.Duration().String())
	case 6, 7:
		if h.haveEM {
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
