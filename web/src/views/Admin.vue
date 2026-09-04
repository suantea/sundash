<template>
  <div class="admin-page">
    <div class="page-bg">
      <div class="bg-orb bg-orb-1"></div>
      <div class="bg-orb bg-orb-2"></div>
    </div>

    <header class="page-header">
      <button class="back-btn" @click="$router.push('/')">
        <Icon icon="mdi:chevron-left" :width="18" :height="18" />
        <span>{{ t("common.back") }}</span>
      </button>
      <h1>{{ t("admin.title") }}</h1>
      <div style="width: 70px;"></div>
    </header>

    <div class="admin-content">
      <!-- ===== 1. System Settings ===== -->
      <div class="section-header">
        <div>
          <h2>{{ t("admin.systemSettings") }}</h2>
          <p>{{ t("admin.manageRegistrationDefaults") }}</p>
        </div>
      </div>

      <div class="settings-card">
        <div class="setting-row">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(52,199,89,0.1);">
              <Icon icon="mdi:account-plus" :width="18" :height="18" color="#34C759" />
            </div>
            <div>
              <div class="setting-title">{{ t("admin.allowRegistration") }}</div>
              <div class="setting-desc">{{ t("admin.allowRegistrationDesc") }}</div>
            </div>
          </div>
          <n-switch v-model:value="globalSettings.allow_registration" @update:value="saveGlobalSetting('allow_registration', $event)" />
        </div>
        <div class="setting-divider"></div>
        <div class="setting-row">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(255,149,0,0.1);">
              <Icon icon="mdi:shield-check" :width="18" :height="18" color="#FF9500" />
            </div>
            <div>
              <div class="setting-title">{{ t("admin.requireApproval") }}</div>
              <div class="setting-desc">{{ t("admin.requireApprovalDesc") }}</div>
            </div>
          </div>
          <n-switch v-model:value="globalSettings.require_approval" @update:value="saveGlobalSetting('require_approval', $event)" />
        </div>
        <div class="setting-divider"></div>
        <div class="setting-row">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
              <Icon icon="mdi:image" :width="18" :height="18" color="#007AFF" />
            </div>
            <div>
              <div class="setting-title">{{ t("admin.defaultWallpaperType") }}</div>
              <div class="setting-desc">{{ t("admin.defaultWallpaperTypeDesc") }}</div>
            </div>
          </div>
          <div class="seg-control">
            <button :class="['seg-btn', { active: globalSettings.default_wallpaper_type === '' }]" @click="saveGlobalSetting('default_wallpaper_type', '')">{{ t("common.default") }}</button>
            <button :class="['seg-btn', { active: globalSettings.default_wallpaper_type === 'bing' }]" @click="saveGlobalSetting('default_wallpaper_type', 'bing')">{{ t("settings.bingWallpaper") }}</button>
            <button :class="['seg-btn', { active: globalSettings.default_wallpaper_type === 'custom' }]" @click="saveGlobalSetting('default_wallpaper_type', 'custom')">{{ t("common.custom") }}</button>
          </div>
        </div>
        <div v-if="globalSettings.default_wallpaper_type === 'custom'" class="setting-sub">
          <input v-model="globalSettings.default_wallpaper_url" class="bg-url-input" :placeholder="t('admin.enterBackgroundUrl')" @change="saveGlobalSetting('default_wallpaper_url', globalSettings.default_wallpaper_url)" />
        </div>
        <div class="setting-divider"></div>
        <div class="setting-row" style="flex-direction: column; align-items: stretch; gap: 12px;">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(90,200,250,0.1);">
              <Icon icon="mdi:sync" :width="18" :height="18" color="#5AC8FA" />
            </div>
            <div>
              <div class="setting-title">{{ t('admin.bookmarkSync') }}</div>
              <div class="setting-desc">{{ t('admin.bookmarkSyncDesc') }}</div>
            </div>
          </div>
          <div class="setting-toggle-row">
            <span>{{ t('admin.bookmarkSyncEnabled') }}</span>
            <n-switch v-model:value="globalSettings.bmsync_enabled" @update:value="saveGlobalSetting('bmsync_enabled', $event)" />
          </div>
          <div v-if="globalSettings.bmsync_enabled" class="setting-sub">
            <input v-model="globalSettings.bmsync_server_url" class="bg-url-input" :placeholder="t('admin.bmsyncServerUrlPlaceholder')" @change="saveGlobalSetting('bmsync_server_url', globalSettings.bmsync_server_url)" />
            <input v-model="globalSettings.bmsync_token" type="password" class="bg-url-input" :placeholder="t('admin.bmsyncTokenPlaceholder')" @change="saveGlobalSetting('bmsync_token', globalSettings.bmsync_token)" />
          </div>
          <div v-else class="setting-sub disabled-hint">
            {{ t('admin.bookmarkSyncDisabledHint') }}
          </div>
        </div>
        <div class="setting-divider"></div>
        <div class="setting-row" style="flex-direction: column; align-items: stretch; gap: 12px;">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
              <Icon icon="mdi:magnify" :width="18" :height="18" color="#007AFF" />
            </div>
            <div>
              <div class="setting-title">{{ t('admin.searchEngine') }}</div>
              <div class="setting-desc">{{ t('admin.searchEngineDesc') }}</div>
            </div>
          </div>
          <div class="seg-control">
            <button
              v-for="eng in searchEngines"
              :key="eng.name"
              :class="['seg-btn', { active: globalSettings.search_engine === eng.name }]"
              @click="saveGlobalSetting('search_engine', eng.name)"
            >{{ eng.name }}</button>
          </div>
        </div>
      </div>

      <!-- ===== 2. Site Configuration ===== -->
      <div class="section-header">
        <div>
          <h2>{{ t("admin.siteConfig") }}</h2>
          <p>{{ t("admin.manageWebsiteConfig") }}</p>
        </div>
        <button class="add-btn" @click="saveAllSiteConfig" :disabled="siteSaving">
          <Icon icon="mdi:content-save" :width="16" :height="16" />
          {{ siteSaving ? t('common.saving') : t('admin.saveConfig') }}
        </button>
      </div>

      <div class="settings-card">
        <div class="setting-row">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
              <Icon icon="mdi:text-short" :width="18" :height="18" color="#007AFF" />
            </div>
            <div>
              <div class="setting-title">{{ t('settings.siteTitle') }}</div>
              <div class="setting-desc">{{ t('admin.siteTitleDesc') }}</div>
            </div>
          </div>
          <input v-model="siteConfig.site_title" class="mini-text-input" :placeholder="t('admin.siteTitlePlaceholder')" />
        </div>
        <div class="setting-divider"></div>
        <div class="setting-row">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(255,149,0,0.1);">
              <Icon icon="mdi:icon" :width="18" :height="18" color="#FF9500" />
            </div>
            <div>
              <div class="setting-title">{{ t('admin.siteIcon') }}</div>
              <div class="setting-desc">{{ t('admin.siteIconDesc') }}</div>
            </div>
          </div>
          <input v-model="siteConfig.site_icon_url" class="mini-text-input wide" :placeholder="t('admin.siteIconPlaceholder')" />
        </div>
        <div class="setting-divider"></div>
        <div class="setting-row">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(52,199,89,0.1);">
              <Icon icon="mdi:cloud" :width="18" :height="18" color="#34C759" />
            </div>
            <div>
              <div class="setting-title">{{ t('admin.cdnUrl') }}</div>
              <div class="setting-desc">{{ t('admin.cdnUrlDesc') }}</div>
            </div>
          </div>
          <input v-model="siteConfig.site_cdn_url" class="mini-text-input wide" :placeholder="t('admin.cdnPlaceholder')" />
        </div>
        <div class="setting-divider"></div>
        <div class="setting-row" style="flex-direction: column; align-items: stretch; gap: 12px;">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(175,82,222,0.1);">
              <Icon icon="mdi:chart-line" :width="18" :height="18" color="#AF52DE" />
            </div>
            <div>
              <div class="setting-title">{{ t('admin.analyticsCode') }}</div>
              <div class="setting-desc">{{ t('admin.analyticsCodeDesc') }}</div>
            </div>
          </div>
          <textarea v-model="siteConfig.site_analytics_code" class="code-input" :placeholder="t('admin.analyticsCodePlaceholder')" rows="3"></textarea>
        </div>
        <div class="setting-divider"></div>
        <div class="setting-row" style="flex-direction: column; align-items: stretch; gap: 12px;">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
              <Icon icon="mdi:code-tags" :width="18" :height="18" color="#007AFF" />
            </div>
            <div>
              <div class="setting-title">{{ t('admin.customHead') }}</div>
              <div class="setting-desc">{{ t('admin.customHeadDesc') }}</div>
            </div>
          </div>
          <textarea v-model="siteConfig.site_custom_head" class="code-input" :placeholder="t('admin.customHeadPlaceholder')" rows="3"></textarea>
        </div>
        <div class="setting-divider"></div>
        <div class="setting-row" style="flex-direction: column; align-items: stretch; gap: 12px;">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(52,199,89,0.1);">
              <Icon icon="mdi:code-tags" :width="18" :height="18" color="#34C759" />
            </div>
            <div>
              <div class="setting-title">{{ t('admin.customFooter') }}</div>
              <div class="setting-desc">{{ t('admin.customFooterDesc') }}</div>
            </div>
          </div>
          <textarea v-model="siteConfig.site_custom_footer" class="code-input" :placeholder="t('admin.customFooterCodePlaceholder')" rows="3"></textarea>
        </div>
      </div>

      <!-- ===== 3. Pending Users ===== -->
      <div class="section-header">
        <div>
          <h2>{{ t("admin.pendingUsers") }}</h2>
          <p>{{ t("admin.approveNewUsers") }}</p>
        </div>
        <span v-if="pendingUsers.length" class="badge">{{ pendingUsers.length }}</span>
      </div>

      <div class="table-card pending-table">
        <div v-for="user in pendingUsers" :key="user.id" class="table-row">
          <div class="td" style="flex: 2.5;">
            <span class="u-avatar pending">{{ user.username[0].toUpperCase() }}</span>
            <span class="u-name">{{ user.username }}</span>
          </div>
          <div class="td" style="flex: 2; color: var(--sd-text-secondary);">{{ user.display_name || '-' }}</div>
          <div class="td" style="flex: 2; gap: 6px;">
            <button class="row-btn approve" @click="approveUser(user)">
              <Icon icon="mdi:check" :width="14" :height="14" />
              {{ t("admin.approve") }}
            </button>
            <button class="row-btn danger" @click="rejectUser(user)">
              <Icon icon="mdi:close" :width="14" :height="14" />
              {{ t("admin.reject") }}
            </button>
          </div>
        </div>
        <div v-if="pendingUsers.length === 0" class="table-empty">{{ t('admin.noPendingUsers') }}</div>
      </div>

      <!-- ===== 3. All Users ===== -->
      <div class="section-header">
        <div>
          <h2>{{ t("admin.userManagement") }}</h2>
          <p>{{ t("admin.manageAllUsers") }}</p>
        </div>
        <button class="add-btn" @click="showAddUser = true">
          <Icon icon="mdi:plus" :width="16" :height="16" />
          {{ t("admin.addUser") }}
        </button>
      </div>

      <div class="table-card users-table">
        <div class="table-head">
          <div class="th" style="flex: 2.5;">{{ t('admin.user') }}</div>
          <div class="th" style="flex: 2;">{{ t('admin.displayName') }}</div>
          <div class="th" style="flex: 1;">{{ t('user.role') }}</div>
          <div class="th" style="flex: 1;">{{ t('common.status') }}</div>
          <div class="th" style="flex: 2.5;">{{ t('common.actions') }}</div>
        </div>
        <div v-for="user in users" :key="user.id" class="table-row">
          <div class="td" style="flex: 2.5;">
            <span class="u-avatar">{{ user.username[0].toUpperCase() }}</span>
            <span class="u-name">{{ user.username }}</span>
          </div>
          <div class="td" style="flex: 2; color: var(--sd-text-secondary);">{{ user.display_name || '-' }}</div>
          <div class="td" style="flex: 1;">
            <span :class="['role-pill', user.role === 'admin' ? 'is-admin' : 'is-user']">
              {{ user.role === 'admin' ? t('admin.roleAdmin') : t('admin.roleUser') }}
            </span>
          </div>
          <div class="td" style="flex: 1;">
            <span :class="['status-pill', 'is-' + (user.status || 'approved')]">
              {{ statusLabel(user.status) }}
            </span>
          </div>
          <div class="td td-actions">
            <button class="row-btn" @click="toggleRole(user)">
              {{ user.role === 'admin' ? t('admin.setAsUser') : t('admin.setAsAdmin') }}
            </button>
            <button class="row-btn" @click="viewUserPanel(user)">
              <Icon icon="mdi:view-grid" :width="14" :height="14" />
              {{ t("admin.bookmarks") }}
            </button>
            <button class="row-btn danger" @click="deleteUser(user)"> {{ t('common.delete') }}</button>
            <button class="row-btn" @click="openEditDisplayName(user)">
              <Icon icon="mdi:rename-box" :width="14" :height="14" />
              {{ t("admin.rename") }}
            </button>
            <button class="row-btn" @click="openResetPassword(user)" style="color: #FF9500; border-color: rgba(255,149,0,0.3);">
              <Icon icon="mdi:key-variant" :width="14" :height="14" />
              {{ t("admin.resetPassword") }}
            </button>
          </div>
        </div>
        <div v-if="users.length === 0" class="table-empty">{{ t('admin.noUsers') }}</div>
      </div>

      <!-- ===== Backup/Restore Section ===== -->
      <div class="section-header">
        <div>
          <h2>{{ t('admin.backupRestore') || '备份与恢复' }}</h2>
          <p>{{ t('admin.backupRestoreDesc') || '下载完整数据库备份或从备份文件恢复' }}</p>
        </div>
      </div>
      <div class="settings-card">
        <div class="setting-row">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(0,122,255,0.1);">
              <Icon icon="mdi:download" :width="18" :height="18" color="#007AFF" />
            </div>
            <div>
              <div class="setting-title">{{ t('admin.downloadBackup') || '下载备份' }}</div>
              <div class="setting-desc">{{ t('admin.backupDownloadDesc') || '将当前数据库完整备份为 .db 文件下载到本地' }}</div>
            </div>
          </div>
          <button class="add-btn" @click="handleBackup">
            <Icon icon="mdi:download" :width="16" :height="16" />
            {{ t('admin.downloadBackup') || '下载备份' }}
          </button>
        </div>
        <div class="setting-divider"></div>
        <div class="setting-row">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(255,149,0,0.1);">
              <Icon icon="mdi:upload" :width="18" :height="18" color="#FF9500" />
            </div>
            <div>
              <div class="setting-title">{{ t('admin.restoreBackup') || '从备份恢复' }}</div>
              <div class="setting-desc">{{ t('admin.restoreBackupDesc') || '从 .db 或 .sqlite 备份文件恢复，将覆盖现有数据' }}</div>
            </div>
          </div>
          <button class="add-btn add-btn-warning" @click="restoreInput?.click()">
            <Icon icon="mdi:upload" :width="16" :height="16" />
            {{ t('admin.restoreBackup') || '从备份恢复' }}
          </button>
          <input ref="restoreInput" type="file" accept=".db,.sqlite" style="display:none" @change="handleRestore" />
        </div>
      </div>

      <!-- ===== AI Auto-Categorize ===== -->
      <div class="section-header">
        <div>
          <h2>{{ t('admin.autoCategorize') || 'AI 智能归类' }}</h2>
          <p>{{ t('admin.autoCategorizeDesc') || '基于关键词分析，自动将「暂停使用」和「其他」中的书签归类到合适的分组' }}</p>
        </div>
      </div>
      <div class="settings-card">
        <div class="setting-row">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(175,82,222,0.1);">
              <Icon icon="mdi:auto-fix" :width="18" :height="18" color="#AF52DE" />
            </div>
            <div>
              <div class="setting-title">{{ t('admin.startCategorize') || '开始分析' }}</div>
              <div class="setting-desc">{{ t('admin.autoCategorizeDesc') || '基于关键词分析，自动将「暂停使用」和「其他」中的书签归类到合适的分组' }}</div>
            </div>
          </div>
          <button class="add-btn add-btn-ai" @click="handleSuggestCategorize" :disabled="categorizeLoading">
            <Icon icon="mdi:auto-fix" :width="16" :height="16" />
            {{ categorizeLoading ? '分析中...' : (t('admin.startCategorize') || '开始分析') }}
          </button>
        </div>
        <div v-if="categorizeSuggestions.length > 0" class="setting-divider"></div>
        <div v-if="categorizeSuggestions.length > 0" class="setting-row">
          <div class="setting-left">
            <div class="setting-icon" style="background: rgba(52,199,89,0.1);">
              <Icon icon="mdi:check-all" :width="18" :height="18" color="#34C759" />
            </div>
            <div>
              <div class="setting-title">{{ t('admin.applyCategorize') || '应用建议' }} ({{ categorizeSuggestions.length }})</div>
              <div class="setting-desc">{{ t('admin.applyCategorizeDesc') || '确认后将书签移动到建议的分组' }}</div>
            </div>
          </div>
          <button class="add-btn add-btn-success" @click="handleApplyCategorize">
            <Icon icon="mdi:check-all" :width="16" :height="16" />
            {{ t('admin.applyCategorize') || '应用建议' }} ({{ categorizeSuggestions.length }})
          </button>
        </div>
        <div v-if="categorizeSummary" class="categorize-summary">
          <div v-for="(count, group) in categorizeSummary" :key="group" class="summary-pill">
            {{ group }}: {{ count }}
          </div>
        </div>
      </div>
    </div>

    <!-- ===== User Bookmarks Modal ===== -->
    <n-modal v-model:show="showUserPanel" preset="card" :title="`${viewingUser?.username || ''} {{ t('admin.sBookmarks') }}`" style="max-width: 700px;" :bordered="false">
      <div v-if="viewingGroups.length === 0" class="table-empty">{{ t('admin.noBookmarksForUser') }}</div>
      <div v-for="group in viewingGroups" :key="group.id" class="view-group">
        <div class="view-group-name">
          <Icon icon="mdi:folder" :width="16" :height="16" color="#007AFF" />
          {{ group.name }}
          <span class="view-card-count">{{ group.cards?.length || 0 }} {{ t('admin.bookmarkCount') }}</span>
        </div>
        <div v-if="group.cards && group.cards.length" class="view-cards">
          <div v-for="card in group.cards" :key="card.id" class="view-card">
            <span class="view-card-title">{{ card.title }}</span>
            <span class="view-card-url">{{ card.url }}</span>
          </div>
        </div>
      </div>
    </n-modal>

    <!-- ===== Add User Modal ===== -->
    <n-modal v-model:show="showAddUser" preset="dialog" :title="t('admin.addUser')" :positive-text="t('admin.create')" :negative-text="t('common.cancel')"
      :loading="modalLoading" @positive-click="handleCreateUser">
      <n-form label-placement="top">
        <n-form-item :label="t('login.username')">
          <n-input v-model:value="newUser.username" :placeholder="t('admin.enterUsername')" @keyup.enter="handleCreateUser" />
        </n-form-item>
        <n-form-item :label="t('login.password')">
          <n-input v-model:value="newUser.password" type="password" :placeholder="t('admin.enterPassword')" show-password-on="click" @keyup.enter="handleCreateUser" />
        </n-form-item>
        <n-form-item :label="t('admin.displayName')">
          <n-input v-model:value="newUser.displayName" :placeholder="t('admin.enterDisplayName')" @keyup.enter="handleCreateUser" />
        </n-form-item>
      </n-form>
    </n-modal>

    <!-- ===== Reset Password Modal ===== -->
    <n-modal v-model:show="showResetPassword" preset="dialog" :title="t('admin.resetUserPassword')" :positive-text="t('admin.reset')" :negative-text="t('common.cancel')"
      :loading="resetPasswordLoading" @positive-click="handleResetPassword">
      <div style="margin-bottom: 12px; color: var(--sd-text-secondary); font-size: 14px;">
        {{ t("admin.setPasswordForUser") }} <strong>{{ resettingUser?.username }}</strong> {{ t('admin.setNewPassword') }}
      </div>
      <n-input v-model:value="newResetPassword" type="password" :placeholder="t('admin.enterNewPassword')" show-password-on="click" @keyup.enter="handleResetPassword" />
    </n-modal>

    <!-- ===== Edit Display Name Modal ===== -->
    <n-modal v-model:show="showEditDisplayName" preset="dialog" :title="t('admin.editDisplayName')" :positive-text="t('common.save')" :negative-text="t('common.cancel')"
      :loading="editDisplayNameLoading" @positive-click="handleEditDisplayName">
      <div style="margin-bottom: 12px; color: var(--sd-text-secondary); font-size: 14px;">
        {{ t('admin.editUser') }} <strong>{{ editingDisplayNameUser?.username }}</strong> {{ t('admin.sDisplayName') }}
      </div>
      <n-input v-model:value="newDisplayName" :placeholder="t('admin.enterNewDisplayName')" @keyup.enter="handleEditDisplayName" />
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { Icon } from '@iconify/vue'
import { useMessage, NSwitch, NModal, NForm, NFormItem, NInput } from 'naive-ui'
import { api } from '../api'
import type { User, PanelGroup } from '../types'

