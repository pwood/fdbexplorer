package panels

import (
	"github.com/pwood/fdbexplorer/data/fdb"
	"github.com/pwood/fdbexplorer/output/ui/components"
	"github.com/pwood/fdbexplorer/output/ui/data/process"
	"github.com/pwood/fdbexplorer/output/ui/views"
	"github.com/rivo/tview"
)

type BackupsPanel struct {
	flex             *tview.Flex
	instancesContent *components.DataTable[fdb.BackupInstance]
	tagsContent      *components.DataTable[fdb.BackupTag]
}

func NewBackups() *BackupsPanel {
	instancesContent := components.NewDataTable[fdb.BackupInstance](
		[]components.ColumnDef[fdb.BackupInstance]{
			views.ColumnBackupInstanceId, views.ColumnBackupInstanceVersion,
			views.ColumnBackupInstanceConfiguredWorkers, views.ColumnBackupInstanceUsedMemory,
			views.ColumnBackupInstanceRecentTransfer, views.ColumnBackupInstanceRecentOperations,
		})

	tagsContent := components.NewDataTable[fdb.BackupTag](
		[]components.ColumnDef[fdb.BackupTag]{
			views.ColumnBackupTagId, views.ColumnBackupStatus, views.ColumnBackupRunning,
			views.ColumnBackupRestorable, views.ColumnBackupSecondsBehind,
			views.ColumnBackupRestorableVersion, views.ColumnBackupRangeBytes, views.ColumnBackupLogBytes,
		})

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(tview.NewTable().SetContent(instancesContent).SetFixed(1, 0).SetSelectable(false, false), 0, 1, false)
	flex.AddItem(tview.NewTable().SetContent(tagsContent).SetFixed(1, 0).SetSelectable(false, false), 0, 1, false)

	return &BackupsPanel{flex: flex, instancesContent: instancesContent, tagsContent: tagsContent}
}

func (p *BackupsPanel) Root() tview.Primitive { return p.flex }

func (p *BackupsPanel) Update(u process.Update) {
	views.UpdateBackupInstances(p.instancesContent.Update)(u)
	views.UpdateBackupTags(p.tagsContent.Update)(u)
}
