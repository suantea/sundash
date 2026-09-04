<template>
  <n-drawer :show="show" @update:show="$emit('update:show', $event)" :width="drawerWidth" placement="right">
    <n-drawer-content :title="$t('settings.title')" :native-scrollbar="false">
      <div class="drawer-settings">
        <!-- ===== 1. 语言 ===== -->
        <section class="settings-section">
          <div class="section-label">{{ $t('settings.language') }}</div>
          <div class="settings-card">
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
                  <Icon icon="mdi:translate" :width="18" :height="18" color="#007AFF" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.language') }}</div>
                  <div class="setting-desc">{{ $t('settings.languageDesc') }}</div>
                </div>
              </div>
              <div class="seg-control">
                <button :class="['seg-btn', { active: locale === 'zh-CN' }]" @click="setLocale('zh-CN')">中文</button>
                <button :class="['seg-btn', { active: locale === 'en' }]" @click="setLocale('en')">English</button>
              </div>
            </div>
          </div>
        </section>

        <!-- ===== 2. 页面 ===== -->
        <section class="settings-section">
          <div class="section-label">{{ $t('settings.site') }}</div>
          <div class="settings-card">
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
                  <Icon icon="mdi:lan" :width="18" :height="18" color="#007AFF" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.networkMode') }}</div>
                  <div class="setting-desc">{{ $t('settings.networkModeDesc') }}</div>
                </div>
              </div>
              <div class="seg-control">
                <button :class="['seg-btn', { active: appStore.networkMode === 'internal' }]" @click="setNetwork('internal')">{{ $t('settings.networkInternal') }}</button>
                <button :class="['seg-btn', { active: appStore.networkMode === 'external' }]" @click="setNetwork('external')">{{ $t('settings.networkExternal') }}</button>
              </div>
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(255,149,0,0.1);">
                  <Icon icon="mdi:image" :width="18" :height="18" color="#FF9500" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.logoType') }}</div>
                  <div class="setting-desc">{{ $t('settings.logoTypeDesc') }}</div>
                </div>
              </div>
              <div class="seg-control">
                <button :class="['seg-btn', { active: appStore.logoType === 'text' }]" @click="appStore.setLogoType('text')">{{ $t('settings.logoText') }}</button>
                <button :class="['seg-btn', { active: appStore.logoType === 'image' }]" @click="appStore.setLogoType('image')">{{ $t('settings.logoImage') }}</button>
              </div>
            </div>
            <div v-if="appStore.logoType === 'image'" class="setting-sub">
              <input v-model="logoImageInput" class="bg-url-input" :placeholder="$t('settings.logoImageUrlPlaceholder')" @change="appStore.setLogoImageUrl(logoImageInput)" />
            </div>
            <div v-if="appStore.logoType === 'text'" class="setting-sub">
              <input v-model="logoTextInput" class="bg-url-input" :placeholder="$t('settings.logoTextPlaceholder')" @change="appStore.setLogoText(logoTextInput)" />
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
                  <Icon icon="mdi:text-short" :width="18" :height="18" color="#007AFF" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.siteTitle') }}</div>
                  <div class="setting-desc">{{ $t('settings.siteTitleDesc') }}</div>
                </div>
              </div>
              <input v-model="siteTitleInput" class="mini-text-input" :placeholder="$t('settings.siteTitlePlaceholder')" @change="updateSiteTitle" />
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row" style="flex-direction: column; align-items: stretch; gap: 12px;">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(52,199,89,0.1);">
                  <Icon icon="mdi:code-tags" :width="18" :height="18" color="#34C759" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.footerHtml') }}</div>
                  <div class="setting-desc">{{ $t('settings.footerHtmlDesc') }}</div>
                </div>
              </div>
              <textarea v-model="footerInput" class="footer-input" :placeholder="$t('settings.footerHtmlPlaceholder')" rows="2" @change="appStore.setFooterHtml(footerInput)"></textarea>
            </div>
          </div>
        </section>

        <!-- ===== 2. 外观 ===== -->
        <section class="settings-section">
          <div class="section-label">{{ $t('settings.appearance') }}</div>
          <div class="settings-card">
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(255,149,0,0.1);">
                  <Icon icon="mdi:theme-light-dark" :width="18" :height="18" color="#FF9500" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.theme') }}</div>
                  <div class="setting-desc">{{ $t('settings.themeDesc') }}</div>
                </div>
              </div>
              <div class="seg-control">
                <button :class="['seg-btn', { active: appStore.themeMode === 'light' }]" @click="appStore.setTheme('light')">{{ $t('settings.themeLight') }}</button>
                <button :class="['seg-btn', { active: appStore.themeMode === 'dark' }]" @click="appStore.setTheme('dark')">{{ $t('settings.themeDark') }}</button>
                <button :class="['seg-btn', { active: appStore.themeMode === 'system' }]" @click="appStore.setTheme('system')">{{ $t('settings.themeSystem') }}</button>
              </div>
              <!-- Custom themes -->
              <div v-if="customThemes.length" class="theme-list">
                <div v-for="theme in customThemes" :key="theme.id"
                  :class="['theme-chip', { active: activeThemeId === theme.id }]"
                  @click="applyTheme(theme)">
                  <span class="theme-dot" :style="{ background: extractPrimaryColor(theme.css_content) }"></span>
                  <span class="theme-chip-name">{{ theme.name }}</span>
                  <button v-if="!theme.is_builtin" class="theme-chip-del" @click.stop="deleteTheme(theme.id)" title="删除">
                    <Icon icon="mdi:close" :width="12" :height="12" />
                  </button>
                </div>
              </div>
              <button class="apply-btn theme-add-btn" @click="openThemeEditor">
                <Icon icon="mdi:plus" :width="14" :height="14" />
                {{ $t('settings.createTheme') || '新建主题' }}
              </button>
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(52,199,89,0.1);">
                  <Icon icon="mdi:image-outline" :width="18" :height="18" color="#34C759" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.wallpaper') }}</div>
                  <div class="setting-desc">{{ $t('settings.wallpaperDesc') }}</div>
                </div>
              </div>
              <div class="seg-control">
                <button :class="['seg-btn', { active: appStore.wallpaperType === 'default' }]" @click="setWallpaperTypeWithSync('default')">{{ $t('settings.wallpaperDefault') }}</button>
                <button :class="['seg-btn', { active: appStore.wallpaperType === 'gradient' }]" @click="setWallpaperTypeWithSync('gradient')">{{ $t('settings.wallpaperGradient') }}</button>
                <button :class="['seg-btn', { active: appStore.wallpaperType === 'bing' }]" @click="setWallpaperTypeWithSync('bing')">{{ $t('settings.wallpaperBing') }}</button>
                <button :class="['seg-btn', { active: appStore.wallpaperType === 'custom' }]" @click="setWallpaperTypeWithSync('custom')">{{ $t('settings.wallpaperCustom') }}</button>
              </div>
            </div>
            <div v-if="appStore.wallpaperType === 'bing'" class="setting-sub">
              <button class="apply-btn" @click="fetchBingWallpaper">{{ $t('settings.refreshBingWallpaper') }}</button>
              <span v-if="appStore.wallpaperUrl" class="bing-preview">
                <img :src="appStore.wallpaperUrl" :alt="$t('settings.preview')" class="bing-preview-img" />
              </span>
            </div>
            <div v-if="appStore.wallpaperType === 'custom'" class="setting-sub">
              <input v-model="customUrl" class="bg-url-input" :placeholder="$t('settings.wallpaperUrlPlaceholder')" />
              <button class="apply-btn" @click="applyCustomBackground">{{ $t('settings.apply') }}</button>
            </div>
            <template v-if="appStore.wallpaperType !== 'default'">
              <div class="setting-divider"></div>
              <div class="setting-row">
                <div class="setting-left">
                  <div class="setting-icon" style="background: rgba(255,149,0,0.1);">
                    <Icon icon="mdi:blur" :width="18" :height="18" color="#FF9500" />
                  </div>
                  <div>
                    <div class="setting-title">{{ $t('settings.wallpaperBlur') }}</div>
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
                    <Icon icon="mdi:opacity" :width="18" :height="18" color="#007AFF" />
                  </div>
                  <div>
                    <div class="setting-title">{{ $t('settings.wallpaperOpacity') }}</div>
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
                  <Icon icon="mdi:palette" :width="18" :height="18" color="#007AFF" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.primaryColor') }}</div>
                  <div class="setting-desc">{{ $t('settings.primaryColorDesc') }}</div>
                </div>
              </div>
              <div class="color-picker-wrap">
                <input type="color" :value="appStore.primaryColor || '#007AFF'" @input="appStore.setPrimaryColor(($event.target as HTMLInputElement).value)" class="color-input" />
                <button v-if="appStore.primaryColor" class="reset-btn" @click="appStore.setPrimaryColor('')">{{ $t('settings.reset') }}</button>
              </div>
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(255,149,0,0.1);">
                  <Icon icon="mdi:format-color-fill" :width="18" :height="18" color="#FF9500" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.borderColor') }}</div>
                  <div class="setting-desc">{{ $t('settings.borderColorDesc') }}</div>
                </div>
              </div>
              <div class="color-picker-wrap">
                <input type="color" :value="appStore.borderColor || '#000000'" @input="appStore.setBorderColor(($event.target as HTMLInputElement).value)" class="color-input" />
                <button v-if="appStore.borderColor" class="reset-btn" @click="appStore.setBorderColor('')">{{ $t('settings.reset') }}</button>
              </div>
            </div>
          </div>
        </section>

        <!-- ===== 3. 布局 ===== -->
        <section class="settings-section">
          <div class="section-label">{{ $t('settings.layout') }}</div>
          <div class="settings-card">
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
                  <Icon icon="mdi:arrow-expand-horizontal" :width="18" :height="18" color="#007AFF" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.maxWidth') }}</div>
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
                  <Icon icon="mdi:format-horizontal-align-left" :width="18" :height="18" color="#34C759" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.paddingX') }}</div>
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
                  <Icon icon="mdi:format-vertical-align-top" :width="18" :height="18" color="#AF52DE" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.paddingTop') }}</div>
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
                  <Icon icon="mdi:format-vertical-align-bottom" :width="18" :height="18" color="#FF9500" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.paddingBottom') }}</div>
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
          <div class="section-label">{{ $t('settings.widgets') }}</div>
          <div class="settings-card">
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(52,199,89,0.1);">
                  <Icon icon="mdi:clock-outline" :width="18" :height="18" color="#34C759" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.clock') }}</div>
                  <div class="setting-desc">{{ $t('settings.clockDesc') }}</div>
                </div>
              </div>
              <n-switch v-model:value="appStore.clockShow" @update:value="appStore.setClockShow" />
            </div>
            <div v-if="appStore.clockShow" class="setting-sub" style="flex-direction: column; align-items: stretch; gap: 4px;">
              <div class="setting-row" style="padding: 6px 0;">
                <span class="mini-label">{{ $t('settings.showSeconds') }}</span>
                <n-switch v-model:value="appStore.clockShowSeconds" @update:value="appStore.setClockShowSeconds" size="small" />
              </div>
              <div class="setting-row" style="padding: 6px 0;">
                <span class="mini-label">{{ $t('settings.clockFormat') }}</span>
                <div class="seg-control">
                  <button :class="['seg-btn', { active: appStore.clockFormat === '24' }]" @click="appStore.setClockFormat('24')">{{ $t('settings.hours24') }}</button>
                  <button :class="['seg-btn', { active: appStore.clockFormat === '12' }]" @click="appStore.setClockFormat('12')">{{ $t('settings.hours12') }}</button>
                </div>
              </div>
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(175,82,222,0.1);">
                  <Icon icon="mdi:server-network" :width="18" :height="18" color="#AF52DE" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.systemStatus') }}</div>
                  <div class="setting-desc">{{ $t('settings.systemStatusDesc') }}</div>
                </div>
              </div>
              <n-switch v-model:value="appStore.showSystemStatus" @update:value="appStore.setShowSystemStatus" />
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(88,86,214,0.1);">
                  <Icon icon="mdi:monitor-dashboard" :width="18" :height="18" color="#5856D6" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.systemMonitor') }}</div>
                  <div class="setting-desc">{{ $t('settings.systemMonitorDesc') }}</div>
                </div>
              </div>
              <n-switch v-model:value="appStore.showSystemMonitor" @update:value="appStore.setShowSystemMonitor" />
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(255,149,0,0.1);">
                  <Icon icon="mdi:weather-partly-cloudy" :width="18" :height="18" color="#FF9500" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.weatherWidget') }}</div>
                </div>
              </div>
              <n-switch v-model:value="appStore.showWeatherWidget" @update:value="appStore.setShowWeatherWidget" />
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(52,199,89,0.1);">
                  <Icon icon="mdi:note-text-outline" :width="18" :height="18" color="#34C759" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.memoWidget') }}</div>
                </div>
              </div>
              <n-switch v-model:value="appStore.showMemo" @update:value="appStore.setShowMemo" />
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(255,69,58,0.1);">
                  <Icon icon="mdi:rss" :width="18" :height="18" color="#FF453A" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.rssWidget') }}</div>
                </div>
              </div>
              <n-switch v-model:value="appStore.showRSSWidget" @update:value="appStore.setShowRSSWidget" />
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(255,149,0,0.1);">
                  <Icon icon="mdi:format-size" :width="18" :height="18" color="#FF9500" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.cardLabelSize') }}</div>
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
                  <Icon icon="mdi:resize" :width="18" :height="18" color="#007AFF" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.cardItemSize') }}</div>
                  <div class="setting-desc">{{ parseInt(appStore.cardItemSize) || 0 > 0 ? '+' + appStore.cardItemSize + 'px' : $t('settings.cardDefault') }}</div>
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
                  <Icon icon="mdi:view-grid" :width="18" :height="18" color="#AF52DE" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.cardsPerRow') }}</div>
                  <div class="setting-desc">{{ parseInt(appStore.cardsPerRow) || 5 }} {{ $t('settings.cardsPerRowUnit') }}</div>
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
                  <Icon icon="mdi:eye-off" :width="18" :height="18" color="#34C759" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.groupCardTransparent') }}</div>
                  <div class="setting-desc">{{ $t('settings.groupCardTransparentDesc') }}</div>
                </div>
              </div>
              <n-switch v-model:value="appStore.groupCardTransparent" @update:value="appStore.setGroupCardTransparent" />
            </div>
          </div>
        </section>

        <!-- ===== 5. 数据 ===== -->
        <section class="settings-section">
          <div class="section-label">{{ $t('settings.data') }}</div>
          <div class="settings-card">
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
                  <Icon icon="mdi:export" :width="18" :height="18" color="#007AFF" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.exportData') }}</div>
                  <div class="setting-desc">{{ $t('settings.exportDataDesc') }}</div>
                </div>
              </div>
              <button class="action-btn" @click="handleExport">{{ $t('settings.export') }}</button>
            </div>
            <div class="setting-divider"></div>
            <div class="setting-row">
              <div class="setting-left">
                <div class="setting-icon" style="background: rgba(52,199,89,0.1);">
                  <Icon icon="mdi:import" :width="18" :height="18" color="#34C759" />
                </div>
                <div>
                  <div class="setting-title">{{ $t('settings.importData') }}</div>
                  <div class="setting-desc">{{ $t('settings.importDataDesc') }}</div>
                </div>
              </div>
              <div style="display: flex; gap: 8px;">
                <button class="action-btn" @click="importTemplate">{{ $t('settings.importTemplate') }}</button>
                <button class="action-btn" @click="triggerImport">{{ $t('settings.import') }}</button>
              </div>
            </div>
          </div>
        </section>
      </div>

      <template #footer>
        <div class="drawer-footer">
          <button v-if="isAdmin" class="admin-entry-btn" @click="goToAdmin">
            <Icon icon="mdi:shield-cog" :width="14" :height="14" />
            <span>{{ $t('settings.goToAdmin') }}</span>
            <Icon icon="mdi:chevron-right" :width="14" :height="14" class="admin-entry-arrow" />
          </button>
          <span class="footer-text">{{ $t('settings.autoSave') }}</span>
        </div>
      </template>
    </n-drawer-content>
  </n-drawer>
  <input ref="fileInput" type="file" accept=".json" style="display: none" @change="handleImport" />
