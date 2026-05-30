package algorithms

import (
	"reflect"
	"testing"

	al "github.com/talaamm/cpu-scheduler-visualizer/algorithms"
	"github.com/talaamm/cpu-scheduler-visualizer/core"
	"github.com/talaamm/cpu-scheduler-visualizer/simulation"
)

type segment struct {
	start int
	end   int
	pid   string
}

func cpuSegments(result simulation.SimulationResult) []segment {
	var segments []segment
	for i, entry := range result.Timeline {
		end := result.TotalTime
		if i+1 < len(result.Timeline) {
			end = result.Timeline[i+1].Time
		}
		// if entry.ProcessID == "IDLE" {
		// 	continue
		// }
		segments = append(segments, segment{
			start: entry.Time,
			end:   end,
			pid:   entry.ProcessID,
		})
	}
	return segments
}

func TestFCFS_Table(t *testing.T) {
	cases := []struct {
		name      string
		processes []core.Process
		expected  []segment
	}{
		{
			name: "WithIOBursts",
			processes: []core.Process{
				{PID: "P1", ArrivalTime: 0, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 5}, {Type: core.IOBurst, Duration: 3}, {Type: core.CPUBurst, Duration: 4}}},
				{PID: "P2", ArrivalTime: 1, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 3}, {Type: core.IOBurst, Duration: 2}, {Type: core.CPUBurst, Duration: 2}}},
				{PID: "P3", ArrivalTime: 2, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 2}, {Type: core.IOBurst, Duration: 1}, {Type: core.CPUBurst, Duration: 3}}},
			},
			expected: []segment{{0, 5, "P1"}, {5, 8, "P2"}, {8, 10, "P3"}, {10, 14, "P1"}, {14, 16, "P2"}, {16, 19, "P3"}},
		},
		{
			name: "IdleBeforeFirstArrival",
			processes: []core.Process{
				{PID: "P1", ArrivalTime: 5, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 3}}},
			},
			expected: []segment{{0, 5, "IDLE"}, {5, 8, "P1"}},
		},
		{
			name: "CPUIdleDuringIO",
			processes: []core.Process{
				{PID: "P1", ArrivalTime: 0, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 2}, {Type: core.IOBurst, Duration: 5}, {Type: core.CPUBurst, Duration: 2}}},
				{PID: "P2", ArrivalTime: 0, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 2}, {Type: core.IOBurst, Duration: 5}, {Type: core.CPUBurst, Duration: 2}}},
			},
			expected: []segment{{0, 2, "P1"}, {2, 4, "P2"}, {4, 7, "IDLE"}, {7, 9, "P1"}, {9, 11, "P2"}},
		},
		{
			name: "SameArrivalOrder",
			processes: []core.Process{
				{PID: "P1", ArrivalTime: 0, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 3}}},
				{PID: "P2", ArrivalTime: 0, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 2}}},
				{PID: "P3", ArrivalTime: 0, Bursts: []core.Burst{{Type: core.CPUBurst, Duration: 1}}},
			},
			expected: []segment{{0, 3, "P1"}, {3, 5, "P2"}, {5, 6, "P3"}},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result := al.NewFCFS().Run(tc.processes)
			got := cpuSegments(result)
			if !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("%s: segments mismatch\nwant: %#v\n got: %#v", tc.name, tc.expected, got)
			}
		})
	}
}

/*
--- PASS: TestFCFS_Table (0.00s)
    --- PASS: TestFCFS_Table/WithIOBursts (0.00s)
    --- PASS: TestFCFS_Table/CPUIdleDuringIO (0.00s)
    --- PASS: TestFCFS_Table/SameArrivalOrder (0.00s)
    --- PASS: TestFCFS_Table/IdleBeforeFirstArrival (0.00s)
PASS
ok      github.com/talaamm/cpu-scheduler-visualizer/tests        0.313s
*/
