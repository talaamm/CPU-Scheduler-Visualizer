<template>
  <div class="queue-panel">
    <div class="queue-title">Ready Queue</div>
    <TransitionGroup name="queue" tag="div" class="queue-slots">
      <div
        v-for="pid in pids"
        :key="pid"
        class="queue-chip"
        :style="{ '--chip-color': colorFor(pid) }"
      >
        <span class="chip-dot" />
        <span class="chip-pid">{{ pid }}</span>
      </div>
      <div v-if="pids.length === 0" key="__empty" class="queue-empty">Empty</div>
    </TransitionGroup>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  pids: string[]
  colorFor: (pid: string) => string
}>()
</script>

<style scoped>
.queue-panel {
  background: #1A1D27;
  border: 1px solid #2E3348;
  border-radius: 10px;
  padding: 14px;
}

.queue-title {
  font-size: 10px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #475569;
  margin-bottom: 12px;
}

.queue-slots {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-height: 60px;
}

.queue-empty {
  color: #475569;
  font-size: 11px;
  font-style: italic;
  text-align: center;
  padding: 16px 0;
}

.queue-chip {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #222535;
  border-radius: 6px;
  border: 1px solid #2E3348;
  padding: 6px 10px;
}

.chip-dot {
  width: 8px; height: 8px; border-radius: 50%;
  background: var(--chip-color);
  box-shadow: 0 0 6px color-mix(in srgb, var(--chip-color) 50%, transparent);
  flex-shrink: 0;
}

.chip-pid {
  font-family: 'JetBrains Mono', monospace;
  font-size: 12px;
  font-weight: 700;
  color: var(--chip-color);
}

.queue-enter-active { transition: all 0.3s ease-out; }
.queue-leave-active { transition: all 0.2s ease-in; position: absolute; }
.queue-enter-from { opacity: 0; transform: translateX(-16px); }
.queue-leave-to { opacity: 0; transform: translateX(16px); }
.queue-move { transition: transform 0.25s ease; }
</style>
