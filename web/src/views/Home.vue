<template>
  <div class="home-page" :class="{ 'dark-wallpaper': isDarkBg }">
    <!-- Wallpaper background (fade in after image loads to prevent flicker) -->
    <div 
      class="wallpaper-bg"
      :class="{
        'wallpaper-bing': appStore.wallpaperType === 'bing',
        'wallpaper-gradient': appStore.wallpaperType === 'gradient',
        'wallpaper-custom': appStore.wallpaperType === 'custom'
      }"
      :style="{ ...wallpaperStyle, opacity: wallpaperReady ? (wallpaperStyle.opacity || 1) : 0, transition: 'opacity 0.4s ease' }"
    >
      <!-- Copyright overlay -->
      <div 
        v-if="appStore.showWallpaperCopyright && wallpaperCopyright" 
        class="wallpaper-copyright"
      >
        {{ wallpaperCopyright }}
      </div>
    </div>

    <!-- Animated background (hidden when wallpaper is active) -->
    <div class="home-bg" v-show="appStore.wallpaperType === 'default'">
      <div class="bg-orb bg-orb-1"></div>
      <div class="bg-orb bg-orb-2"></div>
      <div class="bg-orb bg-orb-3"></div>
    </div>

    <!-- Header -->
    <header class="home-header">
      <!-- Center content: Logo + Divider + Clock -->
      <div class="header-center-content">
        <div class="header-logo">
          <span class="brand-text">{{ appStore.logoText }}</span>
        </div>
        <div class="header-divider">|</div>
        <div class="header-clock">
          <ClockDisplay :compact="true" />
        </div>
      </div>

      <!-- Right actions -->
      <div class="header-actions">
        <button class="hdr-btn" @click="appStore.toggleNetwork()" :title="appStore.networkMode === 'internal' ? '内网模式' : '外网模式'">
          <Icon :icon="appStore.networkMode === 'internal' ? 'mdi:lan' : 'mdi:wan'" :size="17" />
        </button>
        <button class="hdr-btn" @click="toggleTheme" :title="appStore.isDark ? '浅色模式' : '深色模式'">
          <Icon :icon="appStore.isDark ? 'mdi:weather-sunny' : 'mdi:weather-night'" :size="17" />
        </button>
        <div class="hdr-divider"></div>
        <button class="hdr-btn" @click="router.push('/bookmarks')" title="书签同步">
          <Icon icon="mdi:bookmark-multiple-outline" :size="17" />
        </button>
        <button class="hdr-btn accent" @click="showAddGroup = true" title="新建分组">
          <Icon icon="mdi:plus" :size="17" />
        </button>
        <button class="hdr-btn" @click="showSettings = true" title="设置">
          <Icon icon="mdi:cog-outline" :size="17" />
        </button>
        <n-dropdown :options="userMenuOptions" @select="handleUserMenu" trigger="click">
          <button class="avatar-btn">
            <span>{{ (userStore.user?.display_name || userStore.user?.username || 'U')[0].toUpperCase() }}</span>
          </button>
        </n-dropdown>
      </div>
    </header>

    <!-- Search bar -->
    <div class="home-search">
      <SearchBar />
    </div>

    <!-- Main Content -->
    <main class="home-content">
      <!-- System Status -->
      <SystemStatus />

      <!-- System Monitor -->
      <SystemMonitor />

      <!-- Weather Widget -->
      <WeatherWidget />

      <!-- Memo Widget -->
      <MemoWidget />

      <!-- RSS Widget (可选：设置里开启) -->
      <RSSWidget v-if="appStore.showRSSWidget" />

      <!-- Empty State -->
      <div v-if="panelStore.groups.length === 0 && !panelStore.loading" class="empty-state">
        <div class="empty-icon">
          <svg width="56" height="56" viewBox="0 0 48 48" fill="none">
            <rect x="4" y="4" width="18" height="18" rx="5" fill="#007AFF"/>
            <rect x="26" y="4" width="18" height="18" rx="5" fill="#007AFF" opacity="0.4"/>
            <rect x="4" y="26" width="18" height="18" rx="5" fill="#007AFF" opacity="0.4"/>
            <rect x="26" y="26" width="18" height="18" rx="5" fill="#007AFF" opacity="0.15"/>
          </svg>
        </div>
        <h2>欢迎使用 SunDash</h2>
        <p>创建你的第一个分组，开始组织你的导航</p>
        <button class="cta-btn" @click="showAddGroup = true">
          <Icon icon="mdi:plus" :size="18" />
          <span>新建分组</span>
        </button>
      </div>

      <!-- Card Groups -->
      <div v-else class="groups-container">
        <div v-for="group in visibleGroups" :key="group.id" class="group-wrapper">
          <div class="group-card" :class="{ 'card-transparent': appStore.groupCardTransparent }">
            <div class="group-head">
              <div class="group-name" @click="toggleCollapse(group.id)" style="cursor: pointer;">
                <Icon :icon="isCollapsed(group.id) ? 'mdi:chevron-right' : 'mdi:chevron-down'" :size="18" class="collapse-icon" />
                <h3>{{ group.name }}</h3>
                <span class="group-badge">{{ group.cards?.length || 0 }}</span>
              </div>
              <div class="group-tools">
                <button class="tool-btn" @click="openAddCard(group.id)" title="添加卡片">
                  <Icon icon="mdi:plus" :size="15" />
                </button>
                <button class="tool-btn" @click="toggleHideGroup(group.id)" :title="hiddenGroupIds.has(group.id) ? '显示分组' : '隐藏分组'">
                  <Icon :icon="hiddenGroupIds.has(group.id) ? 'mdi:eye-off' : 'mdi:eye-outline'" :size="15" />
                </button>
                <button class="tool-btn" @click="openEditGroup({ id: group.id, name: group.name })" title="$t('home.editGroup')">
                  <Icon icon="mdi:pencil-outline" :size="15" />
                </button>
                <button class="tool-btn danger" @click="confirmDeleteGroup(group.id)" title="$t('home.deleteGroup')">
                  <Icon icon="mdi:delete-outline" :size="15" />
                </button>
              </div>
            </div>

            <div v-show="!isCollapsed(group.id)" class="group-cards-wrap">
              <div class="group-cards" :style="gridStyle">
                <CardItem
                  v-for="card in groupCards(group)"
                  :key="card.id"
                  :card="card"
                  :class="{ 'card-is-hidden': hiddenCardIds.has(card.id) }"
                  @click="openCardUrl(card)"
                  @edit="openEditCard(card)"
                  @delete="confirmDeleteCard(card.id)"
                  @hide="toggleHideCard(card.id)"
                />
              </div>
              <div v-if="hiddenCardsCount(group) > 0" class="hidden-cards-bar">
                <button class="hidden-cards-btn" @click="toggleShowHiddenCards(group.id)">
                  <Icon :icon="showHiddenCardsGroup.has(group.id) ? 'mdi:eye-off' : 'mdi:eye-outline'" :size="13" />
                  <span>{{ hiddenCardsCount(group) }} 个书签已隐藏</span>
                </button>
              </div>
            </div>
          </div>
        </div>

        <div v-if="hiddenGroupIds.size > 0" class="hidden-groups-bar">
          <button class="hidden-groups-btn" @click="showHiddenGroups = !showHiddenGroups">
            <Icon :icon="showHiddenGroups ? 'mdi:eye-off' : 'mdi:eye-outline'" :size="14" />
            <span>{{ hiddenGroupIds.size }} 个分组已隐藏</span>
            <Icon :icon="showHiddenGroups ? 'mdi:chevron-up' : 'mdi:chevron-down'" :size="14" />
          </button>
          <div v-if="showHiddenGroups" class="hidden-groups-list">
            <button v-for="gid in hiddenGroupIds" :key="gid" class="hidden-group-item" @click="toggleHideGroup(gid)">
              <Icon icon="mdi:eye" :size="14" />
              <span>{{ getGroupName(gid) }}</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <footer class="home-footer">
        <div v-if="appStore.footerHtml" class="footer-content" v-text="appStore.footerHtml"></div>
        <div class="copyright-notice">SunDash · Open Source · MIT License</div>
      </footer>
    </main>

    <!-- Modals -->
    <n-modal v-model:show="showAddGroup" preset="dialog" title="新建分组" positive-text="创建" negative-text="取消"
      :loading="modalLoading" @positive-click="handleCreateGroup">
      <n-input v-model:value="newGroupName" placeholder="$t('home.groupName')" autofocus />
    </n-modal>

    <n-modal v-model:show="showEditGroup" preset="dialog" title="$t('home.editGroup')" positive-text="保存" negative-text="取消"
      :loading="modalLoading" @positive-click="handleUpdateGroup">
      <n-input v-model:value="editGroupName" placeholder="$t('home.groupName')" />
    </n-modal>

    <CardEditor
      v-model:show="showCardEditor"
      :card="editingCard"
      :groups="panelStore.groups"
      :default-group-id="defaultCardGroupId"
      @save="handleSaveCard"
      @apply-color-all="handleApplyColorAll"
      @apply-bg-color-all="handleApplyBgColorAll"
    />

    <n-modal v-model:show="showDeleteConfirm" preset="dialog" type="warning"
      :title="deleteTarget.type === 'group' ? $t('home.deleteGroup') : $t('home.deleteCard')"
      :content="deleteTarget.type === 'group' ? '确定要删除该分组及其所有卡片吗？此操作无法撤销。' : '确定要删除该卡片吗？此操作无法撤销。'"
      positive-text="删除" negative-text="取消" :loading="modalLoading"
      @positive-click="handleDelete" />

    <SettingsDrawer v-model:show="showSettings" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, h } from 'vue'
