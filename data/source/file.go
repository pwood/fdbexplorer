package source

import (
	"fmt"
	"github.com/pwood/fdbexplorer/data"
	"os"
	"time"
)

func NewFile(ch chan data.State, fn string) *File {
	return &File{ch: ch, fn: fn}
}

type File struct {
	ch chan data.State
	fn string
}

func (f *File) Run() {
	start := time.Now()

	file, err := os.Open(f.fn)
	defer func(f *os.File) {
		_ = f.Close()
	}(file)

	if err != nil {
		f.ch <- data.State{
			Status: fmt.Sprintf("failed to open input file: %s", err.Error()),
		}
		return
	}

	cs, err := parseFDBStatusJSON(file)

	if err != nil {
		f.ch <- data.State{
			Status: err.Error(),
		}
		return
	}

	f.ch <- data.State{
		Status:       "Successfully read input file.",
		Duration:     time.Now().Sub(start),
		Live:         false,
		Interval:     0,
		ClusterState: cs,
	}
}