const { t } = useI18n()

// 搜索引擎列表
const searchEngines = [
  { name: 'Baidu', url: 'https://www.baidu.com/s?wd=' },
  { name: 'Google', url: 'https://www.google.com/search?q=' },
  { name: 'Bing', url: 'https://www.bing.com/search?q=' },
  { name: 'DuckDuckGo', url: 'https://duckduckgo.com/?q=' },
]

const message = useMessage()
const users = ref<User[]>([])
const pendingUsers = ref<User[]>([])
const loading = ref(false)
const showAddUser = ref(false)
const modalLoading = ref(false)
const newUser = ref({ username: '', password: '', displayName: '' })

// User bookmarks viewer
const showUserPanel = ref(false)
const viewingUser = ref<User | null>(null)
const viewingGroups = ref<PanelGroup[]>([])

// Global settings
const globalSettings = reactive({
  allow_registration: false,
  require_approval: false,
  default_wallpaper_type: '',
  default_wallpaper_url: '',
  default_wallpaper_blur: '0',
  default_wallpaper_opacity: '100',
  default_wallpaper_copyright: 'false',
  bmsync_server_url: '',
  bmsync_token: '',
  bmsync_enabled: true,
  search_engine: 'Baidu',
})

// Site configuration
const siteConfig = reactive({
  site_title: '',
  site_icon_url: '',
  site_cdn_url: '',
  site_analytics_code: '',
  site_custom_head: '',
  site_custom_footer: '',
})
const siteSaving = ref(false)

