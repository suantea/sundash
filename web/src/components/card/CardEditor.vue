<template>
  <n-modal :show="show" @update:show="$emit('update:show', $event)" preset="dialog"
    :title="card ? '编辑卡片' : '添加卡片'" positive-text="保存" negative-text="取消"
    :loading="saving" class="card-editor-modal" @positive-click="handleSave">

    <div class="editor-body">
      <!-- ====== Left Column: Form + Preview ====== -->
      <div class="editor-left">
        <!-- Preview -->
        <div class="preview-section">
          <div class="preview-header">
            <n-checkbox v-model:checked="showPreview">效果预览</n-checkbox>
            <n-checkbox v-model:checked="canvasTransparent">透明画布</n-checkbox>
          </div>
          <div v-if="showPreview" class="preview-box" :class="{ 'canvas-transparent': canvasTransparent }">
            <div class="preview-card preview-card-h" :style="{ background: form.bg_color || 'rgba(42,42,42,0.42)' }">
              <div class="preview-icon">
                <Icon v-if="form.icon" :icon="parsedIconName" :size="38" :color="form.icon_color || '#2080f0'" />
                <Icon v-else icon="mdi:compass" :size="38" color="#999" />
              </div>
              <div class="preview-info">
                <div class="preview-title">{{ form.title || '卡片标题' }}</div>
                <div v-if="form.description" class="preview-desc">{{ form.description }}</div>
              </div>
            </div>
          </div>
        </div>

        <!-- Form Fields -->
        <n-form label-placement="left" :show-feedback="false" class="editor-form">
          <div class="form-grid">
            <n-form-item label="分组">
              <n-select v-model:value="form.group_id" :options="groupOptions" placeholder="选择分组" />
            </n-form-item>
            <n-form-item label="标题">
              <n-input v-model:value="form.title" placeholder="卡片标题" />
            </n-form-item>
          </div>
          <div class="form-grid">
            <n-form-item label="链接">
              <n-input v-model:value="form.url" placeholder="https://example.com" />
            </n-form-item>
            <n-form-item label="内网链接">
              <n-input v-model:value="form.url_internal" placeholder="http://192.168.1.x (可选)" />
            </n-form-item>
          </div>
          <div class="form-grid">
            <n-form-item label="描述">
              <n-input v-model:value="form.description" placeholder="简短描述 (可选)" />
            </n-form-item>
            <n-form-item label="打开方式">
              <div class="open-type-row">
                <n-radio-group v-model:value="form.open_type">
                  <n-radio value="new_tab">新标签页</n-radio>
                  <n-radio value="iframe">弹窗</n-radio>
                </n-radio-group>
              </div>
            </n-form-item>
          </div>
        </n-form>
      </div>

      <!-- ====== Right Column: Icon + Colors ====== -->
      <div class="editor-right">
        <div class="right-header">
          <span class="right-label">图标设置</span>
          <a class="iconify-link" href="https://icon-sets.iconify.design/" target="blank" rel="noopener">
            <Icon icon="mdi:open-in-new" :size="11" />
            <span>Iconify 图库</span>
          </a>
        </div>

        <!-- Icon Input + Picker Trigger -->
        <div class="icon-input-row">
          <n-input v-model:value="form.icon" placeholder="mdi:web 或粘贴 URL" size="small" clearable
            @update:value="onIconInput" style="flex:1;" />
          <button type="button" class="fetch-btn" @click="fetchFavicon" :disabled="!form.url || fetchingIcon">
            <Icon :icon="fetchingIcon ? 'mdi:loading' : 'mdi:image-sync-outline'" :size="14" :class="{ 'spin-icon': fetchingIcon }" />
            <span>{{ fetchingIcon ? '...' : 'ICO' }}</span>
          </button>
        </div>

        <!-- Icon Picker Toggle -->
        <button type="button" class="icon-picker-trigger" @click="showIconPicker = !showIconPicker">
          <div class="trigger-icon">
            <Icon v-if="form.icon" :icon="parsedIconName" :size="20" :color="form.icon_color || '#2080f0'" />
            <Icon v-else icon="mdi:compass" :size="20" color="#999" />
          </div>
          <span class="trigger-label">选择图标</span>
          <Icon :icon="showIconPicker ? 'mdi:chevron-up' : 'mdi:chevron-down'" :size="14" class="trigger-arrow" />
        </button>
        <div v-if="showIconPicker" class="icon-popover-content">
          <div v-for="cat in iconCategories" :key="cat.name" class="icon-cat">
            <div class="icon-cat-name">{{ cat.name }}</div>
            <div class="icon-grid">
              <button v-for="icon in cat.icons" :key="icon" type="button"
                class="icon-btn" :class="{ active: form.icon === icon }"
                @click="form.icon = icon; showIconPicker = false">
                <Icon :icon="icon" :size="18" />
              </button>
            </div>
          </div>
        </div>

        <!-- Divider -->
        <div class="right-divider"></div>

        <!-- Colors -->
        <div class="color-group">
          <div class="color-label">图标颜色</div>
          <div class="color-row">
            <button v-for="c in colorPresets" :key="c" type="button"
              class="color-dot" :class="{ active: form.icon_color === c }"
              :style="{ background: c }" @click="form.icon_color = c" />
            <n-color-picker v-model:value="form.icon_color" :show-alpha="true" format="hex" size="small" style="width: 56px;" />
          </div>
          <button type="button" class="apply-btn" @click="applyColorToAll">
            <Icon icon="mdi:apply-edit" :size="12" /> 统一应用到全部
          </button>
        </div>

        <div class="color-group">
          <div class="color-label">背景颜色</div>
          <div class="color-row">
            <button v-for="c in bgColorPresets" :key="c" type="button"
              class="color-dot" :class="{ active: form.bg_color === c }"
              :style="{ background: c || 'transparent' }" @click="form.bg_color = c" />
            <n-color-picker v-model:value="form.bg_color" :show-alpha="true" format="hex" size="small" style="width: 56px;" />
          </div>
          <button type="button" class="apply-btn" @click="applyBgColorToAll">
            <Icon icon="mdi:apply-edit" :size="12" /> 统一应用到全部
          </button>
        </div>
      </div>
    </div>

  </n-modal>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { Icon } from '@iconify/vue'
