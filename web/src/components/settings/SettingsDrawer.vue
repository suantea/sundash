<template>
  <n-drawer :show="show" @update:show="$emit('update:show', $event)" :width="drawerWidth" placement="right">
    <n-drawer-content title="设置" :native-scrollbar="false">
      <div class="drawer-settings">
        <!-- ===== 1. 页面 ===== -->
        <section class="settings-section">
          <div class="section-label">页面</div>
          <div class="settings-card">
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
                  <Icon icon="mdi:lan" :size="18" color="#007AFF" />
                </div>
                <div>
                  <div class="setting-title">网络模式</div>
                  <div class="setting-desc">切换内网或外网链接</div>
                </div>
              </div>
              <div class="seg-control">
                <button :class="['seg-btn', { active: appStore.networkMode === 'internal' }]" @click="setNetwork('internal')">内网</button>
                <button :class="['seg-btn', { active: appStore.networkMode === 'external' }]" @click="setNetwork('external')">外网</button>
              </div>
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(255,149,0,0.1);">
                  <Icon icon="mdi:image" :size="18" color="#FF9500" />
                </div>
                <div>
                  <div class="setting-title">Logo 类型</div>
                  <div class="setting-desc">图片或文字显示</div>
                </div>
              </div>
              <div class="seg-control">
                <button :class="['seg-btn', { active: appStore.logoType === 'text' }]" @click="appStore.setLogoType('text')">文字</button>
                <button :class="['seg-btn', { active: appStore.logoType === 'image' }]" @click="appStore.setLogoType('image')">图片</button>
              </div>
            </div>
            <div v-if="appStore.logoType === 'image'" class="setting-sub">
              <input v-model="logoImageInput" class="bg-url-input" placeholder="输入 Logo 图片 URL" @change="appStore.setLogoImageUrl(logoImageInput)" />
            </div>
            <div v-if="appStore.logoType === 'text'" class="setting-sub">
              <input v-model="logoTextInput" class="bg-url-input" placeholder="输入 Logo 文字" @change="appStore.setLogoText(logoTextInput)" />
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
                  <Icon icon="mdi:text-short" :size="18" color="#007AFF" />
                </div>
                <div>
                  <div class="setting-title">站点标题</div>
                  <div class="setting-desc">浏览器标签页标题</div>
                </div>
              </div>
              <input v-model="siteTitleInput" class="mini-text-input" placeholder="SunDash" @change="updateSiteTitle" />
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row" style="flex-direction: column; align-items: stretch; gap: 12px;">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(52,199,89,0.1);">
                  <Icon icon="mdi:code-tags" :size="18" color="#34C759" />
                </div>
                <div>
                  <div class="setting-title">自定义页脚</div>
                  <div class="setting-desc">支持 HTML 代码</div>
                </div>
              </div>
              <textarea v-model="footerInput" class="footer-input" placeholder="<p>Powered by SunDash</p>" rows="2" @change="appStore.setFooterHtml(footerInput)"></textarea>
            </div>
          </div>
        </section>

        <!-- ===== 2. 外观 ===== -->
        <section class="settings-section">
          <div class="section-label">外观</div>
          <div class="settings-card">
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(255,149,0,0.1);">
                  <Icon icon="mdi:theme-light-dark" :size="18" color="#FF9500" />
                </div>
                <div>
                  <div class="setting-title">主题模式</div>
                  <div class="setting-desc">选择外观风格</div>
                </div>
              </div>
              <div class="seg-control">
                <button :class="['seg-btn', { active: appStore.themeMode === 'light' }]" @click="appStore.setTheme('light')">浅色</button>
                <button :class="['seg-btn', { active: appStore.themeMode === 'dark' }]" @click="appStore.setTheme('dark')">深色</button>
                <button :class="['seg-btn', { active: appStore.themeMode === 'system' }]" @click="appStore.setTheme('system')">自动</button>
              </div>
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(52,199,89,0.1);">
                  <Icon icon="mdi:image-outline" :size="18" color="#34C759" />
                </div>
                <div>
                  <div class="setting-title">背景</div>
                  <div class="setting-desc">选择页面背景样式</div>
                </div>
              </div>
              <div class="seg-control">
                <button :class="['seg-btn', { active: appStore.wallpaperType === 'default' }]" @click="setWallpaperTypeWithSync('default')">默认</button>
                <button :class="['seg-btn', { active: appStore.wallpaperType === 'gradient' }]" @click="setWallpaperTypeWithSync('gradient')">渐变</button>
                <button :class="['seg-btn', { active: appStore.wallpaperType === 'bing' }]" @click="setWallpaperTypeWithSync('bing')">必应</button>
                <button :class="['seg-btn', { active: appStore.wallpaperType === 'custom' }]" @click="setWallpaperTypeWithSync('custom')">自定义</button>
              </div>
            </div>
            <div v-if="appStore.wallpaperType === 'bing'" class="setting-sub">
              <button class="apply-btn" @click="fetchBingWallpaper">刷新必应壁纸</button>
              <span v-if="appStore.wallpaperUrl" class="bing-preview">
                <img :src="appStore.wallpaperUrl" alt="预览" class="bing-preview-img" />
              </span>
            </div>
            <div v-if="appStore.wallpaperType === 'custom'" class="setting-sub">
              <input v-model="customUrl" class="bg-url-input" placeholder="输入背景图片 URL" />
              <button class="apply-btn" @click="applyCustomBackground">应用</button>
            </div>
            <template v-if="appStore.wallpaperType !== 'default'">
              <div class="setting-divider"></div>
              <div class="setting-row">
                <div class="setting-left">
                  <div class="setting-icon" style="background: rgba(255,149,0,0.1);">
                    <Icon icon="mdi:blur" :size="18" color="#FF9500" />
                  </div>
                  <div>
                    <div class="setting-title">模糊</div>
                    <div class="setting-desc">{{ appStore.wallpaperBlur }}px</div>
                  </div>
                </div>
                <div class="slider-control">
                  <input type="range" min="0" max="20" :value="appStore.wallpaperBlur" @input="setWallpaperBlurWithSync(parseInt(($event.target as HTMLInputElement).value))" class="range-slider" />
                </div>
              </div>
              <div class="setting-divider"></div>
              <div class="setting-row">
                <div class="setting-left">
                  <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
                    <Icon icon="mdi:opacity" :size="18" color="#007AFF" />
                  </div>
                  <div>
                    <div class="setting-title">透明度</div>
                    <div class="setting-desc">{{ appStore.wallpaperOpacity }}%</div>
                  </div>
                </div>
                <div class="slider-control">
                  <input type="range" min="10" max="100" :value="appStore.wallpaperOpacity" @input="setWallpaperOpacityWithSync(parseInt(($event.target as HTMLInputElement).value))" class="range-slider" />
                </div>
              </div>
            </template>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
                  <Icon icon="mdi:palette" :size="18" color="#007AFF" />
                </div>
                <div>
                  <div class="setting-title">主题色</div>
                  <div class="setting-desc">自定义强调色</div>
                </div>
              </div>
              <div class="color-picker-wrap">
                <input type="color" :value="appStore.primaryColor || '#007AFF'" @input="appStore.setPrimaryColor(($event.target as HTMLInputElement).value)" class="color-input" />
                <button v-if="appStore.primaryColor" class="reset-btn" @click="appStore.setPrimaryColor('')">重置</button>
              </div>
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(255,149,0,0.1);">
                  <Icon icon="mdi:format-color-fill" :size="18" color="#FF9500" />
                </div>
                <div>
                  <div class="setting-title">边框颜色</div>
                  <div class="setting-desc">自定义边框色</div>
                </div>
              </div>
              <div class="color-picker-wrap">
                <input type="color" :value="appStore.borderColor || '#000000'" @input="appStore.setBorderColor(($event.target as HTMLInputElement).value)" class="color-input" />
                <button v-if="appStore.borderColor" class="reset-btn" @click="appStore.setBorderColor('')">重置</button>
              </div>
            </div>
          </div>
        </section>

        <!-- ===== 3. 布局 ===== -->
        <section class="settings-section">
          <div class="section-label">布局</div>
          <div class="settings-card">
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
                  <Icon icon="mdi:arrow-expand-horizontal" :size="18" color="#007AFF" />
                </div>
                <div>
                  <div class="setting-title">内容宽度</div>
                  <div class="setting-desc">{{ appStore.contentMaxWidth }}</div>
                </div>
              </div>
              <div class="slider-control">
                <input type="range" min="50" max="100" step="5" :value="parseInt(appStore.contentMaxWidth) || 80" @input="appStore.setContentMaxWidth(($event.target as HTMLInputElement).value + '%')" class="range-slider" />
              </div>
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(52,199,89,0.1);">
                  <Icon icon="mdi:format-horizontal-align-left" :size="18" color="#34C759" />
                </div>
                <div>
                  <div class="setting-title">两侧边距</div>
                  <div class="setting-desc">{{ appStore.contentPaddingX }}</div>
                </div>
              </div>
              <div class="slider-control">
                <input type="range" min="0" max="10" step="1" :value="parseInt(appStore.contentPaddingX) || 5" @input="appStore.setContentPaddingX(($event.target as HTMLInputElement).value + '%')" class="range-slider" />
              </div>
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(175,82,222,0.1);">
                  <Icon icon="mdi:format-vertical-align-top" :size="18" color="#AF52DE" />
                </div>
                <div>
                  <div class="setting-title">顶部边距</div>
                  <div class="setting-desc">{{ appStore.contentPaddingTop }}</div>
                </div>
              </div>
              <div class="slider-control">
                <input type="range" min="0" max="20" step="1" :value="parseInt(appStore.contentPaddingTop) || 10" @input="appStore.setContentPaddingTop(($event.target as HTMLInputElement).value + '%')" class="range-slider" />
              </div>
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(255,149,0,0.1);">
                  <Icon icon="mdi:format-vertical-align-bottom" :size="18" color="#FF9500" />
                </div>
                <div>
                  <div class="setting-title">底部边距</div>
                  <div class="setting-desc">{{ appStore.contentPaddingBottom }}</div>
                </div>
              </div>
              <div class="slider-control">
                <input type="range" min="0" max="10" step="1" :value="parseInt(appStore.contentPaddingBottom) || 5" @input="appStore.setContentPaddingBottom(($event.target as HTMLInputElement).value + '%')" class="range-slider" />
              </div>
            </div>
          </div>
        </section>

        <!-- ===== 4. 组件 ===== -->
        <section class="settings-section">
          <div class="section-label">组件</div>
          <div class="settings-card">
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(52,199,89,0.1);">
                  <Icon icon="mdi:clock-outline" :size="18" color="#34C759" />
                </div>
                <div>
                  <div class="setting-title">时钟</div>
                  <div class="setting-desc">显示当前时间</div>
                </div>
              </div>
              <n-switch v-model:value="appStore.clockShow" @update:value="appStore.setClockShow" />
            </div>
            <div v-if="appStore.clockShow" class="setting-sub" style="flex-direction: column; align-items: stretch; gap: 4px;">
              <div class="setting-row" style="padding: 6px 0;">
                <span class="mini-label">显示秒数</span>
                <n-switch v-model:value="appStore.clockShowSeconds" @update:value="appStore.setClockShowSeconds" size="small" />
              </div>
              <div class="setting-row" style="padding: 6px 0;">
                <span class="mini-label">时间格式</span>
                <div class="seg-control">
                  <button :class="['seg-btn', { active: appStore.clockFormat === '24' }]" @click="appStore.setClockFormat('24')">24小时</button>
                  <button :class="['seg-btn', { active: appStore.clockFormat === '12' }]" @click="appStore.setClockFormat('12')">12小时</button>
                </div>
              </div>
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(175,82,222,0.1);">
                  <Icon icon="mdi:server-network" :size="18" color="#AF52DE" />
                </div>
                <div>
                  <div class="setting-title">系统状态</div>
                  <div class="setting-desc">显示系统信息</div>
                </div>
              </div>
              <n-switch v-model:value="appStore.showSystemStatus" @update:value="appStore.setShowSystemStatus" />
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(255,149,0,0.1);">
                  <Icon icon="mdi:format-size" :size="18" color="#FF9500" />
                </div>
                <div>
                  <div class="setting-title">标签大小</div>
                  <div class="setting-desc">{{ appStore.cardLabelSize }}px</div>
                </div>
              </div>
              <div class="slider-control">
                <input type="range" min="10" max="18" step="1" :value="parseInt(appStore.cardLabelSize) || 12" @input="appStore.setCardLabelSize(($event.target as HTMLInputElement).value)" class="range-slider" />
              </div>
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
                  <Icon icon="mdi:resize" :size="18" color="#007AFF" />
                </div>
                <div>
                  <div class="setting-title">卡片大小</div>
                  <div class="setting-desc">{{ parseInt(appStore.cardItemSize) || 0 > 0 ? '+' + appStore.cardItemSize + 'px' : '默认' }}</div>
                </div>
              </div>
              <div class="slider-control">
                <input type="range" min="-4" max="12" step="1" :value="parseInt(appStore.cardItemSize) || 0" @input="appStore.setCardItemSize(($event.target as HTMLInputElement).value)" class="range-slider" />
              </div>
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(175,82,222,0.1);">
                  <Icon icon="mdi:view-grid" :size="18" color="#AF52DE" />
                </div>
                <div>
                  <div class="setting-title">每行书签数</div>
                  <div class="setting-desc">{{ parseInt(appStore.cardsPerRow) || 5 }} 个</div>
                </div>
              </div>
              <div class="seg-control">
                <button :class="['seg-btn', { active: cardsPerRowVal === 3 }]" @click="appStore.setCardsPerRow('3')">3</button>
                <button :class="['seg-btn', { active: cardsPerRowVal === 4 }]" @click="appStore.setCardsPerRow('4')">4</button>
                <button :class="['seg-btn', { active: cardsPerRowVal === 5 }]" @click="appStore.setCardsPerRow('5')">5</button>
                <button :class="['seg-btn', { active: cardsPerRowVal === 6 }]" @click="appStore.setCardsPerRow('6')">6</button>
              </div>
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(52,199,89,0.1);">
                  <Icon icon="mdi:eye-off" :size="18" color="#34C759" />
                </div>
                <div>
                  <div class="setting-title">分组背景透明</div>
                  <div class="setting-desc">隐藏分组卡片背景</div>
                </div>
              </div>
              <n-switch v-model:value="appStore.groupCardTransparent" @update:value="appStore.setGroupCardTransparent" />
            </div>
          </div>
        </section>

        <!-- ===== 5. 数据 ===== -->
        <section class="settings-section">
          <div class="section-label">数据</div>
          <div class="settings-card">
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
                  <Icon icon="mdi:export" :size="18" color="#007AFF" />
                </div>
                <div>
                  <div class="setting-title">导出数据</div>
                  <div class="setting-desc">备份所有分组和卡片</div>
                </div>
              </div>
              <button class="action-btn" @click="handleExport">导出</button>
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(52,199,89,0.1);">
                  <Icon icon="mdi:import" :size="18" color="#34C759" />
                </div>
                <div>
                  <div class="setting-title">导入数据</div>
                  <div class="setting-desc">从备份文件恢复</div>
                </div>
              </div>
              <div style="display: flex; gap: 8px;">
                <button class="action-btn" @click="importTemplate">模板</button>
                <button class="action-btn" @click="triggerImport">导入</button>
              </div>
            </div>
          </div>
        </section>
      </div>

      <template #footer>
        <div class="drawer-footer">
          <span class="footer-text">设置会自动保存</span>
        </div>
      </template>
    </n-drawer-content>
  </n-drawer>
  <input ref="fileInput" type="file" accept=".json" style="display: none" @change="handleImport" />
