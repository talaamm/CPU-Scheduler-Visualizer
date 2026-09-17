import { ref } from 'vue'

export interface TooltipState {
  visible: boolean
  x: number
  y: number
  header: string
  body: string
}

/**
 * Manages tooltip visibility and position.
 * Pass the returned handlers to mouseenter/mousemove/mouseleave on segments.
 */
export function useTooltip() {
  const tooltip = ref<TooltipState>({
    visible: false,
    x: 0,
    y: 0,
    header: '',
    body: '',
  })

  function show(e: MouseEvent, header: string, body: string) {
    tooltip.value = { visible: true, x: e.clientX + 14, y: e.clientY - 10, header, body }
  }

  function move(e: MouseEvent) {
    if (tooltip.value.visible) {
      tooltip.value.x = e.clientX + 14
      tooltip.value.y = e.clientY - 10
    }
  }

  function hide() {
    tooltip.value.visible = false
  }

  return { tooltip, show, move, hide }
}
