# Integration Guide

## Running the Backend

```bash
# From your project root (cpu-scheduler-visualizer/)
# Add the REST server file to cmd/server/main.go
# Then:
go run cmd/server/main.go
# API available at http://localhost:8080
```

## Generated Files

### backend/cmd/server/main.go
Complete Go REST API server. Drop into your existing project at cmd/server/main.go
(same path used in the Project-Plan.md structure). Import paths already match
your go.mod module: github.com/talaamm/cpu-scheduler-visualizer

Endpoints implemented:
- POST /api/simulate   — run one algorithm
- POST /api/compare    — run several algorithms on the same workload
- GET  /api/algorithms — metadata for all 6 algorithms
- GET  /api/presets    — 4 built-in scenarios (classic, io_heavy, starvation, burst)
- GET  /health         — liveness check

CORS is wide open (Access-Control-Allow-Origin: *) for local dev — tighten
this before deploying.

### frontend/src/types/index.ts
Full TypeScript type definitions for all API contracts.

### frontend/src/api/scheduler.ts  
API client layer + derived data helpers:
- api.simulate() — POST /api/simulate
- api.compare()  — POST /api/compare
- api.getAlgorithms() — GET /api/algorithms  
- api.getPresets()    — GET /api/presets
- toGanttSegments()   — converts timeline to Gantt bars
- buildSimulationFrames() — reconstructs per-tick state for animation
- explainDecision()   — generates "why?" tooltips

### frontend/src/stores/index.ts
Pinia stores:
- useBuilderStore — process list, algorithm, quantum
- useSimulationStore — result, animation state, playback
- useCompareStore — multi-algo comparison
- useMetaStore — cached algorithm/preset metadata

### frontend/index.html  
STANDALONE complete app — runs without a build step.
Opens in browser directly. Has client-side simulation fallback
so it works even without the Go backend running.

## Vue 3 Component Structure (for build-step version)

All files below are generated and working — copy `frontend/` into your repo,
run `npm install`, then `npm run dev`. Note: the entry HTML for the Vite build
is `index-vite.html` (rename to `index.html` if you don't need the standalone
demo file alongside it — they can't both be named `index.html`).

src/
├── components/
│   ├── builder/
│   │   ├── ProcessCard.vue        ✓ expandable process editor w/ burst chips
│   │   └── AlgorithmSelector.vue  ✓ radio group + quantum input
│   ├── gantt/
│   │   └── GanttChart.vue         ✓ SVG Gantt, animated cursor, hover tooltips
│   ├── simulation/
│   │   ├── LiveView.vue           ✓ composes CPU/Ready/IO/Log
│   │   ├── ReadyQueue.vue         ✓ animated process chips
│   │   ├── IOQueue.vue            ✓ IO wait display
│   │   └── EventLog.vue           ✓ timestamped decision log
│   ├── metrics/
│   │   ├── KPIGrid.vue            ✓ 4 key metric cards
│   │   └── ProcessTable.vue       ✓ per-process metrics table
│   ├── compare/
│   │   ├── CompareView.vue        ✓ algorithm picker + run + charts
│   │   ├── BarChart.vue           ✓ horizontal comparison bars
│   │   └── RankingCard.vue        ✓ medal-ranked results
│   └── ui/
│       ├── PlaybackBar.vue        ✓ play/pause/scrub controls
│       ├── ViewTabs.vue           ✓ Gantt/Live/Metrics tab switcher
│       └── LeftPanel.vue          ✓ sidebar shell
├── composables/
│   ├── usePlayback.ts             ✓ play/pause/seek logic (standalone, optional)
│   ├── useKeyboard.ts             ✓ Space, ←, →, R shortcuts
│   └── useTooltip.ts              ✓ hover tooltip positioning (standalone, optional)
├── stores/index.ts                ✓ Pinia: builder, simulation, compare, meta
├── types/index.ts                 ✓ full API + UI type definitions
├── api/scheduler.ts                ✓ HTTP client + Gantt/frame derivation helpers
├── App.vue                         ✓ root layout wiring everything together
└── main.ts                         ✓ Vue + Pinia bootstrap

Note: GanttChart.vue and LiveView.vue use computed properties inline rather
than the usePlayback/useTooltip composables (the simulation store already
owns playback state). The composables are provided as clean, reusable
alternatives if you refactor playback out of the store later.

## Keyboard Shortcuts
Space        — Play / Pause
ArrowRight   — Step forward one tick
ArrowLeft    — Step back one tick
R            — Reset to t=0

## Environment Variables
VITE_API_URL=http://localhost:8080   # default
