package main

import (
	"flag"
	"github.com/pwood/fdbexplorer/data"
	"github.com/pwood/fdbexplorer/data/source"
	"github.com/pwood/fdbexplorer/ui"
	"os"
	"time"
)

func main() {
	defaultClusterFile, found := os.LookupEnv("FDB_CLUSTER_FILE")
	if !found {
		defaultClusterFile = "/etc/foundationdb/fdb.cluster"
	}

	clusterFile := flag.String("cluster-file", defaultClusterFile, "Location of FoundationDB cluster file.")
	interval := flag.Duration("interval", 10*time.Second, "Interval for polling FoundationDB for status.")
	inputFile := flag.String("input-file", "", "Location of an output of 'status json' to explore, will not connect to FoundationDB.")

	flag.Parse()

	stateCh := make(chan data.State)
	defer close(stateCh)

	if len(*inputFile) > 0 {
		go source.NewFile(stateCh, *inputFile).Run()
	} else {
		go source.NewFDB(stateCh, *clusterFile, *interval).Run()
	}

	ui.New(stateCh).Run()
}