</template>

<script setup lang="ts">
import { ref, watch, computed, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Icon } from '@iconify/vue'
import { useMessage, NDrawer, NDrawerContent, NSwitch } from 'naive-ui'
import { useRouter } from 'vue-router'
import { useAppStore } from '../../stores/app'
import { useUserStore } from '../../stores/user'
import { usePanelStore } from '../../stores/panel'
import { importTemplates } from '../../data/importTemplates'

defineProps<{ show: boolean }>()
const emit = defineEmits<{ 'update:show': [val: boolean] }>()

const { t, locale } = useI18n()

function setLocale(lang: string) {
  locale.value = lang
  localStorage.setItem('sundash-locale', lang)
}
const appStore = useAppStore()
const panelStore = usePanelStore()
const userStore = useUserStore()
const router = useRouter()

// 仅管理员能进入管理后台（/admin 路由有权限守卫）
const isAdmin = computed(() => userStore.user?.role === 'admin')

function goToAdmin() {
  emit('update:show', false)
  router.push('/admin')
}

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
      message.success(t('settings.wallpaperUpdated'))
    }
  } catch (e) {
    message.error(t('settings.wallpaperFailed'))
  }
}

function applyCustomBackground() {
  if (!customUrl.value.trim()) {
    message.error(t('settings.enterImageUrl'))
    return
  }
  appStore.setWallpaperUrl(customUrl.value.trim())
  message.success(t('settings.customBgApplied'))
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
  message.success(t('settings.exportSuccess'))
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
    message.success(t('settings.settingsRestored'))
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
    message.success(t('settings.importSuccess', { n: data.groups.length }))
  }
  // Sun-Panel 导出格式（{ icons: [{ title, children: [{ icon: { src/text }, url, ... }] }] }）
  if (data.icons && Array.isArray(data.icons)) {
    let totalCards = 0
    for (const group of data.icons) {
      const newGroup = await panelStore.createGroup(group.title)
      if (!group.children || !Array.isArray(group.children)) continue

      // 收集需要获取 favicon 的 URL
      const urlsToFetch: string[] = []
      const cardDataList: Array<{ child: any; iconObj: any }> = []

      for (const child of group.children) {
        const iconObj: Record<string, any> = child.icon || {}
        let cardIcon = ''
        let cardIconColor = '#2a2a2a6b'
        // itemType 0=文字图标(src为空)，1=图片图标(src有值)，3=图标文本(text=iconify key)
        if (iconObj.itemType === 1 && iconObj.src) {
          cardIcon = iconObj.src
        } else if (iconObj.itemType === 3 && iconObj.text) {
          cardIcon = iconObj.text // iconify key
        } else if (iconObj.itemType === 2 && iconObj.src) {
          cardIcon = iconObj.src
        } else if (iconObj.itemType === 0 && iconObj.text) {
          cardIcon = iconObj.text
        }
        if (iconObj.backgroundColor) cardIconColor = iconObj.backgroundColor

        cardDataList.push({ child, iconObj })
        // 如果没有 icon 且有 URL，收集用于获取 favicon
        if (!cardIcon && child.url) {
          urlsToFetch.push(child.url)
        }
      }

      // 批量获取 favicon
      let faviconResults: any[] = []
      if (urlsToFetch.length > 0) {
        try {
          const { api } = await import('../../api')
          const res = await api.post('favicons', { urls: urlsToFetch })
          faviconResults = res.data?.results || []
        } catch {
          faviconResults = []
        }
      }

      // 创建卡片
      for (let i = 0; i < cardDataList.length; i++) {
        const { child, iconObj } = cardDataList[i]
        let cardIcon = ''
        let cardIconColor = '#2a2a2a6b'
        if (iconObj.itemType === 1 && iconObj.src) {
          cardIcon = iconObj.src
        } else if (iconObj.itemType === 3 && iconObj.text) {
          cardIcon = iconObj.text
        } else if (iconObj.itemType === 2 && iconObj.src) {
          cardIcon = iconObj.src
        } else if (iconObj.itemType === 0 && iconObj.text) {
          cardIcon = iconObj.text
        }
        if (iconObj.backgroundColor) cardIconColor = iconObj.backgroundColor

        // 如果原来没有 icon，尝试用批量获取的结果
        if (!cardIcon && faviconResults[i]) {
          const favicon = faviconResults[i]
          cardIcon = favicon.icon_name || favicon.favicon_url || ''
        }

        await panelStore.createCard({
          group_id: newGroup.id,
          title: child.title,
          url: child.url || '',
          url_internal: child.lanUrl || '',
          icon: cardIcon,
          icon_color: cardIconColor,
          bg_color: child.backgroundColor || '',
          description: child.description || '',
          open_type: 'new_tab',
        })
        totalCards++
      }
    }
    await panelStore.fetchPanel()
    message.success(t('settings.sunPanelImportSuccess', { n: data.icons.length, cards: totalCards }))
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
    message.error(t('settings.importFailed'))
  } finally {
    input.value = ''
  }
}

