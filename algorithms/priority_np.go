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
	return simulation.Run(processes, r.Name(), simulation.SchedulingPolicy{
		SelectNext: func(s *simulation.Simulation) *core.Process {
			return highestPriority(s.ReadyQueue)
		},
	})
}
