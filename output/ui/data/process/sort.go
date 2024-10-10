package process

import (
	"net/netip"
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

func (p *SortControl) Sort(i Process, j Process) bool {
	iKey := ""
	jKey := ""

	switch p.i {
	case SortRole:
		if len(i.FDBData.Roles) > 0 {
			iKey = i.FDBData.Roles[0].Role
		}

		if len(j.FDBData.Roles) > 0 {
			jKey = j.FDBData.Roles[0].Role
		}
	case SortClass:
		iKey = i.FDBData.Class
		jKey = j.FDBData.Class
	case SortUptime:
		return i.FDBData.Uptime < j.FDBData.Uptime
	case SortSelected:
		if !i.Metadata.Selected {
			iKey = "_"
		}

		if !j.Metadata.Selected {
			jKey = "_"
		}
	case SortExcluded:
		if !i.FDBData.Excluded {
			iKey = "_"
		}

		if !j.FDBData.Excluded {
			jKey = "_"
		}
	}

	comp := strings.Compare(iKey, jKey)

	if comp < 0 {
		return true
	} else if comp > 0 {
		return false
	} else {
		iAddrPort, _ := netip.ParseAddrPort(i.FDBData.Address)
		jAddrPort, _ := netip.ParseAddrPort(j.FDBData.Address)

		return iAddrPort.Compare(jAddrPort) < 0
	}
}
