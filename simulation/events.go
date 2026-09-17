package simulation

// EventType identifies what kind of scheduling event occurred at a given tick.
type EventType string

const (
	EventArrived   EventType = "ARRIVED"
	EventSelected  EventType = "SELECTED"
	EventPreempted EventType = "PREEMPTED"
	EventBurstDone EventType = "BURST_DONE" // finished a CPU burst, back in ready queue
	EventIOStart   EventType = "IO_START"
	EventIODone    EventType = "IO_DONE"
	EventCompleted EventType = "COMPLETED"
)

// Event is a human-readable, algorithm-aware record of a scheduling decision
// or state transition, used to explain *why* something happened rather than
// just *that* it happened.
type Event struct {
	Time    int       `json:"time"`
	Type    EventType `json:"type"`
	PID     string    `json:"pid"`
	Message string    `json:"message"`
}

func (s *Simulation) emit(typ EventType, pid string, message string) {
	s.Events = append(s.Events, Event{
		Time:    s.Time,
		Type:    typ,
		PID:     pid,
		Message: message,
	})
}

// IOQueueSnapshotEntry is the visualization-facing view of a process waiting on I/O.
type IOQueueSnapshotEntry struct {
	PID       string `json:"pid"`
	Remaining int    `json:"remaining"`
}

// StateSnapshot captures the full, real system state for a single time tick:
// exactly what a visualizer needs to render the CPU, ready queue, and I/O
// queue without having to reverse-engineer it from the compacted timeline.
type StateSnapshot struct {
	Time              int                    `json:"time"`
	Running           string                 `json:"running"`             // "" means CPU is idle
	RunningRemaining  int                    `json:"running_remaining"`   // remaining time in the running process's current CPU burst
	RunningBurstTotal int                    `json:"running_burst_total"` // total duration of that burst, for progress %
	ReadyQueue        []string               `json:"ready_queue"`
	IOQueue           []IOQueueSnapshotEntry `json:"io_queue"`
}

func (s *Simulation) captureSnapshot() {
	running := ""
	remaining := 0
	burstTotal := 0
	if s.RunningProcess != nil {
		running = s.RunningProcess.PID
		remaining = s.RunningProcess.RemainingBurstTime
		burstTotal = s.RunningProcess.Bursts[s.RunningProcess.CurrentBurstIndex].Duration
	}

	ready := make([]string, len(s.ReadyQueue))
	for i, p := range s.ReadyQueue {
		ready[i] = p.PID
	}

	ioq := make([]IOQueueSnapshotEntry, len(s.IOQueue))
	for i, e := range s.IOQueue {
		ioq[i] = IOQueueSnapshotEntry{PID: e.Process.PID, Remaining: e.CompletionTime - s.Time}
	}

	s.Snapshots = append(s.Snapshots, StateSnapshot{
		Time:              s.Time,
		Running:           running,
		RunningRemaining:  remaining,
		RunningBurstTotal: burstTotal,
		ReadyQueue:        ready,
		IOQueue:           ioq,
	})
}