function importTemplate() {
  const tmpl = importTemplates[0]
  const cardCount = tmpl.groups.reduce((sum, g) => sum + (g.cards?.length || 0), 0)
  message.info(t('settings.importTemplateInfo', { name: tmpl.name, groups: tmpl.groups.length, cards: cardCount }))
  doImport({ groups: tmpl.groups }).catch(() => message.error(t('settings.templateImportFailed')))
}

// Theme management
interface ThemeItem { id: string; name: string; description: string; css_content: string; is_builtin: boolean }
const customThemes = ref<ThemeItem[]>([])
const activeThemeId = ref(localStorage.getItem('sundash-active-theme') || '')

async function fetchThemes() {
  try {
    const res = await api.get('themes')
    customThemes.value = res.data || []
  } catch { /* ignore */ }
}

function applyTheme(theme: ThemeItem) {
  let el = document.getElementById('sundash-custom-theme')
  if (!el) {
    el = document.createElement('style')
    el.id = 'sundash-custom-theme'
    document.head.appendChild(el)
  }
  if (activeThemeId.value === theme.id) {
    // Toggle off
    activeThemeId.value = ''
    el.textContent = ''
    localStorage.removeItem('sundash-active-theme')
  } else {
    el.textContent = theme.css_content
    activeThemeId.value = theme.id
    localStorage.setItem('sundash-active-theme', theme.id)
    message.success(`已应用主题「${theme.name}」`)
  }
}

