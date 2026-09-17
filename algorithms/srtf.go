package algorithms

import (
	"strconv"

	"github.com/talaamm/cpu-scheduler-visualizer/core"
	"github.com/talaamm/cpu-scheduler-visualizer/simulation"
)

type SRTFscheduler struct{}

func NewSRTFscheduler() *SRTFscheduler {
	return &SRTFscheduler{}
}

func (r *SRTFscheduler) Name() string {
	return "Shortest_Remaining_Time_First"
}

func (r *SRTFscheduler) Run(processes []core.Process) simulation.SimulationResult {
	return simulation.Run(processes, r.Name(), simulation.SchedulingPolicy{
		SelectNext: func(s *simulation.Simulation) *core.Process {
			return shortestRemaining(s.ReadyQueue)
		},
		ShouldPreempt: func(s *simulation.Simulation) *core.Process {
			next := shortestRemaining(s.ReadyQueue)
			if next == nil || s.RunningProcess == nil {
				return nil
			}
			if next.RemainingBurstTime < s.RunningProcess.RemainingBurstTime {
				return next
			}
			return nil
		},
		Explain: func(s *simulation.Simulation, selected, previous *core.Process) string {
			if previous != nil {
				return selected.PID + " preempted " + previous.PID + " — remaining time " + strconv.Itoa(selected.RemainingBurstTime) + " < " + strconv.Itoa(previous.RemainingBurstTime)
			}
			return selected.PID + " has the shortest remaining time (" + strconv.Itoa(selected.RemainingBurstTime) + " units) among ready processes"
		},
	})
}
