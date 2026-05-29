# CPU Scheduler Visualizer

## 📖 Overview

This project is an interactive CPU Scheduling Simulator and Visualization tool that demonstrates how different CPU scheduling algorithms behave under the same workload.

It allows users to define processes with multiple CPU bursts, configure scheduling parameters, and compare the performance of different scheduling algorithms through visual and statistical outputs.

The goal is to bridge theoretical Operating Systems concepts with practical performance insights.

---

## 🎯 Objectives

- Implement core CPU scheduling algorithms in C.
- Simulate realistic process execution with multiple CPU bursts
- Support configurable scheduling parameters (e.g., Round Robin quantum)
- Provide performance comparison between all algorithms for a given scenario
- Visualize scheduling behavior using charts and timelines
- Analyze tradeoffs between fairness, throughput, and response time

---

## 🧠 Key Concepts Covered

- Process lifecycle (Ready, Running, Waiting, Terminated)
- CPU burst and I/O burst modeling
- CPU scheduling algorithms
- Context switching
- Performance metrics:
  - Average waiting time
  - Turnaround time
  - Response time
  - CPU utilization
- Throughput and fairness tradeoffs

---

## ⚙️ Supported Scheduling Algorithms

The simulator will implement:

1. First Come First Serve (FCFS)
2. Shortest Job First (SJF - Non-preemptive)
3. Shortest Remaining Time First (SRTF)
4. Round Robin (RR)
5. Priority Scheduling (pree-emptivee / non-preemptive)
6. Priority with RR
7. Multi-level Queue (Future Implementation)
8. Multi-level feedback Queue  (Future Implementation)

---

## 🧩 System Features

### 1. Process Builder

Users can define processes with:

- Process ID
- Arrival Time
- Multiple CPU bursts
- I/O bursts

Example:

```bash
P1:
Arrival: 0
Bursts: [CPU 5, IO 3, CPU 4]
```

---

### 2. Scheduling Configuration

- Select scheduling algorithm
- Set Round Robin time quantum (if applicable)
- Choose simulation speed (for visualization mode REAL TIME ANALOGY?)

---

### 3. Simulation Engine (C Backend)

The backend will simulate:

- Time progression
- Process state transitions
- Queue management
- CPU allocation per time unit

Output:

- Execution timeline
- Process state history
- Scheduling decisions log

---

### 4. Output Data Format

The C engine will export results as:

- JSON (preferred)
- or CSV (fallback)?

Example structure:

```json
{
  "timeline": [
    {"time": 0, "process": "P1"},
    {"time": 1, "process": "P1"},
    {"time": 2, "process": "P2"}
  ],
  "metrics": {
    "P1": {
      "waiting_time": 5,
      "turnaround_time": 12
    }
  }
}
```

---

## 🧭 System Modes

The simulator supports two main operating modes to improve usability and flexibility:

### 1. Single Algorithm Mode

This mode allows users to select one scheduling algorithm and run a full simulation.

#### Purpose

- Educational understanding of a single algorithm
- Step-by-step visualization of execution
- Focus on intuition rather than comparison

#### Features

- Select one algorithm:
  - FCFS
  - SJF
  - SRTF
  - Round Robin
  - Priority Scheduling
- Configure parameters (e.g., time quantum for RR)
- Run simulation
- View:
  - Gantt chart
  - Process execution timeline
  - Performance metrics (for that algorithm only)

---

### 2. Comparative Analysis Mode

This mode runs multiple scheduling algorithms on the same process set.

#### Purpose

- Performance comparison
- System behavior analysis
- Educational insight into tradeoffs

#### Features

- Run all selected algorithms on identical input
- Automatically compute and compare:
  - Average Waiting Time
  - Turnaround Time
  - Response Time
  - CPU Utilization
- Generate comparative visualizations:
  - Side-by-side Gantt charts
  - Bar charts for performance metrics
- Highlight best-performing algorithm per metric

---

## 🔁 Mode Selection Flow (Frontend Design)

