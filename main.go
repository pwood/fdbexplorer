package main

import (
	"flag"
	"fmt"
	"github.com/carlmjohnson/versioninfo"
	"github.com/pwood/fdbexplorer/data"
	"github.com/pwood/fdbexplorer/input"
	"github.com/pwood/fdbexplorer/output"
	"os"
)

func main() {
	header()

	ch := make(chan data.State)
	defer close(ch)

	if in := input.Select(ch); in == nil {
		usage()
	} else {
		go in.Run()
	}

	if out := output.Select(ch); out == nil {
		usage()
	} else {
		out.Run()
	}
}

func header() {
	fmt.Printf("fdbexplorer %s (%s)\n\n", versioninfo.Version, versioninfo.Short())
}

func usage() {
	flag.PrintDefaults()
	os.Exit(1)
}