</template>

<script setup lang="ts">
import { ref, watch, computed, onMounted, onUnmounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useMessage, NDrawer, NDrawerContent, NSwitch } from 'naive-ui'
import { useAppStore } from '../../stores/app'
import { usePanelStore } from '../../stores/panel'
import { importTemplates } from '../../data/importTemplates'

defineProps<{ show: boolean }>()
const emit = defineEmits<{ 'update:show': [val: boolean] }>()

const appStore = useAppStore()
const panelStore = usePanelStore()

// Responsive drawer width
const windowWidth = ref(window.innerWidth)
onMounted(() => {
  const onResize = () => { windowWidth.value = window.innerWidth }
  window.addEventListener('resize', onResize)
  onUnmounted(() => window.removeEventListener('resize', onResize))
})
const drawerWidth = computed(() => windowWidth.value <= 480 ? windowWidth.value : 420)
const cardsPerRowVal = computed(() => parseInt(appStore.cardsPerRow) || 5)
const message = useMessage()
const fileInput = ref<HTMLInputElement>()

// Local inputs
const customUrl = ref('')
const logoImageInput = ref(appStore.logoImageUrl)
const logoTextInput = ref(appStore.logoText)
const siteTitleInput = ref(appStore.siteTitle)
const footerInput = ref(appStore.footerHtml)

