package main

import (
	"flag"
	"time"
)

type State struct {
	Status       string
	Duration     time.Duration
	Interval     time.Duration
	Live         bool
	ClusterState StatusJSON
}

func main() {
	clusterFile := flag.String("cluster-file", "/etc/foundationdb/fdb.cluster", "Location of FoundationDB cluster file.")
	interval := flag.Duration("interval", 10*time.Second, "Interval for polling FoundationDB for status.")
	inputFile := flag.String("input-file", "", "Location of an output of 'status json' to explore, will not connect to FoundationDB.")

	flag.Parse()

	stateCh := make(chan State)
	defer close(stateCh)

	if len(*inputFile) > 0 {
		go handleDataLocal(stateCh, *inputFile)
	} else {
		go handleDataRemote(stateCh, *clusterFile, *interval)
	}

	v := View{ch: stateCh}
	v.run()
}
