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

// A fetch() network failure (backend unreachable, CORS blocked, etc.) throws
// a generic TypeError whose message is just "Failed to fetch" — not useful
// to someone who forgot to start the Go server. Give them something
// actionable instead.
function unreachableBackendError(): APIError {
  return new APIError(
    0,
    `Can't reach the backend at ${BASE_URL}. Make sure it's running: go run ./backend/cmd/server`,
  )
}

async function post<T>(path: string, body: unknown): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(body),
  }).catch(() => {
    throw unreachableBackendError()
  })

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'Unknown error' }))
    throw new APIError(res.status, err.error ?? `HTTP ${res.status}`)
  }

  return res.json()
}

async function get<T>(path: string): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`).catch(() => {
    throw unreachableBackendError()
  })

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

import type { GanttSegment, SimulationFrame, ProcessResult } from '@/types'

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
 * Reconstruct simulation frames for step-by-step animation directly from the
 * engine's real per-tick snapshots and events — no guessing about queue
 * membership or burst progress.
 */
export function buildSimulationFrames(
  result: SimulationResult,
  processes: ProcessResult[],
): SimulationFrame[] {
  const { snapshots, events, total_time } = result

  const pidIndex: Record<string, number> = {}
  processes.forEach((p, i) => {
    pidIndex[p.pid] = i
  })
  const colorFor = (pid: string) => processColor(pidIndex[pid] ?? 0)

  const snapshotAt = new Map(snapshots.map((s) => [s.time, s]))

  const frames: SimulationFrame[] = []
  let last = snapshots[0]

  for (let t = 0; t <= total_time; t++) {
    const snap = snapshotAt.get(t) ?? last
    if (snap) last = snap

    const completedPids = processes
      .filter((p) => p.completed && p.completion_time <= t)
      .map((p) => p.pid)

    const eventLog = events
      .filter((e) => e.time <= t)
      .slice(-30)
      .reverse()
      .map((e) => ({
        time: e.time,
        text: e.message,
        color: e.type === 'COMPLETED' ? '#22C55E' : colorFor(e.pid),
      }))

    frames.push({
      time: t,
      runningPid: snap?.running || null,
      runningRemaining: snap?.running_remaining ?? 0,
      runningBurstTotal: snap?.running_burst_total ?? 0,
      readyQueue: snap?.ready_queue ?? [],
      ioQueue: (snap?.io_queue ?? []).map((e) => ({ pid: e.pid, remainingIO: e.remaining })),
      completedPids,
      eventLog,
    })
  }

  return frames
}

/**
 * Look up the engine's own explanation for why a Gantt segment's process was
 * scheduled at that time (sourced from the algorithm's real selection logic,
 * not reconstructed after the fact).
 */
export function explainDecision(
  segment: GanttSegment,
  algorithm: string,
  result: SimulationResult,
): string {
  if (segment.isIdle) {
    return 'CPU is idle — no processes are in the ready queue at this time. All active processes are waiting on I/O.'
  }

  const match = result.events.find(
    (e) =>
      e.time === segment.start &&
      e.pid === segment.pid &&
      (e.type === 'SELECTED' || e.type === 'PREEMPTED'),
  )
  if (match) return match.message

  return `${segment.pid} continues running at t=${segment.start} under ${algorithm.replace(/_/g, ' ')}.`
}

export { APIError }
