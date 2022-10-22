package data

import (
	"github.com/pwood/fdbexplorer/data/fdb"
	"time"
)

type State struct {
	Status       string
	Duration     time.Duration
	Interval     time.Duration
	Live         bool
	ClusterState fdb.Root
}
