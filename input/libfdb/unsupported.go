//go:build !(cgo && amd64 && (linux || darwin))

package libfdb

import (
	"github.com/pwood/fdbexplorer/data"
	"time"
)

func NewFDB(_ chan data.State, _ time.Duration) (*FDB, bool) {
	return nil, false
}

type FDB struct {
}

func (f *FDB) Run() {
}
