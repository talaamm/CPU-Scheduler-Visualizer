package simulation

import "../core"

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

func (s *Simulation) FirstReady() *core.Process {
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