<template>
  <div class="settings-page">
    <div class="page-bg">
      <div class="bg-orb bg-orb-1"></div>
      <div class="bg-orb bg-orb-2"></div>
    </div>

    <header class="page-header">
      <button class="back-btn" @click="$router.push('/')">
        <Icon icon="mdi:chevron-left" :size="18" />
        <span>返回</span>
      </button>
      <h1>设置</h1>
      <div style="width: 70px;"></div>
    </header>

    <div class="settings-content">
      <!-- Network -->
      <section class="settings-section">
        <div class="section-label">网络</div>
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
        </div>
      </section>

      <!-- Appearance -->
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
                <div class="setting-title">自定义背景</div>
                <div class="setting-desc">选择页面背景样式</div>
              </div>
            </div>
            <div class="seg-control">
              <button :class="['seg-btn', { active: appStore.wallpaperType === 'default' }]" @click="setWallpaperTypeWithSync('default')">默认</button>
              <button :class="['seg-btn', { active: appStore.wallpaperType === 'gradient' }]" @click="setWallpaperTypeWithSync('gradient')">渐变</button>
              <button :class="['seg-btn', { active: appStore.wallpaperType === 'bing' }]" @click="setWallpaperTypeWithSync('bing')">每日必应</button>
              <button :class="['seg-btn', { active: appStore.wallpaperType === 'custom' }]" @click="setWallpaperTypeWithSync('custom')">自定义</button>
            </div>
          </div>
          <div v-if="appStore.wallpaperType === 'bing'" class="setting-sub">
            <button class="apply-btn" @click="fetchBingWallpaper">刷新必应壁纸</button>
            <span v-if="appStore.wallpaperUrl" class="bing-preview">
              <img :src="appStore.wallpaperUrl" alt="必应壁纸预览" class="bing-preview-img" />
            </span>
          </div>
          <div v-if="appStore.wallpaperType === 'custom'" class="setting-sub">
            <input v-model="customUrl" class="bg-url-input" placeholder="输入背景图片 URL" />
            <button class="apply-btn" @click="applyCustomBackground">应用</button>
          </div>
          <div v-if="appStore.wallpaperType !== 'default'" class="setting-row">
            <div class="setting-left">
              <div class="setting-icon" style="background: rgba(255,149,0,0.1);">
                <Icon icon="mdi:blur" :size="18" color="#FF9500" />
              </div>
              <div>
                <div class="setting-title">背景模糊</div>
                <div class="setting-desc">调整背景模糊程度</div>
              </div>
            </div>
            <div class="slider-control">
              <input 
                type="range" 
                min="0" 
                max="20" 
                :value="appStore.wallpaperBlur" 
                @input="setWallpaperBlurWithSync(parseInt(($event.target as HTMLInputElement).value))"
                class="range-slider"
              />
              <span class="slider-value">{{ appStore.wallpaperBlur }}px</span>
            </div>
          </div>
          <div v-if="appStore.wallpaperType !== 'default'" class="setting-row">
            <div class="setting-left">
              <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
                <Icon icon="mdi:opacity" :size="18" color="#007AFF" />
              </div>
              <div>
                <div class="setting-title">背景透明度</div>
                <div class="setting-desc">调整背景透明度</div>
              </div>
            </div>
            <div class="slider-control">
              <input 
                type="range" 
                min="10" 
                max="100" 
                :value="appStore.wallpaperOpacity" 
                @input="setWallpaperOpacityWithSync(parseInt(($event.target as HTMLInputElement).value))"
                class="range-slider"
              />
              <span class="slider-value">{{ appStore.wallpaperOpacity }}%</span>
            </div>
          </div>
          <div v-if="appStore.wallpaperType !== 'default'" class="setting-row">
            <div class="setting-left">
              <div class="setting-icon" style="background: rgba(88,86,214,0.1);">
                <Icon icon="mdi:text" :size="18" color="#5856D6" />
              </div>
              <div>
                <div class="setting-title">显示版权信息</div>
                <div class="setting-desc">在壁纸上显示版权信息</div>
              </div>
            </div>
            <n-switch v-model:value="copyrightSwitch" @update:value="setCopyrightWithSync" />
          </div>
        </div>
      </section>

      <!-- Components -->
      <section class="settings-section">
        <div class="section-label">组件</div>
        <div class="settings-card">
          <!-- 搜索栏 -->
          <div class="setting-row">
            <div class="setting-left">
              <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
                <Icon icon="mdi:magnify" :size="18" color="#007AFF" />
              </div>
              <div>
                <div class="setting-title">搜索栏</div>
                <div class="setting-desc">在首页显示搜索栏</div>
              </div>
            </div>
            <n-switch v-model:value="showSearchBar" @update:value="toggleSearchBar" />
          </div>
          <div class="setting-divider"></div>
          <!-- 时钟 -->
          <div class="setting-row">
            <div class="setting-left">
              <div class="setting-icon" style="background: rgba(52,199,89,0.1);">
                <Icon icon="mdi:clock-outline" :size="18" color="#34C759" />
              </div>
              <div>
                <div class="setting-title">时钟组件</div>
                <div class="setting-desc">在首页显示当前时间</div>
              </div>
            </div>
            <n-switch v-model:value="appStore.clockShow" @update:value="appStore.setClockShow" />
          </div>
          <div v-if="appStore.clockShow" class="setting-sub" style="flex-direction: column; align-items: stretch; gap: 8px;">
            <div class="setting-row" style="padding: 8px 0;">
              <span class="mini-label">显示秒数</span>
              <n-switch v-model:value="appStore.clockShowSeconds" @update:value="appStore.setClockShowSeconds" size="small" />
            </div>
            <div class="setting-row" style="padding: 8px 0;">
              <span class="mini-label">时间格式</span>
              <div class="seg-control">
                <button :class="['seg-btn', { active: appStore.clockFormat === '24' }]" @click="appStore.setClockFormat('24')">24小时</button>
                <button :class="['seg-btn', { active: appStore.clockFormat === '12' }]" @click="appStore.setClockFormat('12')">12小时</button>
              </div>
            </div>
          </div>
          <div class="setting-divider"></div>
          <!-- 系统状态 -->
          <div class="setting-row">
            <div class="setting-left">
              <div class="setting-icon" style="background: rgba(175,82,222,0.1);">
                <Icon icon="mdi:server-network" :size="18" color="#AF52DE" />
              </div>
              <div>
                <div class="setting-title">系统状态</div>
                <div class="setting-desc">在首页显示系统信息</div>
              </div>
            </div>
            <n-switch v-model:value="appStore.showSystemStatus" @update:value="appStore.setShowSystemStatus" />
          </div>
        </div>
      </section>

      <!-- Logo -->
      <section class="settings-section">
        <div class="section-label">Logo</div>
        <div class="settings-card">
          <div class="setting-row">
            <div class="setting-left">
              <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
                <Icon icon="mdi:image" :size="18" color="#007AFF" />
              </div>
              <div>
                <div class="setting-title">Logo 类型</div>
                <div class="setting-desc">选择图片或文字显示</div>
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
        </div>
      </section>

      <!-- Card Style -->
      <section class="settings-section">
        <div class="section-label">卡片风格</div>
        <div class="settings-card">
          <div class="setting-row">
            <div class="setting-left">
              <div class="setting-icon" style="background: rgba(255,149,0,0.1);">
                <Icon icon="mdi:card-multiple" :size="18" color="#FF9500" />
              </div>
              <div>
                <div class="setting-title">卡片样式</div>
                <div class="setting-desc">选择卡片的外观风格</div>
              </div>
            </div>
            <div class="seg-control">
              <button :class="['seg-btn', { active: appStore.cardStyle === 'default' }]" @click="appStore.setCardStyle('default')">默认</button>
              <button :class="['seg-btn', { active: appStore.cardStyle === 'round' }]" @click="appStore.setCardStyle('round')">圆形</button>
              <button :class="['seg-btn', { active: appStore.cardStyle === 'square' }]" @click="appStore.setCardStyle('square')">方形</button>
            </div>
          </div>
        </div>
      </section>

      <!-- Layout -->
      <section class="settings-section">
        <div class="section-label">布局</div>
        <div class="settings-card">
          <div class="setting-row">
            <div class="setting-left">
              <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
                <Icon icon="mdi:arrow-expand-horizontal" :size="18" color="#007AFF" />
              </div>
              <div>
                <div class="setting-title">内容最大宽度</div>
                <div class="setting-desc">调整内容区域宽度</div>
              </div>
            </div>
            <div class="slider-control">
              <input type="range" min="800" max="1400" step="50" :value="appStore.contentMaxWidth" @input="appStore.setContentMaxWidth(String(parseInt(($event.target as HTMLInputElement).value)))" class="range-slider" />
              <span class="slider-value">{{ appStore.contentMaxWidth }}px</span>
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
                <div class="setting-desc">调整页面左右边距</div>
              </div>
            </div>
            <div class="slider-control">
              <input type="range" min="16" max="64" :value="appStore.contentPaddingX" @input="appStore.setContentPaddingX(String(parseInt(($event.target as HTMLInputElement).value)))" class="range-slider" />
              <span class="slider-value">{{ appStore.contentPaddingX }}px</span>
            </div>
          </div>
          <div class="setting-divider"></div>
          <div class="setting-row">
            <div class="setting-left">
              <div class="setting-icon" style="background: rgba(175,82,222,0.1);">
                <Icon icon="mdi:format-vertical-align-top" :size="18" color="#AF52DE" />
              </div>
              <div>
                <div class="setting-title">顶部/底部边距</div>
                <div class="setting-desc">调整内容上下边距</div>
              </div>
            </div>
            <div class="layout-padding">
              <div class="padding-item">
                <span class="mini-label">上</span>
                <input type="number" min="16" max="64" :value="appStore.contentPaddingTop" @change="appStore.setContentPaddingTop(String(parseInt(($event.target as HTMLInputElement).value)))" class="num-input" />
              </div>
              <div class="padding-item">
                <span class="mini-label">下</span>
                <input type="number" min="16" max="128" :value="appStore.contentPaddingBottom" @change="appStore.setContentPaddingBottom(String(parseInt(($event.target as HTMLInputElement).value)))" class="num-input" />
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- Colors -->
      <section class="settings-section">
        <div class="section-label">颜色</div>
        <div class="settings-card">
          <div class="setting-row">
            <div class="setting-left">
              <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
                <Icon icon="mdi:palette" :size="18" color="#007AFF" />
              </div>
              <div>
                <div class="setting-title">主题色</div>
                <div class="setting-desc">自定义主题强调色</div>
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
                <div class="setting-desc">自定义卡片边框色</div>
              </div>
            </div>
            <div class="color-picker-wrap">
              <input type="color" :value="appStore.borderColor || '#000000'" @input="appStore.setBorderColor(($event.target as HTMLInputElement).value)" class="color-input" />
              <button v-if="appStore.borderColor" class="reset-btn" @click="appStore.setBorderColor('')">重置</button>
            </div>
          </div>
        </div>
      </section>

      <!-- Site -->
      <section class="settings-section">
        <div class="section-label">站点</div>
        <div class="settings-card">
          <div class="setting-row">
            <div class="setting-left">
              <div class="setting-icon" style="background: rgba(88,86,214,0.1);">
                <Icon icon="mdi:format-title" :size="18" color="#5856D6" />
              </div>
              <div>
                <div class="setting-title">站点标题</div>
                <div class="setting-desc">自定义浏览器标签页标题</div>
              </div>
            </div>
            <input v-model="siteTitleInput" class="site-input" placeholder="SunDash" @change="updateSiteTitle" />
          </div>
          <div class="setting-divider"></div>
          <div class="setting-row">
            <div class="setting-left">
              <div class="setting-icon" style="background: rgba(52,199,89,0.1);">
                <Icon icon="mdi:login" :size="18" color="#34C759" />
              </div>
              <div>
                <div class="setting-title">登录验证码</div>
                <div class="setting-desc">登录时需要输入验证码</div>
              </div>
            </div>
            <n-switch v-model:value="appStore.enableCaptcha" @update:value="appStore.setEnableCaptcha" />
          </div>
        </div>
      </section>

      <!-- Footer -->
      <section class="settings-section">
        <div class="section-label">页脚</div>
        <div class="settings-card">
          <div class="setting-row" style="flex-direction: column; align-items: stretch; gap: 12px;">
            <div class="setting-left">
              <div class="setting-icon" style="background: rgba(175,82,222,0.1);">
                <Icon icon="mdi:page-layout-footer" :size="18" color="#AF52DE" />
              </div>
              <div>
                <div class="setting-title">自定义页脚</div>
                <div class="setting-desc">在首页底部显示自定义内容（支持纯文本）</div>
              </div>
            </div>
            <textarea v-model="footerInput" class="footer-input" placeholder="输入页脚内容..." rows="3" @change="appStore.setFooterHtml(footerInput)" />
          </div>
        </div>
      </section>

      <!-- Data -->
      <section class="settings-section">
        <div class="section-label">数据</div>
        <div class="settings-card">
          <div class="setting-row">
            <div class="setting-left">
              <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
                <Icon icon="mdi:export-variant" :size="18" color="#007AFF" />
              </div>
              <div>
                <div class="setting-title">导出数据</div>
                <div class="setting-desc">导出面板配置为 JSON</div>
              </div>
            </div>
            <button class="do-btn" @click="handleExport">导出</button>
          </div>
          <div class="setting-divider"></div>
          <div class="setting-row">
            <div class="setting-left">
              <div class="setting-icon" style="background: rgba(52,199,89,0.1);">
                <Icon icon="mdi:import" :size="18" color="#34C759" />
              </div>
              <div>
                <div class="setting-title">导入数据</div>
                <div class="setting-desc">从 JSON 文件导入配置</div>
              </div>
            </div>
            <div style="display: flex; gap: 8px;">
              <button class="do-btn" @click="importTemplate">模板</button>
              <button class="do-btn" @click="triggerImport">导入</button>
            </div>
          </div>
          <input ref="fileInput" type="file" accept=".json" style="display: none;" @change="handleImport" />
        </div>
      </section>

      <!-- Danger Zone -->
      <section class="settings-section">
        <div class="section-label">危险操作</div>
        <div class="settings-card">
          <div class="setting-row clickable" @click="handleResetSettings">
            <div class="setting-left">
              <div class="setting-icon" style="background: rgba(255,149,0,0.1);">
                <Icon icon="mdi:restart" :size="18" color="#FF9500" />
              </div>
              <div>
                <div class="setting-title" style="color: #FF9500;">重置所有设置</div>
                <div class="setting-desc">恢复所有个性化设置为默认值</div>
              </div>
            </div>
            <Icon icon="mdi:chevron-right" :size="16" color="var(--sd-text-tertiary)" />
          </div>
        </div>
      </section>

      <!-- Account -->
      <section class="settings-section">
        <div class="section-label">账户</div>
        <div class="settings-card">
          <div class="setting-row">
            <div class="setting-left">
              <div class="setting-icon" style="background: rgba(88,86,214,0.1);">
                <Icon icon="mdi:account-circle-outline" :size="18" color="#5856D6" />
              </div>
              <div>
                <div class="setting-title">当前用户</div>
                <div class="setting-desc">{{ userStore.user?.display_name || userStore.user?.username }} · {{ userStore.user?.role === 'admin' ? '管理员' : '用户' }}</div>
              </div>
            </div>
          </div>
          <div class="setting-divider"></div>
          <div class="setting-row clickable" @click="userStore.logout()">
            <div class="setting-left">
              <div class="setting-icon" style="background: rgba(255,59,48,0.1);">
                <Icon icon="mdi:logout" :size="18" color="#FF3B30" />
              </div>
              <div>
                <div class="setting-title" style="color: #FF3B30;">退出登录</div>
              </div>
            </div>
            <Icon icon="mdi:chevron-right" :size="16" color="var(--sd-text-tertiary)" />
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Icon } from '@iconify/vue'
import { useMessage, useDialog, NSwitch } from 'naive-ui'
import { useAppStore } from '../stores/app'
import { useUserStore } from '../stores/user'
import { usePanelStore } from '../stores/panel'
import { importTemplates } from '../data/importTemplates'

