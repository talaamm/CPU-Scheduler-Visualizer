package algorithms

import (
	"reflect"
	"testing"

	al "github.com/talaamm/cpu-scheduler-visualizer/algorithms"
	"github.com/talaamm/cpu-scheduler-visualizer/core"
)

// Reference trace verified against tests/run-example/round_robin.txt (Q=2).
func TestRoundRobin_ClassicTextbook(t *testing.T) {
	processes := []core.Process{
		{PID: "P1", ArrivalTime: 0, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 5}, {Type: core.IOBurst, Duration: 3}, {Type: core.CPUBurst, Duration: 4}}},
		{PID: "P2", ArrivalTime: 1, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 3}, {Type: core.IOBurst, Duration: 2}, {Type: core.CPUBurst, Duration: 2}}},
		{PID: "P3", ArrivalTime: 2, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 2}, {Type: core.IOBurst, Duration: 1}, {Type: core.CPUBurst, Duration: 3}}},
	}
	expected := []segment{
		{0, 2, "P1"}, {2, 4, "P2"}, {4, 6, "P3"}, {6, 8, "P1"}, {8, 9, "P2"},
		{9, 11, "P3"}, {11, 12, "P1"}, {12, 14, "P2"}, {14, 15, "P3"}, {15, 19, "P1"},
	}

	result := al.NewRoundRobin(2).Run(processes)
	got := cpuSegments(result)
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("segments mismatch\nwant: %#v\n got: %#v", expected, got)
	}
}

// A process preempted by quantum expiry must go to the back of the ready
// queue, not cut in front of a process that arrived during its slice.
func TestRoundRobin_PreemptedProcessGoesToBackOfQueue(t *testing.T) {
	processes := []core.Process{
		{PID: "P1", ArrivalTime: 0, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 4}}},
		{PID: "P2", ArrivalTime: 1, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 2}}},
	}
	// P1 runs [0:2), P2 becomes ready at t=1 (before P1's quantum ends at
	// t=2) so it must run next, then P1 resumes — not the other way around.
	expected := []segment{{0, 2, "P1"}, {2, 4, "P2"}, {4, 6, "P1"}}

	result := al.NewRoundRobin(2).Run(processes)
	got := cpuSegments(result)
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("segments mismatch\nwant: %#v\n got: %#v", expected, got)
	}
}
