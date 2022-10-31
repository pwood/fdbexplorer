package output

import (
	"github.com/pwood/fdbexplorer/data"
	"github.com/pwood/fdbexplorer/output/ui"
)

type Output interface {
	Run()
}

func Select(ch chan data.State) Output {
	return ui.New(ch)
}