onMounted(() => {
  fetchUsers()
  fetchSettings()
})

function statusLabel(status?: string) {
  if (status === 'pending') return t('admin.statusPending')
  if (status === 'rejected') return t('admin.statusRejected')
  return t('admin.statusApproved')
}

async function fetchUsers() {
  loading.value = true
  try {
    const res = await api.get('users')
    const allUsers = res.data || []
    users.value = (Array.isArray(allUsers) ? allUsers : []).filter((u: User) => u.status !== 'pending')
    pendingUsers.value = (Array.isArray(allUsers) ? allUsers : []).filter((u: User) => u.status === 'pending')
  } catch {
    message.error(t('admin.failedToLoadUsers'))
  } finally {
    loading.value = false
  }
}

async function fetchSettings() {
  try {
    const res = await api.get('admin/settings')
    if (res.data) {
      Object.assign(globalSettings, {
        allow_registration: res.data.allow_registration === 'true',
        require_approval: res.data.require_approval === 'true',
        default_wallpaper_type: res.data.default_wallpaper_type || '',
        default_wallpaper_url: res.data.default_wallpaper_url || '',
        default_wallpaper_blur: res.data.default_wallpaper_blur || '0',
        default_wallpaper_opacity: res.data.default_wallpaper_opacity || '100',
        default_wallpaper_copyright: res.data.default_wallpaper_copyright || 'false',
      })
      globalSettings.bmsync_server_url = res.data.bmsync_server_url || ''
      globalSettings.bmsync_token = res.data.bmsync_token || ''
      globalSettings.bmsync_enabled = res.data.bmsync_enabled !== 'false'
      globalSettings.search_engine = res.data.search_engine || 'Baidu'
      siteConfig.site_title = res.data.site_title || ''
      siteConfig.site_icon_url = res.data.site_icon_url || ''
      siteConfig.site_cdn_url = res.data.site_cdn_url || ''
      siteConfig.site_analytics_code = res.data.site_analytics_code || ''
      siteConfig.site_custom_head = res.data.site_custom_head || ''
      siteConfig.site_custom_footer = res.data.site_custom_footer || ''
    }
  } catch {
    // ignore
  }
}

