//go:build cgo && amd64 && (linux || darwin)

package libfdb

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"os"
)

var clusterFile *string

func init() {
	defaultClusterFile, found := os.LookupEnv("FDB_CLUSTER_FILE")
	if !found {
		defaultClusterFile = "/etc/foundationdb/fdb.cluster"
	}

	clusterFile = flag.String("cluster-file", defaultClusterFile, "Location of FoundationDB cluster file, environment variable FDB_CLUSTER_FILE also obeyed.")
}

func NewFDB() (*FDB, bool) {
	fdb.MustAPIVersion(710)

	f := &FDB{clusterFile: *clusterFile}
	f.db = fdb.MustOpenDatabase(f.clusterFile)

	return f, true
}

type FDB struct {
	clusterFile string
	db          fdb.Database
}

func (f *FDB) Status() (json.RawMessage, error) {
	if d, err := f.db.Transact(func(transaction fdb.Transaction) (interface{}, error) {
		return transaction.Get(fdb.Key("\xff\xff/status/json")).MustGet(), nil
	}); err != nil {
		return nil, fmt.Errorf("foundationdb err: %w", err)
	} else {
		return d.([]byte), nil
	}
}
