package data

import (
	"time"
)

type State struct {
	Err      error
	Duration time.Duration
	Interval time.Duration
	Data     []byte
}