async function saveGlobalSetting(key: string, value: any) {
  const strValue = typeof value === 'boolean' ? (value ? 'true' : 'false') : String(value)
  try {
    await api.put('admin/settings', { settings: { [key]: strValue } })
    message.success(t('common.settingsSaved'))
  } catch {
    message.error(t('common.saveFailed'))
  }
}

async function saveAllSiteConfig() {
  siteSaving.value = true
  try {
    await api.put('admin/settings', {
      settings: {
        site_title: siteConfig.site_title,
        site_icon_url: siteConfig.site_icon_url,
        site_cdn_url: siteConfig.site_cdn_url,
        site_analytics_code: siteConfig.site_analytics_code,
        site_custom_head: siteConfig.site_custom_head,
        site_custom_footer: siteConfig.site_custom_footer,
      }
    })
    message.success(t('admin.siteConfigSaved'))
  } catch {
    message.error(t('common.saveFailed'))
  } finally {
    siteSaving.value = false
  }
}

async function handleCreateUser() {
  if (!newUser.value.username || !newUser.value.password) return false
  modalLoading.value = true
  try {
    await api.post('auth/register', {
      username: newUser.value.username,
      password: newUser.value.password,
      display_name: newUser.value.displayName,
    })
    message.success(t('admin.userCreated'))
    showAddUser.value = false
    newUser.value = { username: '', password: '', displayName: '' }
    await fetchUsers()
  } catch (e: any) {
    message.error(e.response?.data?.error || t('admin.createFailed'))
  } finally {
    modalLoading.value = false
  }
  return true
}

