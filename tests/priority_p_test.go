package algorithms

import (
	"reflect"
	"testing"

	al "github.com/talaamm/cpu-scheduler-visualizer/algorithms"
	"github.com/talaamm/cpu-scheduler-visualizer/core"
)

// Reference trace verified against tests/run-example/priority_p.txt.
// Priorities: P1=3, P2=1, P3=2 (lower value = higher priority).
func TestPriorityPreemptive_ClassicTextbook(t *testing.T) {
	processes := []core.Process{
		{PID: "P1", ArrivalTime: 0, Priority: 3, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 5}, {Type: core.IOBurst, Duration: 3}, {Type: core.CPUBurst, Duration: 4}}},
		{PID: "P2", ArrivalTime: 1, Priority: 1, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 3}, {Type: core.IOBurst, Duration: 2}, {Type: core.CPUBurst, Duration: 2}}},
		{PID: "P3", ArrivalTime: 2, Priority: 2, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 2}, {Type: core.IOBurst, Duration: 1}, {Type: core.CPUBurst, Duration: 3}}},
	}
	expected := []segment{
		{0, 1, "P1"}, {1, 4, "P2"}, {4, 6, "P3"}, {6, 8, "P2"}, {8, 11, "P3"},
		{11, 15, "P1"}, {15, 18, "IDLE"}, {18, 22, "P1"},
	}

	result := al.NewPriorityPreempScheduler().Run(processes)
	got := cpuSegments(result)
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("segments mismatch\nwant: %#v\n got: %#v", expected, got)
	}
}

// A higher-priority process arriving mid-burst must immediately preempt
// the running lower-priority process — that's what makes this preemptive.
func TestPriorityPreemptive_PreemptsOnHigherPriorityArrival(t *testing.T) {
	processes := []core.Process{
		{PID: "Low", ArrivalTime: 0, Priority: 5, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 6}}},
		{PID: "High", ArrivalTime: 2, Priority: 1, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 2}}},
	}
	expected := []segment{{0, 2, "Low"}, {2, 4, "High"}, {4, 8, "Low"}}

	result := al.NewPriorityPreempScheduler().Run(processes)
	got := cpuSegments(result)
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("segments mismatch\nwant: %#v\n got: %#v", expected, got)
	}
}
