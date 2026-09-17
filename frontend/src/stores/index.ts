import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api, toGanttSegments, buildSimulationFrames, processColor } from '@/api/scheduler'
import type {
  Process,
  AlgorithmID,
  SimulationResult,
  CompareResult,
  AlgorithmMeta,
  Preset,
  GanttSegment,
  SimulationFrame,
} from '@/types'

// ─── Process Builder Store ────────────────────────────────────────────────────

export const useBuilderStore = defineStore('builder', () => {
  const processes = ref<Process[]>([
    {
      pid: 'P1',
      arrival_time: 0,
      priority: 3,
      bursts: [
        { type: 'CPU', duration: 5 },
        { type: 'IO', duration: 3 },
        { type: 'CPU', duration: 4 },
      ],
    },
    {
      pid: 'P2',
      arrival_time: 1,
      priority: 1,
      bursts: [
        { type: 'CPU', duration: 3 },
        { type: 'IO', duration: 2 },
        { type: 'CPU', duration: 2 },
      ],
    },
    {
      pid: 'P3',
      arrival_time: 2,
      priority: 2,
      bursts: [
        { type: 'CPU', duration: 2 },
        { type: 'IO', duration: 1 },
        { type: 'CPU', duration: 3 },
      ],
    },
  ])

  const selectedAlgorithm = ref<AlgorithmID>('FCFS')
  const quantum = ref(2)
  const expandedPid = ref<string | null>(null)

  const processColors = computed(() =>
    Object.fromEntries(processes.value.map((p, i) => [p.pid, processColor(i)])),
  )

  function addProcess() {
    const n = processes.value.length + 1
    processes.value.push({
      pid: `P${n}`,
      arrival_time: 0,
      priority: n,
      bursts: [{ type: 'CPU', duration: 3 }],
    })
  }

  function removeProcess(pid: string) {
    processes.value = processes.value.filter((p) => p.pid !== pid)
  }

  function updateProcess(pid: string, updates: Partial<Process>) {
    const proc = processes.value.find((p) => p.pid === pid)
    if (proc) Object.assign(proc, updates)
  }

  function addBurst(pid: string, type: 'CPU' | 'IO') {
    const proc = processes.value.find((p) => p.pid === pid)
    if (proc) proc.bursts.push({ type, duration: type === 'CPU' ? 3 : 2 })
  }

  function removeBurst(pid: string, index: number) {
    const proc = processes.value.find((p) => p.pid === pid)
    if (proc) proc.bursts.splice(index, 1)
  }

  function updateBurstDuration(pid: string, index: number, duration: number) {
    const proc = processes.value.find((p) => p.pid === pid)
    if (proc && proc.bursts[index]) proc.bursts[index].duration = Math.max(1, duration)
  }

  function loadPreset(preset: Preset) {
    processes.value = preset.processes.map((p) => ({ ...p }))
  }

  function toggleExpanded(pid: string) {
    expandedPid.value = expandedPid.value === pid ? null : pid
  }

  return {
    processes,
    selectedAlgorithm,
    quantum,
    expandedPid,
    processColors,
    addProcess,
    removeProcess,
    updateProcess,
    addBurst,
    removeBurst,
    updateBurstDuration,
    loadPreset,
    toggleExpanded,
  }
})

// ─── Simulation Store ─────────────────────────────────────────────────────────