async function toggleRole(user: User) {
  const newRole = user.role === 'admin' ? 'user' : 'admin'
  try {
    await api.put(`users/${user.id}`, { role: newRole })
    user.role = newRole
    message.success(t('admin.roleUpdated'))
  } catch {
    message.error(t('admin.updateFailed'))
  }
}

async function deleteUser(user: User) {
  try {
    await api.delete(`users/${user.id}`)
    message.success(t('admin.userDeleted'))
    await fetchUsers()
  } catch {
    message.error(t('admin.deleteFailed'))
  }
}

async function approveUser(user: User) {
  try {
    await api.post(`users/${user.id}/approve`)
    message.success(`${user.username} ${t('admin.approved')}`)
    pendingUsers.value = pendingUsers.value.filter(u => u.id !== user.id)
    await fetchUsers()
  } catch (e: any) {
    message.error(e.response?.data?.error || t('admin.approvalFailed'))
  }
}

async function rejectUser(user: User) {
  try {
    await api.post(`users/${user.id}/reject`)
    message.success(`${user.username} ${t('admin.rejected')}`)
    pendingUsers.value = pendingUsers.value.filter(u => u.id !== user.id)
  } catch (e: any) {
    message.error(e.response?.data?.error || t('admin.operationFailed'))
  }
}

