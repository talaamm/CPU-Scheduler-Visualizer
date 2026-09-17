<template>
  <div class="ranking-card">
    <div class="chart-title">Ranking by Average Waiting Time</div>
    <div
      v-for="(r, i) in ranked"
      :key="r.algorithm"
      class="rank-row"
    >
      <span class="rank-medal">{{ medal(i) }}</span>
      <span class="rank-name" :style="{ color: colorFor(r.algorithm) }">{{ nameFor(r.algorithm) }}</span>
      <span class="rank-stat">
        wait: {{ r.average_waiting_time.toFixed(2) }} ·
        tat: {{ r.average_turnaround_time.toFixed(2) }} ·
        cpu: {{ r.cpu_utilization.toFixed(0) }}%
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { SimulationResult } from '@/types'

const props = defineProps<{
  results: SimulationResult[]
  colorFor: (algorithm: string) => string
  nameFor: (algorithm: string) => string
}>()

const ranked = computed(() =>
  [...props.results].sort((a, b) => a.average_waiting_time - b.average_waiting_time),
)

function medal(i: number) {
  return ['🥇', '🥈', '🥉'][i] ?? `${i + 1}.`
}
</script>

<style scoped>
.ranking-card {
  background: #1A1D27;
  border: 1px solid #2E3348;
  border-radius: 10px;
  padding: 16px;
  grid-column: span 2;
}

.chart-title {
  font-size: 11px;
  font-weight: 600;
  color: #94A3B8;
  margin-bottom: 14px;
}

.rank-row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 0;
  border-bottom: 1px solid #2E3348;
  font-size: 12px;
}
.rank-row:last-child { border-bottom: none; }

.rank-medal { font-size: 16px; width: 24px; }
.rank-name { flex: 1; font-weight: 500; }
.rank-stat { font-family: 'JetBrains Mono', monospace; color: #475569; font-size: 11px; }
</style>
