<template>
  <div class="card-item" :style="cardItemStyle" @click="$emit('click')" @dblclick.stop @contextmenu.prevent="onContextMenu">
    <div class="card-icon-box" :style="iconBoxStyle">
      <img v-if="isIconUrl && card.icon" :src="card.icon" class="card-icon-img" :alt="card.title" />
      <Icon v-else :icon="resolvedIconName" :width="36" :height="36" :color="card.icon_color || iconFallbackColor" />
    </div>
    <div class="card-info">
      <div class="card-title" :style="cardTitleStyle">{{ card.title }}</div>
      <div v-if="card.description && !isNarrow" class="card-desc">{{ card.description }}</div>
    </div>
    <div v-if="card.url_internal && !isNarrow" class="card-net-dot" :title="t('home.internalAccess')"></div>
    
    <!-- Context Menu -->
    <teleport to="body">
      <div v-if="showContextMenu" class="card-context-menu" :style="contextMenuStyle" @click.stop>
        <div class="ctx-menu-item" @click="handleAction('open')">
          <Icon icon="mdi:open-in-new" :width="16" :height="16" />
          <span>{{ t('home.visitLink') }}</span>
        </div>
        <div v-if="card.url_internal" class="ctx-menu-item" @click="handleAction('open-internal')">
          <Icon icon="mdi:lan" :width="16" :height="16" />
          <span>{{ t('home.visitInternal') }}</span>
        </div>
        <div class="ctx-menu-item" @click="handleAction('copy')">
          <Icon icon="mdi:content-copy" :width="16" :height="16" />
          <span>{{ t('home.copyExternalLink') }}</span>
        </div>
        <div v-if="card.url_internal" class="ctx-menu-item" @click="handleAction('copy-internal')">
          <Icon icon="mdi:content-copy" :width="16" :height="16" />
          <span>{{ t('home.copyInternalLink') }}</span>
        </div>
        <div class="ctx-menu-divider"></div>
        <div class="ctx-menu-item" @click="handleAction('hide')">
          <Icon icon="mdi:eye-off" :width="16" :height="16" />
          <span>{{ t('common.hide') }}</span>
        </div>
        <div class="ctx-menu-item" @click="handleAction('edit')">
          <Icon icon="mdi:pencil-outline" :width="16" :height="16" />
          <span>{{ t('home.editBookmark') }}</span>
        </div>
        <div class="ctx-menu-item danger" @click="handleAction('delete')">
          <Icon icon="mdi:delete-outline" :width="16" :height="16" />
          <span>{{ t('home.deleteBookmark') }}</span>
        </div>
      </div>
    </teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Icon } from '@iconify/vue'
import { useAppStore } from '../../stores/app'
import type { Card } from '../../types'

const { t } = useI18n()
const appStore = useAppStore()
const props = defineProps<{ card: Card }>()

// 窄屏检测（≤480px 时卡片竖排）
const isNarrow = ref(window.innerWidth <= 768)
window.addEventListener('resize', () => { isNarrow.value = window.innerWidth <= 768 })

const iconFallbackColor = computed(() => appStore.primaryColor || '#007AFF')

// Check if icon value is an image URL (not an Iconify icon name)
const isIconUrl = computed(() => {
  const icon = props.card.icon || ''
  return icon.startsWith('http://') || icon.startsWith('https://') || icon.startsWith('//') || icon.startsWith('/')
})

// Get the resolved icon name for Iconify
const resolvedIconName = computed(() => {
  const icon = props.card.icon || ''
  if (!icon) return defaultIcon.value
  // If it's a URL, return default icon
  if (isIconUrl.value) return defaultIcon.value
  return icon
})
const emit = defineEmits<{
  click: []
  edit: []
  delete: []
  hide: []
}>()

// Context menu state
const showContextMenu = ref(false)
const contextMenuStyle = ref<Record<string, string>>({})

