<template>
  <div v-if="appStore.clockShow" class="clock-display" :class="{ 'clock-seconds': appStore.clockShowSeconds, 'clock-compact': compact }">
    <div class="clock-time">
      <span class="clock-hours">{{ hours }}</span>
      <span class="clock-separator">:</span>
      <span class="clock-minutes">{{ minutes }}</span>
      <template v-if="appStore.clockShowSeconds">
        <span class="clock-separator">:</span>
        <span class="clock-seconds-value">{{ seconds }}</span>
      </template>
      <span v-if="appStore.clockFormat === '12'" class="clock-period">{{ period }}</span>
    </div>
    <div class="clock-date">{{ dateStr }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '../../stores/app'

const { t } = useI18n()

const props = defineProps<{ compact?: boolean }>()
const appStore = useAppStore()

const now = ref(new Date())
let timer: ReturnType<typeof setInterval> | null = null

onMounted(() => {
  timer = setInterval(() => {
    now.value = new Date()
  }, 1000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})

const hours = computed(() => {
  const h = now.value.getHours()
  if (appStore.clockFormat === '12') {
    return h === 0 ? '12' : h > 12 ? String(h - 12).padStart(2, '0') : String(h).padStart(2, '0')
  }
  return String(h).padStart(2, '0')
})

const minutes = computed(() => String(now.value.getMinutes()).padStart(2, '0'))
const seconds = computed(() => String(now.value.getSeconds()).padStart(2, '0'))

const period = computed(() => now.value.getHours() >= 12 ? 'PM' : 'AM')

const dateStr = computed(() => {
  const days = [t('clock.sunday'), t('clock.monday'), t('clock.tuesday'), t('clock.wednesday'), t('clock.thursday'), t('clock.friday'), t('clock.saturday')]
  const d = now.value
  return `${d.getMonth() + 1}/${d.getDate()} ${days[d.getDay()]}`
})
</script>

<style scoped>
.clock-display {
  text-align: center;
  padding: 20px 0;
}

.clock-compact {
  padding: 0;
}

.clock-compact .clock-time {
  font-size: 28px;
  font-weight: 600;
}

.clock-compact .clock-separator {
  margin: 0 6px;
}

.clock-time {
  font-size: 64px;
  font-weight: 200;
  letter-spacing: -2px;
  color: var(--sd-text-primary);
  line-height: 1;
  font-variant-numeric: tabular-nums;
}

.clock-separator {
  opacity: 0.5;
  margin: 0 2px;
}

.clock-seconds-value {
  font-size: 48px;
  opacity: 0.7;
}

.clock-period {
  font-size: 24px;
  font-weight: 400;
  opacity: 0.6;
  margin-left: 8px;
}

.clock-date {
  font-size: 16px;
  color: var(--sd-text-secondary);
  margin-top: 8px;
  font-weight: 400;
}

/* Responsive */
@media (max-width: 640px) {
  .clock-time {
    font-size: 48px;
  }
  .clock-seconds-value {
    font-size: 36px;
  }
}

@media (max-width: 480px) {
  .clock-display {
    padding: 10px 0;
  }
  .clock-time {
    font-size: 36px;
    letter-spacing: -1px;
  }
  .clock-compact .clock-time {
    font-size: 22px;
  }
  .clock-seconds-value {
    font-size: 28px;
  }
  .clock-period {
    font-size: 18px;
    margin-left: 4px;
  }
  .clock-date {
    font-size: 14px;
    margin-top: 4px;
  }
}
</style>
