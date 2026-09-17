<template>
  <div class="gantt-wrap" ref="wrapRef">
    <div class="gantt-scroll">
      <svg
        :width="svgWidth"
        :height="svgHeight"
        class="gantt-svg"
        @mouseleave="tooltip.visible = false"
      >
        <!-- Time grid lines -->
        <g class="grid">
          <line
            v-for="t in gridTicks"
            :key="t"
            :x1="LEFT + t * pxPerUnit"
            :x2="LEFT + t * pxPerUnit"
            y1="0"
            :y2="svgHeight - AXIS_H"
            stroke="#1e2235"
            stroke-width="1"
          />
        </g>

        <!-- Per-process rows -->
        <g v-for="(proc, rowIdx) in result.processes" :key="proc.pid">
          <!-- Row background on hover -->
          <rect
            :x="LEFT"
            :y="rowIdx * (ROW_H + ROW_GAP)"
            :width="totalTime * pxPerUnit"
            :height="ROW_H"
            fill="transparent"
            class="row-bg"
          />

          <!-- PID label -->
          <text
            :x="LEFT - 8"
            :y="rowIdx * (ROW_H + ROW_GAP) + ROW_H / 2 + 4"
            text-anchor="end"
            class="pid-label"
            :fill="procColor(rowIdx)"
          >{{ proc.pid }}</text>

          <!-- Segments for this process -->
          <g v-for="(seg, si) in ganttSegments.filter(s => s.pid === proc.pid)" :key="si">
            <rect
              :x="LEFT + seg.start * pxPerUnit"
              :y="rowIdx * (ROW_H + ROW_GAP)"
              :width="Math.max(1, (seg.end - seg.start) * pxPerUnit - 1)"
              :height="ROW_H"
              :fill="procColor(rowIdx)"
              :fill-opacity="seg.end <= currentTime ? 0.88 : 0.15"
              rx="4"
              class="seg-bar"
              style="cursor:pointer"
              @mouseenter="onSegHover($event, seg, proc.pid)"
              @mousemove="onSegMove($event)"
              @mouseleave="tooltip.visible = false"
            />
          </g>
        </g>

        <!-- CPU timeline row -->
        <g :transform="`translate(0, ${cpuRowY})`">
          <text
            :x="LEFT - 8"
            :y="ROW_H / 2 + 4"
            text-anchor="end"
            class="cpu-label"
          >CPU</text>

          <g v-for="(seg, si) in ganttSegments" :key="si">
            <rect
              :x="LEFT + seg.start * pxPerUnit"
              y="0"
              :width="Math.max(1, (seg.end - seg.start) * pxPerUnit - 1)"
              :height="CPU_ROW_H"
              :fill="seg.isIdle ? '#1e2235' : procColorByPid(seg.pid)"
              :fill-opacity="seg.end <= currentTime ? (seg.isIdle ? 0.8 : 0.85) : 0.18"
              rx="3"
              class="seg-bar"
              style="cursor:pointer"
              @mouseenter="onSegHover($event, seg, seg.pid)"
              @mousemove="onSegMove($event)"
              @mouseleave="tooltip.visible = false"
            />
            <text
              v-if="(seg.end - seg.start) * pxPerUnit > 20"
              :x="LEFT + seg.start * pxPerUnit + (seg.end - seg.start) * pxPerUnit / 2"
              :y="CPU_ROW_H / 2 + 4"
              text-anchor="middle"
              class="seg-label"
              :fill="seg.isIdle ? '#475569' : '#fff'"
              :fill-opacity="seg.end <= currentTime ? 0.9 : 0.3"
            >{{ seg.isIdle ? '—' : seg.pid }}</text>
          </g>
        </g>

        <!-- Axis ticks -->
        <g :transform="`translate(0, ${axisY})`">
          <line :x1="LEFT" :x2="LEFT + totalTime * pxPerUnit" y1="0" y2="0" stroke="#2E3348" />
          <g v-for="t in axisTicks" :key="t">
            <line
              :x1="LEFT + t * pxPerUnit"
              :x2="LEFT + t * pxPerUnit"
              y1="0" y2="5"
              stroke="#2E3348"
            />
            <text
              :x="LEFT + t * pxPerUnit"
              y="16"
              text-anchor="middle"
              class="axis-tick"
            >{{ t }}</text>
          </g>
        </g>

        <!-- Time cursor -->
        <line
          :x1="LEFT + currentTime * pxPerUnit"
          :x2="LEFT + currentTime * pxPerUnit"
          y1="0"
          :y2="axisY"
          class="time-cursor"
          :style="{ transition: 'x1 0.08s, x2 0.08s' }"
        />
        <circle
          :cx="LEFT + currentTime * pxPerUnit"
          :cy="0"
          r="4"
          class="cursor-dot"
        />
      </svg>
    </div>

    <!-- Tooltip -->
    <Teleport to="body">
      <div
        v-if="tooltip.visible"
        class="gantt-tooltip"
        :style="{ left: tooltip.x + 'px', top: tooltip.y + 'px' }"
      >
        <div class="tt-header">{{ tooltip.header }}</div>
        <div class="tt-body">{{ tooltip.body }}</div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue'
