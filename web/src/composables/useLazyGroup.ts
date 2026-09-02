import { ref, onMounted, onUnmounted, type Ref } from 'vue'

/**
 * useLazyGroup — IntersectionObserver-based lazy rendering for card groups.
 * Returns a `visible` ref that becomes true once the element enters the viewport.
 * Cards inside a group won't render until the group scrolls into view.
 */
export function useLazyGroup(elRef: Ref<HTMLElement | null>, options?: { rootMargin?: string; threshold?: number }) {
  const visible = ref(false)
  let observer: IntersectionObserver | null = null

  onMounted(() => {
    if (!elRef.value) return
    observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          visible.value = true
          // Once visible, disconnect — no need to keep observing
          observer?.disconnect()
          observer = null
        }
      },
      {
        rootMargin: options?.rootMargin ?? '200px',
        threshold: options?.threshold ?? 0,
      }
    )
    observer.observe(elRef.value)
  })

  onUnmounted(() => {
    observer?.disconnect()
  })

  return { visible }
}
