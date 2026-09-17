import { onMounted, onUnmounted } from 'vue'

interface KeyboardHandlers {
  onPlay?: () => void
  onStepForward?: () => void
  onStepBack?: () => void
  onReset?: () => void
}

/**
 * Binds keyboard shortcuts for playback control.
 * Space → play/pause, ←/→ → step, R → reset.
 * Ignored when focus is inside an <input> or <select>.
 */
export function useKeyboard(handlers: KeyboardHandlers) {
  function onKeyDown(e: KeyboardEvent) {
    const target = e.target as HTMLElement
    if (target.tagName === 'INPUT' || target.tagName === 'SELECT' || target.tagName === 'TEXTAREA') return

    switch (e.code) {
      case 'Space':
        e.preventDefault()
        handlers.onPlay?.()
        break
      case 'ArrowRight':
        e.preventDefault()
        handlers.onStepForward?.()
        break
      case 'ArrowLeft':
        e.preventDefault()
        handlers.onStepBack?.()
        break
      case 'KeyR':
        handlers.onReset?.()
        break
    }
  }

  onMounted(() => window.addEventListener('keydown', onKeyDown))
  onUnmounted(() => window.removeEventListener('keydown', onKeyDown))
}
