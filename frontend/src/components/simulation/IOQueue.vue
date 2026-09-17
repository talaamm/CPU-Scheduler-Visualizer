<template>
  <div class="queue-panel">
    <div class="queue-title">I/O Queue</div>
    <TransitionGroup name="queue" tag="div" class="queue-slots">
      <div v-for="entry in entries" :key="entry.pid" class="io-chip">
        <span class="io-dot" />
        <span class="io-pid">{{ entry.pid }}</span>
        <span class="io-remaining" v-if="entry.remainingIO != null">{{ entry.remainingIO }}t left</span>
      </div>
      <div v-if="entries.length === 0" key="__empty" class="queue-empty">Empty</div>
    </TransitionGroup>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  entries: Array<{ pid: string; remainingIO?: number }>
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

.io-chip {
  display: flex;
  align-items: center;
  gap: 8px;
  background: #222535;
  border-radius: 6px;
  border: 1px solid #2E3348;
  padding: 6px 10px;
}

.io-dot {
  width: 8px; height: 8px; border-radius: 50%;
  background: #2DD4BF;
  box-shadow: 0 0 6px rgba(45,212,191,0.5);
  flex-shrink: 0;
}

.io-pid {
  font-family: 'JetBrains Mono', monospace;
  font-size: 12px;
  font-weight: 700;
  color: #2DD4BF;
}

.io-remaining {
  font-size: 9px;
  color: #475569;
  margin-left: auto;
  font-family: 'JetBrains Mono', monospace;
}

.queue-enter-active { transition: all 0.3s ease-out; }
.queue-leave-active { transition: all 0.2s ease-in; position: absolute; }
.queue-enter-from { opacity: 0; transform: translateY(-10px); }
.queue-leave-to { opacity: 0; transform: translateY(10px); }
.queue-move { transition: transform 0.25s ease; }
</style>
