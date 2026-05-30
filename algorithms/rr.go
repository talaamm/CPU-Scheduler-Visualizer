package algorithms

import (
	"github.com/talaamm/cpu-scheduler-visualizer/core"
	"github.com/talaamm/cpu-scheduler-visualizer/simulation"
)

type RoundRobinScheduler struct {
	Quantum int
}

func NewRoundRobin(q int) *RoundRobinScheduler {
	return &RoundRobinScheduler{
		Quantum: q,
	}
}

func (r *RoundRobinScheduler) Name() string {
	return "Round_Robin"
}

func (r *RoundRobinScheduler) Run(processes []core.Process) simulation.SimulationResult {
	var result simulation.SimulationResult

	// TODO:
	// RR implementation

	return result
}
