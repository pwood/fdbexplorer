package main

import (
	"flag"
	"fmt"
	"github.com/pwood/fdbexplorer/data"
	"github.com/pwood/fdbexplorer/input"
	"github.com/pwood/fdbexplorer/output"
	"os"
	"runtime/debug"
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
	bi, _ := debug.ReadBuildInfo()
	version := bi.Main.Version
	sum := bi.Main.Sum

	fmt.Printf("fdbexplorer %s (%s)\n\n", version, sum)
}

func usage() {
	flag.PrintDefaults()
	os.Exit(1)
}
