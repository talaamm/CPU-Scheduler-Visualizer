package algorithms

import "github.com/talaamm/cpu-scheduler-visualizer/core"

func shortestRemaining(ready []*core.Process) *core.Process {
	return bestReady(ready, func(best, candidate *core.Process) bool {
		return candidate.RemainingBurstTime < best.RemainingBurstTime
	})
}

func highestPriority(ready []*core.Process) *core.Process {
	return bestReady(ready, func(best, candidate *core.Process) bool {
		return candidate.Priority < best.Priority
	})
}

func bestReady(ready []*core.Process, better func(best, candidate *core.Process) bool) *core.Process {
	if len(ready) == 0 {
		return nil
	}

	best := ready[0]
	for _, candidate := range ready[1:] {
		if better(best, candidate) {
			best = candidate
		}
	}
	return best
}
