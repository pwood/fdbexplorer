package main

import (
	"fmt"
	"os"
	"time"
)

func handleDataFile(ch chan State, inputFile string) {
	start := time.Now()

	f, err := os.Open(inputFile)
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	if err != nil {
		ch <- State{
			Status: fmt.Sprintf("failed to open input file: %s", err.Error()),
		}
		return
	}

	cs, err := handleDataParse(f)

	if err != nil {
		ch <- State{
			Status: err.Error(),
		}
		return
	}

	ch <- State{
		Status:       "Successfully read input file.",
		Duration:     time.Now().Sub(start),
		Live:         false,
		Interval:     0,
		ClusterState: cs,
	}
}