// Sync local inputs when store changes
watch(() => appStore.logoImageUrl, v => logoImageInput.value = v)
watch(() => appStore.logoText, v => logoTextInput.value = v)
watch(() => appStore.siteTitle, v => siteTitleInput.value = v)
watch(() => appStore.footerHtml, v => footerInput.value = v)

function setNetwork(mode: 'internal' | 'external') {
  appStore.networkMode = mode
  localStorage.setItem('sundash-network', mode)
}

// Wallpaper sync
async function saveWallpaperSetting(key: string, value: string) {
  try {
    await panelStore.updateSetting(key, value)
  } catch (e) {
    console.error('Failed to save setting:', e)
  }
}

function setWallpaperTypeWithSync(type: 'default' | 'gradient' | 'bing' | 'custom') {
  appStore.setWallpaperType(type)
  saveWallpaperSetting('wallpaper_type', type)
}

function setWallpaperBlurWithSync(blur: number) {
  appStore.setWallpaperBlur(blur)
  saveWallpaperSetting('wallpaper_blur', blur.toString())
}

function setWallpaperOpacityWithSync(opacity: number) {
  appStore.setWallpaperOpacity(opacity)
  saveWallpaperSetting('wallpaper_opacity', opacity.toString())
}

async function fetchBingWallpaper() {
  try {
    const { api } = await import('../../api')
    const res = await api.get('wallpaper/bing')
    if (res.data && res.data.images && res.data.images[0]) {
      appStore.setWallpaperUrl(res.data.images[0].url)
      message.success('壁纸已更新')
    }
  } catch (e) {
    message.error('获取壁纸失败')
  }
}

