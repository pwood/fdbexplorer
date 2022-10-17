package main

type StatusJSON struct {
	Cluster Cluster `json:"cluster"`
}

type Cluster struct {
	Processes map[string]Process `json:"processes"`
}

type Process struct {
	Address  string   `json:"address"`
	Excluded bool     `json:"excluded"`
	Locality Locality `json:"locality"`
	Class    string   `json:"class_type"`
	Roles    []Role   `json:"roles"`
	CPU      CPU      `json:"cpu"`
	Disk     Disk     `json:"disk"`
	Memory   Memory   `json:"memory"`
}

type Locality struct {
	DataHall  string `json:"data_hall"`
	DCID      string `json:"dcid"`
	MachineID string `json:"machineid"`
}

type Role struct {
	Role string `json:"role"`
}

type CPU struct {
	UsageCores float64 `json:"usage_cores"`
}

type Disk struct {
	Busy       float64 `json:"busy"`
	FreeBytes  int     `json:"free_bytes"`
	TotalBytes int     `json:"total_bytes"`
}

type Memory struct {
	AvailableBytes int `json:"available_bytes"`
	UsedBytes      int `json:"used_bytes"`
}
