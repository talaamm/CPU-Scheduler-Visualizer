package algorithms

import (
	"github.com/talaamm/cpu-scheduler-visualizer/core"
	"github.com/talaamm/cpu-scheduler-visualizer/simulation"
)

type SRTFscheduler struct{}

func NewSRTFscheduler() *SRTFscheduler {
	return &SRTFscheduler{}
}

func (r *SRTFscheduler) Name() string {
	return "Shortest_Remaining_Time_First"
}

func (r *SRTFscheduler) Run(processes []core.Process) simulation.SimulationResult {
	var result simulation.SimulationResult

	// TODO:
	// SRTF implementation

	return result
}
