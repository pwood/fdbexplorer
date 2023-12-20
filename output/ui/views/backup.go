package views

import (
	"fmt"
	"github.com/pwood/fdbexplorer/data/fdb"
	"github.com/pwood/fdbexplorer/output/ui/components"
	"github.com/pwood/fdbexplorer/output/ui/data/process"
	"sort"
	"strings"
	"time"
)

func UpdateBackupInstances(f func(instance []fdb.BackupInstance)) func(process.Update) {
	return func(dsu process.Update) {
		var instances []fdb.BackupInstance

		for _, instance := range dsu.Root.Cluster.Layers.Backup.Instances {
			instances = append(instances, instance)
		}

		sort.Slice(instances, func(i, j int) bool {
			return strings.Compare(instances[i].Id, instances[j].Id) < 0
		})

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

func UpdateBackupTags(f func(instance []fdb.BackupTag)) func(process.Update) {
	return func(dsu process.Update) {
		var tags []fdb.BackupTag

		for id, tag := range dsu.Root.Cluster.Layers.Backup.Tags {
			tag.Id = id
			tags = append(tags, tag)
		}

		sort.Slice(tags, func(i, j int) bool {
			return strings.Compare(tags[i].Id, tags[j].Id) < 0
		})

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

func UpdateDrBackupInstances(f func(instance []fdb.DRBackupInstance)) func(process.Update) {
	return func(dsu process.Update) {
		var instances []fdb.DRBackupInstance

		for _, instance := range dsu.Root.Cluster.Layers.DRBackup.Instances {
			instances = append(instances, instance)
		}

		sort.Slice(instances, func(i, j int) bool {
			return strings.Compare(instances[i].Id, instances[j].Id) < 0
		})

		f(instances)
	}
}

func UpdateDrBackupDestInstances(f func(instance []fdb.DRBackupInstance)) func(process.Update) {
	return func(dsu process.Update) {
		var instances []fdb.DRBackupInstance

		for _, instance := range dsu.Root.Cluster.Layers.DRBackupDest.Instances {
			instances = append(instances, instance)
		}

		sort.Slice(instances, func(i, j int) bool {
			return strings.Compare(instances[i].Id, instances[j].Id) < 0
		})

		f(instances)
	}
}

var ColumnDRBackupInstanceId = components.ColumnImpl[fdb.DRBackupInstance]{
	ColName: "Instance Id",
	DataFn: func(instance fdb.DRBackupInstance) string {
		return instance.Id
	},
}

var ColumnDRBackupInstanceVersion = components.ColumnImpl[fdb.DRBackupInstance]{
	ColName: "Version",
	DataFn: func(instance fdb.DRBackupInstance) string {
		return instance.Version
	},
}

var ColumnDRBackupInstanceLastUpdated = components.ColumnImpl[fdb.DRBackupInstance]{
	ColName: "Last Updated",
	DataFn: func(instance fdb.DRBackupInstance) string {
		return time.Unix(int64(instance.LastUpdated), 0).String()
	},
}

var ColumnDRBackupInstanceMemoryUsage = components.ColumnImpl[fdb.DRBackupInstance]{
	ColName: "Memory Usage",
	DataFn: func(instance fdb.DRBackupInstance) string {
		return Convert(float64(instance.MemoryUsage), 1, None)
	},
}

var ColumnDRBackupInstanceResidentSize = components.ColumnImpl[fdb.DRBackupInstance]{
	ColName: "Resident Size",
	DataFn: func(instance fdb.DRBackupInstance) string {
		return Convert(float64(instance.ResidentSize), 1, None)
	},
}

var ColumnDRBackupInstanceMainThreadCPU = components.ColumnImpl[fdb.DRBackupInstance]{
	ColName: "Main Thread CPU",
	DataFn: func(instance fdb.DRBackupInstance) string {
		return fmt.Sprintf("%0.0fs", instance.MainThreadCPUSeconds)
	},
}

var ColumnDRBackupInstanceProcessCPU = components.ColumnImpl[fdb.DRBackupInstance]{
	ColName: "Process CPU",
	DataFn: func(instance fdb.DRBackupInstance) string {
		return fmt.Sprintf("%0.0fs", instance.ProcessCPUSeconds)
	},
}

func UpdateDrBackupTags(f func(instance []fdb.DRBackupTag)) func(process.Update) {
	return func(dsu process.Update) {
		var instances []fdb.DRBackupTag

		for tag, instance := range dsu.Root.Cluster.Layers.DRBackup.Tags {
			instance.Id = tag
			instances = append(instances, instance)
		}

		sort.Slice(instances, func(i, j int) bool {
			return strings.Compare(instances[i].Id, instances[j].Id) < 0
		})

		f(instances)
	}
}

func UpdateDrBackupDestTags(f func(instance []fdb.DRBackupTag)) func(process.Update) {
	return func(dsu process.Update) {
		var instances []fdb.DRBackupTag

		for tag, instance := range dsu.Root.Cluster.Layers.DRBackupDest.Tags {
			instance.Id = tag
			instances = append(instances, instance)
		}

		sort.Slice(instances, func(i, j int) bool {
			return strings.Compare(instances[i].Id, instances[j].Id) < 0
		})

		f(instances)
	}
}

var ColumnDRBackupTagId = components.ColumnImpl[fdb.DRBackupTag]{
	ColName: "Tag",
	DataFn: func(instance fdb.DRBackupTag) string {
		return instance.Id
	},
}

var ColumnDRBackupTagState = components.ColumnImpl[fdb.DRBackupTag]{
	ColName: "State",
	DataFn: func(instance fdb.DRBackupTag) string {
		return Titlify(instance.BackupState)
	},
}

var ColumnDRBackupTagRunning = components.ColumnImpl[fdb.DRBackupTag]{
	ColName: "Running",
	DataFn: func(instance fdb.DRBackupTag) string {
		return Boolify(instance.RunningBackup)
	},
}

var ColumnDRBackupTagRestorable = components.ColumnImpl[fdb.DRBackupTag]{
	ColName: "Restorable",
	DataFn: func(instance fdb.DRBackupTag) string {
		return Boolify(instance.RunningBackup)
	},
}

var ColumnDRBackupTagSecondsBehind = components.ColumnImpl[fdb.DRBackupTag]{
	ColName: "Seconds Behind",
	DataFn: func(instance fdb.DRBackupTag) string {
		return fmt.Sprintf("%0.1f", instance.SecondsBehind)
	},
}

var ColumnDRBackupTagLogBytes = components.ColumnImpl[fdb.DRBackupTag]{
	ColName: "Log Bytes Written",
	DataFn: func(instance fdb.DRBackupTag) string {
		return Convert(float64(instance.MutationLogBytesWritten), 1, None)
	},
}

var ColumnDRBackupTagRangeBytes = components.ColumnImpl[fdb.DRBackupTag]{
	ColName: "Range Bytes Written",
	DataFn: func(instance fdb.DRBackupTag) string {
		return Convert(float64(instance.RangeBytesWritten), 1, None)
	},
}

var ColumnDRBackupTagMutationStream = components.ColumnImpl[fdb.DRBackupTag]{
	ColName: "Mutation Stream Id",
	DataFn: func(instance fdb.DRBackupTag) string {
		return instance.MutationSteamId
	},
}