import { NModal, NCheckbox, NForm, NFormItem, NSelect, NInput, NRadioGroup, NRadioButton, NColorPicker, NButton, NCollapse, NCollapseItem } from 'naive-ui'
import { useAppStore } from '../../stores/app'
import type { Card, PanelGroup } from '../../types'

const appStore = useAppStore()
const showPreview = ref(true)
const canvasTransparent = ref(false)

const props = defineProps<{
  show: boolean
  card: Card | null
  groups: PanelGroup[]
  defaultGroupId: string
}>()

const emit = defineEmits<{
  'update:show': [value: boolean]
  save: [data: any]
  'apply-color-all': [data: { cardIds: string[]; icon_color: string }]
  'apply-bg-color-all': [data: { cardIds: string[]; bg_color: string }]
}>()

const saving = ref(false)
const fetchingIcon = ref(false)
const showIconPicker = ref(false)

const form = ref({
  group_id: '',
  title: '',
  url: '',
  url_internal: '',
  icon: '',
  icon_color: '#2080f0',
  bg_color: '',
  description: '',
  open_type: 'new_tab' as string,
})

// Parse icon name from Iconify URL or return as-is
// Supports: https://icon-sets.iconify.design/mdi/web.svg -> mdi:web
//           https://icon-sets.iconify.design/mdi/web -> mdi:web
//           mdi:web -> mdi:web
const parsedIconName = computed(() => {
  const val = form.value.icon.trim()
  if (!val) return 'mdi:compass'
  // Already iconify format
  if (val.includes(':')) return val
  // Try to parse Iconify URL
  const urlMatch = val.match(/iconify\.design\/([\w-]+)\/([\w-]+)/)
  if (urlMatch) return `${urlMatch[1]}:${urlMatch[2]}`
  // Try generic URL with path like /mdi/web
  const pathMatch = val.match(/\/([\w-]+)\/([\w-]+)(?:\.svg)?(?:\?|$)/)
  if (pathMatch && !pathMatch[1].includes('.')) return `${pathMatch[1]}:${pathMatch[2]}`
  return val
})

