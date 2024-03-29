package input

import (
	"encoding/json"
	"github.com/pwood/fdbexplorer/input/file"
	"github.com/pwood/fdbexplorer/input/libfdb"
	"github.com/pwood/fdbexplorer/input/url"
)

type StatusProvider interface {
	Status() (json.RawMessage, error)
}

type ExclusionManager interface {
	IncludeProcess(includeKey string) error
	ExcludeProcess(excludeKey string) error
	ExcludedProcesses() ([]string, error)
	ExclusionInProgressProcesses() ([]string, error)
}

func Select() StatusProvider {
	if src, ok := file.NewFile(); ok {
		return src
	}

	if src, ok := url.NewURL(); ok {
		return src
	}

	if src, ok := libfdb.NewFDB(); ok {
		return src
	}

	return nil
}