import { useI18n } from 'vue-i18n'
import { Icon } from '@iconify/vue'
import { useMessage, type DropdownOption, NDropdown, NModal, NInput } from 'naive-ui'
import { useAppStore } from '../stores/app'
import { useUserStore } from '../stores/user'
import { usePanelStore } from '../stores/panel'
import { useWallpaper } from '../composables/useWallpaper'
import { useGroupVisibility } from '../composables/useGroupVisibility'
import router from '../router'
import CardEditor from '../components/card/CardEditor.vue'
import CardItem from '../components/card/CardItem.vue'
import ClockDisplay from '../components/clock/ClockDisplay.vue'
import SystemStatus from '../components/status/SystemStatus.vue'
import SystemMonitor from '../components/status/SystemMonitor.vue'
import SearchBar from '../components/search/SearchBar.vue'
import SettingsDrawer from '../components/settings/SettingsDrawer.vue'
import RSSWidget from '../components/rss/RSSWidget.vue'
import type { Card } from '../types'

const appStore = useAppStore()
const userStore = useUserStore()
const panelStore = usePanelStore()
const message = useMessage()
const { t } = useI18n()

// Grid style based on cardsPerRow setting
const gridStyle = computed(() => {
  const count = parseInt(appStore.cardsPerRow) || 5
  return { gridTemplateColumns: `repeat(${count}, 1fr)` }
})

