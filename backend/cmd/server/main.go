package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/talaamm/cpu-scheduler-visualizer/algorithms"
	"github.com/talaamm/cpu-scheduler-visualizer/core"
)

var validAlgorithms = map[string]bool{
	"FCFS": true, "SJF": true, "SRTF": true,
	"RR": true, "Round_Robin": true,
	"Priority_NP": true, "Priority_P": true,
}

// validateProcesses rejects inputs that would produce a nonsensical or
// hanging simulation: missing/duplicate PIDs, negative arrival times,
// non-positive burst durations, unknown burst types, and processes that
// don't start with a CPU burst (the engine assumes the first burst is what
// gets scheduled onto the CPU).
func validateProcesses(procs []ProcessDTO) error {
	seen := make(map[string]bool, len(procs))
	for _, p := range procs {
		if p.PID == "" {
			return fmt.Errorf("process is missing a pid")
		}
		if seen[p.PID] {
			return fmt.Errorf("duplicate pid %q", p.PID)
		}
		seen[p.PID] = true

		if p.ArrivalTime < 0 {
			return fmt.Errorf("%s: arrival_time must be >= 0", p.PID)
		}
		if len(p.Bursts) == 0 {
			continue // a burst-less process completes immediately on arrival; allowed
		}
		if p.Bursts[0].Type != string(core.CPUBurst) {
			return fmt.Errorf("%s: first burst must be type CPU", p.PID)
		}
		for i, b := range p.Bursts {
			if b.Type != string(core.CPUBurst) && b.Type != string(core.IOBurst) {
				return fmt.Errorf("%s: burst %d has invalid type %q (must be CPU or IO)", p.PID, i, b.Type)
			}
			if b.Duration < 1 {
				return fmt.Errorf("%s: burst %d duration must be >= 1", p.PID, i)
			}
		}
	}
	return nil
}

func validateAlgorithm(id string) error {
	if !validAlgorithms[id] {
		return fmt.Errorf("unknown algorithm %q", id)
	}
	return nil
}

// ─── Request / Response DTOs ────────────────────────────────────────────────

type BurstDTO struct {
	Type     string `json:"type"`
	Duration int    `json:"duration"`
}

type ProcessDTO struct {
	PID         string     `json:"pid"`
	ArrivalTime int        `json:"arrival_time"`
	Priority    int        `json:"priority"`
	Bursts      []BurstDTO `json:"bursts"`
}

type SimulateRequest struct {
	Processes []ProcessDTO `json:"processes"`
	Algorithm string       `json:"algorithm"`
	Quantum   int          `json:"quantum"`
}

type CompareRequest struct {
	Processes  []ProcessDTO `json:"processes"`
	Algorithms []string     `json:"algorithms"`
	Quantum    int          `json:"quantum"`
}

type AlgorithmMeta struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	SupportsQuantum    bool   `json:"supports_quantum"`
	SupportsPreemption bool   `json:"supports_preemption"`
	SupportsPriority   bool   `json:"supports_priority"`
}

type PresetProcess struct {
	PID         string     `json:"pid"`
	ArrivalTime int        `json:"arrival_time"`
	Priority    int        `json:"priority"`
	Bursts      []BurstDTO `json:"bursts"`
}

type Preset struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Processes   []PresetProcess `json:"processes"`
}

// ─── DTO helpers ─────────────────────────────────────────────────────────────

func toCoreProcess(d ProcessDTO) core.Process {
	bursts := make([]core.Burst, len(d.Bursts))
	for i, b := range d.Bursts {
		bursts[i] = core.Burst{
			Type:     core.BurstType(b.Type),
			Duration: b.Duration,
		}
	}
	return core.Process{
		PID:         d.PID,
		ArrivalTime: d.ArrivalTime,
		Priority:    d.Priority,
		Bursts:      bursts,
		State:       core.StateNew,
	}
}

func toCoreProcesses(dtos []ProcessDTO) []core.Process {
	procs := make([]core.Process, len(dtos))
	for i, d := range dtos {
		procs[i] = toCoreProcess(d)
	}
	return procs
}

func schedulerFor(algorithm string, quantum int) algorithms.Scheduler {
	switch algorithm {
	case "FCFS":
		return algorithms.NewFCFS()
	case "SJF":
		return algorithms.NewSJFscheduler()
	case "SRTF":
		return algorithms.NewSRTFscheduler()
	case "RR", "Round_Robin":
		q := quantum
		if q <= 0 {
			q = 2
		}
		return algorithms.NewRoundRobin(q)
	case "Priority_NP":
		return algorithms.NewPriorityNONprScheduler()
	case "Priority_P":
		return algorithms.NewPriorityPreempScheduler()
	default:
		return algorithms.NewFCFS()
	}
}

