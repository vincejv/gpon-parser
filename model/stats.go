package model

type DeviceStats struct {
	MemoryUsage     float64   `json:"memoryUsage"`
	CpuUsage        float64   `json:"cpuUsage"`
	CpuDtlUsage     []float64 `json:"cpuDtlUsage"`
	DeviceModel     string    `json:"deviceModel"`
	ModelSerial     string    `json:"modelSerial"`
	SoftwareVersion string    `json:"softwareVersion"`
	Uptime          int64     `json:"uptime"`
}

type OpticalStats struct {
	RxPower       float64 `json:"rxPower"`
	TxPower       float64 `json:"txPower"`
	Temperature   float64 `json:"temperature"`
	SupplyVoltage float64 `json:"supplyVoltage"`
	BiasCurrent   float64 `json:"biasCurrent"`
}

type AllStats struct {
	DeviceStats  *DeviceStats  `json:"deviceStats"`
	OpticalStats *OpticalStats `json:"opticalStats"`
}
