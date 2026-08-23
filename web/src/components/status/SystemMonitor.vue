<template>
  <div v-if="appStore.showSystemMonitor" class="system-monitor" :class="{ 'expanded': isExpanded }">
    <!-- Compact bar -->
    <div class="monitor-bar" @click="isExpanded = !isExpanded">
      <div class="monitor-item cpu" :title="'CPU: ' + stats?.cpu.usage_percent.toFixed(1) + '%'">
        <Icon icon="mdi:cpu-64-bit" :size="13" />
        <span>{{ stats?.cpu.usage_percent.toFixed(0) }}%</span>
      </div>
      <div class="monitor-item memory" :title="'内存: ' + stats?.memory.used_percent.toFixed(1) + '%'">
        <Icon icon="mdi:memory" :size="13" />
        <span>{{ stats?.memory.used_percent.toFixed(0) }}%</span>
      </div>
      <div class="monitor-item disk" :title="'磁盘: ' + stats?.disk.used_percent.toFixed(1) + '%'">
        <Icon icon="mdi:harddisk" :size="13" />
        <span>{{ stats?.disk.used_percent.toFixed(0) }}%</span>
      </div>
      <div class="monitor-item net" :title="'网络 ↑' + formatBytes(stats?.network.bytes_sent || 0) + ' ↓' + formatBytes(stats?.network.bytes_recv || 0)">
        <Icon icon="mdi:swap-vertical" :size="13" />
        <span class="net-speed">
          <Icon icon="mdi:arrow-up" :size="10" class="net-up" />
          {{ formatSpeed(stats?.network.bytes_sent || 0) }}
        </span>
      </div>
      <div class="monitor-expand">
        <Icon :icon="isExpanded ? 'mdi:chevron-up' : 'mdi:chevron-down'" :size="14" />
      </div>
    </div>

    <!-- Expanded detail panel -->
    <Transition name="slide">
      <div v-if="isExpanded" class="monitor-detail">
        <div class="detail-grid">
          <!-- CPU Detail -->
          <div class="detail-card">
            <div class="detail-header">
              <Icon icon="mdi:cpu-64-bit" :size="16" />
              <span>CPU</span>
            </div>
            <div class="detail-value">{{ stats?.cpu.usage_percent.toFixed(1) }}%</div>
            <div class="detail-bar">
              <div class="detail-bar-fill" :style="{ width: Math.min(stats?.cpu.usage_percent || 0, 100) + '%' }" :class="getUsageClass(stats?.cpu.usage_percent || 0)"></div>
            </div>
            <div class="detail-meta">{{ stats?.cpu.coreCount }} 核 · {{ stats?.cpu.model_name }}</div>
          </div>

          <!-- Memory Detail -->
          <div class="detail-card">
            <div class="detail-header">
              <Icon icon="mdi:memory" :size="16" />
              <span>内存</span>
            </div>
            <div class="detail-value">{{ stats?.memory.used_percent.toFixed(1) }}%</div>
            <div class="detail-bar">
              <div class="detail-bar-fill" :style="{ width: Math.min(stats?.memory.used_percent || 0, 100) + '%' }" :class="getUsageClass(stats?.memory.used_percent || 0)"></div>
            </div>
            <div class="detail-meta">{{ formatBytes(stats?.memory.used || 0) }} / {{ formatBytes(stats?.memory.total || 0) }}</div>
          </div>

          <!-- Disk Detail -->
          <div class="detail-card">
            <div class="detail-header">
              <Icon icon="mdi:harddisk" :size="16" />
              <span>磁盘</span>
            </div>
            <div class="detail-value">{{ stats?.disk.used_percent.toFixed(1) }}%</div>
            <div class="detail-bar">
              <div class="detail-bar-fill" :style="{ width: Math.min(stats?.disk.used_percent || 0, 100) + '%' }" :class="getUsageClass(stats?.disk.used_percent || 0)"></div>
            </div>
            <div class="detail-meta">{{ formatBytes(stats?.disk.used || 0) }} / {{ formatBytes(stats?.disk.total || 0) }}</div>
          </div>

          <!-- Uptime -->
          <div class="detail-card">
            <div class="detail-header">
              <Icon icon="mdi:clock-outline" :size="16" />
              <span>运行时间</span>
            </div>
            <div class="detail-value uptime">{{ formatUptime(stats?.host.uptime_seconds || 0) }}</div>
            <div class="detail-meta">{{ stats?.host.hostname }} · {{ stats?.host.os }}</div>
          </div>
        </div>

        <!-- Partitions -->
        <div v-if="stats?.disk.partitions?.length" class="partitions">
          <div class="partition-label">分区详情</div>
          <div class="partition-list">
            <div v-for="p in stats.disk.partitions" :key="p.mount_point" class="partition-item">
              <span class="partition-mount">{{ p.mount_point }}</span>
              <div class="partition-bar">
                <div class="partition-bar-fill" :style="{ width: Math.min(p.used_percent, 100) + '%' }" :class="getUsageClass(p.used_percent)"></div>
              </div>
              <span class="partition-percent">{{ p.used_percent.toFixed(0) }}%</span>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useAppStore } from '../../stores/app'
