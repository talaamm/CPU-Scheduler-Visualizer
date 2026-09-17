package algorithms

import (
	"strconv"

	"github.com/talaamm/cpu-scheduler-visualizer/core"
	"github.com/talaamm/cpu-scheduler-visualizer/simulation"
)

type FCFSscheduler struct{}

/*In Go, interfaces are implemented implicitly,
we do not use an implements keyword like in Java
If a struct implements all the methods declared by an interface,
==> that struct is considered to be an implementation of that interface.*/

func NewFCFS() *FCFSscheduler { // constructor
	return &FCFSscheduler{}
}

func (f *FCFSscheduler) Name() string { // method to satisfy Scheduler interface
	return "FCFS"
}

func (f *FCFSscheduler) Run(processes []core.Process) simulation.SimulationResult {

	// FCFS implementation
	return simulation.Run(processes, "FCFS", simulation.SchedulingPolicy{
		SelectNext: func(s *simulation.Simulation) *core.Process {
			return s.FirstReady()
		},
		Explain: func(s *simulation.Simulation, selected, previous *core.Process) string {
			return selected.PID + " is next in arrival order (FCFS), arrived at t=" + strconv.Itoa(selected.ArrivalTime)
		},
	})
}
