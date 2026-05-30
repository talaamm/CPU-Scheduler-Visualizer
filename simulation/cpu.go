package simulation

import "github.com/talaamm/cpu-scheduler-visualizer/core"

func (s *Simulation) ExecuteCPU(quantum int) {
	p := s.RunningProcess
	if p == nil { // CPU is idle
		s.Timeline = append(s.Timeline, TimelineEntry{
			Time:      s.Time,
			ProcessID: "IDLE",
			Event:     "IDLE",
		})
		return
	}

	if !p.Started {
		p.Started = true
		p.StartTime = s.Time
		p.ResponseTime = s.Time - p.ArrivalTime
	}

	p.State = core.StateRunning
	p.RemainingBurstTime--
	s.CurrentQuantumUsed++
	s.CPUBusyTime++
	s.Timeline = append(s.Timeline, TimelineEntry{
		Time:      s.Time,
		ProcessID: p.PID,
		Event:     "RUNNING",
	})

	if p.RemainingBurstTime == 0 {
		s.CompleteCurrentBurstAt(p, s.Time+1)
		return
	}

	if quantum > 0 && s.CurrentQuantumUsed == quantum { // for rr
		s.RunningProcess = nil
		s.CurrentQuantumUsed = 0
		p.State = core.StateReady
		s.EnqueuePendingReady(p)
	}
}

func (s *Simulation) AssignCPU(p *core.Process) {
	p.State = core.StateRunning
	s.RunningProcess = p
	s.CurrentQuantumUsed = 0
}

func (s *Simulation) CPUIdle() bool {
	return s.RunningProcess == nil
}
