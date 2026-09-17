package algorithms

import (
	"strconv"

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
		Explain: func(s *simulation.Simulation, selected, previous *core.Process) string {
			if previous != nil {
				return selected.PID + " preempted " + previous.PID + " — priority " + strconv.Itoa(selected.Priority) + " is higher than " + strconv.Itoa(previous.Priority) + " (lower value = higher priority)"
			}
			return selected.PID + " has the highest priority (value " + strconv.Itoa(selected.Priority) + ", lower = higher) among ready processes"
		},
	})
}