async function viewUserPanel(user: User) {
  viewingUser.value = user
  try {
    const res = await api.get(`users/${user.id}/panel`)
    viewingGroups.value = res.data.groups || []
    showUserPanel.value = true
  } catch {
    message.error(t('admin.failedToLoadBookmarks'))
  }
}

// Edit display name
const showEditDisplayName = ref(false)
const editingDisplayNameUser = ref<User | null>(null)
const newDisplayName = ref('')
const editDisplayNameLoading = ref(false)

function openEditDisplayName(user: User) {
  editingDisplayNameUser.value = user
  newDisplayName.value = user.display_name || ''
  showEditDisplayName.value = true
}

async function handleEditDisplayName() {
  if (!editingDisplayNameUser.value) return false
  if (!newDisplayName.value.trim()) {
    message.warning(t('admin.displayNameCannotBeEmpty'))
    return false
  }
  editDisplayNameLoading.value = true
  try {
    await api.put(`users/${editingDisplayNameUser.value.id}`, {
      display_name: newDisplayName.value.trim(),
    })
    message.success(`${editingDisplayNameUser.value.username} ${t('admin.sDisplayNameUpdated')}`)
    showEditDisplayName.value = false
    await fetchUsers()
  } catch (e: any) {
    message.error(e.response?.data?.error || t('admin.updateFailed'))
    return false
  } finally {
    editDisplayNameLoading.value = false
  }
  return true
}

// Reset password
const showResetPassword = ref(false)
const resettingUser = ref<User | null>(null)
const newResetPassword = ref('')
const resetPasswordLoading = ref(false)

function openResetPassword(user: User) {
  resettingUser.value = user
  newResetPassword.value = ''
  showResetPassword.value = true
}

async function handleResetPassword() {
  if (!resettingUser.value) return false
  if (!newResetPassword.value || newResetPassword.value.length < 6) {
    message.warning(t('admin.passwordMin6Chars'))
    return false
  }
  resetPasswordLoading.value = true
  try {
    await api.post(`users/${resettingUser.value.id}/reset-password`, {
      new_password: newResetPassword.value,
    })
    message.success(`${resettingUser.value.username} ${t('admin.sPasswordReset')}`)
    showResetPassword.value = false
  } catch (e: any) {
    message.error(e.response?.data?.error || t('admin.resetFailed'))
    return false
  } finally {
    resetPasswordLoading.value = false
  }
  return true
}

// Backup / Restore
const restoreInput = ref<HTMLInputElement | null>(null)

function handleBackup() {
  const token = localStorage.getItem('sundash-token')
  if (!token) return
  const url = '/api/admin/backup'
  // Create a temporary link with auth header via fetch + blob
  fetch(url, { headers: { Authorization: `Bearer ${token}` } })
    .then(res => {
      if (!res.ok) throw new Error('Backup failed')
      return res.blob()
    })
    .then(blob => {
      const a = document.createElement('a')
      a.href = URL.createObjectURL(blob)
      a.download = `sundash-backup-${new Date().toISOString().slice(0,19).replace(/[T:]/g,'-')}.db`
      a.click()
      URL.revokeObjectURL(a.href)
      message.success(t('admin.backupSuccess') || '备份下载成功')
    })
    .catch(() => message.error(t('admin.backupFailed') || '备份失败'))
}