// Wallpaper state & background luminance (extracted composable)
const { wallpaperCopyright, wallpaperReady, isDarkBg, fetchWallpaper, wallpaperStyle } = useWallpaper()

// Sync layout values to CSS variables
function syncLayoutVars() {
  const root = document.documentElement
  root.style.setProperty('--sd-content-max-width', appStore.contentMaxWidth)
  root.style.setProperty('--sd-content-padding-x', appStore.contentPaddingX)
  root.style.setProperty('--sd-content-padding-top', appStore.contentPaddingTop)
  root.style.setProperty('--sd-content-padding-bottom', appStore.contentPaddingBottom)
}
syncLayoutVars()
watch(() => [appStore.contentMaxWidth, appStore.contentPaddingX, appStore.contentPaddingTop, appStore.contentPaddingBottom], syncLayoutVars)


// Group collapse/hide state (extracted composable)
const {
  hiddenGroupIds, hiddenCardIds, showHiddenGroups, showHiddenCardsGroup,
  isCollapsed, toggleCollapse, toggleHideGroup, toggleHideCard,
  visibleGroups, getGroupName, hiddenCardsCount, toggleShowHiddenCards, groupCards,
} = useGroupVisibility()

// Fetch data — one bootstrap request (settings+profile+panels) instead of
// three round-trips; wallpaper runs after settings since it depends on the
// wallpaper type setting.
onMounted(async () => {
  try {
    const { api } = await import('../api')
    const res = await api.get('bootstrap')
    appStore.applyServerSettings(res.data.settings || {})
    userStore.setUser(res.data.profile)
    panelStore.setGroups(res.data.panels?.groups || [])
  } catch (e) {
    // Fallback: fetch individually if the bootstrap endpoint is unavailable.
    console.error('Bootstrap failed, falling back to individual requests:', e)
    await Promise.all([
      appStore.loadSettingsFromServer(),
      userStore.fetchProfile(),
      panelStore.fetchPanel(),
    ])
  }
  // Re-sync layout CSS variables after server settings loaded
  syncLayoutVars()
  await fetchWallpaper()
})

// Theme
function toggleTheme() {
  appStore.setTheme(appStore.isDark ? 'light' : 'dark')
}

// User menu
const userMenuOptions = computed<DropdownOption[]>(() => {
  const items: DropdownOption[] = [
    { label: '个人设置', key: 'profile', icon: () => h(Icon, { icon: 'mdi:account-cog' }) },
    { label: '设置', key: 'settings', icon: () => h(Icon, { icon: 'mdi:cog' }) },
  ]
  if (userStore.user?.role === 'admin') {
    items.push({ label: '管理面板', key: 'admin', icon: () => h(Icon, { icon: 'mdi:shield-cog' }) })
  }
  items.push(
    { type: 'divider', key: 'd1' },
    { label: '退出登录', key: 'logout', icon: () => h(Icon, { icon: 'mdi:logout' }) },
  )
  return items
})

function handleUserMenu(key: string) {
  switch (key) {
    case 'profile':
      router.push('/profile')
      break
    case 'settings': showSettings.value = true; break
    case 'admin':
      router.push('/admin')
      break
    case 'logout': userStore.logout(); return
  }
}

// Group management
const showAddGroup = ref(false)
const showEditGroup = ref(false)
const showSettings = ref(false)
const newGroupName = ref('')
const editGroupName = ref('')
const editingGroupId = ref('')
const modalLoading = ref(false)

function openEditGroup(group: { id: string; name: string }) {
  editingGroupId.value = group.id
  editGroupName.value = group.name
  showEditGroup.value = true
}

