package output

import (
	"github.com/pwood/fdbexplorer/data"
	"github.com/pwood/fdbexplorer/output/http"
	"github.com/pwood/fdbexplorer/output/ui"
)

type Output interface {
	Run()
}

func Select(ch chan data.State) Output {
	if out, ok := http.NewHTTP(ch); ok {
		return out
	}

	return ui.New(ch)
}
