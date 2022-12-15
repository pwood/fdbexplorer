package ui

import "time"

type IntervalControl struct {
	i int
}

var durations = []time.Duration{5 * time.Second, 3 * time.Second, 1 * time.Second, 10 * time.Second}

func (i *IntervalControl) Next() {
	i.i++
	if i.i >= len(durations) {
		i.i = 0
	}
}

func (i *IntervalControl) Duration() time.Duration {
	return durations[i.i]
}
