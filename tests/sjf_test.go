package algorithms

import (
	"reflect"
	"testing"

	al "github.com/talaamm/cpu-scheduler-visualizer/algorithms"
	"github.com/talaamm/cpu-scheduler-visualizer/core"
)

// Reference trace verified against tests/run-example/sjf.txt.
func TestSJF_ClassicTextbook(t *testing.T) {
	processes := []core.Process{
		{PID: "P1", ArrivalTime: 0, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 5}, {Type: core.IOBurst, Duration: 3}, {Type: core.CPUBurst, Duration: 4}}},
		{PID: "P2", ArrivalTime: 1, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 3}, {Type: core.IOBurst, Duration: 2}, {Type: core.CPUBurst, Duration: 2}}},
		{PID: "P3", ArrivalTime: 2, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 2}, {Type: core.IOBurst, Duration: 1}, {Type: core.CPUBurst, Duration: 3}}},
	}
	expected := []segment{{0, 5, "P1"}, {5, 7, "P3"}, {7, 10, "P2"}, {10, 13, "P3"}, {13, 15, "P2"}, {15, 19, "P1"}}

	result := al.NewSJFscheduler().Run(processes)
	got := cpuSegments(result)
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("segments mismatch\nwant: %#v\n got: %#v", expected, got)
	}
}

// Two processes become ready while the CPU is busy, in arrival order
// Mid-then-Short. SJF must pick by burst length (Short first), not by
// ready-queue insertion / arrival order (which would pick Mid first).
func TestSJF_SelectsByBurstLengthNotArrivalOrder(t *testing.T) {
	processes := []core.Process{
		{PID: "Long", ArrivalTime: 0, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 10}}},
		{PID: "Mid", ArrivalTime: 1, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 6}}},
		{PID: "Short", ArrivalTime: 2, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 2}}},
	}
	expected := []segment{{0, 10, "Long"}, {10, 12, "Short"}, {12, 18, "Mid"}}

	result := al.NewSJFscheduler().Run(processes)
	got := cpuSegments(result)
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("segments mismatch\nwant: %#v\n got: %#v", expected, got)
	}
}
