package file

import (
	"flag"
	"fmt"
	"github.com/pwood/fdbexplorer/data"
	"io"
	"os"
	"time"
)

var inputFile *string

func init() {
	inputFile = flag.String("input-file", "", "Location of an output of 'status json' to explore, will not connect to FoundationDB.")
}

func NewFile(ch chan data.State) (*File, bool) {
	if len(*inputFile) == 0 {
		return nil, false
	}

	return &File{ch: ch, fn: *inputFile}, true
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
			Err: fmt.Errorf("failed to open input file: %w", err),
		}
		return
	}

	d, err := io.ReadAll(file)

	if err != nil {
		f.ch <- data.State{
			Err: fmt.Errorf("failed to read input file: %w", err),
		}
		return
	}

	f.ch <- data.State{
		Duration: time.Now().Sub(start),
		Interval: 0,
		Data:     d,
	}
}