function applyCustomBackground() {
  if (!customUrl.value.trim()) {
    message.error('请输入图片 URL')
    return
  }
  appStore.setWallpaperUrl(customUrl.value.trim())
  message.success('自定义背景已应用')
}

function updateSiteTitle() {
  const title = siteTitleInput.value.trim() || 'SunDash'
  appStore.setSiteTitle(title)
}

function handleExport() {
  // Collect all sundash-* settings from localStorage
  const settings: Record<string, string> = {}
  for (let i = 0; i < localStorage.length; i++) {
    const key = localStorage.key(i)
    if (key && key.startsWith('sundash-')) {
      settings[key] = localStorage.getItem(key) || ''
    }
  }
  const data = {
    version: '1.0.0',
    exported_at: new Date().toISOString(),
    groups: panelStore.groups,
    settings,
  }
  const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `sundash-backup-${new Date().toISOString().slice(0, 10)}.json`
  a.click()
  URL.revokeObjectURL(url)
  message.success('数据已导出')
}

function triggerImport() {
  fileInput.value?.click()
}

// Shared import logic: restores settings and appends groups/cards (used by
// both file import and the built-in template).
async function doImport(data: any) {
  if (data.settings && typeof data.settings === 'object') {
    for (const [key, value] of Object.entries(data.settings)) {
      if (key.startsWith('sundash-') && typeof value === 'string') {
        localStorage.setItem(key, value)
      }
    }
    await appStore.loadSettingsFromServer()
    message.success('设置已恢复')
  }
  if (data.groups && Array.isArray(data.groups)) {
    for (const group of data.groups) {
      const newGroup = await panelStore.createGroup(group.name)
      if (group.cards && Array.isArray(group.cards)) {
        for (const card of group.cards) {
          await panelStore.createCard({
            group_id: newGroup.id,
            title: card.title,
            url: card.url,
            url_internal: card.url_internal || '',
            icon: card.icon || '',
            icon_color: card.icon_color || '',
            bg_color: card.bg_color || '',
            description: card.description || '',
            open_type: card.open_type || 'new_tab',
          })
        }
      }
    }
    await panelStore.fetchPanel()
    message.success(`导入成功：${data.groups.length} 个分组`)
  }
}

