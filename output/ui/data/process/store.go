package process

import (
	"github.com/pwood/fdbexplorer/data/fdb"
	"sort"
)

type notifiable struct {
	updateFn func([]Process)
	filterFn func(Process) bool
}
type Store struct {
	notifiables []notifiable
	sortFn      func(i Process, j Process) bool

	store map[string]*Metadata
	data  []Process
}

func NewStore(sortFn func(Process, Process) bool) *Store {
	return &Store{
		store:  make(map[string]*Metadata),
		sortFn: sortFn,
	}
}
func (m *Store) AddNotifiable(updateFn func([]Process), filterFn func(Process) bool) {
	m.notifiables = append(m.notifiables, notifiable{
		updateFn: updateFn,
		filterFn: filterFn,
	})
}

func (m *Store) Update(u Update) {
	var processes []Process

	for _, proc := range u.Root.Cluster.Processes {
		meta := m.findOrCreate(proc.Address)
		meta.Update(proc)
		meta.ExclusionInProgress = false
		processes = append(processes, Process{FDBData: proc, Metadata: meta})
	}

	for _, excluding := range u.ExclusionInProgress {
		meta := m.findOrCreate(excluding)
		meta.ExclusionInProgress = true
	}

	for _, excluded := range u.ExcludedProcesses {
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

	m.data = processes
	m.notify()
}

func (m *Store) Sort() {
	m.notify()
}

func (m *Store) notify() {
	sort.Slice(m.data, func(i, j int) bool {
		return m.sortFn(m.data[i], m.data[j])
	})

	for _, n := range m.notifiables {
		var processes []Process

		for _, p := range m.data {
			if n.filterFn(p) {
				processes = append(processes, p)
			}
		}

		n.updateFn(processes)
	}
}

func (m *Store) ClearSelected() {
	for _, pm := range m.store {
		pm.Selected = false
	}
}

func (m *Store) findOrCreate(id string) *Metadata {
	pd, ok := m.store[id]

	if !ok {
		pd = &Metadata{}
		m.store[id] = pd
	}

	return pd
}
