package algorithms

import (
	"github.com/talaamm/cpu-scheduler-visualizer/core"
	"github.com/talaamm/cpu-scheduler-visualizer/simulation"
)

type SJFscheduler struct{}

func NewSJFscheduler() *SJFscheduler {
	return &SJFscheduler{}
}

func (r *SJFscheduler) Name() string {
	return "Shortest_Job_First"
}

func (r *SJFscheduler) Run(processes []core.Process) simulation.SimulationResult {
	var result simulation.SimulationResult

	// TODO:
	// SJF implementation

	return result
}
