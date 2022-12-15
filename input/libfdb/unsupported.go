//go:build !(cgo && amd64 && (linux || darwin))

package libfdb

import (
	"encoding/json"
)

func NewFDB() (*FDB, bool) {
	return nil, false
}

type FDB struct {
}

func (f *FDB) Status() (json.RawMessage, error) {
	return nil, nil
}
