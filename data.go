package main

import (
	"encoding/json"
	"fmt"
	"github.com/pwood/fdbexplorer/statusjson"
	"io"
)

func handleDataParse(r io.Reader) (statusjson.Root, error) {
	cs := statusjson.Root{}

	if err := json.NewDecoder(r).Decode(&cs); err != nil {
		return statusjson.Root{}, fmt.Errorf("failed to parse state: %w", err)
	}

	return cs, nil
}