function onContextMenu(e: MouseEvent) {
  e.preventDefault()
  e.stopPropagation()
  
  // Calculate position
  const x = e.clientX
  const y = e.clientY
  
  // Get viewport size
  const viewportWidth = window.innerWidth
  const viewportHeight = window.innerHeight
  
  // Menu dimensions (approximate)
  const menuWidth = 180
  const menuHeight = 260
  
  // Adjust position to keep menu in viewport
  let posX = x
  let posY = y
  
  if (x + menuWidth > viewportWidth) {
    posX = viewportWidth - menuWidth - 10
  }
  if (y + menuHeight > viewportHeight) {
    posY = viewportHeight - menuHeight - 10
  }
  
  contextMenuStyle.value = {
    left: `${posX}px`,
    top: `${posY}px`,
  }
  
  showContextMenu.value = true
}

function closeContextMenu() {
  showContextMenu.value = false
}

function handleAction(key: string) {
  closeContextMenu()
  switch (key) {
    case 'open':
      window.open(props.card.url, '_blank')
      break
    case 'open-internal':
      window.open(props.card.url_internal, '_blank')
      break
    case 'copy':
      navigator.clipboard.writeText(props.card.url)
      break
    case 'copy-internal':
      navigator.clipboard.writeText(props.card.url_internal || '')
      break
    case 'edit':
      emit('edit')
      break
    case 'hide':
      emit('hide')
      break
    case 'delete':
      emit('delete')
      break
  }
}

// Close menu on outside click
function onDocumentClick() {
  closeContextMenu()
}

// Close menu on scroll
function onDocumentScroll() {
  closeContextMenu()
}

onMounted(() => {
  document.addEventListener('click', onDocumentClick)
  document.addEventListener('scroll', onDocumentScroll, true)
})

onUnmounted(() => {
  document.removeEventListener('click', onDocumentClick)
  document.removeEventListener('scroll', onDocumentScroll, true)
})

// Better default icon based on card color
const defaultIcon = computed(() => {
  const color = props.card.icon_color || iconFallbackColor.value
  if (color.includes('FF') || color.includes('ff')) return 'mdi:star-four-points'
  if (color.includes('34C759') || color.includes('34c759')) return 'mdi:leaf'
  if (color.includes('5856D6') || color.includes('5856d6')) return 'mdi:shape'
  return 'mdi:bookmark'
})

// Icon box background style
const iconBoxStyle = computed(() => {
  return {}
})

// Card item background + size style
const cardItemStyle = computed(() => {
  const style: Record<string, string> = {}
  if (props.card.bg_color) {
    style.background = props.card.bg_color
  }
  const sizeOffset = parseInt(appStore.cardItemSize) || 0
  if (sizeOffset !== 0) {
    style.padding = `${10 + sizeOffset}px ${14 + sizeOffset}px`
  }
  return style
})

// Card title font size style
const cardTitleStyle = computed(() => {
  return { fontSize: `${appStore.cardLabelSize}px` }
})
</script>

<style scoped>
.card-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  border-radius: 14px;
  background: rgba(255,255,255,0.55);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid rgba(0,0,0,0.04);
  cursor: pointer;
  text-decoration: none;
  transition: all 0.2s ease;
  position: relative;
  user-select: none;
  -webkit-user-select: none;
}

/* 窄屏竖排模式 */
@media (max-width: 480px) {
  .card-item {
    flex-direction: column;
    align-items: center;
    text-align: center;
    padding: 10px 4px;
    gap: 6px;
  }

  .card-icon-box {
    width: 44px;
    height: 44px;
  }

  .card-icon-img {
    width: 28px;
    height: 28px;
  }

  .card-info {
    width: 100%;
  }

  .card-title {
    font-size: 12px;
    white-space: normal;
    overflow: visible;
    text-overflow: clip;
    word-break: break-word;
    line-height: 1.4;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
}

:root[data-theme="dark"] .card-item {
  background: rgba(44,44,46,0.5);
  border-color: rgba(255,255,255,0.06);
}

.card-item:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0,0,0,0.06);
  border-color: rgba(0,122,255,0.12);
  background: rgba(255,255,255,0.75);
}

