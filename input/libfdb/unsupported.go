//go:build cgo && arm64 && linux

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
