<template>
  <div v-if="appStore.showWeatherWidget" class="weather-widget" :class="{ 'expanded': isExpanded }">
    <!-- Compact bar -->
    <div class="weather-bar" @click="isExpanded = !isExpanded">
      <div class="weather-item temp" :title="t('weather.temperature') + ': ' + weather?.temperature.toFixed(1) + '°C'">
        <Icon icon="mdi:thermometer" :size="13" />
        <span>{{ weather?.temperature.toFixed(0) }}°C</span>
      </div>
      <div class="weather-item wind" :title="t('weather.wind') + ': ' + (weather?.windspeed || 0).toFixed(1) + ' km/h'">
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
              <span>{{ t('weather.temperature') }}</span>
            </div>
            <div class="detail-value">{{ weather?.temperature.toFixed(1) }}°C</div>
            <div class="detail-meta">{{ weatherCondition }}</div>
          </div>

          <!-- Wind Detail -->
          <div class="detail-card">
            <div class="detail-header">
              <Icon icon="mdi:weather-windy" :size="16" />
              <span>{{ t('weather.wind') }}</span>
            </div>
            <div class="detail-value">{{ (weather?.windspeed || 0).toFixed(1) }} km/h</div>
            <div class="detail-meta">{{ t('weather.windDirection') }}: {{ weather?.winddirection?.toFixed(0) }}°</div>
          </div>

          <!-- Location -->
          <div class="detail-card">
            <div class="detail-header">
              <Icon icon="mdi:map-marker" :size="16" />
              <span>{{ t('weather.location') }}</span>
            </div>
            <div class="detail-value">{{ weather?.latitude.toFixed(2) }}, {{ weather?.longitude.toFixed(2) }}</div>
            <div class="detail-meta">{{ weather?.timezone }}</div>
          </div>

          <!-- Update Time -->
          <div class="detail-card">
            <div class="detail-header">
              <Icon icon="mdi:clock-outline" :size="16" />
              <span>{{ t('weather.updated') }}</span>
            </div>
            <div class="detail-value update-time">{{ formatTime(weather?.time) }}</div>
            <div class="detail-meta">{{ t('weather.dataSource') }}: Open-Meteo</div>
          </div>
        </div>

        <!-- Weather description -->
        <div v-if="weather" class="weather-description">
          <div class="description-title">{{ t('weather.condition') }}</div>
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
import { useI18n } from 'vue-i18n'
import { api as axios } from '@/api'

const appStore = useAppStore()
const { t } = useI18n()
const weather = ref<any>(null)
const isExpanded = ref(false)
let timer: ReturnType<typeof setInterval> | null = null

async function fetchWeather() {
  try {
    const token = localStorage.getItem('sundash-token')
    const res = await axios.get('weather', {
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
  if (code === undefined) return t('weather.unknown')
  // WMO Weather interpretation codes (https://open-meteo.com/en/docs)
  const weatherCodes: Record<number, string> = {
    0: t('weather.sunny'),
    1: t('weather.mainlyClear'),
    2: t('weather.partlyCloudy'),
    3: t('weather.cloudy'),
    45: t('weather.fog'),
    48: t('weather.freezingFog'),
    51: t('weather.drizzleLight'),
    53: t('weather.drizzleModerate'),
    55: t('weather.drizzleDense'),
    56: t('weather.freezingDrizzleLight'),
    57: t('weather.freezingDrizzleDense'),
    61: t('weather.rainLight'),
    63: t('weather.rainModerate'),
    65: t('weather.rainHeavy'),
    66: t('weather.freezingRainLight'),
    67: t('weather.freezingRainHeavy'),
    71: t('weather.snowLight'),
    73: t('weather.snowModerate'),
    75: t('weather.snowHeavy'),
    77: t('weather.snowGrains'),
    80: t('weather.showersLight'),
    81: t('weather.showersModerate'),
    82: t('weather.showersViolent'),
    85: t('weather.snowShowersLight'),
    86: t('weather.snowShowersHeavy'),
    95: t('weather.thunderstorm'),
    96: t('weather.thunderstormHailLight'),
    99: t('weather.thunderstormHailHeavy')
  }
  return weatherCodes[code] || t('weather.unknown')
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