async function handleCreateGroup() {
  if (!newGroupName.value.trim()) return false
  modalLoading.value = true
  try {
    await panelStore.createGroup(newGroupName.value.trim())
    newGroupName.value = ''
    showAddGroup.value = false
    message.success('分组已创建')
  } catch (e: any) {
    message.error(e.response?.data?.error || '创建失败')
  } finally {
    modalLoading.value = false
  }
  return true
}

async function handleUpdateGroup() {
  if (!editGroupName.value.trim()) return false
  modalLoading.value = true
  try {
    await panelStore.updateGroup(editingGroupId.value, { name: editGroupName.value.trim() })
    showEditGroup.value = false
    message.success('分组已更新')
  } catch (e: any) {
    message.error(e.response?.data?.error || '更新失败')
  } finally {
    modalLoading.value = false
  }
  return true
}

// Card management
const showCardEditor = ref(false)
const editingCard = ref<Card | null>(null)
const defaultCardGroupId = ref('')

function openAddCard(groupId: string) {
  editingCard.value = null
  defaultCardGroupId.value = groupId
  showCardEditor.value = true
}

function openEditCard(card: Card) {
  editingCard.value = card
  defaultCardGroupId.value = card.group_id
  showCardEditor.value = true
}

async function handleSaveCard(data: any) {
  modalLoading.value = true
  try {
    if (editingCard.value) {
      await panelStore.updateCard(editingCard.value.id, data)
      message.success('卡片已更新')
    } else {
      await panelStore.createCard(data)
      message.success('卡片已创建')
    }
    showCardEditor.value = false
  } catch (e: any) {
    message.error(e.response?.data?.error || '操作失败')
  } finally {
    modalLoading.value = false
  }
}

// Apply color to all cards in group
async function handleApplyColorAll(data: { cardIds: string[]; icon_color: string }) {
  try {
    await panelStore.batchUpdateCardColors({
      cardIds: data.cardIds,
      icon_color: data.icon_color,
    })
    message.success('图标颜色已应用到全部卡片')
  } catch (e: any) {
    message.error('批量更新失败')
  }
}

// Apply background color to all cards in group
async function handleApplyBgColorAll(data: { cardIds: string[]; bg_color: string }) {
  try {
    await panelStore.batchUpdateCardColors({
      cardIds: data.cardIds,
      bg_color: data.bg_color,
    })
    message.success('背景颜色已应用到全部卡片')
  } catch (e: any) {
    message.error('批量更新失败')
  }
}

// Delete
const showDeleteConfirm = ref(false)
const deleteTarget = ref<{ type: 'group' | 'card'; id: string }>({ type: 'group', id: '' })

function confirmDeleteGroup(groupId: string) {
  deleteTarget.value = { type: 'group', id: groupId }
  showDeleteConfirm.value = true
}

function confirmDeleteCard(cardId: string) {
  deleteTarget.value = { type: 'card', id: cardId }
  showDeleteConfirm.value = true
}

async function handleDelete() {
  modalLoading.value = true
  try {
    if (deleteTarget.value.type === 'group') {
      await panelStore.deleteGroup(deleteTarget.value.id)
      message.success('分组已删除')
    } else {
      await panelStore.deleteCard(deleteTarget.value.id)
      message.success('卡片已删除')
    }
    showDeleteConfirm.value = false
  } catch (e: any) {
    message.error(e.response?.data?.error || '删除失败')
  } finally {
    modalLoading.value = false
  }
  return true
}

// Reorder
async function handleReorder(groupId: string, cardIds: string[]) {
  const cardOrders = cardIds.map((id, index) => ({
    id,
    group_id: groupId,
    sort_order: index + 1,
  }))
  const groupOrders = panelStore.groups.map((g, index) => ({
    id: g.id,
    sort_order: index + 1,
  }))
  try {
    await panelStore.reorder(groupOrders, cardOrders)
  } catch (e: any) {
    console.error('Reorder failed:', e)
  }
}

// Open card URL
function openCardUrl(card: Card) {
  window.open(card.url, '_blank')
}
</script>

<style scoped>
/* === Page === */
.home-page {
  min-height: 100vh;
  background: var(--sd-bg-root);
  position: relative;
}

/* === Wallpaper Background === */
.wallpaper-bg {
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
}

.wallpaper-bing,
.wallpaper-custom {
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
}

.wallpaper-gradient {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.wallpaper-copyright {
  position: absolute;
  bottom: 20px;
  right: 20px;
  padding: 8px 16px;
  background: rgba(0, 0, 0, 0.5);
  color: white;
  font-size: 12px;
  border-radius: 6px;
  backdrop-filter: blur(10px);
  z-index: 10;
}

/* === Background orbs === */
.home-bg {
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: 1;
  overflow: hidden;
}

.bg-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  opacity: 0.4;
  animation: floatOrb 20s ease-in-out infinite;
}

[data-theme="dark"] .bg-orb {
  opacity: 0.15;
}

