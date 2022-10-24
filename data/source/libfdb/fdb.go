//go:build cgo && amd64 && (linux || darwin)

package libfdb

import (
	"flag"
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/pwood/fdbexplorer/data"
	"os"
	"time"
)

var clusterFile *string

func init() {
	defaultClusterFile, found := os.LookupEnv("FDB_CLUSTER_FILE")
	if !found {
		defaultClusterFile = "/etc/foundationdb/fdb.cluster"
	}

	clusterFile = flag.String("cluster-file", defaultClusterFile, "Location of FoundationDB cluster file, environment variable FDB_CLUSTER_FILE also obeyed.")
}

func NewFDB(ch chan data.State, interval time.Duration) (*FDB, bool) {
	return &FDB{ch: ch, clusterFile: *clusterFile, interval: interval}, true
}

type FDB struct {
	ch          chan data.State
	clusterFile string
	interval    time.Duration
}

func (f *FDB) Run() {
	fdb.MustAPIVersion(710)

	timer := time.NewTicker(f.interval)
	db := fdb.MustOpenDatabase(f.clusterFile)

	nowCh := make(chan struct{}, 1)
	nowCh <- struct{}{}

	for {
		select {
		case <-nowCh:
			f.poll(db, f.ch)
		case <-timer.C:
			f.poll(db, f.ch)
		}
	}
}

func (f *FDB) poll(db fdb.Database, ch chan data.State) {
	start := time.Now()

	d, err := db.Transact(func(transaction fdb.Transaction) (interface{}, error) {
		return transaction.Get(fdb.Key("\xff\xff/status/json")).MustGet(), nil
	})

	if err != nil {
		ch <- data.State{
			Err: fmt.Errorf("foundationdb err: %w", err),
		}
		return
	}

	ch <- data.State{
		Duration: time.Now().Sub(start),
		Interval: f.interval,
		Data:     d.([]byte),
	}
}