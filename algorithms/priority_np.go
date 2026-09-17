package algorithms

import (
	"strconv"

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
	return simulation.Run(processes, r.Name(), simulation.SchedulingPolicy{
		SelectNext: func(s *simulation.Simulation) *core.Process {
			return highestPriority(s.ReadyQueue)
		},
		Explain: func(s *simulation.Simulation, selected, previous *core.Process) string {
			return selected.PID + " has the highest priority (value " + strconv.Itoa(selected.Priority) + ", lower = higher) among ready processes"
		},
	})
}
