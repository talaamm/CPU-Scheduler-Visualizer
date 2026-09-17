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
import { processColor, buildSimulationFrames } from '@/api/scheduler'
import ReadyQueue from './ReadyQueue.vue'
import IOQueue from './IOQueue.vue'
import EventLog from './EventLog.vue'

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

// Frames are derived from the engine's real per-tick snapshots/events — the
// ready queue, I/O queue, and burst progress below all reflect actual
// simulation state, not a post-hoc guess.
const frames = computed(() => buildSimulationFrames(props.result, props.result.processes))
const frame = computed(() => frames.value[Math.min(props.currentTime, frames.value.length - 1)] ?? null)

const runningPid = computed(() => frame.value?.runningPid ?? null)
const readyPids = computed(() => frame.value?.readyQueue ?? [])
const ioEntries = computed(() =>
  (frame.value?.ioQueue ?? []).map((e) => ({ pid: e.pid, remainingIO: e.remainingIO })),
)

const burstProgress = computed(() => {
  const f = frame.value
  if (!f || !f.runningBurstTotal) return 0
  const elapsed = f.runningBurstTotal - f.runningRemaining
  return Math.min(100, Math.max(0, (elapsed / f.runningBurstTotal) * 100))
})

const logEntries = computed(() => frame.value?.eventLog.slice(0, 10) ?? [])
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