import type { SimulationResult, GanttSegment } from '@/types'
import { toGanttSegments, processColor, explainDecision } from '@/api/scheduler'

const props = defineProps<{
  result: SimulationResult
  currentTime: number
}>()

// Layout constants
const LEFT = 50
const ROW_H = 26
const ROW_GAP = 6
const CPU_ROW_H = 20
const AXIS_H = 28
const CPU_SECTION_GAP = 14

const wrapRef = ref<HTMLElement>()

const ganttSegments = computed(() => toGanttSegments(props.result))
const totalTime = computed(() => props.result.total_time)

// Pixel density: fit to available width with min/max
const pxPerUnit = computed(() => {
  const availW = (wrapRef.value?.clientWidth ?? 700) - LEFT - 20
  return Math.max(16, Math.min(44, availW / totalTime.value))
})

const numRows = computed(() => props.result.processes.length)
const cpuRowY = computed(() => numRows.value * (ROW_H + ROW_GAP) + CPU_SECTION_GAP)
const axisY   = computed(() => cpuRowY.value + CPU_ROW_H + 6)
const svgHeight = computed(() => axisY.value + AXIS_H)
const svgWidth  = computed(() => LEFT + totalTime.value * pxPerUnit.value + 20)

const gridTicks = computed(() => {
  const step = Math.ceil(totalTime.value / 20)
  const ticks: number[] = []
  for (let t = 0; t <= totalTime.value; t += step) ticks.push(t)
  return ticks
})

const axisTicks = computed(() => {
  const step = Math.max(1, Math.ceil(totalTime.value / 20))
  const ticks: number[] = []
  for (let t = 0; t <= totalTime.value; t += step) ticks.push(t)
  return ticks
})

// Color helpers
const pidIndex = computed(() => {
  const map: Record<string, number> = {}
  props.result.processes.forEach((p, i) => { map[p.pid] = i })
  return map
})
function procColor(idx: number) { return processColor(idx) }
function procColorByPid(pid: string) {
  return processColor(pidIndex.value[pid] ?? 0)
}

// Tooltip
const tooltip = reactive({ visible: false, x: 0, y: 0, header: '', body: '' })

function onSegHover(e: MouseEvent, seg: GanttSegment, _pid: string) {
  tooltip.header = seg.isIdle
    ? `IDLE  t=${seg.start}–${seg.end}`
    : `${seg.pid}  ·  t=${seg.start} → ${seg.end}  (${seg.end - seg.start} units)`
  tooltip.body = explainDecision(seg, props.result.algorithm, props.result)
  tooltip.visible = true
  tooltip.x = e.clientX + 14
  tooltip.y = e.clientY - 8
}

function onSegMove(e: MouseEvent) {
  tooltip.x = e.clientX + 14
  tooltip.y = e.clientY - 8
}
</script>

<style scoped>
.gantt-wrap { flex: 1; overflow: hidden; display: flex; flex-direction: column; }
.gantt-scroll { flex: 1; overflow: auto; padding: 20px 20px 16px; }
.gantt-svg { display: block; overflow: visible; }

.pid-label, .cpu-label {
  font-family: 'JetBrains Mono', monospace;
  font-size: 11px; font-weight: 700;
}
.cpu-label { fill: #475569; }

.seg-bar { transition: fill-opacity 0.25s; }
.seg-bar:hover { filter: brightness(1.15); }

.seg-label { font-family: 'JetBrains Mono', monospace; font-size: 9px; font-weight: 700; pointer-events: none; }

.axis-tick { font-family: 'JetBrains Mono', monospace; font-size: 9px; fill: #475569; }

.time-cursor {
  stroke: #4F8EF7;
  stroke-width: 1.5;
  filter: drop-shadow(0 0 4px #4F8EF7);
  pointer-events: none;
}
.cursor-dot {
  fill: #4F8EF7;
  filter: drop-shadow(0 0 4px #4F8EF7);
  pointer-events: none;
}

.gantt-tooltip {
  position: fixed;
  background: #222535;
  border: 1px solid #2E3348;
  border-radius: 9px;
  padding: 10px 13px;
  max-width: 270px;
  pointer-events: none;
  z-index: 9999;
  box-shadow: 0 8px 28px rgba(0,0,0,0.55);
}
.tt-header {
  font-family: 'JetBrains Mono', monospace;
  font-size: 11px; font-weight: 700;
  color: #F1F5F9; margin-bottom: 5px;
}
.tt-body { font-size: 11px; color: #94A3B8; line-height: 1.5; }
</style>
