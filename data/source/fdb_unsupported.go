//go:build !(cgo && amd64 && (linux || darwin))

package source

import (
	"github.com/pwood/fdbexplorer/data"
	"time"
)

func NewFDB(_ chan data.State, _ string, _ time.Duration) *FDB {
	return &FDB{}
}

type FDB struct {
}

func (f *FDB) Run() {
	panic("fdbexplorer compiled without CGO or for platform that has no official FoundationDB library.")
}
