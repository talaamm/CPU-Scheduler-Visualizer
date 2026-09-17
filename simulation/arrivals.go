package simulation

import "github.com/talaamm/cpu-scheduler-visualizer/core"

func (s *Simulation) HandleArrivals() {
	for _, p := range s.Processes {
		if p.State == core.StateNew && p.ArrivalTime <= s.Time {
			if len(p.Bursts) == 0 {
				s.finishProcess(p, s.Time)
				continue
			}
			p.CurrentBurstIndex = 0
			p.RemainingBurstTime = p.Bursts[0].Duration
			s.EnqueueReady(p)
			s.emit(EventArrived, p.PID, p.PID+" arrived and entered the ready queue")
		}
	}
}
