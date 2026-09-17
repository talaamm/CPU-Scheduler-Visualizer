import type {
  SimulateRequest,
  SimulationResult,
  CompareRequest,
  CompareResult,
  AlgorithmMeta,
  Preset,
} from '@/types'

const BASE_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

class APIError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message)
    this.name = 'APIError'
  }
}

async function post<T>(path: string, body: unknown): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  })

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'Unknown error' }))
    throw new APIError(res.status, err.error ?? `HTTP ${res.status}`)
  }

  return res.json()
}

async function get<T>(path: string): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`)

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'Unknown error' }))
    throw new APIError(res.status, err.error ?? `HTTP ${res.status}`)
  }

  return res.json()
}

// ─── Public API ───────────────────────────────────────────────────────────────

export const api = {
  /**
   * Run a single scheduling algorithm simulation.
   * @example
   * const result = await api.simulate({ processes, algorithm: 'FCFS' })
   */
  simulate(req: SimulateRequest): Promise<SimulationResult> {
    return post<SimulationResult>('/api/simulate', req)
  },

  /**
   * Run multiple algorithms on the same workload for comparison.
   * @example
   * const { results } = await api.compare({ processes, algorithms: ['FCFS', 'SJF', 'RR'] })
   */
  compare(req: CompareRequest): Promise<CompareResult> {
    return post<CompareResult>('/api/compare', req)
  },

  /**
   * Fetch all available scheduling algorithms with metadata.
   */
  getAlgorithms(): Promise<{ algorithms: AlgorithmMeta[] }> {
    return get('/api/algorithms')
  },

  /**
   * Fetch built-in scenario presets.
   */
  getPresets(): Promise<{ presets: Preset[] }> {
    return get('/api/presets')
  },
}

// ─── Derived data helpers (pure, no HTTP) ────────────────────────────────────

import type {
  GanttSegment,
  ProcessGanttRow,
  SimulationFrame,
  TimelineEntry,
  ProcessResult,
} from '@/types'

/** Process colors — cycle through palette for N processes */
const PROCESS_COLORS = [
  '#4F8EF7', // blue
  '#9B6DFF', // purple
  '#2DD4BF', // teal
  '#F59E0B', // amber
  '#F472B6', // pink
  '#22C55E', // green
  '#FB923C', // orange
  '#A78BFA', // violet
]

export function processColor(index: number): string {
  return PROCESS_COLORS[index % PROCESS_COLORS.length]
}

/**
 * Convert flat timeline to Gantt segments [start, end, pid].
 */
export function toGanttSegments(result: SimulationResult): GanttSegment[] {
  const segments: GanttSegment[] = []
  const { timeline, total_time } = result

  for (let i = 0; i < timeline.length; i++) {
    const entry = timeline[i]
    const nextTime = i + 1 < timeline.length ? timeline[i + 1].time : total_time

    segments.push({
      pid: entry.process_id,
      start: entry.time,
      end: nextTime,
      isIdle: entry.process_id === 'IDLE',
    })
  }

  return segments
}

/**
 * Build per-process Gantt rows showing CPU, IO wait, and idle stretches.
 * This requires the full ProcessResult array for IO burst durations.
 */
export function toProcessGanttRows(result: SimulationResult): ProcessGanttRow[] {
  const { processes, timeline, total_time } = result

  return processes.map((proc, idx) => {
    const color = processColor(idx)
    const segments: ProcessGanttRow['segments'] = []

    // Find all CPU segments for this process
    const cpuEntries = timeline.filter((e) => e.process_id === proc.pid)

    for (let t = proc.arrival_time; t < total_time; ) {
      const cpuAt = cpuEntries.find((e) => {
        const nextIdx = timeline.indexOf(e) + 1
        const end = nextIdx < timeline.length ? timeline[nextIdx].time : total_time
        return e.time <= t && t < end
      })

      if (cpuAt) {
        const idx2 = timeline.indexOf(cpuAt)
        const end = idx2 + 1 < timeline.length ? timeline[idx2 + 1].time : total_time
        segments.push({ start: t, end, type: 'cpu' })
        t = end
      } else {
        // Check if in IO or just waiting
        const burst = getIOBurstAt(proc, t)
        if (burst) {
          segments.push({ start: t, end: t + 1, type: 'io' })
        } else if (t >= proc.arrival_time) {
          segments.push({ start: t, end: t + 1, type: 'waiting' })
        }
        t++
      }
    }

    return { pid: proc.pid, color, segments: mergeAdjacent(segments) }
  })
}

function getIOBurstAt(proc: ProcessResult, _t: number): boolean {
  // Simplified: just checks if process has IO bursts (full impl needs IO queue tracking)
  return proc.bursts.some((b) => b.type === 'IO')
}

function mergeAdjacent(segs: ProcessGanttRow['segments']): ProcessGanttRow['segments'] {
  if (segs.length === 0) return segs
  const merged = [segs[0]]
  for (let i = 1; i < segs.length; i++) {
    const last = merged[merged.length - 1]
    if (last.type === segs[i].type && last.end === segs[i].start) {
      last.end = segs[i].end
    } else {
      merged.push(segs[i])
    }
  }
  return merged
}

/**
 * Reconstruct simulation frames for step-by-step animation.
 * Each frame represents the system state at time t.
 */
export function buildSimulationFrames(
  result: SimulationResult,
  processes: ProcessResult[],
): SimulationFrame[] {
  const frames: SimulationFrame[] = []
  const { timeline, total_time } = result

  for (let t = 0; t <= total_time; t++) {
    // Find running process at time t
    let runningPid: string | null = null
    for (let i = timeline.length - 1; i >= 0; i--) {
      if (timeline[i].time <= t) {
        runningPid = timeline[i].process_id === 'IDLE' ? null : timeline[i].process_id
        break
      }
    }

    // Determine ready queue (arrived, not running, not completed yet)
    const readyQueue: string[] = []
    const ioQueue: Array<{ pid: string; remainingIO: number }> = []
    const completedPids: string[] = []
    const eventLog: string[] = []

    for (const proc of processes) {
      if (proc.arrival_time > t) continue
      if (proc.completion_time <= t && proc.completed) {
        completedPids.push(proc.pid)
        continue
      }
      if (proc.pid === runningPid) continue
      readyQueue.push(proc.pid)
    }

    // Event log for this tick
    const tickEntry = timeline.find((e) => e.time === t)
    if (tickEntry) {
      if (tickEntry.process_id === 'IDLE') {
        eventLog.push(`t=${t}  CPU idle — no processes ready`)
      } else {
        eventLog.push(`t=${t}  ${tickEntry.process_id} running on CPU`)
      }
    }

    frames.push({ time: t, runningPid, readyQueue, ioQueue, completedPids, eventLog })
  }

  return frames
}

/**
 * Generate a "why was this decision made?" explanation for a Gantt segment.
 */
export function explainDecision(
  segment: GanttSegment,
  algorithm: string,
  result: SimulationResult,
): string {
  if (segment.isIdle) {
    return 'CPU is idle — no processes are in the ready queue at this time. All active processes are waiting on I/O.'
  }

  const pid = segment.pid
  const proc = result.processes.find((p) => p.pid === pid)
  if (!proc) return ''

  switch (algorithm) {
    case 'FCFS':
      return `${pid} was selected because it arrived before all other ready processes (FCFS order).`
    case 'SJF':
      return `${pid} was selected because it had the shortest next CPU burst among ready processes.`
    case 'SRTF':
      return `${pid} was selected because it had the least remaining burst time (${proc.bursts[0]?.duration ?? '?'} units) at t=${segment.start}.`
    case 'Round_Robin':
      return `${pid} received its time quantum at t=${segment.start}. Round Robin ensures each process gets equal CPU time.`
    case 'Priority_Non_Preemptive':
    case 'Priority_Preemptive':
      return `${pid} was selected because it had the highest priority (value=${proc.priority}) among ready processes. Lower number = higher priority.`
    default:
      return `${pid} was scheduled at t=${segment.start}.`
  }
}

export { APIError }