:root[data-theme="dark"] .card-item:hover {
  box-shadow: 0 6px 20px rgba(0,0,0,0.2);
  background: rgba(55,55,58,0.65);
}

.card-item:active {
  transform: translateY(0) scale(0.98);
  transition-duration: 0.1s;
}

.card-icon-box {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition: transform 0.2s ease;
  overflow: hidden;
}

.card-item:hover .card-icon-box {
  transform: scale(1.05);
}

.card-icon-img {
  width: 36px;
  height: 36px;
  object-fit: contain;
  border-radius: 6px;
  flex-shrink: 0;
}

.card-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.card-title {
  font-size: 12px;
  font-weight: 500;
  color: var(--sd-text-primary);
  line-height: 1.3;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-desc {
  font-size: 11px;
  color: var(--sd-text-tertiary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-net-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #34C759;
  flex-shrink: 0;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.card-item:hover .card-net-dot {
  opacity: 1;
}

/* Tablet横排: title 换行，图标略缩 */
@media (max-width: 768px) {
  .card-item {
    gap: 8px;
    padding: 8px 10px;
  }

  .card-icon-box {
    width: 40px;
    height: 40px;
  }

  .card-icon-img {
    width: 36px;
    height: 36px;
  }

  .card-title {
    white-space: normal;
    overflow: hidden;
    text-overflow: clip;
    word-break: break-word;
    line-height: 1.4;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    font-size: 12px;
  }
}

/* 窄屏竖排模式 */
@media (max-width: 480px) {
  .card-item {
    padding: 10px 6px;
    gap: 6px;
    border-radius: 12px;
  }

  .card-icon-box {
    width: 40px;
    height: 40px;
    border-radius: 10px;
  }

  .card-icon-img {
    width: 36px;
    height: 36px;
  }

  .card-title {
    font-size: 12px;
    white-space: normal;
    overflow: hidden;
    text-overflow: clip;
    word-break: break-word;
    line-height: 1.4;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
  }

  .card-desc {
    font-size: 11px;
  }

  .card-style-round .card-item {
    padding: 10px 6px;
    min-width: 0;
    max-width: none;
    width: 100%;
  }
}

/* End Scoped Styles */
</style>

<!--
  Context Menu Styles (unscoped)
  These styles target teleported elements rendered outside this component's scope.
-->
<style>
/* Context Menu */
.card-context-menu {
  position: fixed;
  z-index: 99999;
  min-width: 170px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(0, 0, 0, 0.08);
  border-radius: 12px;
  padding: 4px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.12), 0 2px 8px rgba(0, 0, 0, 0.06);
  animation: ctxMenuIn 0.15s ease-out;
}

:root[data-theme="dark"] .card-context-menu {
  background: rgba(44, 44, 46, 0.95);
  border-color: rgba(255, 255, 255, 0.1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3), 0 2px 8px rgba(0, 0, 0, 0.2);
}

@keyframes ctxMenuIn {
  from {
    opacity: 0;
    transform: scale(0.95);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.ctx-menu-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  border-radius: 8px;
  font-size: 13px;
  font-family: var(--sd-font);
  color: var(--sd-text-primary);
  cursor: pointer;
  transition: all 0.15s ease;
  user-select: none;
}

.ctx-menu-item:hover {
  background: var(--sd-primary-light, rgba(0, 122, 255, 0.08));
  color: var(--sd-primary, #007AFF);
}

:root[data-theme="dark"] .ctx-menu-item:hover {
  background: var(--sd-primary-medium, rgba(0, 122, 255, 0.15));
  color: #4DA6FF;
}

.ctx-menu-item.danger:hover {
  background: rgba(255, 59, 48, 0.08);
  color: #FF3B30;
}

:root[data-theme="dark"] .ctx-menu-item.danger:hover {
  background: rgba(255, 59, 48, 0.15);
  color: #FF6B6B;
}

.ctx-menu-divider {
  height: 1px;
  background: rgba(0, 0, 0, 0.06);
  margin: 4px 8px;
}

:root[data-theme="dark"] .ctx-menu-divider {
  background: rgba(255, 255, 255, 0.08);
}
</style>
