<template>
  <div v-if="appStore.showWeatherWidget" class="weather-widget" :class="{ 'expanded': isExpanded }">
    <!-- Compact bar -->
    <div class="weather-bar" @click="isExpanded = !isExpanded">
      <div class="weather-item temp" :title="'温度: ' + weather?.temperature.toFixed(1) + '°C'">
        <Icon icon="mdi:thermometer" :size="13" />
        <span>{{ weather?.temperature.toFixed(0) }}°C</span>
      </div>
      <div class="weather-item wind" :title="'风速: ' + (weather?.windspeed || 0).toFixed(1) + ' km/h'">
        <Icon icon="mdi:weather-windy" :size="13" />
        <span>{{ (weather?.windspeed || 0).toFixed(0) }} km/h</span>
      </div>
      <div class="weather-item condition">
        <Icon :icon="weatherIcon" :size="13" />
        <span>{{ weatherCondition }}</span>
      </div>
      <div class="weather-expand">
        <Icon :icon="isExpanded ? 'mdi:chevron-up' : 'mdi:chevron-down'" :size="14" />
      </div>
    </div>

    <!-- Expanded detail panel -->
    <Transition name="slide">
      <div v-if="isExpanded" class="weather-detail">
        <div class="detail-grid">
          <!-- Temperature Detail -->
          <div class="detail-card">
            <div class="detail-header">
              <Icon icon="mdi:thermometer" :size="16" />
              <span>温度</span>
            </div>
            <div class="detail-value">{{ weather?.temperature.toFixed(1) }}°C</div>
            <div class="detail-meta">{{ weatherCondition }}</div>
          </div>

          <!-- Wind Detail -->
          <div class="detail-card">
            <div class="detail-header">
              <Icon icon="mdi:weather-windy" :size="16" />
              <span>风速</span>
            </div>
            <div class="detail-value">{{ (weather?.windspeed || 0).toFixed(1) }} km/h</div>
            <div class="detail-meta">风向: {{ weather?.winddirection?.toFixed(0) }}°</div>
          </div>

          <!-- Location -->
          <div class="detail-card">
            <div class="detail-header">
              <Icon icon="mdi:map-marker" :size="16" />
              <span>位置</span>
            </div>
            <div class="detail-value">{{ weather?.latitude.toFixed(2) }}, {{ weather?.longitude.toFixed(2) }}</div>
            <div class="detail-meta">{{ weather?.timezone }}</div>
          </div>

          <!-- Update Time -->
          <div class="detail-card">
            <div class="detail-header">
              <Icon icon="mdi:clock-outline" :size="16" />
              <span>更新时间</span>
            </div>
            <div class="detail-value update-time">{{ formatTime(weather?.time) }}</div>
            <div class="detail-meta">数据来源: Open-Meteo</div>
          </div>
        </div>

        <!-- Weather description -->
        <div v-if="weather" class="weather-description">
          <div class="description-title">天气状况</div>
          <div class="description-text">{{ getWeatherDescription(weather?.weathercode) }}</div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useAppStore } from '../../stores/app'
import axios from '@/api'

const appStore = useAppStore()
const weather = ref<any>(null)
const isExpanded = ref(false)
let timer: ReturnType<typeof setInterval> | null = null