// ─── Middleware ───────────────────────────────────────────────────────────────

func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func jsonResponse(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func jsonError(w http.ResponseWriter, status int, msg string) {
	jsonResponse(w, status, map[string]string{"error": msg})
}

// ─── Handlers ─────────────────────────────────────────────────────────────────

func handleSimulate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, http.StatusMethodNotAllowed, "POST only")
		return
	}

	var req SimulateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if len(req.Processes) == 0 {
		jsonError(w, http.StatusBadRequest, "at least one process required")
		return
	}
	if err := validateAlgorithm(req.Algorithm); err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validateProcesses(req.Processes); err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	procs := toCoreProcesses(req.Processes)
	sched := schedulerFor(req.Algorithm, req.Quantum)
	result := sched.Run(procs)

	jsonResponse(w, http.StatusOK, result)
}

func handleCompare(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, http.StatusMethodNotAllowed, "POST only")
		return
	}

	var req CompareRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if len(req.Processes) == 0 {
		jsonError(w, http.StatusBadRequest, "at least one process required")
		return
	}
	if len(req.Algorithms) == 0 {
		jsonError(w, http.StatusBadRequest, "at least one algorithm required")
		return
	}
	for _, alg := range req.Algorithms {
		if err := validateAlgorithm(alg); err != nil {
			jsonError(w, http.StatusBadRequest, err.Error())
			return
		}
	}
	if err := validateProcesses(req.Processes); err != nil {
		jsonError(w, http.StatusBadRequest, err.Error())
		return
	}

	type compareEntry struct {
		Algorithm string      `json:"algorithm"`
		Result    interface{} `json:"result"`
	}

	results := make([]interface{}, 0, len(req.Algorithms))
	for _, alg := range req.Algorithms {
		procs := toCoreProcesses(req.Processes) // fresh copy per algo
		sched := schedulerFor(alg, req.Quantum)
		results = append(results, sched.Run(procs))
	}

	jsonResponse(w, http.StatusOK, map[string]interface{}{"results": results})
}

func handleAlgorithms(w http.ResponseWriter, r *http.Request) {
	algos := []AlgorithmMeta{
		{
			ID:                 "FCFS",
			Name:               "First Come First Served",
			Description:        "Processes are scheduled in the order they arrive. Simple but can cause the convoy effect.",
			SupportsQuantum:    false,
			SupportsPreemption: false,
			SupportsPriority:   false,
		},
		{
			ID:                 "SJF",
			Name:               "Shortest Job First",
			Description:        "Always picks the process with the shortest next CPU burst. Optimal average waiting time but requires burst time knowledge.",
			SupportsQuantum:    false,
			SupportsPreemption: false,
			SupportsPriority:   false,
		},
		{
			ID:                 "SRTF",
			Name:               "Shortest Remaining Time First",
			Description:        "Preemptive version of SJF. Preempts whenever a shorter job arrives. Optimal but causes context switches.",
			SupportsQuantum:    false,
			SupportsPreemption: true,
			SupportsPriority:   false,
		},
		{
			ID:                 "RR",
			Name:               "Round Robin",
			Description:        "Each process gets a fixed time quantum. Fair but turnaround time depends heavily on quantum size.",
			SupportsQuantum:    true,
			SupportsPreemption: true,
			SupportsPriority:   false,
		},
		{
			ID:                 "Priority_NP",
			Name:               "Priority (Non-Preemptive)",
			Description:        "Runs the highest priority process to completion. Low priority processes may starve.",
			SupportsQuantum:    false,
			SupportsPreemption: false,
			SupportsPriority:   true,
		},
		{
			ID:                 "Priority_P",
			Name:               "Priority (Preemptive)",
			Description:        "Preempts whenever a higher priority process arrives. Most responsive to priority but most context switches.",
			SupportsQuantum:    false,
			SupportsPreemption: true,
			SupportsPriority:   true,
		},
	}
	jsonResponse(w, http.StatusOK, map[string]interface{}{"algorithms": algos})
}