// Handle icon input - auto-parse URL
function onIconInput(val: string) {
  if (!val) return
  // If it's an Iconify URL, parse and replace with icon name
  const urlMatch = val.match(/iconify\.design\/([\w-]+)\/([\w-]+)/)
  if (urlMatch) {
    form.value.icon = `${urlMatch[1]}:${urlMatch[2]}`
    return
  }
  // Generic URL path
  const pathMatch = val.match(/\/([\w-]+)\/([\w-]+)(?:\.svg)?(?:\?|$)/)
  if (pathMatch && !pathMatch[1].includes('.') && val.includes('http')) {
    form.value.icon = `${pathMatch[1]}:${pathMatch[2]}`
  }
}

// Fetch favicon from URL via server API
async function fetchFavicon() {
  if (!form.value.url || fetchingIcon.value) return
  fetchingIcon.value = true
  try {
    const token = localStorage.getItem('sundash-token')
    const res = await fetch(`/api/favicon?url=${encodeURIComponent(form.value.url)}`, {
      headers: token ? { Authorization: `Bearer ${token}` } : {},
    })
    if (!res.ok) throw new Error('Failed to fetch')
    const data = await res.json()
    if (data.icon_name) {
      form.value.icon = data.icon_name
    } else if (data.favicon_url) {
      // If only favicon URL available, try to find matching iconify icon
      // For now use mdi:web as fallback
      form.value.icon = 'mdi:web'
    } else {
      form.value.icon = 'mdi:web'
    }
  } catch {
    form.value.icon = 'mdi:web'
  } finally {
    fetchingIcon.value = false
  }
}

const groupOptions = computed(() =>
  props.groups.map(g => ({ label: g.name, value: g.id }))
)

const colorPresets = [
  '#007AFF', '#5856D6', '#AF52DE', '#FF2D55', '#FF9500',
  '#34C759', '#00C7BE', '#FF3B30', '#FFCC00', '#8E8E93',
]

const bgColorPresets = [
  '', 'rgba(255,255,255,0.55)', 'rgba(255,255,255,0.75)', 'rgba(255,255,255,0.9)',
  'rgba(0,0,0,0.3)', 'rgba(0,0,0,0.5)', 'rgba(0,0,0,0.7)',
  'rgba(0,122,255,0.15)', 'rgba(88,86,214,0.15)', 'rgba(52,199,89,0.15)',
]

const iconCategories = [
  {
    name: '常用',
    icons: [
      'mdi:web', 'mdi:earth', 'mdi:compass', 'mdi:star', 'mdi:heart',
      'mdi:bookmark', 'mdi:home', 'mdi:magnify', 'mdi:cog', 'mdi:bell',
    ],
  },
  {
    name: '工具',
    icons: [
      'mdi:wrench', 'mdi:hammer', 'mdi:screwdriver', 'mdi:settings', 'mdi:tune',
      'mdi:console', 'mdi:database', 'mdi:server', 'mdi:cloud', 'mdi:shield',
    ],
  },
  {
    name: '媒体',
    icons: [
      'mdi:play', 'mdi:music', 'mdi:image', 'mdi:video', 'mdi:camera',
      'mdi:film', 'mdi:microphone', 'mdi:headphones', 'mdi:radio', 'mdi:podcast',
    ],
  },
  {
    name: '社交',
    icons: [
      'mdi:account-group', 'mdi:chat', 'mdi:email', 'mdi:bell-ring', 'mdi:share',
      'mdi:message-text', 'mdi:forum', 'mdi:at', 'mdi:rss', 'mdi:link-variant',
    ],
  },
  {
    name: '文件',
    icons: [
      'mdi:folder', 'mdi:file-document', 'mdi:download', 'mdi:upload', 'mdi:archive',
      'mdi:notebook', 'mdi:clipboard-text', 'mdi:file-code', 'mdi:file-pdf', 'mdi:file-image',
    ],
  },
  {
    name: '设备',
    icons: [
      'mdi:monitor', 'mdi:laptop', 'mdi:cellphone', 'mdi:tablet', 'mdi:printer',
      'mdi:router-wireless', 'mdi:lan', 'mdi:usb', 'mdi:harddisk', 'mdi:memory',
    ],
  },
]

