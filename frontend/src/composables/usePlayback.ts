import { ref, computed, onUnmounted } from 'vue'

/**
 * Encapsulates all play/pause/seek/step logic for the simulation timeline.
 * Used by PlaybackBar and the simulation stores.
 */
export function usePlayback(totalTime: ReturnType<typeof computed<number>>) {
  const currentTime = ref(0)
  const isPlaying = ref(false)
  const playSpeed = ref(1) // multiplier

  let intervalId: ReturnType<typeof setInterval> | null = null

  function msPerTick() {
    return Math.round(600 / playSpeed.value)
  }

  function clearTimer() {
    if (intervalId !== null) {
      clearInterval(intervalId)
      intervalId = null
    }
  }

  function play() {
    if (isPlaying.value) return
    if (currentTime.value >= totalTime.value) currentTime.value = 0
    isPlaying.value = true

    intervalId = setInterval(() => {
      if (currentTime.value >= totalTime.value) {
        pause()
        return
      }
      currentTime.value++
    }, msPerTick())
  }

  function pause() {
    isPlaying.value = false
    clearTimer()
  }

  function toggle() {
    isPlaying.value ? pause() : play()
  }

  function reset() {
    pause()
    currentTime.value = 0
  }

  function seekTo(t: number) {
    pause()
    currentTime.value = Math.max(0, Math.min(t, totalTime.value))
  }

  function stepForward() {
    pause()
    currentTime.value = Math.min(currentTime.value + 1, totalTime.value)
  }

  function stepBack() {
    pause()
    currentTime.value = Math.max(currentTime.value - 1, 0)
  }

  function setSpeed(v: number) {
    playSpeed.value = v
    if (isPlaying.value) {
      clearTimer()
      intervalId = setInterval(() => {
        if (currentTime.value >= totalTime.value) { pause(); return }
        currentTime.value++
      }, msPerTick())
    }
  }

  onUnmounted(clearTimer)

  return {
    currentTime,
    isPlaying,
    playSpeed,
    play,
    pause,
    toggle,
    reset,
    seekTo,
    stepForward,
    stepBack,
    setSpeed,
  }
}