async function handleRestore(e: Event) {
  const input = e.target as HTMLInputElement
  const file = input.files?.[0]
  if (!file) return
  if (!confirm(t('admin.restoreConfirm') || '恢复将覆盖当前所有数据，确定继续？')) {
    input.value = ''
    return
  }
  const formData = new FormData()
  formData.append('file', file)
  try {
    await api.post('/admin/restore', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
    message.success(t('admin.restoreSuccess') || '恢复成功，请刷新页面')
  } catch (err: any) {
    message.error(err.response?.data?.error || t('admin.restoreFailed') || '恢复失败')
  } finally {
    input.value = ''
  }
}

// AI Auto-Categorize
const categorizeLoading = ref(false)
const categorizeSuggestions = ref<any[]>([])
const categorizeSummary = ref<Record<string, number> | null>(null)

async function handleSuggestCategorize() {
  categorizeLoading.value = true
  categorizeSuggestions.value = []
  categorizeSummary.value = null
  try {
    const res = await api.get('admin/suggest-categorize')
    categorizeSuggestions.value = res.data.suggestions || []
    categorizeSummary.value = res.data.summary || null
    if (categorizeSuggestions.value.length === 0) {
      message.info(t('admin.noSuggestions') || '没有需要归类的书签')
    } else {
      message.success(`发现 ${categorizeSuggestions.value.length} 条归类建议`)
    }
  } catch (e: any) {
    message.error(e.response?.data?.error || t('admin.categorizeFailed') || '分析失败')
  } finally {
    categorizeLoading.value = false
  }
}

async function handleApplyCategorize() {
  if (categorizeSuggestions.value.length === 0) return
  const ops = categorizeSuggestions.value.map((s: any) => ({
    action: 'move_card',
    card_id: s.card_id,
    target_group_id: s.target_id,
  }))
  try {
    await api.post('panels/batch', { operations: ops })
    message.success(`已归类 ${ops.length} 个书签`)
    categorizeSuggestions.value = []
    categorizeSummary.value = null
    await fetchUsers() // refresh
  } catch (e: any) {
    message.error(e.response?.data?.error || t('admin.applyFailed') || '应用失败')
  }
}
</script>

<style scoped>
.admin-page {
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
  position: sticky;
  top: 0;
  z-index: 100;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: var(--sd-header-height);
  padding: 0 var(--sd-space-6);
  background: rgba(255,255,255,0.8);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border-bottom: 1px solid rgba(0,0,0,0.06);
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
.admin-content {
  position: relative;
  z-index: 1;
  max-width: 900px;
  margin: 0 auto;
  padding: var(--sd-space-6) var(--sd-space-6) var(--sd-space-16);
}

.section-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  margin-bottom: var(--sd-space-5);
}

.section-header + .section-header {
  margin-top: var(--sd-space-8);
}

.settings-card + .section-header,
.table-card + .section-header {
  margin-top: var(--sd-space-8);
}

.section-header h2 {
  font-size: 22px;
  font-weight: 600;
  margin: 0 0 4px;
  color: var(--sd-text-primary);
}

.section-header p {
  font-size: 14px;
  color: var(--sd-text-secondary);
  margin: 0;
}

.badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 24px;
  height: 24px;
  padding: 0 8px;
  background: #FF3B30;
  color: white;
  border-radius: 12px;
  font-size: 13px;
  font-weight: 600;
}

.add-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 0 18px;
  height: 38px;
  background: var(--sd-primary);
  color: white;
  border: none;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  font-family: var(--sd-font);
  cursor: pointer;
  transition: all 0.2s ease;
}

.add-btn:hover {
  background: var(--sd-primary-hover);
  transform: translateY(-1px);
}

.add-btn-warning {
  background: var(--sd-warning, #FF9500);
}

.add-btn-warning:hover {
  background: #D97F00;
}

.add-btn-ai {
  background: #AF52DE;
}

.add-btn-ai:hover {
  background: #9636C8;
}

.add-btn-success {
  background: var(--sd-success, #34C759);
}

.add-btn-success:hover {
  background: #28A745;
}

/* Settings Card */
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
  flex-direction: column;
  align-items: stretch;
  gap: 8px;
}

.setting-toggle-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 14px 4px;
  font-size: 13px;
  color: var(--sd-text-secondary);
}

.disabled-hint {
  padding: 8px 14px 14px;
  font-size: 12px;
  color: var(--sd-text-tertiary);
  line-height: 1.6;
}

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

