package simulation

import "github.com/talaamm/cpu-scheduler-visualizer/core"

type SchedulingPolicy struct {
	SelectNext    func(*Simulation) *core.Process
	ShouldPreempt func(*Simulation) *core.Process // if not given assume no preemption
	Quantum       int                             // if not given assume no time slicing
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
		}
	}

	if s.CPUIdle() && policy.SelectNext != nil { // select next w/out arg  just to check it exists
		next := policy.SelectNext(s) // calls - invokes that function
		if next != nil {
			s.RemoveFromReady(next)
			s.AssignCPU(next)
		}
	}
	s.ExecuteCPU(policy.Quantum) // if quantum is 0 it will just ignore it and run until burst completion or preemption
	s.Time++
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
		p.Completed = true
		p.State = core.StateTerminated
		p.CompletionTime = finishedAt
		s.CompletedProcesses++
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
	} else {
		// go back to ready queue
		p.State = core.StateReady
		p.RemainingBurstTime = next.Duration
		s.EnqueueReady(p)
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

func (s *Simulation) CompleteCurrentBurst(p *core.Process) {
	s.CompleteCurrentBurstAt(p, s.Time+1)
}
