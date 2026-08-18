<template>
  <div v-if="appStore.showSystemStatus" class="system-status">
    <div class="status-item">
      <Icon icon="mdi:folder-outline" :size="14" />
      <span>{{ groupCount }} 分组</span>
    </div>
    <div class="status-divider"></div>
    <div class="status-item">
      <Icon icon="mdi:bookmark-outline" :size="14" />
      <span>{{ cardCount }} 卡片</span>
    </div>
    <div v-if="hiddenCount > 0" class="status-divider"></div>
    <div v-if="hiddenCount > 0" class="status-item hidden">
      <Icon icon="mdi:eye-off-outline" :size="14" />
      <span>{{ hiddenCount }} 已隐藏</span>
    </div>
    <div class="status-divider"></div>
    <div class="status-item">
      <Icon icon="mdi:clock-outline" :size="14" />
      <span>{{ currentTime }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useAppStore } from '../../stores/app'
import { usePanelStore } from '../../stores/panel'

const appStore = useAppStore()
const panelStore = usePanelStore()

const groupCount = computed(() => panelStore.groups.length)
const cardCount = computed(() => {
  return panelStore.groups.reduce((sum, g) => sum + (g.cards?.length || 0), 0)
})

// Hidden cards count
const hiddenCount = computed(() => {
  try {
    const hidden = new Set(JSON.parse(localStorage.getItem('sundash-hidden-cards') || '[]'))
    return hidden.size
  } catch {
    return 0
  }
})

// Live clock
const currentTime = ref('')
let timer: ReturnType<typeof setInterval> | null = null

function updateTime() {
  const now = new Date()
  currentTime.value = now.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

onMounted(() => {
  updateTime()
  timer = setInterval(updateTime, 1000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.system-status {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  padding: 8px 16px;
  font-size: 12px;
  color: var(--sd-text-tertiary);
}

.status-item {
  display: flex;
  align-items: center;
  gap: 4px;
}

.status-item.hidden {
  color: var(--sd-text-secondary);
}

.status-divider {
  width: 1px;
  height: 12px;
  background: var(--sd-border);
}
</style>
