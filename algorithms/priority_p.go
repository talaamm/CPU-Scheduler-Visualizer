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
	return simulation.Run(processes, r.Name(), simulation.SchedulingPolicy{
		SelectNext: func(s *simulation.Simulation) *core.Process {
			return highestPriority(s.ReadyQueue)
		},
		ShouldPreempt: func(s *simulation.Simulation) *core.Process {
			next := highestPriority(s.ReadyQueue)
			if next == nil || s.RunningProcess == nil {
				return nil
			}
			if next.Priority < s.RunningProcess.Priority {
				return next
			}
			return nil
		},
	})
}
