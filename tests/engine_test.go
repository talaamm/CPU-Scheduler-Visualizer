package algorithms

import (
	"testing"

	al "github.com/talaamm/cpu-scheduler-visualizer/algorithms"
	"github.com/talaamm/cpu-scheduler-visualizer/core"
)

// Regression test: a process whose LAST burst is I/O (rather than CPU) must
// still be marked completed when that I/O finishes. Previously the engine
// silently dropped such a process without incrementing CompletedProcesses,
// which made the simulation loop forever.
func TestProcessEndingInIOBurst_Completes(t *testing.T) {
	processes := []core.Process{
		{PID: "P1", ArrivalTime: 0, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 3}, {Type: core.IOBurst, Duration: 2}}},
	}

	result := al.NewFCFS().Run(processes)

	if len(result.Processes) != 1 {
		t.Fatalf("expected 1 process in result, got %d", len(result.Processes))
	}
	p := result.Processes[0]
	if !p.Completed {
		t.Fatalf("expected process to be marked completed")
	}
	if p.State != core.StateTerminated {
		t.Fatalf("expected state TERMINATED, got %v", p.State)
	}
	if p.CompletionTime != 5 { // CPU 3 (t=0..3), IO 2 (t=3..5)
		t.Fatalf("expected completion_time=5, got %d", p.CompletionTime)
	}
}
