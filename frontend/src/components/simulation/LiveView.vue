<template>
  <div class="live-view">
    <div class="live-grid">
      <ReadyQueue :pids="readyPids" :color-for="colorFor" />

      <div class="cpu-panel" :class="{ idle: !runningPid }">
        <div class="cpu-title">CPU</div>
        <template v-if="runningPid">
          <div class="cpu-process" :style="{ color: colorFor(runningPid) }">{{ runningPid }}</div>
          <div class="cpu-sub">t = {{ currentTime }}</div>
          <div class="cpu-burst-bar">
            <div
              class="cpu-burst-fill"
              :style="{ width: burstProgress + '%', background: colorFor(runningPid) }"
            />
          </div>
        </template>
        <div v-else class="cpu-idle-text">IDLE</div>
      </div>

      <IOQueue :entries="ioEntries" />
    </div>

    <EventLog :entries="logEntries" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { SimulationResult } from '@/types'
import { processColor } from '@/api/scheduler'
import ReadyQueue from './ReadyQueue.vue'
import IOQueue from './IOQueue.vue'
import EventLog, { type LogEntry } from './EventLog.vue'

const props = defineProps<{
  result: SimulationResult
  currentTime: number
}>()

const pidIndex = computed(() => {
  const map: Record<string, number> = {}
  props.result.processes.forEach((p, i) => { map[p.pid] = i })
  return map
})
function colorFor(pid: string) {
  return processColor(pidIndex.value[pid] ?? 0)
}

const runningPid = computed<string | null>(() => {
  const { timeline } = props.result
  for (let i = timeline.length - 1; i >= 0; i--) {
    if (timeline[i].time <= props.currentTime) {
      return timeline[i].process_id === 'IDLE' ? null : timeline[i].process_id
    }
  }
  return null
})

// Best-effort IO detection: a process that has appeared on CPU, isn't running
// now, and hasn't completed is assumed to be in IO wait.
function isLikelyInIO(pid: string): boolean {
  const { timeline, total_time, processes } = props.result
  const proc = processes.find((p) => p.pid === pid)
  if (!proc) return false
  if (proc.completion_time <= props.currentTime) return false

  const cpuAppearances = timeline.filter((e) => e.process_id === pid)
  if (cpuAppearances.length === 0) return false

  const lastEntry = cpuAppearances[cpuAppearances.length - 1]
  const idx = timeline.indexOf(lastEntry)
  const lastEnd = idx + 1 < timeline.length ? timeline[idx + 1].time : total_time
  return lastEnd <= props.currentTime
}

const readyPids = computed(() => {
  const t = props.currentTime
  return props.result.processes
    .filter((p) => p.arrival_time <= t && p.completion_time > t && p.pid !== runningPid.value)
    .filter((p) => !isLikelyInIO(p.pid))
    .map((p) => p.pid)
})

const ioEntries = computed(() => {
  const t = props.currentTime
  return props.result.processes
    .filter((p) => p.arrival_time <= t && p.completion_time > t && p.pid !== runningPid.value)
    .filter((p) => isLikelyInIO(p.pid))
    .map((p) => ({ pid: p.pid }))
})

const burstProgress = computed(() => {
  // Rough visual indicator only — real progress needs per-burst remaining time
  return 55
})

const logEntries = computed<LogEntry[]>(() => {
  const { timeline, processes } = props.result
  const t = props.currentTime
  const entries: LogEntry[] = []

  for (const e of timeline) {
    if (e.time > t) break
    if (e.process_id === 'IDLE') {
      entries.push({ time: e.time, text: 'CPU idle — no ready processes', color: '#475569' })
    } else {
      entries.push({
        time: e.time,
        text: `${e.process_id} scheduled on CPU`,
        color: colorFor(e.process_id),
      })
    }
  }

  for (const p of processes) {
    if (p.arrival_time <= t) {
      entries.push({ time: p.arrival_time, text: `${p.pid} arrived`, color: colorFor(p.pid) })
    }
    if (p.completed && p.completion_time <= t) {
      entries.push({ time: p.completion_time, text: `${p.pid} completed`, color: '#22C55E' })
    }
  }

  return entries.sort((a, b) => b.time - a.time).slice(0, 10)
})
</script>

<style scoped>
.live-view { flex: 1; overflow: auto; padding: 20px; }

.live-grid {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  gap: 16px;
  align-items: start;
  margin-bottom: 20px;
}

.cpu-panel {
  background: #1A1D27;
  border: 1px solid #4F8EF7;
  border-radius: 10px;
  padding: 14px;
  min-width: 160px;
  text-align: center;
  box-shadow: 0 0 20px rgba(79,142,247,0.1);
}
.cpu-panel.idle { border-color: #2E3348; box-shadow: none; }

.cpu-title {
  font-size: 10px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #475569;
  margin-bottom: 12px;
}

.cpu-process {
  font-family: 'JetBrains Mono', monospace;
  font-size: 28px;
  font-weight: 800;
  margin-bottom: 4px;
  animation: pulse 1.6s ease-in-out infinite;
}

.cpu-sub { font-size: 10px; color: #475569; }
.cpu-idle-text { color: #475569; font-size: 13px; padding: 10px 0; }

.cpu-burst-bar {
  height: 4px;
  background: #222535;
  border-radius: 2px;
  margin-top: 10px;
  overflow: hidden;
}
.cpu-burst-fill { height: 100%; border-radius: 2px; transition: width 0.3s; }

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.72; }
}
</style>
