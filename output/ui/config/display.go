package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/pwood/fdbexplorer/data/fdb"
)

type Column struct {
	Title string `json:"title"`
	Name  string `json:"name"`
}

type Display struct {
	Columns []Column `json:"columns"`
}

var DefaultColumns = []Column{
	{Title: "Machine", Name: fdb.LocalityMachineID},
	{Title: "Data Hall", Name: fdb.LocalityDataHall},
	{Title: "DC", Name: fdb.LocalityDataCenter},
}

func LoadDisplay() (Display, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return Display{}, fmt.Errorf("config directory: %w", err)
	}

	path := filepath.Join(configDir, "fdbexplorer", "display.json")

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Display{}, nil
	} else if err != nil {
		return Display{}, fmt.Errorf("read %s: %w", path, err)
	}

	var d Display
	if err := json.Unmarshal(data, &d); err != nil {
		return Display{}, fmt.Errorf("parse %s: %w", path, err)
	}

	return d, nil
}

func (d Display) LocalityColumns() []Column {
	return slices.Concat(d.Columns, DefaultColumns)
}