export const useSimulationStore = defineStore('simulation', () => {
  const result = ref<SimulationResult | null>(null)
  const isLoading = ref(false)
  const isSlowToRespond = ref(false) // true once a request has been pending a while — the free-tier backend may be waking from sleep
  const error = ref<string | null>(null)

  // Animation state
  const currentTime = ref(0)
  const isPlaying = ref(false)
  const playbackSpeed = ref(1) // multiplier: 0.5, 1, 2, 4
  const activeView = ref<'gantt' | 'live' | 'metrics' | 'compare'>('gantt')

  // Hover state
  const hoveredSegment = ref<GanttSegment | null>(null)

  const ganttSegments = computed<GanttSegment[]>(() =>
    result.value ? toGanttSegments(result.value) : [],
  )

  const frames = computed<SimulationFrame[]>(() =>
    result.value ? buildSimulationFrames(result.value, result.value.processes) : [],
  )

  const currentFrame = computed<SimulationFrame | null>(
    () => frames.value[Math.min(currentTime.value, frames.value.length - 1)] ?? null,
  )

  const totalTime = computed(() => result.value?.total_time ?? 0)

  let playInterval: ReturnType<typeof setInterval> | null = null

  async function run(builder: ReturnType<typeof useBuilderStore>) {
    pause() // clears any playback interval still running from a previous simulation
    isLoading.value = true
    isSlowToRespond.value = false
    error.value = null
    currentTime.value = 0

    // The free-tier backend spins down after 15 minutes idle and can take
    // 30-60s to wake up on the next request. Surface that after a few
    // seconds so it reads as "waking up", not "broken".
    const slowTimer = setTimeout(() => {
      isSlowToRespond.value = true
    }, 4000)

    try {
      result.value = await api.simulate({
        processes: builder.processes,
        algorithm: builder.selectedAlgorithm,
        quantum: builder.quantum,
      })
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Simulation failed'
    } finally {
      clearTimeout(slowTimer)
      isLoading.value = false
      isSlowToRespond.value = false
    }
  }

  function play() {
    if (isPlaying.value) return
    isPlaying.value = true
    const msPerTick = 600 / playbackSpeed.value

    playInterval = setInterval(() => {
      if (currentTime.value >= totalTime.value) {
        pause()
        return
      }
      currentTime.value++
    }, msPerTick)
  }

  function pause() {
    isPlaying.value = false
    if (playInterval) {
      clearInterval(playInterval)
      playInterval = null
    }
  }

  function reset() {
    pause()
    currentTime.value = 0
  }

  function stepForward() {
    pause()
    currentTime.value = Math.min(currentTime.value + 1, totalTime.value)
  }

  function stepBack() {
    pause()
    currentTime.value = Math.max(currentTime.value - 1, 0)
  }

  function seekTo(t: number) {
    currentTime.value = Math.max(0, Math.min(t, totalTime.value))
  }

  function setView(view: typeof activeView.value) {
    activeView.value = view
  }

  return {
    result,
    isLoading,
    isSlowToRespond,
    error,
    currentTime,
    isPlaying,
    playbackSpeed,
    activeView,
    hoveredSegment,
    ganttSegments,
    frames,
    currentFrame,
    totalTime,
    run,
    play,
    pause,
    reset,
    stepForward,
    stepBack,
    seekTo,
    setView,
  }
})

// ─── Compare Store ─────────────────────────────────────────────────────────────

export const useCompareStore = defineStore('compare', () => {
  const results = ref<SimulationResult[]>([])
  const selectedAlgorithms = ref<AlgorithmID[]>(['FCFS', 'SJF', 'SRTF', 'RR'])
  const isLoading = ref(false)
  const error = ref<string | null>(null)

  const algorithmColors: Record<AlgorithmID, string> = {
    FCFS: '#4F8EF7',
    SJF: '#9B6DFF',
    SRTF: '#2DD4BF',
    RR: '#F59E0B',
    Priority_NP: '#F472B6',
    Priority_P: '#22C55E',
  }

  const ranking = computed(() => {
    return [...results.value].sort(
      (a, b) => a.average_waiting_time - b.average_waiting_time,
    )
  })

  async function runComparison(builder: ReturnType<typeof useBuilderStore>) {
    if (selectedAlgorithms.value.length === 0) return
    isLoading.value = true
    error.value = null

    try {
      const res = await api.compare({
        processes: builder.processes,
        algorithms: selectedAlgorithms.value,
        quantum: builder.quantum,
      })
      results.value = res.results
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Comparison failed'
    } finally {
      isLoading.value = false
    }
  }

  function toggleAlgorithm(id: AlgorithmID) {
    const idx = selectedAlgorithms.value.indexOf(id)
    if (idx >= 0) {
      if (selectedAlgorithms.value.length > 1) {
        selectedAlgorithms.value.splice(idx, 1)
      }
    } else {
      selectedAlgorithms.value.push(id)
    }
  }

  return {
    results,
    selectedAlgorithms,
    isLoading,
    error,
    algorithmColors,
    ranking,
    runComparison,
    toggleAlgorithm,
  }
})

// ─── Metadata Store ────────────────────────────────────────────────────────────

export const useMetaStore = defineStore('meta', () => {
  const algorithms = ref<AlgorithmMeta[]>([])
  const presets = ref<Preset[]>([])

  async function fetchAll() {
    const [algoRes, presetRes] = await Promise.all([api.getAlgorithms(), api.getPresets()])
    algorithms.value = algoRes.algorithms
    presets.value = presetRes.presets
  }

  return { algorithms, presets, fetchAll }
})
