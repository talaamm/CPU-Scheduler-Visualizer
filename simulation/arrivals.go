package simulation

func (s *Simulation) HandleArrivals() {
	for _, p := range s.Processes {
		if p.State == "NEW" && p.ArrivalTime <= s.Time {
			if len(p.Bursts) == 0 {
				p.Completed = true
				p.CompletionTime = s.Time
				s.CompletedProcesses++
				continue
			}
			p.CurrentBurstIndex = 0
			p.RemainingBurstTime = p.Bursts[0].Duration
			s.EnqueueReady(p)
		}
	}
}
