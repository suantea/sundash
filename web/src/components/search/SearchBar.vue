<template>
  <div class="search-bar" :class="{ 'search-dark-bg': isDarkBg }">
    <div class="search-input-wrap" :style="searchWrapStyle">
      <Icon :icon="currentEngine.icon" :size="16" class="search-icon" :style="searchIconStyle" />
      <input
        v-model="query"
        class="search-input"
        :style="searchInputStyle"
        :placeholder="currentEngine.name + ' ' + t('search.placeholder')"
        @keydown.enter="handleSearch"
        @input="onSearchInput">
      <kbd v-if="!query" class="search-kbd">Enter</kbd>
      <button v-else class="search-clear" @click="query = ''; searchResults = []">
        <Icon icon="mdi:close" :size="14" />
      </button>
      <div class="search-divider"></div>
      <n-dropdown :options="engineOptions" @select="selectEngine" trigger="click" placement="bottom-end">
        <button class="engine-select">
          <span>{{ currentEngine.name }}</span>
          <Icon icon="mdi:chevron-down" :size="12" />
        </button>
      </n-dropdown>
    </div>
    <!-- Bookmarks search results -->
    <div v-if="searchResults.length > 0" class="search-results">
      <div v-for="item in searchResults" :key="item.id" class="search-result-item" @click="openBookmark(item.url)">
        <Icon :icon="item.icon || 'mdi:bookmark-outline'" :size="16" class="result-icon" />
        <div class="result-info">
          <div class="result-title">{{ item.title }}</div>
          <div class="result-url">{{ item.url }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { NDropdown } from 'naive-ui'
import { useI18n } from 'vue-i18n'
import { ref, computed, h, onMounted, onUnmounted, watch } from 'vue'
import { Icon } from '@iconify/vue'
import { useAppStore } from '../../stores/app'
import { usePanelStore } from '../../stores/panel'
import { useRouter } from 'vue-router'
import { useToast } from 'vue-toast-notification'
import axios from '@/api'

const { t } = useI18n()
const appStore = useAppStore()
const panelStore = usePanelStore()
const toast = useToast()
const router = useRouter()

const query = ref('')
const isDarkBg = ref(false)

// Search bar styles synced with card settings
const searchWrapStyle = computed(() => {
  const style: Record<string, string> = {}
  if (appStore.bgColor) {
    style['--custom-bg'] = appStore.bgColor
    style.background = appStore.bgColor
    style.border = 'none'
  }
  return style
})

const searchIconStyle = computed(() => {
  const style: Record<string, string> = {}
  if (appStore.textColor) {
    style.color = appStore.textColor
  }
  return style
})

const searchInputStyle = computed(() => {
  const style: Record<string, string> = {}
  if (appStore.textColor) {
    style.color = appStore.textColor
  }
  return style
})

const engines = [
  { name: 'Google', url: 'https://www.google.com/search?q=', icon: 'mdi:google' },
  { name: 'Bing', url: 'https://www.bing.com/search?q=', icon: 'mdi:microsoft-bing' },
  { name: 'Baidu', url: 'https://www.baidu.com/s?wd=', icon: 'mdi:alpha-b-circle' },
  { name: 'DuckDuckGo', url: 'https://duckduckgo.com/?q=', icon: 'mdi:duck' },
]

// Initialize from saved setting
const savedEngineIndex = engines.findIndex(e => e.name === appStore.searchEngine)
const currentEngineIndex = ref(savedEngineIndex >= 0 ? savedEngineIndex : 0)
const currentEngine = ref(engines[currentEngineIndex.value])

const engineOptions = engines.map((e, i) => ({
  label: e.name,
  key: i,
  icon: () => h(Icon, { icon: e.icon, width: 16 }),
}))

// Search results for bookmarks from backend
const searchResults = ref<Array<{ id: string; title: string; url: string; icon?: string }>>([])

// Debounce timer for suggestions
let suggestionTimer: NodeJS.Timeout | null = null

function selectEngine(index: number) {
  currentEngineIndex.value = index
  currentEngine.value = engines[index]
  // Save to user account
  appStore.setSearchEngine(engines[index].name)
}

function handleSearch() {
  if (query.value.trim()) {
    window.open(currentEngine.value.url + encodeURIComponent(query.value.trim()), '_blank')
  }
}

// Search bookmarks via backend API
async function onSearchInput() {
  const q = query.value.trim()
  if (!q) {
    searchResults.value = []
    return
  }

  // Clear previous timer
  if (suggestionTimer) clearTimeout(suggestionTimer)

  // Debounce to avoid too many requests
  suggestionTimer = setTimeout(async () => {
    try {
      const token = localStorage.getItem('sundash-token')
      const res = await axios.get('/api/search', {
        params: { q, limit: 8 },
        headers: token ? { Authorization: `Bearer ${token}` } : {}
      })
      if (res.status === 200) {
        const data = res.data
        // Transform backend results to frontend format
        searchResults.value = data.results.map((r: any) => ({
          id: r.id,
          title: r.title,
          url: r.url,
          icon: r.icon
        }))
      }
    } catch (err: any) {
      console.error('Search error:', err)
      // Fallback to local search if backend fails
      fallbackLocalSearch(q)
    }
  }, 300)
}

// Fallback to local search (original implementation)
function fallbackLocalSearch(q: string) {
  const results: Array<{ id: string; title: string; url: string; icon?: string }> = []
  for (const group of panelStore.groups) {
    if (group.cards) {
      for (const card of group.cards) {
        if (
          card.title.toLowerCase().includes(q) ||
          card.url.toLowerCase().includes(q) ||
          (card.description && card.description.toLowerCase().includes(q))
        ) {
          results.push({
            id: card.id,
            title: card.title,
            url: card.url_internal || card.url,
            icon: card.icon,
          })
        }
      }
    }
  }
  searchResults.value = results.slice(0, 8)
}

// Get search suggestions from backend
async function getSuggestions(prefix: string) {
  if (!prefix) return []
  try {
    const token = localStorage.getItem('sundash-token')
    const res = await axios.get('/api/search/suggestions', {
      params: { q: prefix, limit: 5 },
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    if (res.status === 200) {
      return res.data.suggestions || []
    }
  } catch (err) {
    console.error('Suggestions error:', err)
  }
  return []
}

// Detect wallpaper background luminance
function detectBackground() {
  const wallpaper = document.querySelector('.wallpaper-bg') as HTMLElement
  if (!wallpaper) return

  const style = getComputedStyle(wallpaper)
  const bgImage = style.backgroundImage

  // If there's a wallpaper image, sample its luminance
  if (bgImage && bgImage !== 'none' && bgImage.includes('url')) {
    const urlMatch = bgImage.match(/url\(['"]?(.+?)['"]?\)/)
    if (urlMatch && urlMatch[1]) {
      const img = new Image()
      img.crossOrigin = 'anonymous'
      img.onload = () => {
        try {
          const canvas = document.createElement('canvas')
          const ctx = canvas.getContext('2d')
          canvas.width = 50
          canvas.height = 50
          ctx!.drawImage(img, 0, 0, 50, 50)
          const data = ctx!.getImageData(0, 0, 50, 50).data
          let sum = 0
          for (let i = 0; i < data.length; i += 4) {
            // ITU-R BT.709 luminance
            sum += 0.2126 * data[i] + 0.7152 * data[i + 1] + 0.0722 * data[i + 2]
          }
          const avgLuminance = sum / (data.length / 4)
          isDarkBg.value = avgLuminance < 128
        } catch {
          isDarkBg.value = false
        }
      }
      img.onerror = () => { isDarkBg.value = false }
      img.src = urlMatch[1]
      return
    }
  }

  // Check dark theme
  isDarkBg.value = appStore.isDark
}

let observer: MutationObserver | null = null

onMounted(() => {
  detectBackground()
  // Watch for wallpaper changes
  observer = new MutationObserver(() => detectBackground())
  const wallpaper = document.querySelector('.wallpaper-bg')
  if (wallpaper) {
    observer.observe(wallpaper, { attributes: true, attributeFilter: ['style', 'class'] })
  }
})

onUnmounted(() => {
  observer?.disconnect()
  if (suggestionTimer) clearTimeout(suggestionTimer)
})

// Re-detect when theme changes
watch(() => appStore.isDark, () => {
  if (!appStore.wallpaperUrl) detectBackground()
})
</script>

<style scoped>
.search-bar {
  width: 80%;
  max-width: 880px;
  margin: 0 auto;
  position: relative;
}

.search-input-wrap {
  display: flex;
  align-items: center;
  height: 42px;
  background: var(--custom-bg, rgba(0,0,0,0.04));
  border: 1px solid rgba(0,0,0,0.06);
  border-radius: 12px;
  padding: 0 4px 0 14px;
  transition: all 0.2s ease;
}

/* Dark background detected */
.search-dark-bg .search-input-wrap {
  background: var(--custom-bg, rgba(255,255,255,0.12));
  border-color: var(--custom-bg, rgba(255,255,255,0.15));
}

.search-dark-bg .search-input-wrap:focus-within {
  background: var(--custom-bg, rgba(255,255,255,0.18));
  border-color: var(--custom-bg, rgba(255,255,255,0.3));
  box-shadow: 0 0 0 3px rgba(255,255,255,0.08);
}

.search-dark-bg .search-input,
.search-dark-bg .search-input::placeholder {
  color: rgba(255,255,255,0.9);
}

.search-dark-bg .search-icon {
  color: rgba(255,255,255,0.6);
}

.search-dark-bg .search-kbd {
  background: rgba(255,255,255,0.1);
  color: rgba(255,255,255,0.5);
}

.search-dark-bg .search-clear {
  background: rgba(255,255,255,0.1);
  color: rgba(255,255,255,0.7);
}

.search-dark-bg .search-clear:hover {
  background: rgba(255,255,255,0.18);
  color: white;
}

.search-dark-bg .search-divider {
  background: rgba(255,255,255,0.12);
}

.search-dark-bg .engine-select {
  color: rgba(255,255,255,0.7);
}

.search-dark-bg .engine-select:hover {
  background: rgba(255,255,255,0.08);
  color: white;
}

:root[data-theme="dark"] .search-input-wrap {
  background: rgba(255,255,255,0.06);
  border-color: rgba(255,255,255,0.08);
}

.search-input-wrap:focus-within {
  background: var(--custom-bg, white);
  border-color: var(--sd-primary);
  box-shadow: 0 0 0 3px rgba(0,122,255,0.12);
}

:root[data-theme="dark"] .search-input-wrap:focus-within {
  background: var(--custom-bg, rgba(28,28,30,0.9));
}

.search-icon {
  flex-shrink: 0;
  color: var(--sd-text-tertiary);
  margin-right: 8px;
}

.search-input {
  flex: 1;
  height: 100%;
  border: none;
  outline: none;
  background: transparent;
  font-size: 14px;
  color: var(--sd-text-primary);
  padding: 0 8px;
}

.search-input::placeholder {
  color: var(--sd-text-tertiary);
}

.search-kbd {
  display: inline-block;
  width: 32px;
  height: 20px;
  background: rgba(0,0,0,0.05);
  border-radius: 3px;
  font-size: 10px;
  line-height: 20px;
  text-align: center;
  margin-left: 8px;
  font-family: monospace;
}

.search-clear {
  width: 24px;
  height: 24px;
  border: none;
  background: transparent;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s ease;
  padding: 0;
}

.search-clear:hover {
  background: rgba(0,0,0,0.1);
}

.search-divider {
  width: 1px;
  height: 20px;
  background: var(--sd-border);
  margin: 0 8px;
}

.engine-select {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: var(--sd-text-tertiary);
  padding: 4px 8px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.engine-select:hover {
  background: rgba(0,0,0,0.04);
}

.search-results {
  margin-top: 8px;
  max-height: 200px;
  overflow-y: auto;
  border-radius: 8px;
  border: 1px solid var(--sd-border);
  background: var(--sd-bg-elevated);
}

.search-result-item {
  display: flex;
  align-items: center;
  padding: 10px 12px;
  cursor: pointer;
  transition: background 0.2s ease;
}

.search-result-item:hover {
  background: var(--sd-bg-base);
}

.result-icon {
  flex-shrink: 0;
  margin-right: 10px;
  color: var(--sd-text-tertiary);
}

.result-info {
  flex: 1;
  min-width: 0;
}

.result-title {
  font-size: 13px;
  color: var(--sd-text-primary);
  line-height: 1.4;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.result-url {
  font-size: 11px;
  color: var(--sd-text-tertiary);
  line-height: 1.4;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  margin-top: 2px;
}

/* Scrollbar styling */
.search-results::-webkit-scrollbar {
  width: 6px;
}

.search-results::-webkit-scrollbar-track {
  background: var(--sd-bg-base);
}

.search-results::-webkit-scrollbar-thumb {
  background: var(--sd-border);
  border-radius: 3px;
}

.search-results::-webkit-scrollbar-thumb:hover {
  background: var(--sd-text-tertiary);
}
</style>