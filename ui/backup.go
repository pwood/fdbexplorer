package ui

import (
	"fmt"
	"github.com/pwood/fdbexplorer/data/fdb"
	"github.com/pwood/fdbexplorer/ui/components"
)

func UpdateBackupInstances(f func(instance []fdb.BackupInstance)) func(fdb.Root) {
	return func(root fdb.Root) {
		var instances []fdb.BackupInstance

		for _, instance := range root.Cluster.Layers.Backup.Instances {
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
		return convert(instance.RSSBytes, 1, false)
	},
}

func UpdateBackupTags(f func(instance []fdb.BackupTag)) func(fdb.Root) {
	return func(root fdb.Root) {
		var tags []fdb.BackupTag

		for id, tag := range root.Cluster.Layers.Backup.Tags {
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