```text
User Input
    ↓
Select Mode:
    ├── Single Algorithm Mode
    │       ↓
    │   Choose Algorithm → Run Simulation → Visualize Output
    │
    └── Comparative Mode ? between all or just 2?
            ↓
    Run All (OR SELECTED?) Algorithms → Collect Metrics → Compare → Visualize Results
```

---

🎯 Design Benefit

This dual-mode design allows the system to serve two audiences:

Beginners → understand how each algorithm works individually
Advanced users → analyze and compare performance tradeoffs

It transforms the project from a simple simulator into an interactive educational OS analysis tool.

---

Project Structure

```t
CPU-Scheduler-Visualizer/
│
├── 📂 backend/                          # C simulation engine
│   ├── 📂 algorithms/                   # Each algorithm in its own file
│   │   ├── fcfs.c                       # FCFS implementation
│   │   ├── fcfs.h
│   │   ├── sjf.c                        # SJF (non-preemptive)
│   │   ├── sjf.h
│   │   ├── srtf.c                       # Shortest Remaining Time First
│   │   ├── srtf.h
│   │   ├── round_robin.c                # Round Robin
│   │   ├── round_robin.h
│   │   ├── priority.c                   # Priority Scheduling
│   │   ├── priority.h
│   │   ├── priority_rr.c                # Priority + Round Robin
│   │   ├── priority_rr.h
│   │   ├── multilevel_queue.c           # Multi-level Queue (future)
│   │   └── multilevel_queue.h
│   │
│   ├── 📂 core/                         # Core simulation engine
│   │   ├── process.c                    # Process structure & lifecycle
│   │   ├── process.h
│   │   ├── queue.c                      # Queue management
│   │   ├── queue.h
│   │   ├── scheduler.c                  # Main scheduler dispatcher
│   │   ├── scheduler.h
│   │   └── simulation.c                 # Simulation engine
│   │   └── simulation.h
│   │
│   ├── 📂 utils/                        # Utility functions
│   │   ├── json_output.c                # JSON export
│   │   ├── json_output.h
│   │   ├── csv_output.c                 # CSV export
│   │   ├── csv_output.h
│   │   ├── metrics.c                    # Performance metrics calculation
│   │   ├── metrics.h
│   │   └── config_parser.c              # Parse input configuration
│   │   └── config_parser.h
│   │
│   ├── main.c                           # Entry point
│   ├── Makefile                         # Build configuration
│   └── CMakeLists.txt                   # (Alternative build system)
│
├── 📂 frontend/                         # Web/GUI interface
│   ├── 📂 public/
│   │   ├── index.html
│   │   └── favicon.ico
│   │
│   ├── 📂 src/
│   │   ├── 📂 components/
│   │   │   ├── ProcessBuilder.js        # Define processes UI
│   │   │   ├── AlgorithmSelector.js     # Choose algorithm(s)
│   │   │   ├── ConfigPanel.js           # Set parameters (quantum, etc.)
│   │   │   ├── GanttChart.js            # Visualization
│   │   │   ├── MetricsDisplay.js        # Show performance stats
│   │   │   └── ComparisonPanel.js       # Comparative mode UI
│   │   │
│   │   ├── 📂 api/
│   │   │   └── simulator.js             # API calls to backend
│   │   │
│   │   ├── 📂 utils/
│   │   │   ├── parser.js                # Parse input
│   │   │   └── formatter.js             # Format output data
│   │   │
│   │   ├── App.js                       # Main app component
│   │   └── index.js
│   │
│   ├── package.json
│   └── .gitignore
│
├── 📂 docs/
│   ├── algorithms.md                    # Algorithm explanations
│   ├── api.md                           # Backend API docs
│   ├── setup.md                         # Setup instructions
│   └── examples.md                      # Usage examples
│
├── 📂 tests/                            # Test suite
│   ├── 📂 backend/
│   │   ├── test_fcfs.c
│   │   ├── test_sjf.c
│   │   ├── test_round_robin.c
│   │   └── test_metrics.c
│   │
│   └── 📂 frontend/
│       ├── components.test.js
│       └── api.test.js
│
├── Project Plan.md
├── README.md                            # Project overview
└── .gitignore
```

