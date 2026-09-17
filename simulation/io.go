package simulation

import "github.com/talaamm/cpu-scheduler-visualizer/core"

type IOEntry struct {
	Process        *core.Process
	CompletionTime int
}

// Handles IO completion and moves processes back to ready queue when their IO is done.
// A process whose LAST burst is an I/O burst must be marked completed here —
// otherwise it is silently dropped and the simulation never terminates.
func (s *Simulation) HandleIOCompletion() {
	remaining := s.IOQueue[:0]
	for _, entry := range s.IOQueue {
		if entry.CompletionTime <= s.Time {
			p := entry.Process
			p.CurrentBurstIndex++
			if p.CurrentBurstIndex < len(p.Bursts) {
				p.RemainingBurstTime = p.Bursts[p.CurrentBurstIndex].Duration
				s.EnqueueReady(p)
				s.emit(EventIODone, p.PID, p.PID+" finished I/O and returned to the ready queue")
			} else {
				s.finishProcess(p, s.Time)
			}
			continue
		}
		remaining = append(remaining, entry)
	}
	s.IOQueue = remaining
}

func (s *Simulation) MoveToIO(p *core.Process, d int, startTime ...int) {
	start := s.Time
	if len(startTime) > 0 {
		start = startTime[0]
	}
	p.State = core.StateWaiting
	p.RemainingBurstTime = 0
	s.IOQueue = append(s.IOQueue, &IOEntry{
		Process:        p,
		CompletionTime: start + d,
	})
}
