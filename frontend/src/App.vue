<template>
  <div class="app-shell">
    <!-- Top bar -->
    <div class="topbar">
      <div class="logo">
        <span class="logo-dot" />
        CPU Scheduler Visualizer
      </div>
      <div v-if="sim.result" class="algo-badge-top">{{ sim.result.algorithm.replace(/_/g, ' ') }}</div>
      <div class="spacer" />
      <div class="mode-toggle">
        <button class="mode-btn" :class="{ active: mode === 'single' }" @click="setMode('single')">Single</button>
        <button class="mode-btn" :class="{ active: mode === 'compare' }" @click="setMode('compare')">Compare</button>
      </div>
    </div>

    <div class="app-body">
      <LeftPanel
        :is-loading="sim.isLoading"
        :is-slow="sim.isSlowToRespond"
        :error-message="sim.error"
        @run="handleRun"
      />

      <div class="main-canvas">
        <ViewTabs
          v-if="mode === 'single' && sim.result"
          v-model="sim.activeView"
          :tabs="singleTabs"
        />

        <PlaybackBar
          v-if="mode === 'single' && sim.result"
          :current-time="sim.currentTime"
          :total-time="sim.totalTime"
          :is-playing="sim.isPlaying"
          :play-speed="sim.playbackSpeed"
          @toggle="sim.isPlaying ? sim.pause() : sim.play()"
          @step-forward="sim.stepForward()"
          @step-back="sim.stepBack()"
          @reset="sim.reset()"
          @seek="sim.seekTo($event)"
          @set-speed="sim.playbackSpeed = $event"
        />

        <div class="view-content">
          <template v-if="mode === 'compare'">
            <CompareView />
          </template>

          <template v-else-if="!sim.result">
            <div class="empty-state">
              <div class="empty-icon">⏱</div>
              <div class="empty-title">Ready to simulate</div>
              <div class="empty-sub">Configure processes, pick an algorithm, and click Run.</div>
            </div>
          </template>

          <template v-else>
            <GanttChart
              v-if="sim.activeView === 'gantt'"
              :result="sim.result"
              :current-time="sim.currentTime"
            />
            <LiveView
              v-else-if="sim.activeView === 'live'"
              :result="sim.result"
              :current-time="sim.currentTime"
            />
            <div v-else-if="sim.activeView === 'metrics'" class="metrics-view">
              <KPIGrid :result="sim.result" />
              <ProcessTable :processes="sim.result.processes" />
            </div>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useBuilderStore, useSimulationStore, useMetaStore } from '@/stores'
import { useKeyboard } from '@/composables/useKeyboard'

import LeftPanel from '@/components/ui/LeftPanel.vue'
import ViewTabs from '@/components/ui/ViewTabs.vue'
import PlaybackBar from '@/components/ui/PlaybackBar.vue'
import GanttChart from '@/components/gantt/GanttChart.vue'
import LiveView from '@/components/simulation/LiveView.vue'
import KPIGrid from '@/components/metrics/KPIGrid.vue'
import ProcessTable from '@/components/metrics/ProcessTable.vue'
import CompareView from '@/components/compare/CompareView.vue'

const builder = useBuilderStore()
const sim = useSimulationStore()
const meta = useMetaStore()

const mode = ref<'single' | 'compare'>('single')

const singleTabs = [
  { id: 'gantt', label: 'Gantt Chart' },
  { id: 'live', label: 'Live Simulation' },
  { id: 'metrics', label: 'Metrics' },
]

function setMode(m: 'single' | 'compare') {
  mode.value = m
  sim.pause()
}

async function handleRun() {
  await sim.run(builder)
  if (sim.result) {
    sim.activeView = 'gantt'
    setTimeout(() => sim.play(), 250)
  }
}

useKeyboard({
  onPlay: () => (sim.isPlaying ? sim.pause() : sim.play()),
  onStepForward: () => sim.stepForward(),
  onStepBack: () => sim.stepBack(),
  onReset: () => sim.reset(),
})

onMounted(async () => {
  try {
    await meta.fetchAll()
  } catch {
    // Backend not reachable yet — presets/algorithms fall back to defaults in components
  }
  handleRun() // auto-run on load for an immediate "wow" moment
})
</script>

<style>
:root {
  --bg: #0F1117;
  --surface: #1A1D27;
  --surface2: #222535;
  --border: #2E3348;
}

* { box-sizing: border-box; margin: 0; padding: 0; }

body {
  background: var(--bg);
  color: #F1F5F9;
  font-family: 'Inter', -apple-system, sans-serif;
  font-size: 13px;
}

#app, .app-shell {
  height: 100vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}
</style>

<style scoped>
.topbar {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 0 20px;
  height: 48px;
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
  background: var(--surface);
}

.logo {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 700;
  font-size: 14px;
  letter-spacing: -0.02em;
}
.logo-dot {
  width: 8px; height: 8px; border-radius: 50%;
  background: #4F8EF7;
  box-shadow: 0 0 8px #4F8EF7;
}

.algo-badge-top {
  font-size: 11px;
  color: #475569;
  margin-left: 4px;
}

.spacer { flex: 1; }

.mode-toggle { display: flex; gap: 4px; }
.mode-btn {
  background: none;
  border: 1px solid var(--border);
  border-radius: 6px;
  color: #94A3B8;
  cursor: pointer;
  font-size: 11px;
  font-weight: 500;
  padding: 5px 12px;
  transition: all 0.15s;
}
.mode-btn.active { background: #4F8EF7; border-color: #4F8EF7; color: white; }
.mode-btn:hover:not(.active) { border-color: #94A3B8; color: #F1F5F9; }

.app-body { display: flex; flex: 1; overflow: hidden; }

.main-canvas { flex: 1; display: flex; flex-direction: column; overflow: hidden; }

.view-content { flex: 1; overflow: hidden; display: flex; flex-direction: column; }

.metrics-view { flex: 1; overflow: auto; padding: 20px; }

.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: #475569;
}
.empty-icon { font-size: 48px; opacity: 0.3; }
.empty-title { font-size: 16px; font-weight: 600; color: #94A3B8; }
.empty-sub { font-size: 12px; }
</style>
