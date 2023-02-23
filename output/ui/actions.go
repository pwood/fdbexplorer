package ui

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/input"
	"github.com/pwood/fdbexplorer/output/ui/data/process"
	"github.com/pwood/fdbexplorer/output/ui/views"
)

func manageProcesses(em input.ExclusionManager, s *process.Store, include bool) error {
	selectedProcesses := s.FilterFetch(views.Selected)

	for _, p := range selectedProcesses {
		if include {
			if err := em.IncludeProcess(p.FDBData.Address); err != nil {
				return err
			}
		} else {
			if err := em.ExcludeProcess(p.FDBData.Address); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m *Main) rootAction(event *tcell.EventKey) *tcell.EventKey {
	switch event.Key() {
	case tcell.KeyLeft:
		m.slideShow.Prev()
	case tcell.KeyRight:
		m.slideShow.Next()
	case tcell.KeyF1:
		m.sorter.Next()
		m.processStore.Sort()
	case tcell.KeyF2:
		if filename, err := m.snapshotData(); err != nil {
			m.updateStatus(fmt.Sprintf("Failed to write snapshot: %s", err.Error()), StatusFailure)
		} else {
			m.updateStatus(fmt.Sprintf("Snapshot written: %s", filename), StatusSuccess)
		}
	case tcell.KeyF3:
		m.interval.Next()
	case tcell.KeyF5:
		m.upCh <- struct{}{}
	case tcell.KeyF7:
		if m.em != nil {
			if err := manageProcesses(m.em, m.processStore, true); err != nil {
				m.updateStatus(fmt.Sprintf("Failed to include processes: %s", err.Error()), StatusFailure)
			}
		}
	case tcell.KeyF8:
		if m.em != nil {
			if err := manageProcesses(m.em, m.processStore, false); err != nil {
				m.updateStatus(fmt.Sprintf("Failed to exclude processes: %s", err.Error()), StatusFailure)
			}
		}
	case tcell.KeyESC:
		m.app.Stop()
	case tcell.KeyCtrlL:
		go m.app.Draw()
	case tcell.KeyRune:
		switch event.Rune() {
		case '\\':
			m.processStore.ClearSelected()
			m.processStore.Sort()
		default:
			return event
		}
	default:
		return event
	}
	return nil
}