watch(() => props.show, (val) => {
  if (val) {
    if (props.card) {
      form.value = {
        group_id: props.card.group_id,
        title: props.card.title,
        url: props.card.url,
        url_internal: props.card.url_internal || '',
        icon: props.card.icon || '',
        icon_color: props.card.icon_color || '#2080f0',
        bg_color: (props.card as any).bg_color || '',
        description: props.card.description || '',
        open_type: props.card.open_type || 'new_tab',
      }
    } else {
      form.value = {
        group_id: props.defaultGroupId,
        title: '',
        url: '',
        url_internal: '',
        icon: '',
        icon_color: '#2080f0',
        bg_color: '',
        description: '',
        open_type: 'new_tab',
      }
    }
  }
})

async function handleSave() {
  if (!form.value.title || !form.value.url || !form.value.group_id) {
    return false
  }
  saving.value = true
  try {
    emit('save', { ...form.value })
    // Don't return true - let parent control modal close after API success
  } catch {
    // Error handled by parent
  } finally {
    saving.value = false
  }
}

// Apply current icon color to ALL cards in ALL groups
function applyColorToAll() {
  // Collect card IDs from all groups
  const cardIds: string[] = []
  for (const group of props.groups) {
    if (group.cards) {
      for (const card of group.cards) {
        cardIds.push(card.id)
      }
    }
  }
  if (cardIds.length === 0) return
  
  // Also update global text color for search bar sync
  appStore.setTextColor(form.value.icon_color)
  
  // Emit event to parent for batch update
  emit('apply-color-all', {
    cardIds,
    icon_color: form.value.icon_color,
  })
}

// Apply current background color to ALL cards in ALL groups
function applyBgColorToAll() {
  // Collect card IDs from all groups
  const cardIds: string[] = []
  for (const group of props.groups) {
    if (group.cards) {
      for (const card of group.cards) {
        cardIds.push(card.id)
      }
    }
  }
  if (cardIds.length === 0) return
  
  // Also update global background color for search bar sync
  appStore.setBgColor(form.value.bg_color)
  
  // Emit event to parent for batch update
  emit('apply-bg-color-all', {
    cardIds,
    bg_color: form.value.bg_color,
  })
}
</script>

<style scoped>
/* ===== Modal Width ===== */
.card-editor-modal {
  width: 60vw !important;
  max-width: calc(100vw - 24px);
}
:deep(.n-dialog) {
  width: 60vw !important;
  max-width: calc(100vw - 24px);
  border-radius: 14px !important;
  overflow: hidden;
}
:deep(.n-dialog__content) {
  padding: 0 !important;
  margin: 0 !important;
}
:deep(.n-dialog__title) {
  padding: 12px 20px 10px !important;
  font-size: 14px !important;
  font-weight: 600 !important;
  border-bottom: 1px solid var(--sd-border);
  margin: 0 !important;
}
:deep(.n-dialog__action) {
  padding: 10px 20px !important;
  border-top: 1px solid var(--sd-border) !important;
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin: 0 !important;
}
</style>

<style>
/* Global override for card editor dialog width */
.card-editor-modal.n-dialog,
.card-editor-modal .n-dialog {
  width: 60vw !important;
  max-width: calc(100vw - 24px) !important;
}
/* Icon picker popover z-index above dialog */
</style>

