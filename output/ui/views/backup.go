package views

import (
	"fmt"
	"github.com/pwood/fdbexplorer/data/fdb"
	"github.com/pwood/fdbexplorer/output/ui/components"
	"github.com/pwood/fdbexplorer/output/ui/data"
)

func UpdateBackupInstances(f func(instance []fdb.BackupInstance)) func(data.Update) {
	return func(dsu data.Update) {
		var instances []fdb.BackupInstance

		for _, instance := range dsu.Root.Cluster.Layers.Backup.Instances {
			instances = append(instances, instance)
		}

		f(instances)
	}
}

var ColumnBackupInstanceId = components.ColumnImpl[fdb.BackupInstance]{
	ColName: "Instance Id",
	DataFn: func(instance fdb.BackupInstance) string {
		return instance.Id
	},
}

var ColumnBackupInstanceVersion = components.ColumnImpl[fdb.BackupInstance]{
	ColName: "Version",
	DataFn: func(instance fdb.BackupInstance) string {
		return instance.Version
	},
}

var ColumnBackupInstanceConfiguredWorkers = components.ColumnImpl[fdb.BackupInstance]{
	ColName: "Workers",
	DataFn: func(instance fdb.BackupInstance) string {
		return fmt.Sprintf("%d", instance.ConfiguredWorkers)
	},
}

var ColumnBackupInstanceUsedMemory = components.ColumnImpl[fdb.BackupInstance]{
	ColName: "RAM Usage",
	DataFn: func(instance fdb.BackupInstance) string {
		return Convert(instance.RSSBytes, 1, None)
	},
}

var ColumnBackupInstanceRecentTransfer = components.ColumnImpl[fdb.BackupInstance]{
	ColName: "Recent Transfer",
	DataFn: func(instance fdb.BackupInstance) string {
		return fmt.Sprintf("%s / %s", Convert(instance.BlobStats.Recent.BytesPerSecond, 1, "s"), Convert(instance.BlobStats.Recent.BytesSent, 1, None))
	},
}

var ColumnBackupInstanceRecentOperations = components.ColumnImpl[fdb.BackupInstance]{
	ColName: "Recent Operations",
	DataFn: func(instance fdb.BackupInstance) string {
		return fmt.Sprintf("%d Succeeded / %d Failed", int(instance.BlobStats.Recent.RequestsSuccessful), int(instance.BlobStats.Recent.RequestsFailed))
	},
}

func UpdateBackupTags(f func(instance []fdb.BackupTag)) func(data.Update) {
	return func(dsu data.Update) {
		var tags []fdb.BackupTag

		for id, tag := range dsu.Root.Cluster.Layers.Backup.Tags {
			tag.Id = id
			tags = append(tags, tag)
		}

		f(tags)
	}
}

var ColumnBackupTagId = components.ColumnImpl[fdb.BackupTag]{
	ColName: "Tag",
	DataFn: func(tag fdb.BackupTag) string {
		return tag.Id
	},
}

var ColumnBackupStatus = components.ColumnImpl[fdb.BackupTag]{
	ColName: "Status",
	DataFn: func(tag fdb.BackupTag) string {
		return Titlify(tag.CurrentStatus)
	},
}

var ColumnBackupRunning = components.ColumnImpl[fdb.BackupTag]{
	ColName: "Running?",
	DataFn: func(tag fdb.BackupTag) string {
		return Boolify(tag.RunningBackup)
	},
}

var ColumnBackupRestorable = components.ColumnImpl[fdb.BackupTag]{
	ColName: "Restorable?",
	DataFn: func(tag fdb.BackupTag) string {
		return Boolify(tag.RunningBackupIsRestorable)
	},
}

var ColumnBackupSecondsBehind = components.ColumnImpl[fdb.BackupTag]{
	ColName: "Seconds Behind",
	DataFn: func(tag fdb.BackupTag) string {
		return fmt.Sprintf("%0.1f", tag.LastRestorableSecondsBehind)
	},
}

var ColumnBackupRestorableVersion = components.ColumnImpl[fdb.BackupTag]{
	ColName: "Restorable Version",
	DataFn: func(tag fdb.BackupTag) string {
		return fmt.Sprintf("%d", tag.LastRestorableVersion)
	},
}

var ColumnBackupRangeBytes = components.ColumnImpl[fdb.BackupTag]{
	ColName: "Range Bytes",
	DataFn: func(tag fdb.BackupTag) string {
		return Convert(float64(tag.RangeBytesWritten), 1, None)
	},
}

var ColumnBackupLogBytes = components.ColumnImpl[fdb.BackupTag]{
	ColName: "Log Bytes",
	DataFn: func(tag fdb.BackupTag) string {
		return Convert(float64(tag.MutationLogBytesWritten), 1, None)
	},
}