async function deleteTheme(id: string) {
  try {
    await api.delete(`themes/${id}`)
    customThemes.value = customThemes.value.filter(t => t.id !== id)
    if (activeThemeId.value === id) {
      activeThemeId.value = ''
      const el = document.getElementById('sundash-custom-theme')
      if (el) el.textContent = ''
      localStorage.removeItem('sundash-active-theme')
    }
    message.success('主题已删除')
  } catch { message.error('删除失败') }
}

function extractPrimaryColor(css: string): string {
  const m = css.match(/--sd-primary:\s*(#[0-9a-fA-F]{3,8}|rgba?\([^)]+\))/)
  return m ? m[1] : '#007AFF'
}

function openThemeEditor() {
  const name = prompt('主题名称：')
  if (!name) return
  const css = prompt('CSS 内容（自定义变量，如 :root { --sd-primary: #ff0000 }）：')
  if (!css) return
  api.post('themes', { name, css_content: css }).then(res => {
    customThemes.value.push(res.data)
    message.success('主题已创建')
  }).catch(() => message.error('创建失败'))
}

// Restore active theme on mount
const savedThemeId = localStorage.getItem('sundash-active-theme')
if (savedThemeId) {
  fetchThemes().then(() => {
    const theme = customThemes.value.find(t => t.id === savedThemeId)
    if (theme) {
      const el = document.createElement('style')
      el.id = 'sundash-custom-theme'
      el.textContent = theme.css_content
      document.head.appendChild(el)
      activeThemeId.value = savedThemeId
    }
  })
} else {
  fetchThemes()
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

/* 主题区「新建主题」按钮：在 flex 行内独立占一行，不再与分段控件挤压 */
.theme-add-btn {
  width: 100%;
  justify-content: center;
  margin-top: 8px;
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
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
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
}

.admin-entry-btn {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 4px 10px;
  border: none;
  border-radius: 6px;
  background: transparent;
  color: var(--sd-text-secondary);
  font-size: 12px;
  cursor: pointer;
  transition: background 0.15s ease, color 0.15s ease;
}

.admin-entry-btn:hover {
  background: var(--sd-primary-light);
  color: var(--sd-primary);
}

.admin-entry-arrow {
  opacity: 0.5;
}

.footer-text {
  font-size: 11px;
  color: var(--sd-text-tertiary);
}

/* 窄屏适配 */
@media (max-width: 480px) {
  .drawer-settings {
    padding: 4px 0 16px;
    gap: 12px;
  }

  .section-label {
    font-size: 12px;
    margin-bottom: 6px;
  }

  .settings-card {
    border-radius: 12px;
  }

  .setting-row {
    padding: 10px 12px;
    min-height: 44px;
    flex-wrap: wrap;
  }

  .setting-left {
    gap: 6px;
  }

  .setting-title {
    font-size: 13px;
  }

  .setting-desc {
    font-size: 11px;
  }

  .setting-sub {
    padding: 6px 12px 12px;
  }
}

/* Theme chips */
.theme-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 8px;
  width: 100%;
  flex-shrink: 0;
}

.theme-chip {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  border-radius: 8px;
  background: rgba(0,0,0,0.04);
  cursor: pointer;
  font-size: 12px;
  transition: all 0.15s;
}

.theme-chip:hover {
  background: rgba(0,0,0,0.08);
}

.theme-chip.active {
  background: rgba(0, 122, 255, 0.12);
  color: var(--sd-primary, #007AFF);
}

:root[data-theme="dark"] .theme-chip {
  background: rgba(255,255,255,0.06);
}

:root[data-theme="dark"] .theme-chip:hover {
  background: rgba(255,255,255,0.1);
}

:root[data-theme="dark"] .theme-chip.active {
  background: rgba(0, 122, 255, 0.2);
}

.theme-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;
  flex-shrink: 0;
}

.theme-chip-name {
  max-width: 80px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.theme-chip-del {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  border: none;
  background: transparent;
  color: #999;
  cursor: pointer;
  border-radius: 50%;
  padding: 0;
}

.theme-chip-del:hover {
  background: rgba(244, 67, 54, 0.15);
  color: #f44336;
}
</style>