<style scoped>
/* ===== Two-Column Layout ===== */
.editor-body {
  display: grid;
  grid-template-columns: 1fr 300px;
  min-height: 440px;
  max-height: calc(100vh - 140px);
}

.editor-left {
  padding: 14px 20px;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.editor-right {
  padding: 14px 16px;
  border-left: 1px solid var(--sd-border);
  background: var(--sd-bg-surface);
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

/* ===== Left: Preview ===== */
.preview-section {
  flex-shrink: 0;
}

.preview-header {
  display: flex;
  gap: 14px;
  margin-bottom: 8px;
}

.preview-box {
  background: linear-gradient(135deg, #e8f0fe 0%, #f0f4ff 100%);
  border-radius: 10px;
  padding: 12px;
  display: flex;
  gap: 12px;
  justify-content: center;
  flex-wrap: wrap;
}

.preview-box.canvas-transparent {
  background: transparent;
  border: 1px dashed var(--sd-border);
  background-image: linear-gradient(45deg, #e0e0e0 25%, transparent 25%, transparent 75%, #e0e0e0 75%),
    linear-gradient(45deg, #e0e0e0 25%, transparent 25%, transparent 75%, #e0e0e0 75%);
  background-size: 12px 12px;
  background-position: 0 0, 6px 6px;
  background-color: #f5f5f5;
}

.preview-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  border-radius: 14px;
  background: rgba(42, 42, 42, 0.42);
  color: white;
  width: 179px;
  height: 53.6px;
  box-sizing: border-box;
}

.preview-card-h {
  flex-direction: row;
}



.preview-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.preview-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.preview-title {
  font-weight: 500;
  font-size: 12px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.preview-desc {
  font-size: 11px;
  opacity: 0.75;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}



/* ===== Left: Form ===== */
.editor-form {
  flex: 1;
}

.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 0 14px;
}

.form-grid > :deep(.n-form-item) {
  margin-bottom: 8px;
}

.form-grid > :deep(.n-form-item .n-form-item-label__text) {
  font-size: 12px;
  font-weight: 500;
}

.open-type-row {
  display: flex;
  align-items: center;
  height: 32px;
}

/* ===== Right: Header ===== */
.right-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-shrink: 0;
}

.right-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--sd-text-primary);
}

.iconify-link {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  font-size: 11px;
  color: var(--sd-primary);
  text-decoration: none;
  padding: 2px 8px;
  border-radius: 6px;
  transition: background 0.15s;
}
.iconify-link:hover {
  background: var(--sd-primary-light);
}

/* ===== Right: Icon Input ===== */
.icon-input-row {
  display: flex;
  gap: 6px;
  align-items: center;
  flex-shrink: 0;
}

.fetch-btn {
  display: inline-flex;
  align-items: center;
  gap: 3px;
  padding: 0 10px;
  height: 28px;
  border: 1px solid var(--sd-border);
  border-radius: 6px;
  background: var(--sd-bg-card);
  color: var(--sd-text-secondary);
  font-size: 11px;
  font-family: var(--sd-font);
  cursor: pointer;
  transition: all 0.15s;
  white-space: nowrap;
  flex-shrink: 0;
}
.fetch-btn:hover:not(:disabled) {
  background: var(--sd-primary-light);
  color: var(--sd-primary);
  border-color: var(--sd-primary);
}
.fetch-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.spin-icon {
  animation: spin 1s linear infinite;
}
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* ===== Right: Icon Picker Trigger ===== */
.icon-picker-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 8px 12px;
  border: 1px solid var(--sd-border);
  border-radius: 8px;
  background: var(--sd-bg-card);
  color: var(--sd-text-secondary);
  font-size: 12px;
  font-family: var(--sd-font);
  cursor: pointer;
  transition: all 0.15s;
  flex-shrink: 0;
}
.icon-picker-trigger:hover {
  border-color: var(--sd-primary);
  background: var(--sd-primary-light);
}

.trigger-icon {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  background: rgba(0,122,255,0.06);
  flex-shrink: 0;
}

