package data

import (
	"github.com/pwood/fdbexplorer/data/fdb"
)

type Store struct {
	Metadata map[string]*Metadata
}

func (m *Store) Update(f func([]Process)) func(Update) {
	return func(dsu Update) {
		var processes []Process

		for _, proc := range dsu.Root.Cluster.Processes {
			meta := m.findOrCreate(proc.Address)
			meta.Update(proc)
			meta.ExclusionInProgress = false
			processes = append(processes, Process{FDBData: proc, Metadata: meta})
		}

		for _, excluding := range dsu.ExclusionInProgress {
			meta := m.findOrCreate(excluding)
			meta.ExclusionInProgress = true
		}

		for _, excluded := range dsu.ExcludedProcesses {
			found := false

			for _, p := range processes {
				if p.FDBData.Address == excluded {
					found = true
					break
				}
			}

			if !found {
				processes = append(processes, Process{FDBData: fdb.Process{Address: excluded, Excluded: true}, Metadata: &Metadata{Health: HealthExcludedOnly}})
			}
		}

		f(processes)
	}
}

func (m *Store) ClearSelected() {
	for _, pm := range m.Metadata {
		pm.Selected = false
	}
}

func (m *Store) findOrCreate(id string) *Metadata {
	pd, ok := m.Metadata[id]

	if !ok {
		pd = &Metadata{}
		m.Metadata[id] = pd
	}

	return pd
}
