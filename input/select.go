package input

import (
	"flag"
	"github.com/pwood/fdbexplorer/data"
	"github.com/pwood/fdbexplorer/input/file"
	"github.com/pwood/fdbexplorer/input/libfdb"
	"time"
)

type Input interface {
	Run()
}

func Select(ch chan data.State) Input {
	interval := flag.Duration("interval", 10*time.Second, "Interval for polling FoundationDB for status.")

	flag.Parse()

	if src, ok := file.NewFile(ch); ok {
		return src
	}

	if src, ok := libfdb.NewFDB(ch, *interval); ok {
		return src
	}

	return nil
}