---

### AI Usage

Planning helper: ChatGPT - codex
Coding helper: VS Code Co-pilot: Claude Haiku 4.5 agent

---

# SWITCHING FROM C-LANG TO GO-LANG

Yes. Honestly, for *this specific project* (CPU scheduling simulator with queues/events/processes), Go would probably feel much cleaner and easier to reason about than C — especially for you.

Not because C is bad. C is actually the classic OS language.
But because right now you're still building:

* architecture thinking
* abstraction skills
* simulation logic
* data structure confidence

and C forces you to fight:

* memory management
* pointers
* linked list bugs
* malloc/free
* manual resizing
* string handling
* queue bookkeeping

ALL while also learning scheduling algorithms.

That becomes mentally exhausting.

---

# Why your C code probably feels “too complicated”

Because in C, even SIMPLE concepts become verbose.

Example:
A ready queue in C often becomes:

```c id="gqglfx"
typedef struct Node {
    Process *process;
    struct Node *next;
} Node;
```

then:

* enqueue
* dequeue
* free memory
* NULL checks
* pointer bugs

But conceptually you only wanted:

> “put process into queue”

---

# In Go, the SAME idea becomes:

```go id="dgrn0e"
readyQueue := []*Process{}
```

enqueue:

```go id="u0q1sa"
readyQueue = append(readyQueue, p)
```

dequeue:

```go id="0g9v1z"
p := readyQueue[0]
readyQueue = readyQueue[1:]
```

MUCH closer to how your brain thinks.

---

# For scheduling simulators, Go is REALLY nice because:

## 1. Structs are clean

```go id="t9b1iv"
type Process struct {
    ID           string
    ArrivalTime  int
    Bursts       []Burst
    CurrentBurst int
    Remaining    int
    Priority     int
}
```

Readable immediately.

---

## 2. Dynamic arrays (slices)

No malloc headaches.

```go id="i0m31d"
ready := []*Process{}
ioQueue := []*Process{}
```

---

## 3. Sorting is easy

SJF:

```go id="sv2r44"
sort.Slice(ready, func(i, j int) bool {
    return ready[i].Remaining < ready[j].Remaining
})
```

Done.

---

## 4. Maps are beautiful

```go id="z3xw1g"
processMap := map[string]*Process{}
```

---

## 5. Less code noise

C code often becomes:

* 60% infrastructure
* 40% actual scheduling logic

Go becomes:

* 80% scheduling logic
* 20% infrastructure

which is MUCH better for learning algorithms.

---

# BUT — very important

## Using Go will NOT automatically simplify BAD architecture.

You can still write:

* giant functions
* duplicated schedulers
* messy state handling

in Go.

The BIGGEST improvement comes from architecture, not language.

---

# The ideal design in Go

You’d want something like:

```text id="n6h5sq"
scheduler/
    fcfs.go
    sjf.go
    srtf.go
    rr.go
    priority.go

simulation/
    engine.go
    process.go
    queue.go
```

Where:

* engine handles time/events
* schedulers only choose next process

THAT is the clean part.

---

# Another important thing

If this is for an OS course:

Your professor may EXPECT C.

So before rewriting:
check:

* assignment requirements
* grading language
* restrictions

because some OS courses only allow:

* C
* C++
* maybe Java

---

# My actual recommendation for YOU specifically

Since you already know some Go and are still developing abstraction skills:

## If language is NOT restricted:

Use Go.

You will:

* learn faster
* debug faster
* iterate faster
* understand the scheduler better
* focus on concepts instead of memory bugs

And your code will probably become:

* shorter
* cleaner
* easier to extend

especially for:

* SRTF
* RR
* Priority RR
* IO handling

---

# One warning though

Do NOT let Claude write the whole thing blindly.

Because then you end up with:

* code that works
* but architecture you don’t understand

Instead:

* build small pieces yourself
* ask AI for review/help
* refactor gradually

That’s how your engineering skill grows.