.mini-text-input {
  width: 160px;
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

.mini-text-input.wide {
  width: 260px;
}

.code-input {
  width: 100%;
  padding: 10px;
  border: 1px solid rgba(0,0,0,0.1);
  border-radius: 8px;
  font-size: 12px;
  font-family: var(--sd-font-mono);
  background: var(--sd-bg-card);
  color: var(--sd-text-primary);
  resize: vertical;
  outline: none;
}

.code-input:focus {
  border-color: #007AFF;
}

/* Table */
.table-card {
  background: rgba(255,255,255,0.85);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border-radius: 14px;
  border: 1px solid rgba(0,0,0,0.06);
  overflow: hidden;
}

:root[data-theme="dark"] .table-card {
  background: rgba(28,28,30,0.8);
  border-color: rgba(255,255,255,0.08);
}

.table-head {
  display: flex;
  align-items: center;
  padding: 10px 16px;
  background: rgba(0,0,0,0.03);
  border-bottom: 1px solid rgba(0,0,0,0.06);
}

:root[data-theme="dark"] .table-head {
  background: rgba(255,255,255,0.03);
  border-bottom-color: rgba(255,255,255,0.06);
}

.th {
  font-size: 12px;
  font-weight: 600;
  color: var(--sd-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.3px;
}

.table-row {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid rgba(0,0,0,0.04);
  transition: background 0.15s ease;
}

.table-row:last-child {
  border-bottom: none;
}

.table-row:hover {
  background: rgba(0,0,0,0.02);
}

:root[data-theme="dark"] .table-row:hover {
  background: rgba(255,255,255,0.03);
}

.td {
  font-size: 14px;
  color: var(--sd-text-primary);
  display: flex;
  align-items: center;
  gap: 10px;
}

.u-avatar {
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0,122,255,0.1);
  color: var(--sd-primary);
  border-radius: 8px;
  font-size: 13px;
  font-weight: 600;
  flex-shrink: 0;
}

.u-avatar.pending {
  background: rgba(255,149,0,0.12);
  color: #FF9500;
}

.u-name {
  font-weight: 500;
}

.role-pill, .status-pill {
  display: inline-flex;
  align-items: center;
  padding: 2px 10px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
}

.is-admin {
  background: rgba(255,149,0,0.12);
  color: #FF9500;
}

.is-user {
  background: rgba(0,122,255,0.1);
  color: #007AFF;
}

.status-pill.is-approved {
  background: rgba(52,199,89,0.1);
  color: #34C759;
}

.status-pill.is-pending {
  background: rgba(255,149,0,0.12);
  color: #FF9500;
}

.status-pill.is-rejected {
  background: rgba(255,59,48,0.1);
  color: #FF3B30;
}

.row-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 12px;
  background: transparent;
  border: 1px solid rgba(0,0,0,0.1);
  border-radius: 6px;
  color: var(--sd-text-secondary);
  font-size: 12px;
  font-family: var(--sd-font);
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease;
}

:root[data-theme="dark"] .row-btn {
  border-color: rgba(255,255,255,0.12);
}

.row-btn:hover {
  background: rgba(0,0,0,0.05);
  color: var(--sd-text-primary);
}

.row-btn.approve {
  color: #34C759;
  border-color: rgba(52,199,89,0.3);
}

.row-btn.approve:hover {
  background: rgba(52,199,89,0.08);
}

.row-btn.danger {
  color: #FF3B30;
  border-color: rgba(255,59,48,0.2);
}

.row-btn.danger:hover {
  background: rgba(255,59,48,0.08);
}

.table-empty {
  padding: 48px;
  text-align: center;
  color: var(--sd-text-tertiary);
  font-size: 14px;
}

/* User Bookmarks Viewer */
.view-group {
  margin-bottom: 16px;
}

.view-group-name {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  font-weight: 600;
  color: var(--sd-text-primary);
  margin-bottom: 8px;
}

.view-card-count {
  font-size: 12px;
  font-weight: 400;
  color: var(--sd-text-tertiary);
}

.view-cards {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding-left: 24px;
}

.view-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 6px 10px;
  border-radius: 6px;
  background: rgba(0,0,0,0.02);
}

:root[data-theme="dark"] .view-card {
  background: rgba(255,255,255,0.03);
}

.view-card-title {
  font-size: 13px;
  font-weight: 500;
  color: var(--sd-text-primary);
  min-width: 80px;
}

.view-card-url {
  font-size: 12px;
  color: var(--sd-text-tertiary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}


/* === Responsive === */
@media (max-width: 768px) {
  .admin-content {
    padding: var(--sd-space-4) var(--sd-space-3) var(--sd-space-16);
  }

  .section-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  /* Hide role/status columns on mobile, use card layout (users table only) */
  .users-table .th:nth-child(3),
  .users-table .th:nth-child(4),
  .users-table .td:nth-child(3),
  .users-table .td:nth-child(4) {
    display: none;
  }

  .th, .td {
    font-size: 12px;
  }

  .td-actions {
    gap: 3px;
  }

  .row-btn {
    padding: 3px 6px;
    font-size: 10px;
  }

  .row-btn span {
    display: none;
  }

  .setting-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .mini-text-input,
  .mini-text-input.wide {
    width: 100%;
  }
}

@media (max-width: 480px) {
  .admin-content {
    padding: 12px 10px 80px;
  }

  .admin-page header.page-header {
    padding: 12px 10px;
  }

  .section-header {
    padding: 0 0 8px;
  }

  .section-header h2 {
    font-size: 15px;
  }

  .section-header p {
    font-size: 12px;
  }

  .settings-card {
    border-radius: 12px;
  }

  .setting-row {
    padding: 10px 12px;
  }

  .setting-icon {
    width: 32px;
    height: 32px;
  }

  .setting-icon :deep(svg) {
    width: 16px;
    height: 16px;
  }

  .setting-title {
    font-size: 13px;
  }

  .setting-desc {
    font-size: 11px;
  }

  /* 表格卡片化 */
  .data-table {
    display: block;
    overflow-x: auto;
  }

  .table-head .th,
  .table-row .td {
    min-width: 80px;
    padding: 8px 6px;
    font-size: 11px;
  }

  /* 弹窗窄屏适配 */
  .n-modal {
    width: calc(100vw - 24px) !important;
    max-width: 400px;
  }
}

.categorize-summary {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  padding: 0 14px 14px;
}

.summary-pill {
  padding: 4px 12px;
  border-radius: 12px;
  background: rgba(175, 82, 222, 0.12);
  color: #AF52DE;
  font-size: 12px;
  font-weight: 500;
}
</style>