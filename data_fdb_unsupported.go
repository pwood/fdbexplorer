//go:build !(cgo && amd64 && (linux || darwin))

package main

import "time"

func handleDataFDB(_ chan State, _ string, _ time.Duration) {
	panic("fdbexplorer compiled without CGO or for platform that has no official FoundationDB library.")
}
