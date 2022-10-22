package common

import (
	"encoding/json"
	"fmt"
	"github.com/pwood/fdbexplorer/data/fdb"
	"io"
)

func ParseJSON(r io.Reader) (fdb.Root, error) {
	cs := fdb.Root{}

	if err := json.NewDecoder(r).Decode(&cs); err != nil {
		return fdb.Root{}, fmt.Errorf("failed to parse state: %w", err)
	}

	return cs, nil
}