const appStore = useAppStore()
const userStore = useUserStore()
const panelStore = usePanelStore()
const message = useMessage()
const dialog = useDialog()

const fileInput = ref<HTMLInputElement | null>(null)
const customUrl = ref(appStore.wallpaperUrl)
const copyrightSwitch = ref(appStore.showWallpaperCopyright)
const logoImageInput = ref(appStore.logoImageUrl)
const logoTextInput = ref(appStore.logoText)
const siteTitleInput = ref(appStore.siteTitle)
const footerInput = ref(appStore.footerHtml)

// Load settings from backend on mount
onMounted(async () => {
  await panelStore.fetchSettings()
  const savedType = panelStore.getSetting('wallpaper_type')
  const savedUrl = panelStore.getSetting('wallpaper_url')
  const savedBlur = panelStore.getSetting('wallpaper_blur')
  const savedOpacity = panelStore.getSetting('wallpaper_opacity')
  const savedCopyright = panelStore.getSetting('wallpaper_copyright')
  const savedSearch = panelStore.getSetting('show_searchbar')
  const savedClock = panelStore.getSetting('show_clock')
  const savedTitle = panelStore.getSetting('site_title')

  if (savedType) appStore.setWallpaperType(savedType as any)
  if (savedUrl) appStore.setWallpaperUrl(savedUrl)
  if (savedBlur) appStore.setWallpaperBlur(parseInt(savedBlur))
  if (savedOpacity) appStore.setWallpaperOpacity(parseInt(savedOpacity))
  if (savedCopyright) appStore.setShowWallpaperCopyright(savedCopyright === 'true')
  if (savedTitle) {
    siteTitleInput.value = savedTitle
    appStore.setSiteTitle(savedTitle)
  }
  showSearchBar.value = savedSearch !== 'false'
  if (savedClock) appStore.setClockShow(savedClock === 'true')
})

