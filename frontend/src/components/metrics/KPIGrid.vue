<template>
  <div class="kpi-grid">
    <div class="kpi-card">
      <div class="kpi-label">CPU Utilization</div>
      <div class="kpi-value" style="color:#4F8EF7">{{ result.cpu_utilization.toFixed(1) }}%</div>
      <div class="kpi-sub">{{ algorithmLabel }}</div>
    </div>
    <div class="kpi-card">
      <div class="kpi-label">Avg Waiting Time</div>
      <div class="kpi-value" style="color:#9B6DFF">{{ result.average_waiting_time.toFixed(2) }}</div>
      <div class="kpi-sub">time units</div>
    </div>
    <div class="kpi-card">
      <div class="kpi-label">Avg Turnaround</div>
      <div class="kpi-value" style="color:#2DD4BF">{{ result.average_turnaround_time.toFixed(2) }}</div>
      <div class="kpi-sub">time units</div>
    </div>
    <div class="kpi-card">
      <div class="kpi-label">Avg Response Time</div>
      <div class="kpi-value" style="color:#F59E0B">{{ result.average_response_time.toFixed(2) }}</div>
      <div class="kpi-sub">time units</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { SimulationResult } from '@/types'

const props = defineProps<{
  result: SimulationResult
}>()

const algorithmLabel = computed(() => props.result.algorithm.replace(/_/g, ' '))
</script>

<style scoped>
.kpi-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  margin-bottom: 20px;
}

.kpi-card {
  background: #1A1D27;
  border: 1px solid #2E3348;
  border-radius: 10px;
  padding: 16px;
}

.kpi-label {
  font-size: 10px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #475569;
  margin-bottom: 8px;
}

.kpi-value {
  font-family: 'JetBrains Mono', monospace;
  font-size: 24px;
  font-weight: 700;
  margin-bottom: 4px;
}

.kpi-sub { font-size: 10px; color: #475569; }
</style>
