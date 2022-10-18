package main

import (
	"encoding/json"
	"fmt"
	"io"
)

func handleDataParse(r io.Reader) (StatusJSON, error) {
	cs := StatusJSON{}

	if err := json.NewDecoder(r).Decode(&cs); err != nil {
		return StatusJSON{}, fmt.Errorf("failed to parse state: %w", err)
	}

	return cs, nil
}
