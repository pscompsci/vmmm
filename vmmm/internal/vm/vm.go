package vm

type VirtualMachine struct {
	ID              int    `db:"vm_id" json:"id,omitempty"`
	Name            string `db:"name"  json:"name"`
	Parent          string `db:"parent" json:"parent"`
	Network         string `db:"network" json:"network"`
	OperatingSystem string `db:"operating_system" json:"operatingSystem"`
	IPAddress       string `db:"ipAddress" json:"ipAddress"`
	CPU             int32  `db:"cpu" json:"cpu"`
	Memory          int32  `db:"memory" json:"memory"`
	DiskType        string `db:"disk_type" json:"diskType"`
	DiskCapacity    int32  `db:"disk_capacity" json:"diskCapacity"`
	DiskFreeSpace   int32  `db:"disk_freespace" json:"diskFreeSpace"`
	State           string `db:"state" json:"state,omitempty"`
	OverallStatus   string `db:"overall_status" json:"overallStatus"`
}
