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

	store        map[string]*Process
	storeTouched map[string]struct{}

	data []*Process
}

func NewStore(sortFn func(Process, Process) bool) *Store {
	return &Store{
		store:  make(map[string]*Process),
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
	m.storeTouched = make(map[string]struct{})

	for _, proc := range u.Root.Cluster.Processes {
		p, _ := m.findOrCreate(proc.Address)
		p.FDBData = &proc
		p.Metadata.Update(proc)
		p.Metadata.ExclusionInProgress = false
	}

	for _, excluding := range u.ExclusionInProgress {
		p, _ := m.findOrCreate(excluding)
		p.Metadata.ExclusionInProgress = true
	}

	for _, excluded := range u.ExcludedProcesses {
		p, created := m.findOrCreate(excluded)

		if created {
			p.FDBData.Excluded = true
			p.Metadata.Health = HealthExcludedOnly
		}
	}

	var nd []*Process

	for addr := range m.storeTouched {
		nd = append(nd, m.store[addr])
	}

	m.data = nd

	m.notify()
}

func (m *Store) Sort() {
	m.notify()
}

func (m *Store) notify() {
	sort.Slice(m.data, func(i, j int) bool {
		return m.sortFn(*m.data[i], *m.data[j])
	})

	for _, n := range m.notifiables {
		var processes []Process

		for _, p := range m.data {
			if n.filterFn(*p) {
				processes = append(processes, *p)
			}
		}

		n.updateFn(processes)
	}
}

func (m *Store) ClearSelected() {
	for _, pm := range m.store {
		pm.Metadata.Selected = false
	}
}

func (m *Store) FilterFetch(fn func(process Process) bool) []Process {
	var processes []Process

	for _, p := range m.store {
		if fn(*p) {
			processes = append(processes, *p)
		}
	}

	return processes
}

func (m *Store) findOrCreate(id string) (*Process, bool) {
	pd, ok := m.store[id]

	if !ok {
		pd = &Process{
			FDBData:  &fdb.Process{Address: id},
			Metadata: &Metadata{},
		}
		m.store[id] = pd
	}

	m.storeTouched[id] = struct{}{}

	return pd, !ok
}
