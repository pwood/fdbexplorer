package fdb

type Root struct {
	Cluster Cluster `json:"cluster"`
}

type Cluster struct {
	Processes         map[string]Process `json:"processes"`
	DatabaseAvailable bool               `json:"database_available"`
	Workload          Workload           `json:"workload"`
	Messages          []Message          `json:"messages"`
	RecoveryState     RecoveryState      `json:"recovery_state"`
	Data              Data               `json:"data"`
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

type Health int

const (
	HealthCritical Health = iota
	HealthWarning
	HealthNormal
	HealthExcluded
)

type Process struct {
	Health           Health    `json:"-"`
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

func AnnotateProcessHealth(p Process) Process {
	p.Health = HealthNormal

	if p.Excluded || p.UnderMaintenance {
		p.Health = HealthExcluded
	}

	if len(p.Messages) > 0 {
		p.Health = HealthWarning
	}

	if p.Degraded {
		p.Health = HealthCritical
	}

	return p
}

type Message struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Locality struct {
	DataHall  string `json:"data_hall"`
	DCID      string `json:"dcid"`
	MachineID string `json:"machineid"`
}

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
}