async function fetchBingWallpaper() {
  try {
    const { api } = await import('../api')
    const res = await api.get('wallpaper/bing')
    if (res.data && res.data.images && res.data.images[0]) {
      appStore.setWallpaperUrl(res.data.images[0].url)
      message.success('必应壁纸已刷新')
    }
  } catch (e) {
    message.error('获取必应壁纸失败')
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

// Component toggles
const showSearchBar = ref(localStorage.getItem('sundash-show-searchbar') !== 'false')

function toggleSearchBar(val: boolean) {
  localStorage.setItem('sundash-show-searchbar', val.toString())
}

function updateSiteTitle() {
  const title = siteTitleInput.value.trim() || 'SunDash'
  appStore.setSiteTitle(title)
  message.success('站点标题已更新')
}

function setNetwork(mode: 'internal' | 'external') {
  appStore.networkMode = mode
  localStorage.setItem('sundash-network', mode)
}

function handleResetSettings() {
  dialog.warning({
    title: '重置设置',
    content: '即将重置所有个性化设置为默认值。\n\n此操作不会删除你的分组和书签数据。',
    positiveText: '确认重置',
    negativeText: '取消',
    onPositiveClick: () => {
      // Clear all sundash localStorage keys
      const keys = Object.keys(localStorage).filter(k => k.startsWith('sundash-'))
      keys.forEach(k => localStorage.removeItem(k))
      // Reload to apply defaults
      window.location.reload()
    },
  })
}

// Save wallpaper settings to backend
async function saveWallpaperSetting(key: string, value: string) {
  try {
    await panelStore.updateSetting(key, value)
  } catch (e) {
    console.error('Failed to save setting to backend:', e)
  }
}

// Override appStore methods to also persist to backend
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

function setCopyrightWithSync(show: boolean) {
  appStore.setShowWallpaperCopyright(show)
  saveWallpaperSetting('wallpaper_copyright', show.toString())
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

// Import a built-in template (common navigation) with confirmation.
function importTemplate() {
  const t = importTemplates[0]
  const cardCount = t.groups.reduce((sum: number, g: any) => sum + (g.cards?.length || 0), 0)
  dialog.warning({
    title: '确认导入模板',
    content: `即将导入「${t.name}」：${t.groups.length} 个分组、${cardCount} 个书签。\n\n现有数据不会被删除，新数据将追加到当前面板。`,
    positiveText: '确认导入',
    negativeText: '取消',
    onPositiveClick: async () => {
      await doImport({ groups: t.groups })
    },
  })
}

async function handleImport(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  
  try {
    const text = await file.text()
    const data = JSON.parse(text)
    if (!data.groups || !Array.isArray(data.groups)) {
      message.error('无效的备份文件格式')
      return
    }
    
    const groupCount = data.groups.length
    const cardCount = data.groups.reduce((sum: number, g: any) => sum + (g.cards?.length || 0), 0)
    
    dialog.warning({
      title: '确认导入',
      content: `即将导入 ${groupCount} 个分组、${cardCount} 个书签。\n\n现有数据不会被删除，新数据将追加到当前面板。`,
      positiveText: '确认导入',
      negativeText: '取消',
      onPositiveClick: async () => {
        await doImport(data)
      },
    })
  } catch {
    message.error('导入失败：文件格式错误')
  } finally {
    input.value = ''
  }
}

async function doImport(data: any) {
  try {
    // Restore settings if present
    if (data.settings && typeof data.settings === 'object') {
      for (const [key, value] of Object.entries(data.settings)) {
        if (key.startsWith('sundash-') && typeof value === 'string') {
          localStorage.setItem(key, value)
        }
      }
      // Reload settings from server to apply
      await appStore.loadSettingsFromServer()
      message.success('设置已恢复')
    }
    // Import groups and cards
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
  } catch {
    message.error('导入失败')
  }
}


</script>

<style scoped>
.settings-page {
  min-height: 100vh;
  background: var(--sd-bg-root);
  position: relative;
}

.page-bg {
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 0;
  overflow: hidden;
}

.bg-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
}

.bg-orb-1 {
  width: 500px;
  height: 500px;
  background: radial-gradient(circle, rgba(0,122,255,0.07) 0%, transparent 70%);
  top: -150px;
  right: -100px;
  animation: floatOrb 20s ease-in-out infinite;
}

.bg-orb-2 {
  width: 400px;
  height: 400px;
  background: radial-gradient(circle, rgba(88,86,214,0.05) 0%, transparent 70%);
  bottom: -100px;
  left: -80px;
  animation: floatOrb 25s ease-in-out infinite reverse;
}

@keyframes floatOrb {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(20px, -20px); }
}

/* Header */
.page-header {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: var(--sd-header-height);
  padding: 0 var(--sd-space-6);
  background: rgba(255,255,255,0.8);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border-bottom: 1px solid rgba(0,0,0,0.06);
  position: sticky;
  top: 0;
}

:root[data-theme="dark"] .page-header {
  background: rgba(0,0,0,0.7);
  border-bottom-color: rgba(255,255,255,0.08);
}

.back-btn {
  display: flex;
  align-items: center;
  gap: 2px;
  background: none;
  border: none;
  color: var(--sd-primary);
  font-size: 14px;
  font-family: var(--sd-font);
  font-weight: 500;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 8px;
  transition: background 0.15s ease;
}

.back-btn:hover {
  background: rgba(0,122,255,0.08);
}

.page-header h1 {
  font-size: 17px;
  font-weight: 600;
  margin: 0;
  color: var(--sd-text-primary);
}

/* Content */
.settings-content {
  position: relative;
  z-index: 1;
  max-width: 600px;
  margin: 0 auto;
  padding: var(--sd-space-6) var(--sd-space-6) var(--sd-space-16);
  display: flex;
  flex-direction: column;
  gap: var(--sd-space-6);
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
  -webkit-backdrop-filter: blur(20px) saturate(180%);
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
  padding: 14px 16px;
  min-height: 56px;
}

.setting-row.clickable {
  cursor: pointer;
  transition: background 0.15s ease;
}

.setting-row.clickable:hover {
  background: rgba(0,0,0,0.03);
}

:root[data-theme="dark"] .setting-row.clickable:hover {
  background: rgba(255,255,255,0.05);
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
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  flex-shrink: 0;
}

.setting-title {
  font-size: 15px;
  font-weight: 500;
  color: var(--sd-text-primary);
}

.setting-desc {
  font-size: 12px;
  color: var(--sd-text-tertiary);
  margin-top: 1px;
}

.setting-divider {
  height: 1px;
  background: rgba(0,0,0,0.06);
  margin: 0 16px;
}

:root[data-theme="dark"] .setting-divider {
  background: rgba(255,255,255,0.06);
}

.setting-sub {
  padding: 0 16px 14px;
  display: flex;
  gap: 8px;
  align-items: center;
}

.bg-url-input {
  flex: 1;
  height: 34px;
  padding: 0 12px;
  background: rgba(0,0,0,0.04);
  border: 1px solid rgba(0,0,0,0.08);
  border-radius: 8px;
  font-size: 13px;
  font-family: var(--sd-font);
  color: var(--sd-text-primary);
  outline: none;
}

.bg-url-input:focus {
  border-color: var(--sd-primary);
  box-shadow: 0 0 0 3px rgba(0,122,255,0.1);
}

/* Segment Control */
.seg-control {
  display: flex;
  background: rgba(0,0,0,0.04);
  border-radius: 8px;
  padding: 2px;
  flex-shrink: 0;
}

:root[data-theme="dark"] .seg-control {
  background: rgba(255,255,255,0.06);
}

.seg-btn {
  padding: 5px 14px;
  background: transparent;
  border: none;
  border-radius: 6px;
  font-size: 13px;
  font-family: var(--sd-font);
  font-weight: 500;
  color: var(--sd-text-secondary);
  cursor: pointer;
  transition: all 0.2s ease;
  white-space: nowrap;
}

.seg-btn.active {
  background: white;
  color: var(--sd-text-primary);
  box-shadow: 0 1px 3px rgba(0,0,0,0.1);
}

:root[data-theme="dark"] .seg-btn.active {
  background: rgba(56,56,58,0.9);
}

.do-btn {
  padding: 6px 16px;
  background: var(--sd-primary);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 13px;
  font-family: var(--sd-font);
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease;
  flex-shrink: 0;
}

.do-btn:hover {
  background: var(--sd-primary-hover);
}

.apply-btn {
  padding: 6px 14px;
  background: var(--sd-primary);
  color: white;
  border: none;
  border-radius: 8px;
  font-size: 13px;
  font-family: var(--sd-font);
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease;
}

.apply-btn:hover {
  background: var(--sd-primary-hover);
}

.bing-preview {
  display: inline-block;
  width: 120px;
  height: 68px;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid rgba(0,0,0,0.08);
  flex-shrink: 0;
}

.bing-preview-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.slider-control {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
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
  appearance: none;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: var(--sd-primary);
  cursor: pointer;
}

.slider-value {
  font-size: 13px;
  color: var(--sd-text-secondary);
  min-width: 36px;
  text-align: right;
}

.site-input {
  width: 160px;
  height: 34px;
  padding: 0 12px;
  background: rgba(0,0,0,0.04);
  border: 1px solid rgba(0,0,0,0.08);
  border-radius: 8px;
  font-size: 13px;
  font-family: var(--sd-font);
  color: var(--sd-text-primary);
  outline: none;
  flex-shrink: 0;
}

.site-input:focus {
  border-color: var(--sd-primary);
  box-shadow: 0 0 0 3px rgba(0,122,255,0.1);
}

.mini-label {
  font-size: 13px;
  color: var(--sd-text-secondary);
}

.layout-padding {
  display: flex;
  gap: 12px;
  flex-shrink: 0;
}

.padding-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.num-input {
  width: 60px;
  height: 30px;
  padding: 0 8px;
  background: rgba(0,0,0,0.04);
  border: 1px solid rgba(0,0,0,0.08);
  border-radius: 6px;
  font-size: 13px;
  text-align: center;
}

.color-picker-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.color-input {
  width: 36px;
  height: 36px;
  padding: 0;
  border: 2px solid rgba(0,0,0,0.08);
  border-radius: 8px;
  cursor: pointer;
  background: none;
}

.reset-btn {
  padding: 4px 10px;
  background: rgba(0,0,0,0.05);
  border: none;
  border-radius: 6px;
  font-size: 12px;
  color: var(--sd-text-secondary);
  cursor: pointer;
}

.footer-input {
  width: 100%;
  min-height: 80px;
  padding: 12px;
  background: rgba(0,0,0,0.04);
  border: 1px solid rgba(0,0,0,0.08);
  border-radius: 10px;
  font-size: 13px;
  font-family: var(--sd-font);
  color: var(--sd-text-primary);
  resize: vertical;
  outline: none;
}

.footer-input:focus {
  border-color: var(--sd-primary);
  box-shadow: 0 0 0 3px rgba(0,122,255,0.1);
}
</style>