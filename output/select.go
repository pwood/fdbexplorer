package output

import (
	"github.com/pwood/fdbexplorer/input"
	"github.com/pwood/fdbexplorer/output/http"
	"github.com/pwood/fdbexplorer/output/ui"
)

type Output interface {
	Run()
}

func Select(ds input.StatusProvider) Output {
	if out, ok := http.NewHTTP(ds); ok {
		return out
	}

	return ui.New(ds)
}