async function handleImport(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  try {
    const text = await file.text()
    await doImport(JSON.parse(text))
  } catch {
    message.error('导入失败：文件格式错误')
  } finally {
    input.value = ''
  }
}

function importTemplate() {
  const t = importTemplates[0]
  const cardCount = t.groups.reduce((sum, g) => sum + (g.cards?.length || 0), 0)
  message.info(`导入模板「${t.name}」：${t.groups.length} 个分组、${cardCount} 个书签`)
  doImport({ groups: t.groups }).catch(() => message.error('模板导入失败'))
}
</script>

<style scoped>
.drawer-settings {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 4px 0 20px;
}

.section-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--sd-text-secondary);
  margin-bottom: 8px;
  padding-left: 4px;
}

.settings-card {
  background: rgba(255,255,255,0.85);
  backdrop-filter: blur(20px) saturate(180%);
  border-radius: 14px;
  border: 1px solid rgba(0,0,0,0.06);
  overflow: hidden;
}

:root[data-theme="dark"] .settings-card {
  background: rgba(28,28,30,0.8);
  border-color: rgba(255,255,255,0.08);
}

.setting-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 14px;
  min-height: 48px;
}

.setting-left {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  min-width: 0;
}

.setting-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.setting-title {
  font-size: 14px;
  font-weight: 500;
  color: var(--sd-text-primary);
  line-height: 1.3;
}

