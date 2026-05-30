package algorithms

import (
	"cpu-scheduler/core"
	"cpu-scheduler/simulation"
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