async function fetchWeather() {
  try {
    const token = localStorage.getItem('sundash-token')
    const res = await axios.get('/api/weather', {
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    if (res.status === 200) {
      weather.value = res.data
    }
  } catch (err) {
    console.error('Failed to fetch weather:', err)
    // Keep old data or show error? For now, just ignore.
  }
}

function formatTime(timeString: string | undefined): string {
  if (!timeString) return '--:--'
  try {
    const date = new Date(timeString)
    return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  } catch {
    return timeString
  }
}

function getWeatherDescription(code: number | undefined): string {
  if (code === undefined) return '未知'
  // WMO Weather interpretation codes (https://open-meteo.com/en/docs)
  const weatherCodes: Record<number, string> = {
    0: '晴朗',
    1: '主要晴朗',
    2: '部分多云',
    3: '多云',
    45: '雾',
    48: '沉雾',
    51: '轻度毛毛雨',
    53: '中度毛毛雨',
    55: '密集毛毛雨',
    56: '轻度冻毛毛雨',
    57: '密集冻毛毛雨',
    61: '轻度雨',
    63: '中度雨',
    65: '大雨',
    66: '轻度冻雨',
    67: '大冻雨',
    71: '轻度雪',
    73: '中度雪',
    75: '大雪',
    77: '雪粒',
    80: '轻度阵雨',
    81: '中度阵雨',
    82: '猛烈阵雨',
    85: '轻度雪阵雨',
    86: '重度雪阵雨',
    95: '雷暴',
    96: '雷暴伴有轻度冰雹',
    99: '雷暴伴有重度冰雹'
  }
  return weatherCodes[code] || '未知'
}

function getWeatherIcon(code: number | undefined): string {
  if (code === undefined) return 'mdi:weather-cloudy'
  // Map weather codes to icons (using Material Design Icons)
  const iconMap: Record<number, string> = {
    0: 'mdi:weather-sunny',
    1: 'mdi:weather-partlycloudy',
    2: 'mdi:weather-partlycloudy',
    3: 'mdi:weather-cloudy',
    45: 'mdi:weather-fog',
    48: 'mdi:weather-fog',
    51: 'mdi:weather-rainy',
    53: 'mdi:weather-rainy',
    55: 'mdi:weather-rainy',
    56: 'mdi:weather-snowy-rainy',
    57: 'mdi:weather-snowy-rainy',
    61: 'mdi:weather-rainy',
    63: 'mdi:weather-rainy',
    65: 'mdi:weather-pouring',
    66: 'mdi:weather-snowy-rainy',
    67: 'mdi:weather-snowy-rainy',
    71: 'mdi:weather-snowy',
    73: 'mdi:weather-snowy',
    75: 'mdi:weather-snowy',
    77: 'mdi:weather-snowy',
    80: 'mdi:weather-lightning-rainy',
    81: 'mdi:weather-lightning-rainy',
    82: 'mdi:weather-lightning-rainy',
    85: 'mdi:weather-snowy',
    86: 'mdi:weather-snowy',
    95: 'mdi:weather-lightning',
    96: 'mdi:weather-lightning',
    99: 'mdi:weather-lightning'
  }
  return iconMap[code] || 'mdi:weather-cloudy'
}

const weatherIcon = computed(() => getWeatherIcon(weather?.weathercode))
const weatherCondition = computed(() => getWeatherDescription(weather?.weathercode))

onMounted(() => {
  fetchWeather()
  timer = setInterval(fetchWeather, 10 * 60 * 1000) // 10 minutes
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.weather-widget {
  background: var(--sd-bg-elevated);
  border: 1px solid var(--sd-border);
  border-radius: 12px;
  overflow: hidden;
  margin: 0 16px 12px;
  backdrop-filter: blur(10px);
}

.weather-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 8px 14px;
  cursor: pointer;
  user-select: none;
}

.weather-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--sd-text-secondary);
}

.weather-item span {
  font-variant-numeric: tabular-nums;
}

.weather-expand {
  margin-left: auto;
  color: var(--sd-text-tertiary);
}

/* Expanded detail */
.weather-detail {
  padding: 12px 14px 14px;
  border-top: 1px solid var(--sd-border);
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 10px;
}

.detail-card {
  background: var(--sd-bg-subtle);
  border-radius: 8px;
  padding: 10px;
}

.detail-header {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 11px;
  color: var(--sd-text-tertiary);
  margin-bottom: 6px;
}

.detail-value {
  font-size: 20px;
  font-weight: 600;
  color: var(--sd-text-primary);
  font-variant-numeric: tabular-nums;
}

.detail-meta {
  font-size: 10px;
  color: var(--sd-text-tertiary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.description-title {
  font-size: 12px;
  color: var(--sd-text-primary);
  margin-bottom: 4px;
}

.description-text {
  font-size: 13px;
  color: var(--sd-text-secondary);
  line-height: 1.4;
}

/* Transition */
.slide-enter-active,
.slide-leave-active {
  transition: all 0.25s ease;
  overflow: hidden;
}

.slide-enter-from,
.slide-leave-to {
  opacity: 0;
  max-height: 0;
}

.slide-enter-to,
.slide-leave-from {
  opacity: 1;
  max-height: 500px;
}
</style>