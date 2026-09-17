<template>
  <div class="playback-bar">
    <button class="pb-btn" title="Reset (R)" @click="$emit('reset')">↺</button>
    <button class="pb-btn" title="Step back (←)" @click="$emit('stepBack')">←</button>
    <button class="pb-btn" :class="{ active: isPlaying }" title="Play/Pause (Space)" @click="$emit('toggle')">
      {{ isPlaying ? '⏸' : '▶' }}
    </button>
    <button class="pb-btn" title="Step forward (→)" @click="$emit('stepForward')">→</button>

    <span class="pb-time mono">t = {{ currentTime }} / {{ totalTime }}</span>

    <input
      class="pb-scrubber"
      type="range"
      min="0"
      :max="totalTime"
      :value="currentTime"
      @input="$emit('seek', +($event.target as HTMLInputElement).value)"
    />

    <select class="speed-select" :value="playSpeed" @change="$emit('setSpeed', +($event.target as HTMLSelectElement).value)">
      <option :value="0.5">0.5×</option>
      <option :value="1">1×</option>
      <option :value="2">2×</option>
      <option :value="4">4×</option>
    </select>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  currentTime: number
  totalTime: number
  isPlaying: boolean
  playSpeed: number
}>()

defineEmits<{
  toggle: []
  stepForward: []
  stepBack: []
  reset: []
  seek: [t: number]
  setSpeed: [s: number]
}>()
</script>

<style scoped>
.playback-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 20px;
  border-bottom: 1px solid #2E3348;
  background: #1A1D27;
  flex-shrink: 0;
}

.pb-btn {
  background: #222535;
  border: 1px solid #2E3348;
  border-radius: 6px;
  color: #F1F5F9;
  cursor: pointer;
  font-size: 13px;
  padding: 5px 10px;
  transition: all 0.12s;
  line-height: 1;
}
.pb-btn:hover { background: #2E3348; }
.pb-btn.active { background: #4F8EF7; border-color: #4F8EF7; }

.pb-time {
  font-family: 'JetBrains Mono', monospace;
  font-size: 12px;
  color: #94A3B8;
  min-width: 76px;
  flex-shrink: 0;
}

.pb-scrubber {
  flex: 1;
  -webkit-appearance: none;
  appearance: none;
  height: 4px;
  border-radius: 2px;
  background: #222535;
  outline: none;
  cursor: pointer;
}
.pb-scrubber::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 14px; height: 14px;
  border-radius: 50%;
  background: #4F8EF7;
  border: 2px solid #0F1117;
  box-shadow: 0 0 6px rgba(79,142,247,0.5);
}
.pb-scrubber::-moz-range-thumb {
  width: 14px; height: 14px;
  border-radius: 50%;
  background: #4F8EF7;
  border: 2px solid #0F1117;
  box-shadow: 0 0 6px rgba(79,142,247,0.5);
}

.speed-select {
  background: #222535;
  border: 1px solid #2E3348;
  border-radius: 6px;
  color: #94A3B8;
  font-size: 11px;
  padding: 4px 6px;
  outline: none;
  cursor: pointer;
  flex-shrink: 0;
}
</style>
