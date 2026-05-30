package simulation

import "github.com/talaamm/cpu-scheduler-visualizer/core"

func (s *Simulation) EnqueueReady(p *core.Process) {
	if p == nil || p.Completed {
		return
	}
	for _, queued := range s.ReadyQueue {
		if queued == p {
			return
		}
	}
	p.State = core.StateReady
	s.ReadyQueue = append(s.ReadyQueue, p)
}

func (s *Simulation) DequeueReady() *core.Process {
	if len(s.ReadyQueue) == 0 {
		return nil
	}
	p := s.ReadyQueue[0]
	s.ReadyQueue = s.ReadyQueue[1:]
	return p
}

func (s *Simulation) FirstReady() *core.Process { // used for FCFS, we just want to peek at the first process in the ready queue without removing it
	if len(s.ReadyQueue) == 0 {
		return nil
	}
	return s.ReadyQueue[0]
}

func (s *Simulation) RemoveFromReady(p *core.Process) {
	for i, queued := range s.ReadyQueue {
		if queued == p {
			s.ReadyQueue = append(s.ReadyQueue[:i], s.ReadyQueue[i+1:]...)
			return
		}
	}
}

// this function is used to move processes from pending ready to ready at the end of each step,
// this allows us to avoid modifying the ready queue while we are iterating over it in the scheduling policies
// pending ready queue is used to store processes that are ready but should not be considered
// for scheduling until the next step, this is useful for processes that have just completed IO or have been preempted,
// we want to give them a chance to be scheduled in the next step rather than immediately
func (s *Simulation) FlushPendingReady() {
	pending := s.PendingReadyQueue
	s.PendingReadyQueue = nil
	for _, p := range pending {
		s.EnqueueReady(p)
	}
}

func (s *Simulation) EnqueuePendingReady(p *core.Process) {
	if p == nil || p.Completed {
		return
	}
	for _, queued := range s.PendingReadyQueue {
		if queued == p {
			return
		}
	}
	s.PendingReadyQueue = append(s.PendingReadyQueue, p)
}
