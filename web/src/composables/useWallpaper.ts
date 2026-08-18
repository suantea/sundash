import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useAppStore } from '../stores/app'

/**
 * Wallpaper state & background-luminance detection.
 * Extracted from Home.vue to keep the view focused on layout.
 */
export function useWallpaper() {
  const appStore = useAppStore()

  const wallpaperCopyright = ref('')
  const wallpaperReady = ref(false)

  // Preload image to prevent flicker.
  function preloadImage(url: string): Promise<void> {
    return new Promise((resolve) => {
      const img = new Image()
      img.onload = () => resolve()
      img.onerror = () => resolve() // Still resolve to show fallback
      img.src = url
    })
  }

  const fetchWallpaper = async () => {
    if (appStore.wallpaperType === 'bing') {
      try {
        const { api } = await import('../api')
        const res = await api.get('wallpaper/bing')
        if (res.data && res.data.images && res.data.images[0]) {
          const newUrl = res.data.images[0].url
          await preloadImage(newUrl)
          appStore.setWallpaperUrl(newUrl)
          wallpaperCopyright.value = res.data.images[0].copyright
          wallpaperReady.value = true
        }
      } catch (e) {
        console.error('Failed to fetch Bing wallpaper:', e)
        wallpaperReady.value = true
      }
    } else if (appStore.wallpaperType === 'custom' && appStore.wallpaperUrl) {
      await preloadImage(appStore.wallpaperUrl)
      wallpaperReady.value = true
    } else {
      wallpaperReady.value = true
    }
  }

  const wallpaperStyle = computed(() => {
    const style: Record<string, string> = {}

    if (appStore.wallpaperType === 'bing' || appStore.wallpaperType === 'custom') {
      if (appStore.wallpaperUrl) {
        style.backgroundImage = `url(${appStore.wallpaperUrl})`
        style.backgroundSize = 'cover'
        style.backgroundPosition = 'center'
        style.backgroundRepeat = 'no-repeat'
      }
    } else if (appStore.wallpaperType === 'gradient') {
      style.background = 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
    }

    if (appStore.wallpaperBlur > 0) {
      style.filter = `blur(${appStore.wallpaperBlur}px)`
    }
    if (appStore.wallpaperOpacity < 100) {
      style.opacity = (appStore.wallpaperOpacity / 100).toString()
    }
    return style
  })

  // Background luminance detection for text color adaptation.
  const isDarkBg = ref(false)

  function sampleLuminance(img: HTMLImageElement): boolean {
    const canvas = document.createElement('canvas')
    const ctx = canvas.getContext('2d')!
    canvas.width = 50
    canvas.height = 50
    ctx.drawImage(img, 0, 0, 50, 50)
    const data = ctx.getImageData(0, 0, 50, 50).data
    let sum = 0
    for (let i = 0; i < data.length; i += 4) {
      sum += 0.2126 * data[i] + 0.7152 * data[i + 1] + 0.0722 * data[i + 2]
    }
    return (sum / (data.length / 4)) < 128
  }

  async function detectBackground() {
    const wallpaper = document.querySelector('.wallpaper-bg') as HTMLElement
    if (!wallpaper) {
      isDarkBg.value = appStore.isDark
      return
    }
    const style = getComputedStyle(wallpaper)
    const bgImage = style.backgroundImage
    if (bgImage && bgImage !== 'none' && bgImage.includes('url')) {
      const urlMatch = bgImage.match(/url\(['"]?(.+?)['"]?\)/)
      if (urlMatch && urlMatch[1]) {
        const url = urlMatch[1]
        // Step 1: Try direct Canvas sampling (works for same-origin)
        try {
          const img = new Image()
          img.crossOrigin = 'anonymous'
          const loaded = await new Promise<boolean>((resolve) => {
            img.onload = () => resolve(true)
            img.onerror = () => resolve(false)
            img.src = url
          })
          if (loaded) {
            isDarkBg.value = sampleLuminance(img)
            return
          }
        } catch { /* fall through */ }
        // Step 2: CORS failed - try fetch as blob to get same-origin URL
        try {
          const resp = await fetch(url)
          const blob = await resp.blob()
          const blobUrl = URL.createObjectURL(blob)
          const img2 = new Image()
          const loaded2 = await new Promise<boolean>((resolve) => {
            img2.onload = () => resolve(true)
            img2.onerror = () => resolve(false)
            img2.src = blobUrl
          })
          if (loaded2) {
            isDarkBg.value = sampleLuminance(img2)
            URL.revokeObjectURL(blobUrl)
            return
          }
          URL.revokeObjectURL(blobUrl)
        } catch { /* fall through */ }
      }
    }
    isDarkBg.value = appStore.isDark
  }

  let bgObserver: MutationObserver | null = null
  let bgDetectTimer: ReturnType<typeof setTimeout> | null = null
  onMounted(() => {
    detectBackground()
    bgObserver = new MutationObserver(() => {
      if (bgDetectTimer) clearTimeout(bgDetectTimer)
      bgDetectTimer = setTimeout(() => detectBackground(), 200)
    })
    const wallpaper = document.querySelector('.wallpaper-bg')
    if (wallpaper) bgObserver.observe(wallpaper, { attributes: true, attributeFilter: ['style', 'class'] })
  })
  onUnmounted(() => {
    bgObserver?.disconnect()
    if (bgDetectTimer) clearTimeout(bgDetectTimer)
  })
  watch(() => appStore.isDark, () => { if (!appStore.wallpaperUrl) detectBackground() })
  watch(() => appStore.wallpaperUrl, () => { setTimeout(() => detectBackground(), 500) })

  return { wallpaperCopyright, wallpaperReady, isDarkBg, fetchWallpaper, wallpaperStyle }
}
