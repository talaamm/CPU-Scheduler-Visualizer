<template>
  <div class="chart-card">
    <div class="chart-title">{{ title }}</div>
    <div class="bar-chart">
      <div v-for="row in rows" :key="row.label" class="bar-row">
        <span class="bar-label">{{ row.label }}</span>
        <div class="bar-track">
          <div
            class="bar-fill"
            :style="{ width: (row.value / maxValue) * 100 + '%', background: row.color }"
          />
        </div>
        <span class="bar-value">{{ row.display ?? row.value.toFixed(1) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

export interface BarRow {
  label: string
  value: number
  color: string
  display?: string
}

const props = defineProps<{
  title: string
  rows: BarRow[]
}>()

const maxValue = computed(() => Math.max(...props.rows.map((r) => r.value), 1))
</script>

<style scoped>
.chart-card {
  background: #1A1D27;
  border: 1px solid #2E3348;
  border-radius: 10px;
  padding: 16px;
}

.chart-title {
  font-size: 11px;
  font-weight: 600;
  color: #94A3B8;
  margin-bottom: 14px;
}

.bar-chart { display: flex; flex-direction: column; gap: 8px; }

.bar-row { display: flex; align-items: center; gap: 10px; }

.bar-label {
  font-size: 11px;
  font-family: 'JetBrains Mono', monospace;
  width: 110px;
  flex-shrink: 0;
  color: #94A3B8;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.bar-track {
  flex: 1;
  height: 16px;
  background: #222535;
  border-radius: 4px;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  border-radius: 4px;
  transition: width 0.6s cubic-bezier(0.34, 1.2, 0.64, 1);
}

.bar-value {
  font-family: 'JetBrains Mono', monospace;
  font-size: 10px;
  color: #475569;
  width: 40px;
  text-align: right;
  flex-shrink: 0;
}
</style>
