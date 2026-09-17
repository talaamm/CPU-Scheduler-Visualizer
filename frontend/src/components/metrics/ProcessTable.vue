<template>
  <div class="metrics-table">
    <table>
      <thead>
        <tr>
          <th>Process</th>
          <th>Arrival</th>
          <th>Start</th>
          <th>Completion</th>
          <th>Waiting</th>
          <th>Turnaround</th>
          <th>Response</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(p, i) in processes" :key="p.pid">
          <td class="pid-cell" :style="{ color: colorFor(i) }">{{ p.pid }}</td>
          <td>{{ p.arrival_time }}</td>
          <td>{{ p.start_time }}</td>
          <td>{{ p.completion_time }}</td>
          <td>{{ p.waiting_time }}</td>
          <td>{{ p.turnaround_time }}</td>
          <td>{{ p.response_time }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import type { ProcessResult } from '@/types'
import { processColor } from '@/api/scheduler'

defineProps<{
  processes: ProcessResult[]
}>()

function colorFor(idx: number) {
  return processColor(idx)
}
</script>

<style scoped>
.metrics-table {
  background: #1A1D27;
  border: 1px solid #2E3348;
  border-radius: 10px;
  overflow: hidden;
}

table { width: 100%; border-collapse: collapse; }

th {
  background: #222535;
  padding: 10px 14px;
  text-align: left;
  font-size: 10px;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: #475569;
  font-weight: 600;
}

td {
  padding: 10px 14px;
  border-top: 1px solid #2E3348;
  font-family: 'JetBrains Mono', monospace;
  font-size: 12px;
}

tr:hover td { background: #222535; }

.pid-cell { font-weight: 700; }
</style>
