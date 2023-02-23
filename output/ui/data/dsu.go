package data

import "github.com/pwood/fdbexplorer/data/fdb"

type DataSourceUpdate struct {
	Root                fdb.Root
	ExcludedProcesses   []string
	ExclusionInProgress []string
}
