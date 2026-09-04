import { defineStore } from 'pinia'
import { ref, watch } from 'vue'
import { usePreferredDark } from '@vueuse/core'
import { api } from '../api'

// Helper to migrate old asuan- keys to sundash-
function migrateKeys() {
  const keys = [
    'theme', 'network', 'wallpaper-type', 'wallpaper-url', 'wallpaper-blur',
    'wallpaper-opacity', 'wallpaper-copyright', 'clock-show', 'clock-seconds',
    'clock-format', 'logo-type', 'logo-image', 'logo-text', 'card-label-size',
    'card-item-size', 'cards-per-row', 'content-max-width', 'content-padding-x',
    'content-padding-top', 'content-padding-bottom', 'bg-color', 'text-color',
    'border-color', 'primary-color', 'footer-html', 'site-title', 'site-icon',
    'login-bg', 'enable-captcha', 'show-system-status', 'network-auto',
    'group-card-transparent', 'search-engine',
  ]
  for (const k of keys) {
    const old = localStorage.getItem(`asuan-${k}`)
    if (old !== null && localStorage.getItem(`sundash-${k}`) === null) {
      localStorage.setItem(`sundash-${k}`, old)
    }
    localStorage.removeItem(`asuan-${k}`)
  }
  // Also migrate token
  const oldToken = localStorage.getItem('asuan-token')
  if (oldToken && !localStorage.getItem('sundash-token')) {
    localStorage.setItem('sundash-token', oldToken)
  }
  localStorage.removeItem('asuan-token')
}

// Run migration on import
migrateKeys()

// localStorage key helper
const L = (key: string) => `sundash-${key}`

