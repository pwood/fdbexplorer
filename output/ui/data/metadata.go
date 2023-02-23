package data

import (
	"github.com/pwood/fdbexplorer/data/fdb"
	"github.com/pwood/fdbexplorer/output/ui/views"
)

type MetadataStore struct {
	Metadata map[string]*views.ProcessMetadata
}

func (m *MetadataStore) Update(f func([]views.ProcessData)) func(DataSourceUpdate) {
	return func(dsu DataSourceUpdate) {
		var processes []views.ProcessData

		for _, proc := range dsu.Root.Cluster.Processes {
			meta := m.findOrCreateMetadata(proc.Address)
			meta.Update(proc)
			meta.ExclusionInProgress = false
			processes = append(processes, views.ProcessData{Process: proc, Metadata: meta})
		}

		for _, excluding := range dsu.ExclusionInProgress {
			meta := m.findOrCreateMetadata(excluding)
			meta.ExclusionInProgress = true
		}

		for _, excluded := range dsu.ExcludedProcesses {
			found := false

			for _, p := range processes {
				if p.Process.Address == excluded {
					found = true
					break
				}
			}

			if !found {
				processes = append(processes, views.ProcessData{Process: fdb.Process{Address: excluded, Excluded: true}, Metadata: &views.ProcessMetadata{Health: views.HealthExcludedOnly}})
			}
		}

		f(processes)
	}
}

func (m *MetadataStore) ClearSelected() {
	for _, pm := range m.Metadata {
		pm.Selected = false
	}
}

func (m *MetadataStore) findOrCreateMetadata(id string) *views.ProcessMetadata {
	pd, ok := m.Metadata[id]

	if !ok {
		pd = &views.ProcessMetadata{}
		m.Metadata[id] = pd
	}

	return pd
}
