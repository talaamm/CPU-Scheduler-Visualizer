<template>
  <div class="compare-view">
    <div class="compare-algo-picker">
      <button
        v-for="algo in availableAlgorithms"
        :key="algo.id"
        class="algo-pill"
        :class="{ selected: store.selectedAlgorithms.includes(algo.id) }"
        @click="store.toggleAlgorithm(algo.id)"
      >
        {{ algo.name }}
      </button>
    </div>

    <button class="compare-run-btn" :disabled="store.isLoading" @click="store.runComparison(builder)">
      <span v-if="store.isLoading" class="spinner" />
      {{ store.isLoading ? 'Comparing…' : 'Compare Selected' }}
    </button>

    <div v-if="store.error" class="error-banner">{{ store.error }}</div>

    <div v-if="store.results.length === 0" class="empty-state">
      <div class="empty-icon">📊</div>
      <div class="empty-title">Click Compare to run all algorithms</div>
      <div class="empty-sub">Results will appear here.</div>
    </div>

    <div v-else class="compare-charts">
      <BarChart
        title="Avg Waiting Time (lower is better)"
        :rows="waitingRows"
      />
      <BarChart
        title="Avg Turnaround Time (lower is better)"
        :rows="turnaroundRows"
      />
      <BarChart
        title="Avg Response Time"
        :rows="responseRows"
      />
      <BarChart
        title="CPU Utilization (higher is better)"
        :rows="utilizationRows"
      />
      <RankingCard
        :results="store.results"
        :color-for="colorFor"
        :name-for="nameFor"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useCompareStore, useBuilderStore } from '@/stores'
import type { AlgorithmID } from '@/types'
import BarChart, { type BarRow } from './BarChart.vue'
import RankingCard from './RankingCard.vue'

const store = useCompareStore()
const builder = useBuilderStore()

const availableAlgorithms: Array<{ id: AlgorithmID; name: string }> = [
  { id: 'FCFS', name: 'FCFS' },
  { id: 'SJF', name: 'SJF' },
  { id: 'SRTF', name: 'SRTF' },
  { id: 'RR', name: 'Round Robin' },
  { id: 'Priority_NP', name: 'Priority NP' },
  { id: 'Priority_P', name: 'Priority P' },
]

function colorFor(algorithm: string) {
  return store.algorithmColors[algorithm as AlgorithmID] ?? '#4F8EF7'
}
function nameFor(algorithm: string) {
  return availableAlgorithms.find((a) => algorithm.startsWith(a.id))?.name ?? algorithm.replace(/_/g, ' ')
}

function buildRows(metric: (r: typeof store.results[number]) => number, display?: (v: number) => string): BarRow[] {
  return store.results.map((r) => ({
    label: r.algorithm.replace(/_/g, ' '),
    value: metric(r),
    color: colorFor(r.algorithm),
    display: display ? display(metric(r)) : undefined,
  }))
}

const waitingRows = computed(() => buildRows((r) => r.average_waiting_time))
const turnaroundRows = computed(() => buildRows((r) => r.average_turnaround_time))
const responseRows = computed(() => buildRows((r) => r.average_response_time))
const utilizationRows = computed(() => buildRows((r) => r.cpu_utilization, (v) => `${v.toFixed(0)}%`))
</script>

<style scoped>
.compare-view { flex: 1; overflow: auto; padding: 20px; }

.compare-algo-picker {
  display: flex;
  gap: 8px;
  margin-bottom: 20px;
  flex-wrap: wrap;
}

.algo-pill {
  border-radius: 20px;
  border: 1px solid #2E3348;
  color: #94A3B8;
  cursor: pointer;
  font-size: 11px;
  font-weight: 500;
  padding: 5px 12px;
  transition: all 0.15s;
  background: none;
}
.algo-pill.selected {
  border-color: #4F8EF7;
  color: #fff;
  background: rgba(79,142,247,0.15);
}

.compare-run-btn {
  background: #9B6DFF;
  border: none;
  border-radius: 10px;
  color: white;
  cursor: pointer;
  font-size: 12px;
  font-weight: 700;
  padding: 8px 18px;
  transition: all 0.15s;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 8px;
}
.compare-run-btn:hover:not(:disabled) { background: #b07fff; }
.compare-run-btn:disabled { opacity: 0.6; cursor: not-allowed; }

.spinner {
  width: 12px; height: 12px;
  border: 2px solid rgba(255,255,255,0.35);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }

.error-banner {
  background: rgba(248,113,113,0.1);
  border: 1px solid rgba(248,113,113,0.3);
  border-radius: 6px;
  color: #fca5a5;
  font-size: 11px;
  margin-bottom: 16px;
  padding: 8px 12px;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  color: #475569;
  min-height: 220px;
}
.empty-icon { font-size: 40px; opacity: 0.4; }
.empty-title { font-size: 14px; font-weight: 600; color: #94A3B8; }
.empty-sub { font-size: 11px; }

.compare-charts {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}
</style>
