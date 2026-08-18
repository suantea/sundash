<template>
  <n-config-provider :theme="theme" :theme-overrides="themeOverrides">
    <n-loading-bar-provider>
      <n-message-provider>
        <n-notification-provider>
          <n-dialog-provider>
            <router-view v-slot="{ Component }">
              <transition name="page" mode="out-in">
                <component :is="Component" />
              </transition>
            </router-view>
          </n-dialog-provider>
        </n-notification-provider>
      </n-message-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed, watchEffect, onMounted } from 'vue'
import { darkTheme, type GlobalThemeOverrides, NConfigProvider, NLoadingBarProvider, NMessageProvider, NNotificationProvider, NDialogProvider } from 'naive-ui'
import { useAppStore } from './stores/app'

const appStore = useAppStore()

onMounted(() => {
  appStore.fetchSiteConfig()
})

const theme = computed(() => appStore.isDark ? darkTheme : null)

const defaultPrimary = '#007AFF'
function adjustColor(hex: string, amount: number): string {
  hex = hex.replace('#', '')
  if (hex.length === 3) hex = hex.split('').map(c => c + c).join('')
  const num = parseInt(hex, 16)
  let r = Math.min(255, Math.max(0, (num >> 16) + amount))
  let g = Math.min(255, Math.max(0, ((num >> 8) & 0x00FF) + amount))
  let b = Math.min(255, Math.max(0, (num & 0x0000FF) + amount))
  return `#${((1 << 24) + (r << 16) + (g << 8) + b).toString(16).slice(1)}`
}

const themeOverrides = computed<GlobalThemeOverrides>(() => {
  const primary = appStore.primaryColor || defaultPrimary
  return {
    common: {
      primaryColor: primary,
      primaryColorHover: adjustColor(primary, -30),
      primaryColorPressed: adjustColor(primary, -60),
      primaryColorSuppl: primary,
      borderRadius: '8px',
      borderRadiusSmall: '6px',
      fontFamily: '-apple-system, BlinkMacSystemFont, "SF Pro Display", "Helvetica Neue", Roboto, sans-serif',
      fontSize: '15px',
      fontSizeMini: '11px',
      fontSizeTiny: '12px',
      fontSizeSmall: '13px',
      fontSizeMedium: '15px',
      fontSizeLarge: '16px',
      heightMedium: '36px',
      heightLarge: '40px',
    },
    Card: { borderRadius: '16px' },
    Input: {
      borderRadius: '10px',
      heightMedium: '40px',
      heightLarge: '44px',
      boxShadowFocus: `0 0 0 3px ${primary}1F`,
    },
    Button: { borderRadiusMedium: '10px', borderRadiusLarge: '12px', heightLarge: '42px', fontWeight: '500' },
    DataTable: { borderRadius: '16px' },
    Modal: { borderRadius: '20px', titleFontSize: '18px', titleFontWeight: '600' },
    Tag: { borderRadius: '8px' },
    Tabs: { borderRadius: '10px' },
    Dropdown: { borderRadius: '12px', optionBorderRadius: '8px' },
    Message: { borderRadius: '12px' },
  }
})

watchEffect(() => {
  document.documentElement.setAttribute('data-theme', appStore.isDark ? 'dark' : 'light')
  const primary = appStore.primaryColor || defaultPrimary
  document.documentElement.style.setProperty('--sd-primary', primary)
  document.documentElement.style.setProperty('--sd-primary-hover', adjustColor(primary, -30))
  document.documentElement.style.setProperty('--sd-primary-light', `${primary}14`)
  document.documentElement.style.setProperty('--sd-primary-medium', `${primary}26`)
})
</script>
