package simulation

import (
	"strconv"

	"github.com/talaamm/cpu-scheduler-visualizer/core"
)

type SchedulingPolicy struct {
	SelectNext    func(*Simulation) *core.Process
	ShouldPreempt func(*Simulation) *core.Process // if not given assume no preemption
	Quantum       int                             // if not given assume no time slicing

	// Explain produces a human-readable reason for why `selected` was given
	// the CPU. `previous` is nil for a normal selection (CPU was idle) and
	// non-nil when `selected` preempted `previous`. Optional — if nil, a
	// generic message is used instead.
	Explain func(s *Simulation, selected *core.Process, previous *core.Process) string
}

// stores the ENTIRE CURRENT STATE OF THE SYSTEM at any moment in time
type Simulation struct {
	Time               int
	Processes          []*core.Process
	ReadyQueue         []*core.Process
	IOQueue            []*IOEntry
	RunningProcess     *core.Process
	PendingReadyQueue  []*core.Process
	Timeline           []TimelineEntry
	Snapshots          []StateSnapshot
	Events             []Event
	CompletedProcesses int
	CPUBusyTime        int
	CurrentQuantumUsed int
}

func NewSimulation(processes []core.Process) *Simulation {
	copied := make([]*core.Process, len(processes))
	for i := range processes {
		p := processes[i]
		p.Bursts = append([]core.Burst(nil), processes[i].Bursts...)
		p.CurrentBurstIndex = 0
		p.RemainingBurstTime = 0
		p.State = core.StateNew
		p.Started = false
		p.Completed = false
		p.StartTime = 0
		p.CompletionTime = 0
		p.WaitingTime = 0
		p.TurnaroundTime = 0
		p.ResponseTime = 0
		copied[i] = &p
	}

	return &Simulation{
		Time:               0,
		Processes:          copied,
		ReadyQueue:         []*core.Process{},
		IOQueue:            []*IOEntry{},
		RunningProcess:     nil, // not yet we just started, we  have to decide
		PendingReadyQueue:  []*core.Process{},
		Timeline:           []TimelineEntry{},
		Snapshots:          []StateSnapshot{},
		Events:             []Event{},
		CompletedProcesses: 0,
		CPUBusyTime:        0,
		CurrentQuantumUsed: 0,
	}
}

func (s *Simulation) AllProcessesCompleted() bool {
	return len(s.Processes) == s.CompletedProcesses
	// for _, p := range s.Processes {
	// 	if !p.Completed {
	// 		return false
	// 	}

	// }
	// return true
}

func (s *Simulation) Step(policy SchedulingPolicy) {
	s.HandleArrivals()
	s.HandleIOCompletion()
	s.FlushPendingReady()

	if s.RunningProcess != nil && policy.ShouldPreempt != nil { // check if we should preempt the current running process
		next := policy.ShouldPreempt(s)
		if next != nil {
			s.RemoveFromReady(next)
			current := s.RunningProcess
			s.RunningProcess = nil
			current.State = core.StateReady // preempt change state from running to ready
			s.EnqueueReady(current)
			s.AssignCPU(next)
			s.emit(EventPreempted, current.PID, s.explain(policy, next, current))
			s.emit(EventSelected, next.PID, s.explain(policy, next, current))
		}
	}

	if s.CPUIdle() && policy.SelectNext != nil { // select next w/out arg  just to check it exists
		next := policy.SelectNext(s) // calls - invokes that function
		if next != nil {
			s.RemoveFromReady(next)
			s.AssignCPU(next)
			s.emit(EventSelected, next.PID, s.explain(policy, next, nil))
		}
	}
	s.captureSnapshot()
	s.ExecuteCPU(policy.Quantum) // if quantum is 0 it will just ignore it and run until burst completion or preemption
	s.Time++
}

func (s *Simulation) explain(policy SchedulingPolicy, selected *core.Process, previous *core.Process) string {
	if policy.Explain != nil {
		return policy.Explain(s, selected, previous)
	}
	if previous != nil {
		return selected.PID + " preempted " + previous.PID
	}
	return selected.PID + " selected to run"
}

func Run(processes []core.Process, algorithm string, policy SchedulingPolicy) SimulationResult {
	sim := NewSimulation(processes)
	for !sim.AllProcessesCompleted() {
		sim.Step(policy)
	}
	return sim.BuildResult(algorithm)
}

func (s *Simulation) CompleteCurrentBurstAt(p *core.Process, finishedAt int) {
	// 1. Move to next burst
	p.CurrentBurstIndex++

	// 2. Check if process is finished
	if p.CurrentBurstIndex >= len(p.Bursts) {
		s.finishProcess(p, finishedAt)
		s.RunningProcess = nil
		s.CurrentQuantumUsed = 0
		return
	}

	// 3. Get next burst
	next := p.Bursts[p.CurrentBurstIndex]

	// 4. Decide what happens next
	if next.Type == "IO" {
		// move to IO queue
		s.MoveToIO(p, next.Duration, finishedAt)
		s.emit(EventIOStart, p.PID, p.PID+" started an I/O burst ("+strconv.Itoa(next.Duration)+" units)")
	} else {
		// go back to ready queue
		p.State = core.StateReady
		p.RemainingBurstTime = next.Duration
		s.EnqueueReady(p)
		s.emit(EventBurstDone, p.PID, p.PID+" finished a CPU burst and returned to the ready queue")
	}

	// 5. CPU becomes free
	s.RunningProcess = nil
	s.CurrentQuantumUsed = 0
	/*Burst ends →
	  either:
	      → IO wait
	      → OR next CPU burst
	      → OR process finished*/
}

// finishProcess marks a process as terminated, regardless of whether it
// finished on a CPU burst or an I/O burst.
func (s *Simulation) finishProcess(p *core.Process, finishedAt int) {
	p.Completed = true
	p.State = core.StateTerminated
	p.CompletionTime = finishedAt
	s.CompletedProcesses++
	s.emit(EventCompleted, p.PID, p.PID+" completed all bursts")
}

func (s *Simulation) CompleteCurrentBurst(p *core.Process) {
	s.CompleteCurrentBurstAt(p, s.Time+1)
}
