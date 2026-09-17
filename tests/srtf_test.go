package algorithms

import (
	"reflect"
	"testing"

	al "github.com/talaamm/cpu-scheduler-visualizer/algorithms"
	"github.com/talaamm/cpu-scheduler-visualizer/core"
)

// Reference trace verified against tests/run-example/srtf.txt.
func TestSRTF_ClassicTextbook(t *testing.T) {
	processes := []core.Process{
		{PID: "P1", ArrivalTime: 0, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 5}, {Type: core.IOBurst, Duration: 3}, {Type: core.CPUBurst, Duration: 4}}},
		{PID: "P2", ArrivalTime: 1, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 3}, {Type: core.IOBurst, Duration: 2}, {Type: core.CPUBurst, Duration: 2}}},
		{PID: "P3", ArrivalTime: 2, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 2}, {Type: core.IOBurst, Duration: 1}, {Type: core.CPUBurst, Duration: 3}}},
	}
	expected := []segment{
		{0, 1, "P1"}, {1, 4, "P2"}, {4, 6, "P3"}, {6, 8, "P2"}, {8, 11, "P3"},
		{11, 15, "P1"}, {15, 18, "IDLE"}, {18, 22, "P1"},
	}

	result := al.NewSRTFscheduler().Run(processes)
	got := cpuSegments(result)
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("segments mismatch\nwant: %#v\n got: %#v", expected, got)
	}
}

// A process with a shorter remaining burst arriving mid-execution must
// preempt the running process immediately (unlike non-preemptive SJF).
func TestSRTF_PreemptsOnShorterArrival(t *testing.T) {
	processes := []core.Process{
		{PID: "P1", ArrivalTime: 0, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 8}}},
		{PID: "P2", ArrivalTime: 2, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 2}}},
	}
	expected := []segment{{0, 2, "P1"}, {2, 4, "P2"}, {4, 10, "P1"}}

	result := al.NewSRTFscheduler().Run(processes)
	got := cpuSegments(result)
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("segments mismatch\nwant: %#v\n got: %#v", expected, got)
	}
}