.bg-orb-1 {
  width: 600px;
  height: 600px;
  background: linear-gradient(135deg, #007AFF, #5856D6);
  top: -200px;
  right: -100px;
  animation-delay: 0s;
}

.bg-orb-2 {
  width: 500px;
  height: 500px;
  background: linear-gradient(135deg, #5856D6, #FF2D55);
  bottom: -200px;
  left: -100px;
  animation-delay: -7s;
}

.bg-orb-3 {
  width: 400px;
  height: 400px;
  background: linear-gradient(135deg, #34C759, #007AFF);
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation-delay: -14s;
}

@keyframes floatOrb {
  0%, 100% { transform: translate(0, 0) scale(1); }
  25% { transform: translate(30px, -40px) scale(1.05); }
  50% { transform: translate(-20px, 20px) scale(0.95); }
  75% { transform: translate(40px, 30px) scale(1.02); }
}

/* === Header === */
.home-header {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  height: auto;
  padding: 40px 24px 16px;
  margin-top: var(--sd-content-padding-top, 10%);
  background: transparent;
}

/* Center content */
.header-center-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: white;
}

.header-logo .brand-text {
  font-size: 40px;
  font-weight: 700;
  color: white;
  text-shadow: 0 2px 10px rgba(0,0,0,0.3);
  letter-spacing: -0.5px;
}

.header-divider {
  font-size: 20px;
  color: rgba(255,255,255,0.5);
  margin: 0 4px;
}

.header-clock :deep(*) {
  color: white !important;
}

.header-clock :deep(.clock-time) {
  font-size: 28px !important;
  font-weight: 600 !important;
  color: white !important;
  text-shadow: 0 2px 10px rgba(0,0,0,0.3);
}

.header-clock :deep(.clock-date),
.header-clock :deep(.clock-week) {
  color: rgba(255,255,255,0.8) !important;
  text-shadow: 0 1px 6px rgba(0,0,0,0.3);
}

/* Right actions - fixed to page top-right */
.header-actions {
  position: fixed;
  right: 24px;
  top: 20px;
  z-index: 100;
  display: flex;
  align-items: center;
  gap: 4px;
  background: rgba(0,0,0,0.2);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255,255,255,0.15);
  border-radius: 14px;
  padding: 4px;
}

/* Search bar */
.home-search {
  position: relative;
  z-index: 1;
  display: flex;
  justify-content: center;
  padding: 0 24px 32px;
  max-width: 1100px;
  margin: 0 auto;
}

.header-left {
  display: flex;
  align-items: center;
  gap: var(--sd-space-3);
  flex-shrink: 0;
}

.logo-image {
  width: 26px;
  height: 26px;
  border-radius: 6px;
  object-fit: contain;
}

.brand {
  font-size: 18px;
  font-weight: 700;
  color: var(--sd-text-primary);
  letter-spacing: -0.4px;
}

.header-center {
  flex: 1;
  max-width: 520px;
  margin: 0 var(--sd-space-8);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 2px;
  flex-shrink: 0;
}

/* === Header buttons === */
.hdr-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border: none;
  background: transparent;
  color: rgba(255,255,255,0.8);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s var(--sd-ease);
}

.hdr-btn:hover {
  background: rgba(255,255,255,0.15);
  color: white;
  transform: translateY(-1px);
}

.hdr-btn:active {
  transform: scale(0.95) translateY(0);
}

.hdr-btn.accent {
  background: rgba(255,255,255,0.2);
  color: white;
}

.hdr-btn.accent:hover {
  background: rgba(255,255,255,0.3);
}

.hdr-divider {
  width: 1px;
  height: 18px;
  background: rgba(255,255,255,0.3);
  margin: 0 6px;
}

.avatar-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  background: rgba(255,255,255,0.2);
  border-radius: 50%;
  cursor: pointer;
  transition: all 0.2s ease;
  margin-left: 6px;
}

.avatar-btn span {
  color: white;
  font-size: 13px;
  font-weight: 600;
}

.avatar-btn:hover {
  transform: scale(1.08);
  box-shadow: 0 0 0 3px var(--sd-primary-light, rgba(0,122,255,0.15));
}

/* === Content === */
.home-content {
  position: relative;
  z-index: 1;
  max-width: var(--sd-content-max-width, 80%);
  margin: 0 auto;
  padding: 0 var(--sd-content-padding-x, 5%) var(--sd-content-padding-bottom, 5%);
}

/* === Empty State === */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 120px var(--sd-space-6);
  animation: fadeInUp 0.6s cubic-bezier(0.34, 1.56, 0.64, 1) both;
}

@keyframes fadeInUp {
  from { opacity: 0; transform: translateY(20px); }
  to { opacity: 1; transform: translateY(0); }
}

.empty-icon {
  width: 100px;
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: white;
  border-radius: 24px;
  box-shadow: 0 2px 12px rgba(0,0,0,0.06), 0 0 0 1px rgba(0,0,0,0.03);
  margin-bottom: var(--sd-space-6);
}

:root[data-theme="dark"] .empty-icon {
  background: var(--sd-bg-surface);
}

.empty-state h2 {
  font-size: 22px;
  font-weight: 600;
  color: var(--sd-text-primary);
  margin: 0 0 8px;
}

