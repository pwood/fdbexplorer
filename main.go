package main

import (
	"flag"
	"fmt"
	"github.com/carlmjohnson/versioninfo"
	"github.com/pwood/fdbexplorer/input"
	"github.com/pwood/fdbexplorer/output"
	"os"
)

func main() {
	header()

	flag.Parse()

	in := input.Select()
	if in == nil {
		usage()
	}

	if out := output.Select(in); out == nil {
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
