# CPU Scheduler Visualizer

## Overview

CPU Scheduler Visualizer is an interactive Operating Systems educational platform that allows users to create custom workloads, simulate CPU scheduling algorithms, visualize process execution, and compare algorithm performance.

The project aims to bridge Operating Systems theory with practical understanding through real-time visualization and performance analysis.

Users can define processes containing multiple CPU and I/O bursts, select scheduling algorithms, run simulations, and inspect detailed execution timelines and metrics.

---

# Goals

## Educational Goals

* Understand CPU scheduling algorithms
* Visualize process state transitions
* Demonstrate CPU and I/O interaction
* Show effects of scheduling decisions on performance
* Compare scheduling policies under identical workloads

## Technical Goals

* Build a reusable scheduling simulation engine
* Separate scheduling policy from simulation logic
* Create an interactive web-based visualization platform
* Expose simulation functionality through a backend API
* Support future extension to advanced scheduling concepts

---

# Core Concepts

The simulator models:

* Processes
* CPU bursts
* I/O bursts
* Ready queue
* I/O queue
* CPU execution
* Process state transitions

Process states:

* NEW
* READY
* RUNNING
* WAITING
* TERMINATED

---

# System Architecture

## High-Level Architecture

```text
Frontend
    ↓
REST API
    ↓
Simulation Engine
    ↓
Scheduling Policy
```

---

## Component Breakdown

### Frontend

Responsibilities:

* Process creation
* Workload configuration
* Algorithm selection
* Visualization
* Performance comparison

Technology:

* HTML
* Tailwind CSS
* JavaScript
* Chart.js

---

### API Layer

Responsibilities:

* Receive simulation requests
* Validate input
* Execute scheduler
* Return simulation results

Technology:

* Go
* Gin (optional)

---

### Simulation Engine

Responsibilities:

* Time progression
* Process arrivals
* Ready queue management
* I/O queue management
* CPU execution
* State transitions
* Timeline generation
* Metrics collection

The engine is independent of scheduling algorithms.

---

### Scheduling Policies

Responsibilities:

* Decide which READY process receives CPU

Policies:

* FCFS
* SJF
* Round Robin
* Priority Scheduling
* SRTF

Scheduling algorithms contain only selection logic.

---

# Project Structure

```text
cpu-scheduler-visualizer/

cmd/
└── server/
    └── main.go

internal/

├── process/
│   ├── process.go
│   ├── burst.go
│   └── state.go
│
├── simulation/
│   ├── engine.go
│   ├── cpu.go
│   ├── io.go
│   ├── timeline.go
│   ├── metrics.go
│   └── result.go
│
├── scheduler/
│   ├── scheduler.go
│   ├── fcfs.go
│   ├── sjf.go
│   ├── rr.go
│   ├── priority.go
│   └── srtf.go
│
├── api/
│   ├── routes.go
│   ├── handlers.go
│   └── dto.go
│
└── queue/
    ├── ready_queue.go
    └── io_queue.go

web/

├── static/
├── assets/
└── frontend/

README.md
go.mod
```

---

# Data Model

## Process

Fields:

* PID
* Arrival Time
* Bursts
* Current Burst Index
* Remaining Burst Time
* State
* Waiting Time
* Response Time
* Turnaround Time
* Completion Time

---

## Burst

Fields:

* Type

  * CPU
  * IO
* Duration

---

## Timeline Entry

Fields:

* Time
* Process ID
* Event Type

Examples:

* RUNNING
* READY
* WAITING
* TERMINATED
* IDLE

---

# Scheduler Interface

Every algorithm implements:

```go
type Scheduler interface {
    Run(processes []process.Process) simulation.SimulationResult
    Name() string
}
```

This allows all algorithms to be executed interchangeably.

---

# User Modes

## Single Algorithm Mode

Purpose:

Visualize one algorithm in detail.

Workflow:

1. Create processes
2. Select algorithm
3. Run simulation
4. View execution timeline
5. View metrics

Outputs:

* Gantt Chart
* Timeline Animation
* Process Statistics

---

## Comparative Analysis Mode

Purpose:

Compare multiple algorithms using the same workload.

Workflow:

1. Create workload
2. Select algorithms
3. Run simulations
4. Compare results

Outputs:

* Waiting Time Comparison
* Turnaround Time Comparison
* Response Time Comparison
* CPU Utilization Comparison
* Throughput Comparison

---

# Visualization Features

## Process Builder

Users can:

* Add processes
* Add CPU bursts
* Add I/O bursts
* Configure arrival times

Example:

P1:
Arrival = 0

Bursts:
CPU 5
IO 3
CPU 4

---

## Timeline Visualization

Animated execution timeline showing:

* Running process
* CPU idle periods
* State transitions

---

## Gantt Chart

Visual representation of CPU allocation over time.

Example:

| P1 | P1 | P1 | P2 | P2 | P3 |

---

## Queue Visualization

Display:

* Ready Queue
* I/O Queue

As simulation progresses.

---

## Metrics Dashboard

Display:

* Average Waiting Time
* Average Turnaround Time
* Average Response Time
* CPU Utilization
* Throughput

---

## Comparative Charts

Using Chart.js:

* Bar Charts
* Line Charts
* Algorithm Ranking

---

# Metrics

For each process:

* Waiting Time
* Turnaround Time
* Response Time
* Completion Time

System metrics:

* CPU Utilization
* Throughput
* Average Waiting Time
* Average Turnaround Time
* Average Response Time

---

# Future Extensions

## Advanced Scheduling

* Multi-Level Queue
* Multi-Level Feedback Queue
* Lottery Scheduling

## Multiprocessor Support

* Multiple CPUs
* Load Balancing
* CPU Affinity

## Hardware Concepts

* Context Switch Cost
* Cache Penalties
* NUMA Awareness

## Memory Management

Potential future module:

* Paging Visualization
* Page Replacement Algorithms

## Deadlock Simulation

Potential future module:

* Resource Allocation Graphs
* Deadlock Detection
* Deadlock Avoidance

---

# Portfolio Value

This project demonstrates:

* Operating Systems knowledge
* Scheduling algorithms
* Data structures
* Simulation design
* Backend development
* API design
* Software architecture
* Data visualization
* Performance analysis

The project serves as both an educational tool and a systems engineering portfolio piece.

---

## Project Phases

### Phase 1

Process model
Simulation engine
FCFS
RR

### Phase 2

SJF
Priority
SRTF

### Phase 3

REST API

### Phase 4

Frontend visualizer

### Phase 5

Comparative mode

### Phase 6

Deployment + demo video + portfolio write-up
