<template>
  <div class="event-log">
    <div class="event-log-title">Event Log</div>
    <TransitionGroup name="log" tag="div" class="log-entries">
      <div v-for="(entry, i) in entries" :key="entry.time + '-' + entry.text" class="log-entry">
        <span class="log-time">t={{ entry.time }}</span>
        <span :style="{ color: entry.color }">{{ entry.text }}</span>
      </div>
    </TransitionGroup>
  </div>
</template>

<script setup lang="ts">
export interface LogEntry {
  time: number
  text: string
  color: string
}

defineProps<{
  entries: LogEntry[]
}>()
</script>

<style scoped>
.event-log {
  background: #1A1D27;
  border: 1px solid #2E3348;
  border-radius: 10px;
  padding: 14px;
  max-height: 220px;
  overflow-y: auto;
}

.event-log-title {
  font-size: 10px;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: #475569;
  margin-bottom: 10px;
}

.log-entries { display: flex; flex-direction: column-reverse; gap: 0; }

.log-entry {
  display: flex;
  gap: 10px;
  padding: 4px 0;
  border-bottom: 1px solid #2E3348;
  font-size: 11px;
  font-family: 'JetBrains Mono', monospace;
  color: #94A3B8;
}
.log-entry:first-child { border-bottom: none; }

.log-time { color: #475569; flex-shrink: 0; width: 36px; }

.log-enter-active { transition: all 0.25s ease-out; }
.log-enter-from { opacity: 0; transform: translateY(-6px); }
</style>
