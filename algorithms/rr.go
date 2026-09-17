package algorithms

import (
	"strconv"

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
	return simulation.Run(processes, r.Name(), simulation.SchedulingPolicy{
		SelectNext: func(s *simulation.Simulation) *core.Process {
			return s.FirstReady()
		},
		Quantum: r.Quantum,
		Explain: func(s *simulation.Simulation, selected, previous *core.Process) string {
			return selected.PID + " is next in the round-robin rotation and receives a " + strconv.Itoa(r.Quantum) + "-unit time quantum"
		},
	})
}
