package simulation

import "github.com/talaamm/cpu-scheduler-visualizer/core"

type SimulationResult struct { // output of simulation, ready for visualization
	Algorithm string `json:"algorithm"`

	Timeline              []TimelineEntry `json:"timeline"`
	Processes             []core.Process  `json:"processes"`
	TotalTime             int             `json:"total_time"`
	CPUuilization         float64         `json:"cpu_utilization"`
	AverageWaitingTime    float64         `json:"average_waiting_time"`
	AverageTurnaroundTime float64         `json:"average_turnaround_time"`
	AverageResponseTime   float64         `json:"average_response_time"`
}
