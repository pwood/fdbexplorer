//go:build cgo && amd64 && (linux || darwin)

package main

import (
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"strings"
	"time"
)

func handleDataFDB(ch chan State, clusterFile string, interval time.Duration) {
	fdb.MustAPIVersion(710)

	timer := time.NewTicker(interval)
	db := fdb.MustOpenDatabase(clusterFile)

	nowCh := make(chan struct{}, 1)
	nowCh <- struct{}{}

	for {
		select {
		case <-nowCh:
			pollFDB(db, ch)
		case <-timer.C:
			pollFDB(db, ch)
		}
	}
}

func pollFDB(db fdb.Database, ch chan State) {
	start := time.Now()

	data, err := db.Transact(func(transaction fdb.Transaction) (interface{}, error) {
		return transaction.Get(fdb.Key("\xff\xff/status/json")).MustGet(), nil
	})

	if err != nil {
		ch <- State{
			Status: fmt.Sprintf("foundationdb err: %s", err.Error()),
		}
		return
	}

	cs, err := handleDataParse(strings.NewReader(string(data.([]byte))))

	if err != nil {
		ch <- State{
			Status: err.Error(),
		}
		return
	}

	ch <- State{
		Status:       "Successfully read from FDB.",
		Duration:     time.Now().Sub(start),
		Live:         true,
		Interval:     0,
		ClusterState: cs,
	}
}
