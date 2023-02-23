package ui

import (
	"github.com/pwood/fdbexplorer/output/ui/data"
	"strings"
)

type SortControl struct {
	i int
}

func (p *SortControl) Next() {
	p.i++
	if p.i > SortExcluded {
		p.i = 0
	}
}

const (
	SortAddress int = iota
	SortRole
	SortClass
	SortUptime
	SortSelected
	SortExcluded
)

func (p *SortControl) SortName() string {
	switch p.i {
	case SortAddress:
		return "Address"
	case SortRole:
		return "Role"
	case SortClass:
		return "Class"
	case SortUptime:
		return "Uptime"
	case SortSelected:
		return "Selected"
	case SortExcluded:
		return "Excluded"
	default:
		return "Unknown"
	}
}

func (p *SortControl) Sort(i data.Process, j data.Process) bool {
	iKey := i.FDBData.Address
	jKey := j.FDBData.Address

	switch p.i {
	case SortRole:
		iRole := ""
		if len(i.FDBData.Roles) > 0 {
			iRole = i.FDBData.Roles[0].Role
		}

		jRole := ""
		if len(j.FDBData.Roles) > 0 {
			jRole = j.FDBData.Roles[0].Role
		}

		iKey = iRole + iKey
		jKey = jRole + jKey
	case SortClass:
		iKey = i.FDBData.Class + iKey
		jKey = j.FDBData.Class + jKey
	case SortUptime:
		return i.FDBData.Uptime < j.FDBData.Uptime
	case SortSelected:
		if !i.Metadata.Selected {
			iKey = "_" + iKey
		}

		if !j.Metadata.Selected {
			jKey = "_" + jKey
		}
	case SortExcluded:
		if !i.FDBData.Excluded {
			iKey = "_" + iKey
		}

		if !j.FDBData.Excluded {
			jKey = "_" + jKey
		}
	}

	return strings.Compare(iKey, jKey) < 0
}