.empty-state p {
  font-size: 15px;
  color: var(--sd-text-secondary);
  margin: 0 0 28px;
}

.cta-btn {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 0 24px;
  height: 44px;
  background: var(--sd-primary);
  color: white;
  border: none;
  border-radius: 22px;
  font-size: 15px;
  font-weight: 600;
  font-family: var(--sd-font);
  cursor: pointer;
  transition: all 0.2s ease;
}

.cta-btn:hover {
  background: var(--sd-primary-hover);
  transform: translateY(-1px);
  box-shadow: 0 4px 16px rgba(0,122,255,0.3);
}

.cta-btn:active {
  transform: translateY(0);
}

/* === Group Cards === */
.groups-container {
  display: flex;
  flex-direction: column;
  gap: var(--sd-space-6);
}

.group-wrapper {
  animation: fadeInUp 0.5s cubic-bezier(0.34, 1.56, 0.64, 1) both;
}

.group-wrapper:nth-child(2) { animation-delay: 0.05s; }
.group-wrapper:nth-child(3) { animation-delay: 0.1s; }
.group-wrapper:nth-child(4) { animation-delay: 0.15s; }

.group-card {
  background: rgba(255,255,255,0.65);
  backdrop-filter: blur(24px) saturate(180%);
  -webkit-backdrop-filter: blur(24px) saturate(180%);
  border-radius: 20px;
  border: 1px solid rgba(0,0,0,0.06);
  padding: var(--sd-space-5) var(--sd-space-6);
  transition: all 0.3s ease;
  box-shadow: 0 1px 3px rgba(0,0,0,0.03);
}

.group-card.card-transparent {
  background: transparent;
  backdrop-filter: none;
  -webkit-backdrop-filter: none;
  border-color: transparent;
  box-shadow: none;
}

.group-card.card-transparent:hover {
  box-shadow: none;
  border-color: transparent;
}

:root[data-theme="dark"] .group-card.card-transparent:hover {
  box-shadow: none;
  border-color: transparent;
}

:root[data-theme="dark"] .group-card {
  background: rgba(28,28,30,0.7);
  border-color: rgba(255,255,255,0.08);
  box-shadow: 0 1px 3px rgba(0,0,0,0.15);
}

:root[data-theme="dark"] .group-card.card-transparent {
  background: transparent;
  border-color: transparent;
  box-shadow: none;
}

.group-card:hover {
  box-shadow: 0 8px 32px rgba(0,0,0,0.06);
  border-color: var(--sd-primary-light, rgba(0,122,255,0.1));
}

:root[data-theme="dark"] .group-card:hover {
  box-shadow: 0 8px 32px rgba(0,0,0,0.2);
  border-color: var(--sd-primary-medium, rgba(0,122,255,0.15));
}

.group-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--sd-space-4);
}

.group-name {
  display: flex;
  align-items: center;
  gap: 10px;
}

.group-name h3 {
  margin: 0;
  font-size: 17px;
  font-weight: 600;
  color: var(--sd-text-primary);
}

.group-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 22px;
  height: 22px;
  padding: 0 8px;
  background: var(--sd-primary-light, rgba(0,122,255,0.1));
  color: var(--sd-primary);
  border-radius: 11px;
  font-size: 12px;
  font-weight: 600;
}

