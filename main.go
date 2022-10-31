package main

import (
	"flag"
	"github.com/pwood/fdbexplorer/data"
	"github.com/pwood/fdbexplorer/input"
	"github.com/pwood/fdbexplorer/output"
	"os"
)

func main() {
	ch := make(chan data.State)
	defer close(ch)

	if in := input.Select(ch); in == nil {
		flag.PrintDefaults()
		os.Exit(1)
	} else {
		go in.Run()
	}

	if out := output.Select(ch); out == nil {
		flag.PrintDefaults()
		os.Exit(1)
	} else {
		out.Run()
	}
}
