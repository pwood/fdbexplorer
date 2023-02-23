package process

import "github.com/pwood/fdbexplorer/data/fdb"

type Update struct {
	Root                fdb.Root
	ExcludedProcesses   []string
	ExclusionInProgress []string
}

type Process struct {
	FDBData  *fdb.Process
	Metadata *Metadata
}

type Health int

const (
	HealthCritical Health = iota
	HealthWarning
	HealthNormal
	HealthExcluded
	HealthExcludedOnly
)

type Metadata struct {
	Health              Health
	Selected            bool
	ExclusionInProgress bool
}

func (m *Metadata) ToggleSelected() {
	m.Selected = !m.Selected
}

func (m *Metadata) Update(proc fdb.Process) {
	m.Health = HealthNormal

	if proc.Excluded || proc.UnderMaintenance {
		m.Health = HealthExcluded
	}

	if len(proc.Messages) > 0 {
		m.Health = HealthWarning
	}

	if proc.Degraded {
		m.Health = HealthCritical
	}
}
