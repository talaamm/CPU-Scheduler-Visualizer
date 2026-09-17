<template>
  <div class="algo-selector">
    <div
      v-for="algo in algorithms"
      :key="algo.id"
      class="algo-option"
      :class="{ selected: modelValue === algo.id }"
      @click="$emit('update:modelValue', algo.id)"
    >
      <div class="radio">
        <div class="radio-inner" />
      </div>
      <div class="algo-info">
        <span class="algo-name">{{ algo.name }}</span>
        <span class="algo-desc">{{ algo.description }}</span>
      </div>
      <span class="algo-badge" :class="algo.badgeClass">{{ algo.badge }}</span>
    </div>

    <!-- Quantum input — shown only for RR -->
    <Transition name="fade">
      <div v-if="modelValue === 'RR'" class="quantum-row">
        <label>Time Quantum</label>
        <input
          class="quantum-input"
          type="number"
          min="1"
          :value="quantum"
          @change="$emit('update:quantum', +($event.target as HTMLInputElement).value)"
        />
        <span class="quantum-unit">units</span>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import type { AlgorithmID } from '@/types'

defineProps<{
  modelValue: AlgorithmID
  quantum: number
}>()

defineEmits<{
  'update:modelValue': [id: AlgorithmID]
  'update:quantum': [q: number]
}>()

const algorithms = [
  { id: 'FCFS' as AlgorithmID,        name: 'First Come First Served', badge: 'Non-preemptive', badgeClass: '',          description: 'Processes run in arrival order.' },
  { id: 'SJF' as AlgorithmID,         name: 'Shortest Job First',       badge: 'Non-preemptive', badgeClass: '',          description: 'Shortest next burst runs first.' },
  { id: 'SRTF' as AlgorithmID,        name: 'Shortest Remaining Time',  badge: 'Preemptive',     badgeClass: 'preempt',   description: 'Preempts on shorter arrival.' },
  { id: 'RR' as AlgorithmID,          name: 'Round Robin',              badge: 'Quantum',        badgeClass: 'quantum',   description: 'Equal time slices, cyclic.' },
  { id: 'Priority_NP' as AlgorithmID, name: 'Priority Non-Preemptive', badge: 'Priority',       badgeClass: 'priority',  description: 'Highest priority runs to completion.' },
  { id: 'Priority_P' as AlgorithmID,  name: 'Priority Preemptive',     badge: 'Priority',       badgeClass: 'priority',  description: 'Preempts on higher priority arrival.' },
]
</script>

<style scoped>
.algo-selector { display: flex; flex-direction: column; gap: 1px; }

.algo-option {
  display: flex; align-items: center; gap: 8px;
  padding: 7px 8px; border-radius: 7px;
  cursor: pointer; transition: background 0.12s;
}
.algo-option:hover { background: #222535; }
.algo-option.selected { background: rgba(79,142,247,0.08); }

.radio {
  width: 13px; height: 13px; border-radius: 50%;
  border: 2px solid #2E3348; flex-shrink: 0;
  display: flex; align-items: center; justify-content: center;
  transition: border-color 0.15s;
}
.selected .radio { border-color: #4F8EF7; }
.radio-inner {
  width: 5px; height: 5px; border-radius: 50%; background: #4F8EF7;
  transform: scale(0); transition: transform 0.15s;
}
.selected .radio-inner { transform: scale(1); }

.algo-info { flex: 1; min-width: 0; }
.algo-name { display: block; font-size: 12px; font-weight: 500; color: #F1F5F9; }
.algo-desc { display: block; font-size: 10px; color: #475569; margin-top: 1px; }

.algo-badge {
  font-size: 9px; padding: 1px 6px; border-radius: 3px;
  background: rgba(255,255,255,0.05); color: #64748b;
  white-space: nowrap; flex-shrink: 0;
}
.algo-badge.preempt  { background: rgba(155,109,255,0.12); color: #9B6DFF; }
.algo-badge.quantum  { background: rgba(245,158,11,0.12);  color: #F59E0B; }
.algo-badge.priority { background: rgba(244,114,182,0.12); color: #F472B6; }

.quantum-row {
  display: flex; align-items: center; gap: 8px;
  padding: 6px 8px 6px 28px;
  background: rgba(245,158,11,0.05);
  border-radius: 0 0 7px 7px;
  border: 1px solid rgba(245,158,11,0.15);
  margin-top: 2px;
}
label { font-size: 10px; letter-spacing: 0.06em; text-transform: uppercase; color: #64748b; }
.quantum-input {
  background: #1A1D27; border: 1px solid #2E3348;
  border-radius: 5px; color: #F1F5F9;
  font-family: 'JetBrains Mono', monospace;
  font-size: 12px; padding: 3px 6px; width: 48px;
  text-align: center; outline: none;
  transition: border-color 0.15s;
}
.quantum-input:focus { border-color: #F59E0B; }
.quantum-unit { font-size: 10px; color: #475569; }

.fade-enter-active, .fade-leave-active { transition: opacity 0.15s, transform 0.15s; }
.fade-enter-from, .fade-leave-to { opacity: 0; transform: translateY(-4px); }
</style>
