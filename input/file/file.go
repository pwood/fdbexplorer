package file

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

var inputFile *string

func init() {
	inputFile = flag.String("input-file", "", "Location of an output of 'status json' to explore, will not connect to FoundationDB.")
}

func NewFile() (*File, bool) {
	if len(*inputFile) == 0 {
		return nil, false
	}

	return &File{fn: *inputFile}, true
}

type File struct {
	fn string
}

func (f *File) Status() (json.RawMessage, error) {
	file, err := os.Open(f.fn)
	defer func(f *os.File) {
		_ = f.Close()
	}(file)

	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}

	d, err := io.ReadAll(file)

	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %w", err)
	}

	return d, nil
}
