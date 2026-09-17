// ─── Core domain types ──────────────────────────────────────────────────────

export type BurstType = 'CPU' | 'IO'

export interface Burst {
  type: BurstType
  duration: number
}

export interface Process {
  pid: string
  arrival_time: number
  priority: number
  bursts: Burst[]
}

export type AlgorithmID =
  | 'FCFS'
  | 'SJF'
  | 'SRTF'
  | 'RR'
  | 'Priority_NP'
  | 'Priority_P'

// ─── API response types ──────────────────────────────────────────────────────

export interface TimelineEntry {
  time: number
  process_id: string
  event: 'RUNNING' | 'IDLE' | 'READY' | 'WAITING' | 'TERMINATED'
}

/** Ground-truth system state for a single time tick, as computed by the Go engine. */
export interface StateSnapshot {
  time: number
  running: string // "" means CPU is idle
  running_remaining: number
  running_burst_total: number
  ready_queue: string[]
  io_queue: Array<{ pid: string; remaining: number }>
}

export type ScheduleEventType =
  | 'ARRIVED'
  | 'SELECTED'
  | 'PREEMPTED'
  | 'BURST_DONE'
  | 'IO_START'
  | 'IO_DONE'
  | 'COMPLETED'

/** An algorithm-sourced explanation of a scheduling decision or transition. */
export interface ScheduleEvent {
  time: number
  type: ScheduleEventType
  pid: string
  message: string
}

export interface ProcessResult {
  pid: string
  arrival_time: number
  priority: number
  bursts: Burst[]
  start_time: number
  completion_time: number
  waiting_time: number
  turnaround_time: number
  response_time: number
  started: boolean
  completed: boolean
}

export interface SimulationResult {
  algorithm: string
  timeline: TimelineEntry[]
  snapshots: StateSnapshot[]
  events: ScheduleEvent[]
  processes: ProcessResult[]
  total_time: number
  cpu_utilization: number
  average_waiting_time: number
  average_turnaround_time: number
  average_response_time: number
}

export interface CompareResult {
  results: SimulationResult[]
}

export interface AlgorithmMeta {
  id: AlgorithmID
  name: string
  description: string
  supports_quantum: boolean
  supports_preemption: boolean
  supports_priority: boolean
}

export interface Preset {
  id: string
  name: string
  description: string
  processes: Process[]
}

// ─── Request types ───────────────────────────────────────────────────────────

export interface SimulateRequest {
  processes: Process[]
  algorithm: AlgorithmID
  quantum?: number
}

export interface CompareRequest {
  processes: Process[]
  algorithms: AlgorithmID[]
  quantum?: number
}

// ─── Derived UI types ─────────────────────────────────────────────────────────

/** A compacted Gantt segment with start/end times (derived from timeline) */
export interface GanttSegment {
  pid: string
  start: number
  end: number
  isIdle: boolean
}

/** Per-process row for the per-process Gantt view */
export interface ProcessGanttRow {
  pid: string
  color: string
  segments: Array<{
    start: number
    end: number
    type: 'cpu' | 'io' | 'idle' | 'waiting'
  }>
}

/** Simulation state at a specific time tick (for animation scrubbing) */
export interface SimulationFrame {
  time: number
  runningPid: string | null
  runningRemaining: number
  runningBurstTotal: number
  readyQueue: string[]
  ioQueue: Array<{ pid: string; remainingIO: number }>
  completedPids: string[]
  eventLog: Array<{ time: number; text: string; color: string }>
}
