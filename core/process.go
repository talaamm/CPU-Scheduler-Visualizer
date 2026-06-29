package core

type Process struct {
	PID string `json:"pid"`

	ArrivalTime int `json:"arrival_time"`

	Priority int `json:"priority"`

	Bursts []Burst `json:"bursts"`

	CurrentBurstIndex int `json:"current_burst_index"` // which burst process is currently executing

	RemainingBurstTime int `json:"remaining_burst_time"` // preemption easier and for RR/SRTF later

	State ProcessState `json:"state"`

	StartTime      int `json:"start_time"`
	CompletionTime int `json:"completion_time"`

	WaitingTime    int `json:"waiting_time"`
	TurnaroundTime int `json:"turnaround_time"`
	ResponseTime   int `json:"response_time"`

	Started   bool `json:"started"`
	Completed bool `json:"completed"`
}