.setting-desc {
  font-size: 12px;
  color: var(--sd-text-tertiary);
  line-height: 1.3;
}

.setting-divider {
  height: 1px;
  background: rgba(0,0,0,0.05);
  margin: 0 14px;
}

:root[data-theme="dark"] .setting-divider {
  background: rgba(255,255,255,0.06);
}

.setting-sub {
  padding: 8px 14px 14px;
  display: flex;
  align-items: center;
  gap: 8px;
}

/* Controls */
.seg-control {
  display: flex;
  background: rgba(0,0,0,0.04);
  border-radius: 8px;
  padding: 2px;
  gap: 2px;
}

:root[data-theme="dark"] .seg-control {
  background: rgba(255,255,255,0.06);
}

.seg-btn {
  padding: 6px 12px;
  font-size: 12px;
  font-weight: 500;
  font-family: var(--sd-font);
  border: none;
  border-radius: 6px;
  background: transparent;
  color: var(--sd-text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.seg-btn.active {
  background: var(--sd-bg-card);
  color: var(--sd-text-primary);
  box-shadow: 0 1px 3px rgba(0,0,0,0.08);
}

:root[data-theme="dark"] .seg-btn.active {
  background: rgba(44,44,46,0.9);
  box-shadow: 0 1px 3px rgba(0,0,0,0.3);
}

.slider-control {
  display: flex;
  align-items: center;
  gap: 8px;
}

.range-slider {
  width: 100px;
  height: 4px;
  -webkit-appearance: none;
  appearance: none;
  background: rgba(0,0,0,0.1);
  border-radius: 2px;
  outline: none;
}

.range-slider::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #007AFF;
  cursor: pointer;
  box-shadow: 0 1px 4px rgba(0,122,255,0.3);
}

.mini-label {
  font-size: 12px;
  color: var(--sd-text-secondary);
}

.bg-url-input {
  flex: 1;
  height: 32px;
  padding: 0 10px;
  border: 1px solid rgba(0,0,0,0.1);
  border-radius: 8px;
  font-size: 12px;
  font-family: var(--sd-font);
  background: var(--sd-bg-card);
  color: var(--sd-text-primary);
  outline: none;
  transition: border-color 0.2s;
}

.bg-url-input:focus {
  border-color: #007AFF;
}

.apply-btn {
  padding: 6px 14px;
  background: #007AFF;
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 500;
  font-family: var(--sd-font);
  cursor: pointer;
  transition: background 0.15s;
  white-space: nowrap;
}

.apply-btn:hover {
  background: #0066DD;
}

.bing-preview {
  display: block;
  margin-top: 8px;
}

.bing-preview-img {
  width: 100%;
  height: 80px;
  object-fit: cover;
  border-radius: 8px;
}

.color-picker-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
}

.color-input {
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  padding: 0;
}

.reset-btn {
  padding: 4px 8px;
  font-size: 11px;
  font-family: var(--sd-font);
  background: rgba(0,0,0,0.05);
  border: none;
  border-radius: 6px;
  color: var(--sd-text-secondary);
  cursor: pointer;
}

.mini-text-input {
  width: 120px;
  height: 32px;
  padding: 0 10px;
  border: 1px solid rgba(0,0,0,0.1);
  border-radius: 8px;
  font-size: 13px;
  font-family: var(--sd-font);
  background: var(--sd-bg-card);
  color: var(--sd-text-primary);
  outline: none;
}

.mini-text-input:focus {
  border-color: #007AFF;
}

.footer-input {
  width: 100%;
  padding: 10px;
  border: 1px solid rgba(0,0,0,0.1);
  border-radius: 8px;
  font-size: 12px;
  font-family: var(--sd-font);
  background: var(--sd-bg-card);
  color: var(--sd-text-primary);
  resize: vertical;
  outline: none;
}

.footer-input:focus {
  border-color: #007AFF;
}

.action-btn {
  padding: 6px 16px;
  background: rgba(0,122,255,0.1);
  color: #007AFF;
  border: none;
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  font-family: var(--sd-font);
  cursor: pointer;
  transition: background 0.15s;
}

.action-btn:hover {
  background: rgba(0,122,255,0.15);
}

.drawer-footer {
  text-align: center;
}

.footer-text {
  font-size: 12px;
  color: var(--sd-text-tertiary);
}
</style>