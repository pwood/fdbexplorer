package main

import (
	"flag"
	"github.com/pwood/fdbexplorer/data"
	"github.com/pwood/fdbexplorer/data/source"
	"github.com/pwood/fdbexplorer/ui"
	"os"
)

func main() {
	ch := make(chan data.State)
	defer close(ch)

	if src := source.Select(ch); src == nil {
		flag.PrintDefaults()
		os.Exit(1)
	} else {
		go src.Run()
	}

	ui.New(ch).Run()
}