import type { SystemStats } from '../../types/system'

const appStore = useAppStore()
const stats = ref<SystemStats | null>(null)
const isExpanded = ref(false)
let timer: ReturnType<typeof setInterval> | null = null

async function fetchStats() {
  try {
    const token = localStorage.getItem('sundash-token')
    const res = await fetch('/api/system/stats', {
      headers: token ? { Authorization: `Bearer ${token}` } : {}
    })
    if (res.ok) {
      stats.value = await res.json()
    }
  } catch {
    // Silently ignore fetch errors
  }
}

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i]
}

function formatSpeed(bytes: number): string {
  return formatBytes(bytes)
}

function formatUptime(seconds: number): string {
  const days = Math.floor(seconds / 86400)
  const hours = Math.floor((seconds % 86400) / 3600)
  const mins = Math.floor((seconds % 3600) / 60)
  if (days > 0) return `${days}天${hours}小时`
  if (hours > 0) return `${hours}小时${mins}分`
  return `${mins}分钟`
}

function getUsageClass(percent: number): string {
  if (percent >= 90) return 'critical'
  if (percent >= 70) return 'warning'
  return 'normal'
}

onMounted(() => {
  fetchStats()
  timer = setInterval(fetchStats, 5000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>

<style scoped>
.system-monitor {
  background: var(--sd-bg-elevated);
  border: 1px solid var(--sd-border);
  border-radius: 12px;
  overflow: hidden;
  margin: 0 16px 12px;
  backdrop-filter: blur(10px);
}

.monitor-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 8px 14px;
  cursor: pointer;
  user-select: none;
}

.monitor-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--sd-text-secondary);
}

.monitor-item span {
  font-variant-numeric: tabular-nums;
}

.net-speed {
  display: flex;
  align-items: center;
  gap: 1px;
}

.net-up {
  color: var(--sd-color-success);
}

.monitor-expand {
  margin-left: auto;
  color: var(--sd-text-tertiary);
}

/* Expanded detail */
.monitor-detail {
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

.detail-value.uptime {
  font-size: 16px;
}

.detail-bar {
  height: 4px;
  background: var(--sd-bg-base);
  border-radius: 2px;
  margin: 6px 0 4px;
  overflow: hidden;
}

.detail-bar-fill {
  height: 100%;
  border-radius: 2px;
  transition: width 0.5s ease;
}

.detail-bar-fill.normal {
  background: var(--sd-color-success);
}

.detail-bar-fill.warning {
  background: var(--sd-color-warning);
}

.detail-bar-fill.critical {
  background: var(--sd-color-danger);
}

.detail-meta {
  font-size: 10px;
  color: var(--sd-text-tertiary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* Partitions */
.partitions {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px solid var(--sd-border);
}

.partition-label {
  font-size: 11px;
  color: var(--sd-text-tertiary);
  margin-bottom: 8px;
}

.partition-list {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.partition-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 11px;
}

.partition-mount {
  width: 60px;
  flex-shrink: 0;
  color: var(--sd-text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.partition-bar {
  flex: 1;
  height: 4px;
  background: var(--sd-bg-base);
  border-radius: 2px;
  overflow: hidden;
}

.partition-bar-fill {
  height: 100%;
  border-radius: 2px;
  transition: width 0.5s ease;
}

.partition-bar-fill.normal { background: var(--sd-color-success); }
.partition-bar-fill.warning { background: var(--sd-color-warning); }
.partition-bar-fill.critical { background: var(--sd-color-danger); }

.partition-percent {
  width: 32px;
  text-align: right;
  color: var(--sd-text-tertiary);
  font-variant-numeric: tabular-nums;
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
