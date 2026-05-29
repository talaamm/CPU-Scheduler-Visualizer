package simulation

//Visualization-ready timeline structure
type TimelineEntry struct {
	Time      int    `json:"time"`
	ProcessID string `json:"process_id"`
	Event     string `json:"event"`
}
