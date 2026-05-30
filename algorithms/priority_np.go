package algorithms

import (
	"github.com/talaamm/cpu-scheduler-visualizer/core"
	"github.com/talaamm/cpu-scheduler-visualizer/simulation"
)

type PriorityNONprScheduler struct{}

func NewPriorityNONprScheduler() *PriorityNONprScheduler {
	return &PriorityNONprScheduler{}
}

func (r *PriorityNONprScheduler) Name() string {
	return "Priority_Non_Preemptive"
}

func (r *PriorityNONprScheduler) Run(processes []core.Process) simulation.SimulationResult {
	var result simulation.SimulationResult

	// TODO:
	// Priority Non-Preemptive implementation

	return result
}
