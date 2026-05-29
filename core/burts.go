package core

// defines CPU/IO burst

type BurstType string

const (
	CPUBurst BurstType = "CPU"
	IOBurst  BurstType = "IO"
)

type Burst struct {
	Type     BurstType `json:"type"`
	Duration int       `json:"duration"`
}