.group-tools {
  display: flex;
  gap: 2px;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.group-card:hover .group-tools {
  opacity: 1;
}

.tool-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: none;
  background: transparent;
  color: var(--sd-text-tertiary);
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.tool-btn:hover {
  background: rgba(0,0,0,0.05);
  color: var(--sd-text-primary);
}

:root[data-theme="dark"] .tool-btn:hover {
  background: rgba(255,255,255,0.1);
}

.tool-btn.danger:hover {
  background: rgba(255,59,48,0.1);
  color: #FF3B30;
}

/* === Dark Wallpaper - Global Text/Icon Adaptation === */
.dark-wallpaper {
  --sd-text-primary: rgba(255,255,255,0.95);
  --sd-text-secondary: rgba(255,255,255,0.7);
  --sd-text-tertiary: rgba(255,255,255,0.5);
  --sd-border: rgba(255,255,255,0.1);
  --sd-divider: rgba(255,255,255,0.08);
  --sd-glass-bg: rgba(28,28,30,0.55);
  --sd-glass-border: rgba(255,255,255,0.08);
}

/* Group card - dark wallpaper */
.dark-wallpaper .group-card {
  background: rgba(28,28,30,0.45);
  border-color: rgba(255,255,255,0.08);
  box-shadow: 0 1px 3px rgba(0,0,0,0.15);
}

.dark-wallpaper .group-card.card-transparent {
  background: transparent;
  border-color: transparent;
  box-shadow: none;
}

.dark-wallpaper .group-card:hover {
  box-shadow: 0 8px 32px rgba(0,0,0,0.2);
  border-color: rgba(255,255,255,0.15);
}

.dark-wallpaper .group-card.card-transparent:hover {
  box-shadow: none;
  border-color: transparent;
}

.dark-wallpaper .group-badge {
  background: rgba(255,255,255,0.12);
  color: rgba(255,255,255,0.85);
}

.dark-wallpaper .group-head:hover .collapse-icon {
  color: white;
}

.dark-wallpaper .group-action-btn {
  color: rgba(255,255,255,0.5);
}

.dark-wallpaper .group-action-btn:hover {
  background: rgba(255,255,255,0.1);
  color: white;
}

/* Tool buttons */
.dark-wallpaper .tool-btn {
  color: rgba(255,255,255,0.5);
}

.dark-wallpaper .tool-btn:hover {
  background: rgba(255,255,255,0.12);
  color: white;
}

.dark-wallpaper .tool-btn.danger:hover {
  background: rgba(255,59,48,0.25);
  color: #FF6B6B;
}

/* Card items */
.dark-wallpaper .card-item {
  background: rgba(255,255,255,0.08);
  border-color: rgba(255,255,255,0.06);
}

.dark-wallpaper .card-item:hover {
  background: rgba(255,255,255,0.14);
  border-color: rgba(255,255,255,0.12);
  box-shadow: 0 6px 20px rgba(0,0,0,0.2);
}

/* Card add button */
.dark-wallpaper .card-add {
  color: rgba(255,255,255,0.4);
  border-color: rgba(255,255,255,0.12);
}

.dark-wallpaper .card-add:hover {
  color: rgba(255,255,255,0.8);
  border-color: rgba(255,255,255,0.25);
  background: rgba(255,255,255,0.06);
}

/* Hidden items */
.dark-wallpaper .hidden-cards-btn,
.dark-wallpaper .hidden-groups-btn {
  color: rgba(255,255,255,0.4);
}

.dark-wallpaper .hidden-cards-btn:hover,
.dark-wallpaper .hidden-groups-btn:hover {
  background: rgba(255,255,255,0.08);
  color: rgba(255,255,255,0.9);
}

.dark-wallpaper .hidden-group-item {
  border-color: rgba(255,255,255,0.12);
  color: rgba(255,255,255,0.6);
}

.dark-wallpaper .hidden-group-item:hover {
  border-color: rgba(255,255,255,0.3);
  background: rgba(255,255,255,0.06);
  color: white;
}

/* Empty state */
.dark-wallpaper .empty-icon {
  background: rgba(255,255,255,0.08);
}

.dark-wallpaper .empty-icon svg rect {
  fill: rgba(255,255,255,0.6);
}

.dark-wallpaper .cta-btn {
  background: rgba(255,255,255,0.15);
  color: white;
}

.dark-wallpaper .cta-btn:hover {
  background: rgba(255,255,255,0.25);
  box-shadow: 0 4px 16px rgba(0,0,0,0.3);
}

/* Footer */
.dark-wallpaper .home-footer {
  border-top-color: rgba(255,255,255,0.06);
  color: rgba(255,255,255,0.35);
}

/* System status */
.dark-wallpaper .system-status .status-item {
  color: rgba(255,255,255,0.45);
}

.dark-wallpaper .system-status .status-divider {
  background: rgba(255,255,255,0.08);
}

/* SearchBar - dark wallpaper */
.dark-wallpaper :deep(.search-bar) {
  background: transparent;
}

.dark-wallpaper :deep(.search-bar input) {
  color: white;
}

.dark-wallpaper :deep(.search-bar input::placeholder) {
  color: rgba(255,255,255,0.4);
}

/* Naive UI dropdown in dark wallpaper */
.dark-wallpaper :deep(.n-dropdown-menu) {
  background: rgba(40,40,42,0.92) !important;
  backdrop-filter: blur(20px) !important;
  border-color: rgba(255,255,255,0.1) !important;
}

.dark-wallpaper :deep(.n-dropdown-option) {
  color: rgba(255,255,255,0.85) !important;
}

.dark-wallpaper :deep(.n-dropdown-option:hover) {
  background: rgba(255,255,255,0.08) !important;
}

/* === Card Grid === */
.group-cards {
  display: grid;
  gap: 8px;
}

/* === Footer === */
.home-footer {
  margin-top: var(--sd-space-8);
  padding: var(--sd-space-6) 0;
  text-align: center;
  border-top: 1px solid var(--sd-border);
}

.footer-content {
  font-size: 13px;
  color: var(--sd-text-tertiary);
  line-height: 1.6;
}

.copyright-notice {
  margin-top: 12px;
  font-size: 12px;
  color: var(--sd-text-tertiary);
  opacity: 0.6;
}

/* === Collapse === */
.collapse-icon {
  transition: transform 0.2s ease;
  flex-shrink: 0;
  color: var(--sd-text-tertiary);
}

.group-head:hover .collapse-icon {
  color: var(--sd-text-primary);
}

/* === Hidden Groups Bar === */
.hidden-groups-bar {
  margin-top: 4px;
}

.hidden-groups-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  border: none;
  background: transparent;
  color: var(--sd-text-tertiary);
  font-size: 12px;
  font-family: var(--sd-font);
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.hidden-groups-btn:hover {
  background: rgba(0,0,0,0.04);
  color: var(--sd-text-secondary);
}

:root[data-theme="dark"] .hidden-groups-btn:hover {
  background: rgba(255,255,255,0.06);
}

.hidden-groups-list {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  padding: 6px 12px 4px;
}

.hidden-group-item {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 4px 10px;
  border: 1px solid var(--sd-border);
  background: transparent;
  color: var(--sd-text-secondary);
  font-size: 12px;
  font-family: var(--sd-font);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.hidden-group-item:hover {
  border-color: var(--sd-primary);
  color: var(--sd-primary);
  background: var(--sd-primary-light, rgba(0,122,255,0.04));
}

/* === Hidden Cards Bar === */
.group-cards-wrap {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.hidden-cards-bar {
  display: flex;
  justify-content: center;
  padding-top: 2px;
}

.hidden-cards-btn {
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 4px 12px;
  border: none;
  background: transparent;
  color: var(--sd-text-tertiary);
  font-size: 12px;
  font-family: var(--sd-font);
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.hidden-cards-btn:hover {
  background: rgba(0,0,0,0.04);
  color: var(--sd-text-secondary);
}

:root[data-theme="dark"] .hidden-cards-btn:hover {
  background: rgba(255,255,255,0.06);
}

.card-is-hidden {
  opacity: 0.45;
}

/* ==================== Mobile Responsive ==================== */

/* Tablet: <= 768px */
@media (max-width: 768px) {
  .home-header {
    padding: 24px 16px 12px;
    margin-top: calc(8% + 50px);
  }

  .header-center-content {
    gap: 8px;
  }

  .header-logo .brand-text {
    font-size: 24px;
  }

  .header-clock :deep(.clock-time) {
    font-size: 22px !important;
  }

  .header-actions {
    right: 12px;
    top: 16px;
    gap: 2px;
    padding: 3px;
    border-radius: 12px;
  }

  .hdr-btn {
    width: 32px;
    height: 32px;
  }

  .hdr-btn .iconify {
    width: 16px;
    height: 16px;
  }

  .hdr-divider {
    height: 14px;
    margin: 0 3px;
  }

  .avatar-btn {
    width: 30px;
    height: 30px;
    font-size: 12px;
    margin-left: 4px;
  }

  .home-search {
    padding: 0 16px 20px;
  }

  .group-card {
    padding: var(--sd-space-4) var(--sd-space-5);
    border-radius: 16px;
  }

  .group-cards {
    gap: 6px;
  }

  .group-head {
    margin-bottom: var(--sd-space-3);
  }

  .group-name h3 {
    font-size: 15px;
  }

  .group-tools {
    opacity: 1;
  }

  .tool-btn {
    width: 30px;
    height: 30px;
  }

  .empty-state {
    padding: 80px var(--sd-space-4);
  }

  .empty-state h2 {
    font-size: 18px;
  }

  .empty-state p {
    font-size: 14px;
  }
}

/* Phone: <= 480px */
@media (max-width: 480px) {
  .home-header {
    padding: 18px 12px 10px;
    margin-top: calc(6% + 40px);
    flex-direction: column;
    gap: 10px;
  }

  .header-center-content {
    gap: 6px;
  }

  .header-logo .brand-text {
    font-size: 20px;
  }

  .header-divider {
    font-size: 16px;
  }

  .header-clock :deep(.clock-time) {
    font-size: 18px !important;
  }

  .header-clock :deep(.clock-date),
  .header-clock :deep(.clock-week) {
    font-size: 11px !important;
  }

  .header-actions {
    position: relative;
    right: auto;
    top: auto;
    justify-content: center;
  }

  .home-search {
    padding: 0 12px 16px;
  }

  .group-card {
    padding: var(--sd-space-3) var(--sd-space-4);
    border-radius: 14px;
  }

  .group-cards {
    gap: 5px;
  }

  .groups-container {
    gap: var(--sd-space-4);
  }

  .group-name h3 {
    font-size: 14px;
  }

  .group-badge {
    min-width: 18px;
    height: 18px;
    padding: 0 5px;
    font-size: 11px;
  }

  .tool-btn {
    width: 28px;
    height: 28px;
  }

  .tool-btn .iconify {
    width: 14px;
    height: 14px;
  }

  .empty-icon {
    width: 80px;
    height: 80px;
    border-radius: 20px;
  }

  .empty-icon svg {
    width: 40px;
    height: 40px;
  }

  .empty-state {
    padding: 60px var(--sd-space-3);
  }

  .cta-btn {
    height: 40px;
    padding: 0 20px;
    font-size: 14px;
  }

  .hidden-groups-list {
    gap: 4px;
  }

  .hidden-group-item {
    padding: 3px 8px;
    font-size: 11px;
  }
}
</style>