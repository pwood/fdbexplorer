package panels

import (
	"github.com/gdamore/tcell/v2"
	"github.com/pwood/fdbexplorer/data/fdb"
	"github.com/pwood/fdbexplorer/output/ui/components"
	"github.com/pwood/fdbexplorer/output/ui/data/process"
	"github.com/pwood/fdbexplorer/output/ui/views"
	"github.com/rivo/tview"
)

type DRBackupsPanel struct {
	container            *tview.Flex
	instancesContent     *components.DataTable[fdb.DRBackupInstance]
	tagsContent          *components.DataTable[fdb.DRBackupTag]
	destInstancesContent *components.DataTable[fdb.DRBackupInstance]
	destTagsContent      *components.DataTable[fdb.DRBackupTag]
}

func NewDRBackups() *DRBackupsPanel {
	instanceColumns := []components.ColumnDef[fdb.DRBackupInstance]{
		views.ColumnDRBackupInstanceId, views.ColumnDRBackupInstanceLastUpdated,
		views.ColumnDRBackupInstanceProcessCPU, views.ColumnDRBackupInstanceMemoryUsage,
		views.ColumnDRBackupInstanceResidentSize, views.ColumnDRBackupInstanceVersion,
	}
	tagColumns := []components.ColumnDef[fdb.DRBackupTag]{
		views.ColumnDRBackupTagId, views.ColumnDRBackupTagRunning, views.ColumnDRBackupTagRestorable,
		views.ColumnDRBackupTagSecondsBehind, views.ColumnDRBackupTagState,
		views.ColumnDRBackupTagRangeBytes, views.ColumnDRBackupTagLogBytes, views.ColumnDRBackupTagMutationStream,
	}

	instancesContent := components.NewDataTable[fdb.DRBackupInstance](instanceColumns)
	tagsContent := components.NewDataTable[fdb.DRBackupTag](tagColumns)
	destInstancesContent := components.NewDataTable[fdb.DRBackupInstance](instanceColumns)
	destTagsContent := components.NewDataTable[fdb.DRBackupTag](tagColumns)

	sourceFlex := tview.NewFlex()
	sourceFlex.SetDirection(tview.FlexRow)
	sourceFlex.AddItem(tview.NewTable().SetContent(instancesContent).SetFixed(1, 0).SetSelectable(false, false), 0, 2, false)
	sourceFlex.AddItem(tview.NewTable().SetContent(tagsContent).SetFixed(1, 0).SetSelectable(false, false), 0, 1, false)
	sourceFlex.SetBorder(true)
	sourceFlex.SetTitle("'Source' (Local) Cluster")
	sourceFlex.SetTitleColor(tcell.ColorAqua)

	destFlex := tview.NewFlex()
	destFlex.SetDirection(tview.FlexRow)
	destFlex.AddItem(tview.NewTable().SetContent(destInstancesContent).SetFixed(1, 0).SetSelectable(false, false), 0, 2, false)
	destFlex.AddItem(tview.NewTable().SetContent(destTagsContent).SetFixed(1, 0).SetSelectable(false, false), 0, 1, false)
	destFlex.SetBorder(true)
	destFlex.SetTitle("'Destination' (Remote) Cluster")
	destFlex.SetTitleColor(tcell.ColorAqua)

	container := tview.NewFlex()
	container.SetDirection(tview.FlexRow)
	container.AddItem(sourceFlex, 0, 1, false)
	container.AddItem(destFlex, 0, 1, false)

	return &DRBackupsPanel{
		container:            container,
		instancesContent:     instancesContent,
		tagsContent:          tagsContent,
		destInstancesContent: destInstancesContent,
		destTagsContent:      destTagsContent,
	}
}

func (p *DRBackupsPanel) Root() tview.Primitive { return p.container }

func (p *DRBackupsPanel) Update(u process.Update) {
	views.UpdateDrBackupInstances(p.instancesContent.Update)(u)
	views.UpdateDrBackupTags(p.tagsContent.Update)(u)
	views.UpdateDrBackupDestInstances(p.destInstancesContent.Update)(u)
	views.UpdateDrBackupDestTags(p.destTagsContent.Update)(u)
}
