//go:build cgo && ((amd64 && linux) || darwin)

package libfdb

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"os"
	"strings"
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

func (f *FDB) ExcludeProcess(excludeKey string) error {
	if _, err := f.db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		if err := tr.Options().SetSpecialKeySpaceEnableWrites(); err != nil {
			return nil, err
		}

		tr.Set(fdb.Key(fmt.Sprintf("\xff\xff/management/excluded/%s", excludeKey)), []byte{})

		return nil, nil
	}); err != nil {
		return fmt.Errorf("foundationdb err: %w", err)
	} else {
		return nil
	}
}

func (f *FDB) IncludeProcess(includeKey string) error {
	if _, err := f.db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		if err := tr.Options().SetSpecialKeySpaceEnableWrites(); err != nil {
			return nil, err
		}

		tr.Clear(fdb.Key(fmt.Sprintf("\xff\xff/management/excluded/%s", includeKey)))

		return nil, nil
	}); err != nil {
		return fmt.Errorf("foundationdb err: %w", err)
	} else {
		return nil
	}
}

func (f *FDB) ExcludedProcesses() ([]string, error) {
	return f.getProcesses("\xff\xff/management/excluded/")
}

func (f *FDB) ExclusionInProgressProcesses() ([]string, error) {
	return f.getProcesses("\xff\xff/management/in_progress_exclusion/")
}

func (f *FDB) getProcesses(keyPrefix string) ([]string, error) {
	if excluded, err := f.db.Transact(func(tr fdb.Transaction) (interface{}, error) {
		if err := tr.Options().SetAccessSystemKeys(); err != nil {
			return nil, err
		}

		result, err := tr.GetRange(fdb.KeyRange{Begin: fdb.Key(keyPrefix), End: fdb.Key(fmt.Sprintf("%s\xff", keyPrefix))}, fdb.RangeOptions{Mode: fdb.StreamingModeWantAll}).GetSliceWithError()

		if err != nil {
			return nil, err
		}

		var processes []string

		for _, v := range result {
			keyParts := strings.Split(v.Key.String(), "/")
			processes = append(processes, keyParts[len(keyParts)-1])
		}

		return processes, nil
	}); err != nil {
		return nil, fmt.Errorf("foundationdb err: %w", err)
	} else {
		return excluded.([]string), nil
	}
}