export const useAppStore = defineStore('app', () => {
  const prefersDark = usePreferredDark()
  const themeMode = ref<'light' | 'dark' | 'system'>(
    (localStorage.getItem(L('theme')) as 'light' | 'dark' | 'system') || 'system'
  )
  const networkMode = ref<'internal' | 'external'>(
    (localStorage.getItem(L('network')) as 'internal' | 'external') || 'internal'
  )
  const sidebarCollapsed = ref(false)

  // Wallpaper settings
  const wallpaperType = ref<'default' | 'gradient' | 'bing' | 'custom'>(
    (localStorage.getItem(L('wallpaper-type')) as 'default' | 'gradient' | 'bing' | 'custom') || 'default'
  )
  const wallpaperUrl = ref(localStorage.getItem(L('wallpaper-url')) || '')
  const wallpaperBlur = ref(parseInt(localStorage.getItem(L('wallpaper-blur')) || '0'))
  const wallpaperOpacity = ref(parseInt(localStorage.getItem(L('wallpaper-opacity')) || '100'))
  const showWallpaperCopyright = ref(localStorage.getItem(L('wallpaper-copyright')) === 'true')

  const isDark = ref(false)

  function updateDark() {
    isDark.value = themeMode.value === 'system' ? prefersDark.value : themeMode.value === 'dark'
  }

  watch(themeMode, () => {
    localStorage.setItem(L('theme'), themeMode.value)
    updateDark()
  }, { immediate: true })
  watch(prefersDark, updateDark)

  function setTheme(mode: 'light' | 'dark' | 'system') { themeMode.value = mode }
  function toggleNetwork() {
    networkMode.value = networkMode.value === 'internal' ? 'external' : 'internal'
    localStorage.setItem(L('network'), networkMode.value)
  }
  function toggleSidebar() { sidebarCollapsed.value = !sidebarCollapsed.value }

  function setWallpaperType(type: 'default' | 'gradient' | 'bing' | 'custom') { wallpaperType.value = type; localStorage.setItem(L('wallpaper-type'), type) }
  function setWallpaperUrl(url: string) { wallpaperUrl.value = url; localStorage.setItem(L('wallpaper-url'), url) }
  function setWallpaperBlur(blur: number) { wallpaperBlur.value = blur; localStorage.setItem(L('wallpaper-blur'), blur.toString()) }
  function setWallpaperOpacity(opacity: number) { wallpaperOpacity.value = opacity; localStorage.setItem(L('wallpaper-opacity'), opacity.toString()) }
  function setShowWallpaperCopyright(show: boolean) { showWallpaperCopyright.value = show; localStorage.setItem(L('wallpaper-copyright'), show.toString()) }

  // Clock settings
  const clockShow = ref(localStorage.getItem(L('clock-show')) !== 'false')
  const clockShowSeconds = ref(localStorage.getItem(L('clock-seconds')) === 'true')
  const clockFormat = ref<'12' | '24'>((localStorage.getItem(L('clock-format')) as '12' | '24') || '24')

  // Logo settings
  const logoType = ref<'image' | 'text'>((localStorage.getItem(L('logo-type')) as 'image' | 'text') || 'text')
  const logoImageUrl = ref(localStorage.getItem(L('logo-image')) || '')
  const logoText = ref(localStorage.getItem(L('logo-text')) || 'SunDash')

  // Card settings
  const cardLabelSize = ref(localStorage.getItem(L('card-label-size')) || '12')
  const cardItemSize = ref(localStorage.getItem(L('card-item-size')) || '0')
  const cardsPerRow = ref(localStorage.getItem(L('cards-per-row')) || '5')

  // Layout settings
  const migrateLayoutKey = (key: string, fallback: string) => {
    const v = localStorage.getItem(key)
    if (v && !v.includes('%') && !v.includes('calc')) {
      localStorage.removeItem(key)
    }
    return localStorage.getItem(key) || fallback
  }
  const contentMaxWidth = ref(migrateLayoutKey(L('content-max-width'), '80%'))
  const contentPaddingX = ref(migrateLayoutKey(L('content-padding-x'), '5%'))
  const contentPaddingTop = ref(migrateLayoutKey(L('content-padding-top'), '10%'))
  const contentPaddingBottom = ref(migrateLayoutKey(L('content-padding-bottom'), '5%'))

  // Color customization
  const bgColor = ref(localStorage.getItem(L('bg-color')) || '')
  const textColor = ref(localStorage.getItem(L('text-color')) || '')
  const borderColor = ref(localStorage.getItem(L('border-color')) || '')
  const primaryColor = ref(localStorage.getItem(L('primary-color')) || '')

  // Footer
  const footerHtml = ref(localStorage.getItem(L('footer-html')) || '')

  // Site settings
  const siteTitle = ref(localStorage.getItem(L('site-title')) || 'SunDash')
  const siteIconUrl = ref(localStorage.getItem(L('site-icon')) || '')
  const loginBgUrl = ref(localStorage.getItem(L('login-bg')) || '')
  const enableCaptcha = ref(localStorage.getItem(L('enable-captcha')) === 'true')

  // Component settings
  const showSystemStatus = ref(localStorage.getItem(L('show-system-status')) === 'true')
  const showSystemMonitor = ref(localStorage.getItem(L('show-system-monitor')) !== 'false')
  const showWeatherWidget = ref(localStorage.getItem(L('show-weather-widget')) !== 'false')
  const showMemo = ref(localStorage.getItem(L('show-memo')) !== 'false')
  const showRSSWidget = ref(localStorage.getItem(L('show-rss-widget')) === 'true')
  // 书签同步总开关：默认开启；site-config 下发 bmsync_enabled=false 时隐藏主页入口
  const bmsyncEnabled = ref(true)
  const networkModeAuto = ref(localStorage.getItem(L('network-auto')) === 'true')
  const groupCardTransparent = ref(localStorage.getItem(L('group-card-transparent')) === 'true')
  const searchEngine = ref(localStorage.getItem(L('search-engine')) || 'Baidu')
  const cardStyle = ref<'default' | 'round' | 'square'>((localStorage.getItem(L('card-style')) as 'default' | 'round' | 'square') || 'default')

  // Auth settings (from server)
  const allowRegistration = ref(true)
  const requireApproval = ref(false)

  // Helper: localStorage + server sync
  async function setLocal(key: string, value: string) {
    localStorage.setItem(key, value)
    try {
      await api.put('settings', { key, value })
    } catch (e) {
      console.error('Failed to sync setting to server:', key, e)
    }
  }

  // Setter functions
  function setClockShow(show: boolean) { clockShow.value = show; setLocal(L('clock-show'), show.toString()) }
  function setClockShowSeconds(show: boolean) { clockShowSeconds.value = show; setLocal(L('clock-seconds'), show.toString()) }
  function setClockFormat(format: '12' | '24') { clockFormat.value = format; setLocal(L('clock-format'), format) }
  function setLogoType(type: 'image' | 'text') { logoType.value = type; setLocal(L('logo-type'), type) }
  function setLogoImageUrl(url: string) { logoImageUrl.value = url; setLocal(L('logo-image'), url) }
  function setLogoText(text: string) { logoText.value = text; setLocal(L('logo-text'), text) }
  function setCardLabelSize(size: string) { cardLabelSize.value = size; setLocal(L('card-label-size'), size) }
  function setCardItemSize(size: string) { cardItemSize.value = size; setLocal(L('card-item-size'), size) }
  function setCardsPerRow(count: string) { cardsPerRow.value = count; setLocal(L('cards-per-row'), count) }
  function setContentMaxWidth(w: string) { contentMaxWidth.value = w; setLocal(L('content-max-width'), w) }
  function setContentPaddingX(p: string) { contentPaddingX.value = p; setLocal(L('content-padding-x'), p) }
  function setContentPaddingTop(p: string) { contentPaddingTop.value = p; setLocal(L('content-padding-top'), p) }
  function setContentPaddingBottom(p: string) { contentPaddingBottom.value = p; setLocal(L('content-padding-bottom'), p) }
  function setBgColor(color: string) { bgColor.value = color; setLocal(L('bg-color'), color) }
  function setTextColor(color: string) { textColor.value = color; setLocal(L('text-color'), color) }
  function setBorderColor(color: string) { borderColor.value = color; setLocal(L('border-color'), color) }
  function setPrimaryColor(color: string) { primaryColor.value = color; setLocal(L('primary-color'), color) }
  function setFooterHtml(html: string) { footerHtml.value = html; setLocal(L('footer-html'), html) }
  function setSiteTitle(title: string) { siteTitle.value = title; setLocal(L('site-title'), title); document.title = title }
  function setSiteIconUrl(url: string) { siteIconUrl.value = url; setLocal(L('site-icon'), url) }
  function setLoginBgUrl(url: string) { loginBgUrl.value = url; setLocal(L('login-bg'), url) }
  function setEnableCaptcha(enable: boolean) { enableCaptcha.value = enable; setLocal(L('enable-captcha'), enable.toString()) }
  function setShowSystemStatus(show: boolean) { showSystemStatus.value = show; setLocal(L('show-system-status'), show.toString()) }
  function setShowSystemMonitor(show: boolean) { showSystemMonitor.value = show; setLocal(L('show-system-monitor'), show.toString()) }
  function setShowWeatherWidget(show: boolean) { showWeatherWidget.value = show; setLocal(L('show-weather-widget'), show.toString()) }
  function setShowMemo(show: boolean) { showMemo.value = show; setLocal(L('show-memo'), show.toString()) }
  function setShowRSSWidget(show: boolean) { showRSSWidget.value = show; setLocal(L('show-rss-widget'), show.toString()) }
  function setNetworkModeAuto(auto: boolean) { networkModeAuto.value = auto; setLocal(L('network-auto'), auto.toString()) }
  function setGroupCardTransparent(transparent: boolean) { groupCardTransparent.value = transparent; setLocal(L('group-card-transparent'), transparent.toString()) }
  function setSearchEngine(engine: string) { searchEngine.value = engine; setLocal(L('search-engine'), engine) }
  function setCardStyle(style: 'default' | 'round' | 'square') { cardStyle.value = style; setLocal(L('card-style'), style) }

  // Apply server settings to reactive state + localStorage. Shared by
  // loadSettingsFromServer and the bootstrap endpoint.
  function applyServerSettings(serverData: Record<string, string>) {
    const s = serverData || {}

    // Mapping: server key → [ref, localStorage key, transformer]
    const apply = (
      serverKey: string,
      target: { value: any },
      localKey: string,
      transform?: (v: string) => any,
    ) => {
      const val = s[serverKey]
      if (val !== undefined) {
        target.value = transform ? transform(val) : val
        localStorage.setItem(localKey, val)
      }
    }

    apply('sundash-theme', themeMode, L('theme'))
    apply('sundash-network', networkMode, L('network'))
    apply('sundash-wallpaper-type', wallpaperType, L('wallpaper-type'))
    apply('sundash-wallpaper-url', wallpaperUrl, L('wallpaper-url'))
    apply('sundash-wallpaper-blur', wallpaperBlur, L('wallpaper-blur'), (v) => parseInt(v))
    apply('sundash-wallpaper-opacity', wallpaperOpacity, L('wallpaper-opacity'), (v) => parseInt(v))
    apply('sundash-wallpaper-copyright', showWallpaperCopyright, L('wallpaper-copyright'), (v) => v === 'true')
    apply('sundash-clock-show', clockShow, L('clock-show'), (v) => v !== 'false')
    apply('sundash-clock-seconds', clockShowSeconds, L('clock-seconds'), (v) => v === 'true')
    apply('sundash-clock-format', clockFormat, L('clock-format'))
    apply('sundash-logo-type', logoType, L('logo-type'))
    apply('sundash-logo-image', logoImageUrl, L('logo-image'))
    apply('sundash-logo-text', logoText, L('logo-text'))
    apply('sundash-card-label-size', cardLabelSize, L('card-label-size'))
    apply('sundash-card-item-size', cardItemSize, L('card-item-size'))
    apply('sundash-cards-per-row', cardsPerRow, L('cards-per-row'))
    apply('sundash-content-max-width', contentMaxWidth, L('content-max-width'))
    apply('sundash-content-padding-x', contentPaddingX, L('content-padding-x'))
    apply('sundash-content-padding-top', contentPaddingTop, L('content-padding-top'))
    apply('sundash-content-padding-bottom', contentPaddingBottom, L('content-padding-bottom'))
    apply('sundash-bg-color', bgColor, L('bg-color'))
    apply('sundash-text-color', textColor, L('text-color'))
    apply('sundash-border-color', borderColor, L('border-color'))
    apply('sundash-primary-color', primaryColor, L('primary-color'))
    apply('sundash-footer-html', footerHtml, L('footer-html'))
    apply('sundash-login-bg', loginBgUrl, L('login-bg'))
    apply('sundash-enable-captcha', enableCaptcha, L('enable-captcha'), (v) => v === 'true')
    apply('sundash-show-system-status', showSystemStatus, L('show-system-status'), (v) => v === 'true')
    apply('sundash-show-system-monitor', showSystemMonitor, L('show-system-monitor'), (v) => v === 'true')
    apply('sundash-show-weather-widget', showWeatherWidget, L('show-weather-widget'), (v) => v === 'true')
    apply('sundash-show-memo', showMemo, L('show-memo'), (v) => v === 'true')
    apply('sundash-show-rss-widget', showRSSWidget, L('show-rss-widget'), (v) => v === 'true')
    apply('sundash-network-auto', networkModeAuto, L('network-auto'), (v) => v === 'true')
    apply('sundash-group-card-transparent', groupCardTransparent, L('group-card-transparent'), (v) => v === 'true')
    apply('sundash-search-engine', searchEngine, L('search-engine'))

    // Site title: also update document
    if (s['sundash-site-title'] !== undefined) {
      siteTitle.value = s['sundash-site-title']
      localStorage.setItem(L('site-title'), s['sundash-site-title'])
      document.title = s['sundash-site-title']
    }
    apply('sundash-site-icon', siteIconUrl, L('site-icon'))
  }

  async function loadSettingsFromServer() {
    try {
      const res = await api.get<Record<string, string>>('settings')
      applyServerSettings(res.data || {})
    } catch (e) {
      console.error('Failed to load settings from server:', e)
    }
  }

  // Fetch auth settings (registration allowed, approval required)
  async function fetchAuthSettings() {
    try {
      const res = await api.get<{ allow_registration: boolean; require_approval: boolean }>('auth/settings')
      allowRegistration.value = res.data.allow_registration
      requireApproval.value = res.data.require_approval
    } catch (e) {
      console.error('Failed to load auth settings:', e)
    }
  }

  // Fetch and apply site config (title, favicon, etc.)
  async function fetchSiteConfig() {
    try {
      const res = await api.get<Record<string, string>>('site-config')
      const config = res.data || {}
      if (config.site_title) {
        // 同步 store + localStorage + document.title（而不是只改 document.title）
        setSiteTitle(config.site_title)
        // 顶栏 Logo 文字若仍是默认值，跟随站点名（用户后续可到设置里自定义覆盖）
        if (logoText.value === 'SunDash') {
          setLogoText(config.site_title)
        }
      }
      if (config.site_icon_url) {
        const link = document.querySelector("link[rel='icon']") as HTMLLinkElement
        if (link) link.href = config.site_icon_url
      }
      // 书签同步总开关（site-config 下发 bmsync_enabled；缺失默认开启）
      if (config.bmsync_enabled !== undefined) {
        bmsyncEnabled.value = config.bmsync_enabled !== 'false'
      }
    } catch (e) {
      console.error('Failed to load site config:', e)
    }
  }

  return {
    // Theme
    themeMode, isDark, setTheme,
    // Network
    networkMode, sidebarCollapsed, toggleNetwork, toggleSidebar,
    // Wallpaper
    wallpaperType, wallpaperUrl, wallpaperBlur, wallpaperOpacity, showWallpaperCopyright,
    setWallpaperType, setWallpaperUrl, setWallpaperBlur, setWallpaperOpacity, setShowWallpaperCopyright,
    // Clock
    clockShow, clockShowSeconds, clockFormat, setClockShow, setClockShowSeconds, setClockFormat,
    // Logo
    logoType, logoImageUrl, logoText, setLogoType, setLogoImageUrl, setLogoText,
    // Card
    cardLabelSize, setCardLabelSize, cardItemSize, setCardItemSize, cardsPerRow, setCardsPerRow,
    // Layout
    contentMaxWidth, contentPaddingX, contentPaddingTop, contentPaddingBottom,
    setContentMaxWidth, setContentPaddingX, setContentPaddingTop, setContentPaddingBottom,
    // Colors
    bgColor, textColor, borderColor, primaryColor, setBgColor, setTextColor, setBorderColor, setPrimaryColor,
    // Footer
    footerHtml, setFooterHtml,
    // Site
    siteTitle, siteIconUrl, loginBgUrl, enableCaptcha,
    setSiteTitle, setSiteIconUrl, setLoginBgUrl, setEnableCaptcha,
    // Components
    showSystemStatus, showWeatherWidget, showMemo, showRSSWidget, bmsyncEnabled, networkModeAuto, groupCardTransparent, searchEngine, cardStyle,
    showSystemMonitor,
    setShowSystemStatus, setShowWeatherWidget, setShowMemo, setShowRSSWidget, setNetworkModeAuto, setGroupCardTransparent, setSearchEngine, setCardStyle,
    setShowSystemMonitor,
    // Server sync
    loadSettingsFromServer, applyServerSettings,
    // Auth settings
    allowRegistration, requireApproval, fetchAuthSettings,
    // Site config
    fetchSiteConfig,
  }
})