func handlePresets(w http.ResponseWriter, r *http.Request) {
	presets := []Preset{
		{
			ID:          "classic",
			Name:        "Classic Textbook",
			Description: "The standard 3-process example from OS textbooks with CPU and IO bursts.",
			Processes: []PresetProcess{
				{PID: "P1", ArrivalTime: 0, Priority: 3, Bursts: []BurstDTO{{Type: "CPU", Duration: 5}, {Type: "IO", Duration: 3}, {Type: "CPU", Duration: 4}}},
				{PID: "P2", ArrivalTime: 1, Priority: 1, Bursts: []BurstDTO{{Type: "CPU", Duration: 3}, {Type: "IO", Duration: 2}, {Type: "CPU", Duration: 2}}},
				{PID: "P3", ArrivalTime: 2, Priority: 2, Bursts: []BurstDTO{{Type: "CPU", Duration: 2}, {Type: "IO", Duration: 1}, {Type: "CPU", Duration: 3}}},
			},
		},
		{
			ID:          "io_heavy",
			Name:        "IO-Heavy Workload",
			Description: "Processes with long IO waits — ideal for seeing CPU idle periods and IO queue behavior.",
			Processes: []PresetProcess{
				{PID: "P1", ArrivalTime: 0, Priority: 2, Bursts: []BurstDTO{{Type: "CPU", Duration: 2}, {Type: "IO", Duration: 8}, {Type: "CPU", Duration: 2}}},
				{PID: "P2", ArrivalTime: 0, Priority: 1, Bursts: []BurstDTO{{Type: "CPU", Duration: 2}, {Type: "IO", Duration: 8}, {Type: "CPU", Duration: 2}}},
				{PID: "P3", ArrivalTime: 1, Priority: 3, Bursts: []BurstDTO{{Type: "CPU", Duration: 3}, {Type: "IO", Duration: 6}, {Type: "CPU", Duration: 3}}},
			},
		},
		{
			ID:          "starvation",
			Name:        "Starvation Scenario",
			Description: "A low-priority process that may never get CPU time under priority scheduling.",
			Processes: []PresetProcess{
				{PID: "P1", ArrivalTime: 0, Priority: 1, Bursts: []BurstDTO{{Type: "CPU", Duration: 10}}},
				{PID: "P2", ArrivalTime: 1, Priority: 1, Bursts: []BurstDTO{{Type: "CPU", Duration: 8}}},
				{PID: "P3", ArrivalTime: 2, Priority: 1, Bursts: []BurstDTO{{Type: "CPU", Duration: 6}}},
				{PID: "P4", ArrivalTime: 0, Priority: 5, Bursts: []BurstDTO{{Type: "CPU", Duration: 12}}},
			},
		},
		{
			ID:          "burst",
			Name:        "CPU-Burst Heavy",
			Description: "All CPU bursts, no IO — ideal for comparing pure scheduling algorithms cleanly.",
			Processes: []PresetProcess{
				{PID: "P1", ArrivalTime: 0, Priority: 3, Bursts: []BurstDTO{{Type: "CPU", Duration: 8}}},
				{PID: "P2", ArrivalTime: 2, Priority: 1, Bursts: []BurstDTO{{Type: "CPU", Duration: 4}}},
				{PID: "P3", ArrivalTime: 4, Priority: 2, Bursts: []BurstDTO{{Type: "CPU", Duration: 6}}},
				{PID: "P4", ArrivalTime: 5, Priority: 4, Bursts: []BurstDTO{{Type: "CPU", Duration: 2}}},
			},
		},
	}
	jsonResponse(w, http.StatusOK, map[string]interface{}{"presets": presets})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, map[string]string{"status": "ok"})
}

func locateEntryFile() string {
	candidates := []string{
		"cpu-scheduler-visualizer.html",
		filepath.Join("frontend", "index.html"),
		filepath.Join("frontend", "index-vite.html"),
	}

	if cwd, err := os.Getwd(); err == nil {
		for dir := cwd; ; dir = filepath.Dir(dir) {
			for _, name := range candidates {
				path := filepath.Join(dir, name)
				if _, err := os.Stat(path); err == nil {
					return path
				}
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
		}
	}

	if _, file, _, ok := runtime.Caller(0); ok {
		for dir := filepath.Dir(file); ; dir = filepath.Dir(dir) {
			for _, name := range candidates {
				path := filepath.Join(dir, name)
				if _, err := os.Stat(path); err == nil {
					return path
				}
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
		}
	}

	return ""
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	root := locateEntryFile()
	if root == "" {
		http.Error(w, "missing frontend entry file", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, root)
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", serveRoot)
	mux.HandleFunc("/api/simulate", cors(handleSimulate))
	mux.HandleFunc("/api/compare", cors(handleCompare))
	mux.HandleFunc("/api/algorithms", cors(handleAlgorithms))
	mux.HandleFunc("/api/presets", cors(handlePresets))
	mux.HandleFunc("/health", cors(handleHealth))

	return mux
}

// ─── Main ─────────────────────────────────────────────────────────────────────

func main() {
	mux := newMux()

	log.Println("CPU Scheduler API running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
