package algorithms

import (
	"github.com/talaamm/cpu-scheduler-visualizer/core"
	"github.com/talaamm/cpu-scheduler-visualizer/simulation"
)

type PriorityPreempScheduler struct{}

func NewPriorityPreempScheduler() *PriorityPreempScheduler {
	return &PriorityPreempScheduler{}
}

func (r *PriorityPreempScheduler) Name() string {
	return "Priority_Preemptive"
}

func (r *PriorityPreempScheduler) Run(processes []core.Process) simulation.SimulationResult {
	var result simulation.SimulationResult

	// TODO:
	// Priority Preemptive implementation

	return result
}
