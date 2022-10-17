package main

import (
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"sort"
	"strings"
	"sync"
)

type View struct {
	ch  chan State
	app *tview.Application

	pd *ProcessData
}

func (v *View) runData() {
	for s := range v.ch {
		v.pd.m.Lock()

		v.pd.processes = []Process{}

		for _, process := range s.ClusterState.Cluster.Processes {
			v.pd.processes = append(v.pd.processes, process)
		}

		sort.Slice(v.pd.processes, func(i, j int) bool {
			return strings.Compare(v.pd.processes[i].Address, v.pd.processes[j].Address) < 0
		})

		v.pd.m.Unlock()
		v.app.Draw()
	}
}

func (v *View) run() {
	v.pd = &ProcessData{m: &sync.RWMutex{}}

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	locality := &ProcessView{
		pd:      v.pd,
		columns: []string{"IP Address:Port", "Status", "Machine", "Locality", "Class", "Roles"},
		dataFn: func(process Process, column int) *tview.TableCell {

			data := ""

			switch column {
			case 0:
				data = process.Address
			case 1:
				var status []string

				if process.Excluded {
					status = append(status, "Excluded")
				}

				data = strings.Join(status, " ")
			case 2:
				data = process.Locality.MachineID
			case 3:
				data = fmt.Sprintf("%s / %s", process.Locality.DataHall, process.Locality.DCID)
			case 4:
				data = process.Class
			case 5:
				var roles []string

				for _, role := range process.Roles {
					roles = append(roles, role.Role)
				}

				data = strings.Join(roles, ", ")
			}

			return tview.NewTableCell(data)
		},
	}

	usage := &ProcessView{
		pd:      v.pd,
		columns: []string{"IP Address:Port", "Roles", "CPU Usage", "RAM Usage", "Disk Usage"},
		dataFn: func(process Process, column int) *tview.TableCell {
			data := ""

			switch column {
			case 0:
				data = process.Address
			case 1:
				var roles []string

				for _, role := range process.Roles {
					roles = append(roles, role.Role)
				}

				data = strings.Join(roles, ", ")
			case 2:
				data = fmt.Sprintf("%0.1f%%", process.CPU.UsageCores*100)
			case 3:
				memUsage := float64(process.Memory.UsedBytes) / float64(process.Memory.AvailableBytes)

				used := process.Memory.UsedBytes / 1024 / 1024
				available := process.Memory.AvailableBytes / 1024 / 1024

				data = fmt.Sprintf("%0.1f%% (%d MiB of %d MiB)", memUsage*100, used, available)
			case 4:
				usedBytes := process.Disk.TotalBytes - process.Disk.FreeBytes
				diskUsage := float64(usedBytes) / float64(process.Disk.TotalBytes)

				used := float64(usedBytes) / 1024 / 1024 / 1024
				available := float64(process.Disk.TotalBytes) / 1024 / 1024 / 1024

				data = fmt.Sprintf("%0.1f%% (%0.1f GiB of %0.1f GiB)", diskUsage*100, used, available)
			}

			return tview.NewTableCell(data)
		},
	}

	pages := tview.NewPages()

	localityTable := tview.NewTable().SetContent(locality)
	localityTable.SetFixed(1, 0)
	localityTable.SetSelectable(true, false)
	pages.AddPage("locality", localityTable, true, true)

	usageTable := tview.NewTable().SetContent(usage)
	usageTable.SetFixed(1, 0)
	usageTable.SetSelectable(true, false)
	pages.AddPage("usage", usageTable, true, false)

	pageIndex := []string{"locality", "usage"}
	pageNow := 0

	pages.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyLeft:
			pageNow--
			if pageNow < 0 {
				pageNow = len(pageIndex) - 1
			}
		case tcell.KeyRight:
			pageNow++
			if pageNow >= len(pageIndex) {
				pageNow = 0
			}
		default:
			return event
		}

		pages.SwitchToPage(pageIndex[pageNow])
		return nil
	})

	grid := tview.NewGrid().SetRows(5, 0).SetBorders(true)
	grid.AddItem(newPrimitive("Header"), 0, 0, 1, 1, 0, 0, false)
	grid.AddItem(pages, 1, 0, 1, 1, 0, 0, true)

	v.app = tview.NewApplication().SetRoot(grid, true).SetFocus(localityTable)

	go v.runData()

	if err := v.app.Run(); err != nil {
		panic(err)
	}
}