.trigger-label {
  flex: 1;
  text-align: left;
  font-weight: 500;
}

.trigger-arrow {
  opacity: 0.5;
  transition: transform 0.2s;
}

/* ===== Popover Content ===== */
.icon-popover-content {
  max-height: 320px;
  overflow-y: auto;
}

.icon-cat {
  margin-bottom: 8px;
}

.icon-cat:last-child {
  margin-bottom: 0;
}

.icon-cat-name {
  font-size: 10px;
  color: var(--sd-text-tertiary);
  font-weight: 600;
  margin-bottom: 5px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.icon-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.icon-btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid transparent;
  border-radius: 8px;
  background: transparent;
  color: var(--sd-text-secondary);
  cursor: pointer;
  transition: all 0.12s;
}
.icon-btn:hover {
  background: rgba(0,122,255,0.08);
  color: var(--sd-primary);
  border-color: rgba(0,122,255,0.12);
}
.icon-btn.active {
  background: rgba(0,122,255,0.15);
  color: var(--sd-primary);
  border-color: var(--sd-primary);
}

/* ===== Right: Divider ===== */
.right-divider {
  height: 1px;
  background: var(--sd-border);
  margin: 2px 0;
  flex-shrink: 0;
}

/* ===== Right: Colors ===== */
.color-group {
  flex-shrink: 0;
}

.color-label {
  font-size: 11px;
  font-weight: 600;
  color: var(--sd-text-secondary);
  margin-bottom: 6px;
}

.color-row {
  display: flex;
  gap: 5px;
  align-items: flex-start;
  flex-wrap: wrap;
  margin-bottom: 8px;
}

:deep(.color-row .n-color-picker) {
  height: 24px !important;
  width: 56px !important;
  flex-shrink: 0;
}
:deep(.color-row .n-color-picker .n-color-picker-trigger) {
  height: 24px !important;
  width: 56px !important;
  border-radius: 6px;
  overflow: hidden;
}
/* Hide checkerboard transparency pattern */
:deep(.color-row .n-color-picker .n-color-picker-checkboard) {
  display: none !important;
}
/* Hide color value text, show only the color swatch */
:deep(.color-row .n-color-picker .n-color-picker-trigger__value) {
  display: none !important;
}
/* Ensure color fill is visible */
:deep(.color-row .n-color-picker .n-color-picker__fill) {
  position: relative;
  z-index: 1;
  width: 100%;
  height: 100%;
  border-radius: 4px;
}
/* Color picker panel value input */
:deep(.n-color-picker .n-color-picker-value) {
  padding-left: 175px;
}
:deep(.n-color-picker .n-color-picker-panel .n-color-picker-value-input) {
  font-size: 12px !important;
}

.color-dot {
  width: 22px;
  height: 22px;
  border-radius: 50%;
  border: 1.5px solid rgba(0,0,0,0.12);
  cursor: pointer;
  transition: all 0.12s;
  flex-shrink: 0;
}
.color-dot:hover {
  transform: scale(1.15);
}
.color-dot.active {
  border-color: var(--sd-primary);
  box-shadow: 0 0 0 2px var(--sd-bg-card), 0 0 0 3.5px var(--sd-primary);
}

.apply-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border: 1px solid var(--sd-border);
  border-radius: 6px;
  background: var(--sd-bg-card);
  color: var(--sd-text-secondary);
  font-size: 11px;
  font-family: var(--sd-font);
  cursor: pointer;
  transition: all 0.12s;
  white-space: nowrap;
}
.apply-btn:hover {
  background: var(--sd-primary-light);
  color: var(--sd-primary);
  border-color: var(--sd-primary);
}

/* ===== Responsive ===== */
@media (max-width: 900px) {
  :deep(.n-dialog.n-modal) {
    width: calc(100vw - 24px) !important;
  }
  .editor-body {
    grid-template-columns: 1fr;
    max-height: calc(100vh - 100px);
  }
  .editor-right {
    border-left: none;
    border-top: 1px solid var(--sd-border);
  }
}

@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
