package simulation

import "github.com/talaamm/cpu-scheduler-visualizer/core"

func (s *Simulation) BuildResult(algorithm string) SimulationResult {
	processes := make([]core.Process, len(s.Processes))
	var totalWaiting int
	var totalTurnaround int
	var totalResponse int

	for i, p := range s.Processes {
		p.TurnaroundTime = p.CompletionTime - p.ArrivalTime
		p.WaitingTime = p.TurnaroundTime - totalBurstDuration(p)
		processes[i] = *p
		totalWaiting += p.WaitingTime
		totalTurnaround += p.TurnaroundTime
		totalResponse += p.ResponseTime
	}

	count := float64(len(s.Processes))
	if count == 0 {
		count = 1
	}

	totalTime := s.Time
	cpuUtilization := 0.0
	if totalTime > 0 {
		cpuUtilization = float64(s.CPUBusyTime) / float64(totalTime) * 100
	}

	return SimulationResult{
		Algorithm:             algorithm,
		Timeline:              compactTimeline(s.Timeline),
		Processes:             processes,
		TotalTime:             totalTime,
		CPUuilization:         cpuUtilization,
		AverageWaitingTime:    float64(totalWaiting) / count,
		AverageTurnaroundTime: float64(totalTurnaround) / count,
		AverageResponseTime:   float64(totalResponse) / count,
	}
}

func compactTimeline(entries []TimelineEntry) []TimelineEntry {
	if len(entries) == 0 {
		return entries
	}

	compacted := []TimelineEntry{entries[0]}
	for _, entry := range entries[1:] {
		last := &compacted[len(compacted)-1]
		if last.ProcessID == entry.ProcessID && last.Event == entry.Event {
			continue
		}
		compacted = append(compacted, entry)
	}
	return compacted
}

func totalBurstDuration(p *core.Process) int {
	total := 0
	for _, b := range p.Bursts {
		total += b.Duration
	}
	return total
}
