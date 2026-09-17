<template>
  <div class="left-panel">
    <div class="left-panel-scroll">
      <div class="section-label">Processes</div>
      <ProcessCard
        v-for="(proc, idx) in builder.processes"
        :key="proc.pid"
        :process="proc"
        :color="processColor(idx)"
        :is-expanded="builder.expandedPid === proc.pid"
        :show-priority="needsPriority"
        @toggle="builder.toggleExpanded(proc.pid)"
        @remove="builder.removeProcess(proc.pid)"
        @update="(field, value) => builder.updateProcess(proc.pid, { [field]: value })"
        @add-burst="(type) => builder.addBurst(proc.pid, type)"
        @remove-burst="(i) => builder.removeBurst(proc.pid, i)"
        @update-burst="(i, dur) => builder.updateBurstDuration(proc.pid, i, dur)"
      />
      <button class="add-process-btn" @click="builder.addProcess()">＋ Add Process</button>

      <div class="section-label">Presets</div>
      <div class="presets-grid">
        <button
          v-for="preset in meta.presets"
          :key="preset.id"
          class="preset-btn"
          @click="builder.loadPreset(preset)"
        >
          <strong>{{ preset.name }}</strong>{{ preset.description }}
        </button>
      </div>

      <div class="section-label">Algorithm</div>
      <AlgorithmSelector
        v-model="builder.selectedAlgorithm"
        :quantum="builder.quantum"
        @update:quantum="(q) => (builder.quantum = q)"
      />

      <div v-if="errorMessage" class="error-banner">{{ errorMessage }}</div>
    </div>

    <div class="left-panel-footer">
      <button class="run-btn" :disabled="isLoading" @click="$emit('run')">
        <span v-if="isLoading" class="spinner" />
        {{ isLoading ? 'Simulating…' : '▶  Run Simulation' }}
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useBuilderStore, useMetaStore } from '@/stores'
import { processColor } from '@/api/scheduler'
import ProcessCard from './ProcessCard.vue'
import AlgorithmSelector from './AlgorithmSelector.vue'

const builder = useBuilderStore()
const meta = useMetaStore()

defineProps<{
  isLoading: boolean
  errorMessage: string | null
}>()

defineEmits<{
  run: []
}>()

const needsPriority = builder.selectedAlgorithm === 'Priority_NP' || builder.selectedAlgorithm === 'Priority_P'
</script>

<style scoped>
.left-panel {
  width: 280px;
  flex-shrink: 0;
  border-right: 1px solid #2E3348;
  background: #1A1D27;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.left-panel-scroll {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  scrollbar-width: thin;
  scrollbar-color: #2E3348 transparent;
}

.left-panel-footer {
  padding: 12px 16px;
  border-top: 1px solid #2E3348;
}

.section-label {
  font-size: 10px;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: #475569;
  margin-bottom: 8px;
  margin-top: 16px;
}
.section-label:first-child { margin-top: 0; }

.add-process-btn {
  width: 100%;
  background: none;
  border: 1px dashed #2E3348;
  border-radius: 10px;
  color: #475569;
  cursor: pointer;
  font-size: 11px;
  padding: 8px;
  transition: all 0.15s;
  margin-top: 4px;
}
.add-process-btn:hover { border-color: #4F8EF7; color: #4F8EF7; }

.presets-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 6px; }
.preset-btn {
  background: #222535;
  border: 1px solid #2E3348;
  border-radius: 6px;
  color: #94A3B8;
  cursor: pointer;
  font-size: 10px;
  padding: 6px 8px;
  text-align: left;
  transition: all 0.15s;
  line-height: 1.3;
}
.preset-btn:hover { border-color: #4F8EF7; color: #F1F5F9; }
.preset-btn strong { display: block; font-size: 11px; color: #F1F5F9; margin-bottom: 2px; }

.error-banner {
  background: rgba(248,113,113,0.1);
  border: 1px solid rgba(248,113,113,0.3);
  border-radius: 6px;
  color: #fca5a5;
  font-size: 11px;
  margin-top: 12px;
  padding: 8px 10px;
}

.run-btn {
  width: 100%;
  background: #4F8EF7;
  border: none;
  border-radius: 10px;
  color: white;
  cursor: pointer;
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0.02em;
  padding: 11px;
  transition: all 0.15s;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}
.run-btn:hover:not(:disabled) {
  background: #6ba0fa;
  transform: translateY(-1px);
  box-shadow: 0 4px 16px rgba(79,142,247,0.35);
}
.run-btn:disabled { background: #222535; color: #475569; cursor: not-allowed; }

.spinner {
  width: 14px; height: 14px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }
</style>
