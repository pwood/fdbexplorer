package fdb

type Root struct {
	Cluster Cluster `json:"cluster"`
}

type Cluster struct {
	Processes         map[string]Process `json:"processes"`
	DatabaseAvailable bool               `json:"database_available"`
	DatabaseLockState DatabaseLockState  `json:"database_lock_state"`
	Workload          Workload           `json:"workload"`
	Messages          []Message          `json:"messages"`
	RecoveryState     RecoveryState      `json:"recovery_state"`
	Data              Data               `json:"data"`
	Layers            Layers             `json:"layers"`
}

type DatabaseLockState struct {
	Locked bool `json:"locked"`
}

type Layers struct {
	Backup       Backup   `json:"backup"`
	DRBackup     DRBackup `json:"dr_backup"`
	DRBackupDest DRBackup `json:"dr_backup_dest"`
}

type Backup struct {
	Instances map[string]BackupInstance `json:"instances"`
	Tags      map[string]BackupTag      `json:"tags"`
}

type BackupInstance struct {
	Id                string               `json:"id"`
	BlobStats         BackupBlockStatsBlob `json:"blob_stats"`
	RSSBytes          float64              `json:"resident_size"`
	ConfiguredWorkers int                  `json:"configured_workers"`
	Version           string               `json:"version"`
}

type BackupBlockStatsBlob struct {
	Recent BackupBlobStatsIndividual `json:"recent"`
	Total  BackupBlobStatsIndividual `json:"total"`
}

type BackupBlobStatsIndividual struct {
	BytesPerSecond     float64 `json:"bytes_per_second"`
	BytesSent          float64 `json:"bytes_sent"`
	RequestsFailed     float64 `json:"requests_failed"`
	RequestsSuccessful float64 `json:"requests_successful"`
}

type BackupTag struct {
	Id                          string  `json:"-"`
	CurrentContainer            string  `json:"current_container"`
	CurrentStatus               string  `json:"current_status"`
	LastRestorableSecondsBehind float64 `json:"last_restorable_seconds_behind"`
	LastRestorableVersion       int     `json:"last_restorable_version"`
	MutationLogBytesWritten     int     `json:"mutation_log_bytes_written"`
	RangeBytesWritten           int     `json:"range_bytes_written"`
	RunningBackup               bool    `json:"running_backup"`
	RunningBackupIsRestorable   bool    `json:"running_backup_is_restorable"`
}

type DRBackup struct {
	Instances map[string]DRBackupInstance
	Paused    bool `json:"paused"`
	Tags      map[string]DRBackupTag
}

type DRBackupInstance struct {
	ConfiguredWorkers    int     `json:"configured_workers"`
	Id                   string  `json:"id"`
	LastUpdated          float64 `json:"last_updated"`
	MainThreadCPUSeconds float64 `json:"main_thread_cpu_seconds"`
	MemoryUsage          int     `json:"memory_usage"`
	ProcessCPUSeconds    float64 `json:"process_cpu_seconds"`
	ResidentSize         int     `json:"resident_size"`
	Version              string  `json:"version"`
}

type DRBackupTag struct {
	Id                      string  `json:"-"`
	BackupState             string  `json:"backup_state"`
	MutationSteamId         string  `json:"mutation_stream_id"`
	MutationLogBytesWritten int     `json:"mutation_log_bytes_written"`
	RangeBytesWritten       int     `json:"range_bytes_written"`
	RunningBackup           bool    `json:"running_backup"`
	BackupRestorable        bool    `json:"running_backup_is_restorable"`
	SecondsBehind           float64 `json:"seconds_behind"`
}

type Data struct {
	State      State      `json:"state"`
	MovingData MovingData `json:"moving_data"`
}

type State struct {
	Health               bool   `json:"healthy"`
	Name                 string `json:"name"`
	MinReplicasRemaining int    `json:"min_replicas_remaining"`
}

type MovingData struct {
	InFlightBytes int `json:"in_flight_bytes"`
	InQueueBytes  int `json:"in_queue_bytes"`
}

type RecoveryState struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Workload struct {
	Transactions Transactions `json:"transactions"`
	Operations   Operations   `json:"operations"`
	Bytes        Bytes        `json:"bytes"`
}

type Operations struct {
	Reads  Stats `json:"reads"`
	Writes Stats `json:"writes"`
}

type Transactions struct {
	Committed                Stats `json:"committed"`
	Conflicted               Stats `json:"conflicted"`
	RejectedForQueuedTooLong Stats `json:"rejected_for_queued_too_long"`
	Started                  Stats `json:"started"`
}

type Bytes struct {
	Read    Stats `json:"read"`
	Written Stats `json:"written"`
}

type Process struct {
	Address          string    `json:"address"`
	Degraded         bool      `json:"degraded"`
	Excluded         bool      `json:"excluded"`
	Locality         Locality  `json:"locality"`
	Class            string    `json:"class_type"`
	CommandLine      string    `json:"command_line"`
	Roles            []Role    `json:"roles"`
	CPU              CPU       `json:"cpu"`
	Disk             Disk      `json:"disk"`
	Memory           Memory    `json:"memory"`
	Network          Network   `json:"network"`
	Uptime           float64   `json:"uptime_seconds"`
	Version          string    `json:"version"`
	UnderMaintenance bool      `json:"under_maintenance"`
	Messages         []Message `json:"messages"`
}

type Message struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

const (
	LocalityDataHall   = "data_hall"
	LocalityDataCenter = "dcid"
	LocalityMachineID  = "machineid"
	LocalityProcessID  = "processid"
)

type Locality map[string]string

type Role struct {
	Role string `json:"role"`

	// Storage Only
	KVUsedBytes   float64 `json:"kvstore_used_bytes"`
	TotalQueries  Stats   `json:"total_queries"`
	DataLag       Lag     `json:"data_lag"`
	DurabilityLag Lag     `json:"durability_lag"`

	// Log Only
	QueueUsedBytes float64 `json:"queue_disk_used_bytes"`

	// Both
	InputBytes   Stats `json:"input_bytes"`
	DurableBytes Stats `json:"durable_bytes"`
}

type Lag struct {
	Seconds  float64 `json:"seconds"`
	Versions int     `json:"versions"`
}

type CPU struct {
	UsageCores float64 `json:"usage_cores"`
}

type Stats struct {
	Hz        float64 `json:"hz"`
	Counter   float64 `json:"counter"`
	Roughness float64 `json:"roughness"`
}

type Hz struct {
	Hz float64 `json:"hz"`
}

type Disk struct {
	Busy       float64 `json:"busy"`
	FreeBytes  int     `json:"free_bytes"`
	TotalBytes int     `json:"total_bytes"`
	Reads      Hz      `json:"reads"`
	Writes     Hz      `json:"writes"`
}

type Network struct {
	MegabitsSent     Hz `json:"megabits_sent"`
	MegabitsReceived Hz `json:"megabits_received"`
}

type Memory struct {
	AvailableBytes int `json:"available_bytes"`
	UsedBytes      int `json:"used_bytes"`
	RSSBytes       int `json:"rss_bytes"`
}
