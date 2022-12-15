package ui

import "strings"

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

func (p *SortControl) Sort(i ProcessData, j ProcessData) bool {
	iKey := i.Process.Address
	jKey := j.Process.Address

	switch p.i {
	case SortRole:
		iRole := ""
		if len(i.Process.Roles) > 0 {
			iRole = i.Process.Roles[0].Role
		}

		jRole := ""
		if len(j.Process.Roles) > 0 {
			jRole = j.Process.Roles[0].Role
		}

		iKey = iRole + iKey
		jKey = jRole + jKey
	case SortClass:
		iKey = i.Process.Class + iKey
		jKey = j.Process.Class + jKey
	case SortUptime:
		return i.Process.Uptime < j.Process.Uptime
	case SortSelected:
		if !i.Metadata.Selected {
			iKey = "_" + iKey
		}

		if !j.Metadata.Selected {
			jKey = "_" + jKey
		}
	case SortExcluded:
		if !i.Process.Excluded {
			iKey = "_" + iKey
		}

		if !j.Process.Excluded {
			jKey = "_" + jKey
		}
	}

	return strings.Compare(iKey, jKey) < 0
}
