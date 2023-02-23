package ui

import (
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
