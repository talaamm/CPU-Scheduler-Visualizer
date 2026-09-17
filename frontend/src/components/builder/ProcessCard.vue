<template>
  <div class="process-card" :class="{ expanded: isExpanded }" :style="{ '--proc-color': color }">
    <!-- Header row -->
    <div class="process-header" @click="$emit('toggle')">
      <span class="proc-dot" />
      <span class="proc-pid">{{ process.pid }}</span>

      <div class="burst-preview">
        <span
          v-for="(b, i) in process.bursts"
          :key="i"
          class="burst-chip"
          :class="b.type.toLowerCase()"
        >{{ b.type[0] }}{{ b.duration }}</span>
      </div>

      <button class="btn-icon remove" @click.stop="$emit('remove')" title="Remove process">✕</button>
    </div>

    <!-- Expanded body -->
    <Transition name="expand">
      <div v-if="isExpanded" class="process-body">
        <!-- Arrival + Priority -->
        <div class="field-row">
          <div class="field-group">
            <label>Arrival</label>
            <input
              class="field-input"
              type="number"
              min="0"
              :value="process.arrival_time"
              @change="update('arrival_time', +($event.target as HTMLInputElement).value)"
            />
          </div>
          <div v-if="showPriority" class="field-group">
            <label>Priority <span class="hint">(lower = higher)</span></label>
            <input
              class="field-input"
              type="number"
              min="1"
              :value="process.priority"
              @change="update('priority', +($event.target as HTMLInputElement).value)"
            />
          </div>
        </div>

        <!-- Burst editor -->
        <div class="bursts-label">
          <label>Bursts</label>
          <span class="hint">Click duration to edit</span>
        </div>
        <div class="burst-editor-row">
          <div
            v-for="(b, i) in process.bursts"
            :key="i"
            class="burst-edit-chip"
            :class="b.type.toLowerCase()"
          >
            <span class="burst-type-tag">{{ b.type }}</span>
            <input
              class="burst-dur-input"
              type="number"
              min="1"
              :value="b.duration"
              @change="updateBurst(i, +($event.target as HTMLInputElement).value)"
              @click.stop
            />
            <button class="burst-del" @click="$emit('removeBurst', i)" title="Remove burst">×</button>
          </div>

          <button class="add-burst-btn" @click="$emit('addBurst', 'CPU')">+ CPU</button>
          <button class="add-burst-btn io" @click="$emit('addBurst', 'IO')">+ IO</button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import type { Process } from '@/types'

const props = defineProps<{
  process: Process
  color: string
  isExpanded: boolean
  showPriority: boolean
}>()

const emit = defineEmits<{
  toggle: []
  remove: []
  update: [field: string, value: number | string]
  addBurst: [type: 'CPU' | 'IO']
  removeBurst: [index: number]
  updateBurst: [index: number, duration: number]
}>()

function update(field: string, value: number | string) {
  emit('update', field, value)
}

function updateBurst(index: number, duration: number) {
  emit('updateBurst', index, Math.max(1, duration))
}
</script>

<style scoped>
.process-card {
  background: #222535;
  border: 1px solid #2E3348;
  border-radius: 10px;
  margin-bottom: 8px;
  overflow: hidden;
  transition: border-color 0.15s;
}
.process-card:hover { border-color: #3d4460; }
.process-card.expanded { border-color: var(--proc-color, #4F8EF7); }

.process-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 10px;
  cursor: pointer;
  user-select: none;
}

.proc-dot {
  width: 8px; height: 8px; border-radius: 50%;
  background: var(--proc-color);
  box-shadow: 0 0 6px color-mix(in srgb, var(--proc-color) 40%, transparent);
  flex-shrink: 0;
}

.proc-pid {
  font-family: 'JetBrains Mono', monospace;
  font-size: 12px;
  font-weight: 700;
  color: var(--proc-color);
  min-width: 28px;
}

.burst-preview { display: flex; gap: 3px; flex: 1; flex-wrap: wrap; }

.burst-chip {
  font-size: 9px;
  font-family: 'JetBrains Mono', monospace;
  padding: 1px 5px;
  border-radius: 3px;
}
.burst-chip.cpu { background: rgba(79,142,247,0.15); color: #4F8EF7; border: 1px solid rgba(79,142,247,0.3); }
.burst-chip.io  { background: rgba(45,212,191,0.12); color: #2DD4BF; border: 1px solid rgba(45,212,191,0.25); }

.btn-icon {
  background: none; border: none; cursor: pointer;
  border-radius: 4px; line-height: 1;
  padding: 2px 5px;
  font-size: 13px;
  transition: all 0.12s;
}
.btn-icon.remove { color: #475569; }
.btn-icon.remove:hover { color: #f87171; background: rgba(248,113,113,0.1); }

.process-body {
  padding: 10px 12px 12px;
  border-top: 1px solid #2E3348;
}

.field-row { display: flex; gap: 10px; margin-bottom: 10px; }
.field-group { flex: 1; display: flex; flex-direction: column; gap: 4px; }

label {
  font-size: 10px;
  letter-spacing: 0.07em;
  text-transform: uppercase;
  color: #64748b;
}
.hint { text-transform: none; font-size: 9px; color: #475569; letter-spacing: 0; margin-left: 4px; }

.field-input {
  background: #1A1D27;
  border: 1px solid #2E3348;
  border-radius: 6px;
  color: #F1F5F9;
  font-family: 'JetBrains Mono', monospace;
  font-size: 12px;
  padding: 5px 8px;
  outline: none;
  transition: border-color 0.15s;
}
.field-input:focus { border-color: #4F8EF7; }

.bursts-label { display: flex; align-items: baseline; gap: 6px; margin-bottom: 6px; }

.burst-editor-row { display: flex; flex-wrap: wrap; gap: 5px; align-items: center; }

.burst-edit-chip {
  display: flex; align-items: center; gap: 3px;
  border-radius: 6px; padding: 3px 5px 3px 7px;
  font-size: 10px;
  font-family: 'JetBrains Mono', monospace;
}
.burst-edit-chip.cpu { background: rgba(79,142,247,0.1); border: 1px solid rgba(79,142,247,0.28); }
.burst-edit-chip.io  { background: rgba(45,212,191,0.09); border: 1px solid rgba(45,212,191,0.22); }

.burst-type-tag { color: #64748b; font-size: 9px; }

.burst-dur-input {
  background: transparent; border: none;
  color: #F1F5F9; font-family: 'JetBrains Mono', monospace;
  font-size: 12px; width: 26px; text-align: center; outline: none;
}

.burst-del {
  background: none; border: none; cursor: pointer;
  color: #475569; font-size: 12px; padding: 0 2px; line-height: 1;
  transition: color 0.12s;
}
.burst-del:hover { color: #f87171; }

.add-burst-btn {
  background: none;
  border: 1px dashed #2E3348;
  border-radius: 5px;
  color: #475569; cursor: pointer;
  font-size: 10px; padding: 3px 8px;
  transition: all 0.15s;
}
.add-burst-btn:hover { border-color: #4F8EF7; color: #4F8EF7; }
.add-burst-btn.io:hover { border-color: #2DD4BF; color: #2DD4BF; }

/* Expand transition */
.expand-enter-active, .expand-leave-active {
  transition: max-height 0.22s ease, opacity 0.18s ease;
  overflow: hidden;
  max-height: 300px;
}
.expand-enter-from, .expand-leave-to { max-height: 0; opacity: 0; }
</style>
